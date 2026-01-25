package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type config struct {
	OutDir     string
	Package    string
	TypeName   string
	TableName  string
	Namespace  string
	Force      bool
	WithWorker bool
	WithImpex  bool
	WithDebug  bool
}

type templateData struct {
	PackageName      string
	TypeName         string
	TableName        string
	Namespace        string
	TableVar         string
	FieldsVar        string
	DomainFields     string            // Field definitions from .definition file
	FieldDefinitions []FieldDefinition // Parsed field definitions for documentation
}

type FieldDefinition struct {
	Name    string
	Type    string
	Tags    string
	Purpose string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.OutDir, "out", ".", "output directory")
	flag.StringVar(&cfg.Package, "pkg", "", "package name (required)")
	flag.StringVar(&cfg.TypeName, "type", "", "type name (required)")
	flag.StringVar(&cfg.TableName, "table", "", "table name (defaults to type)")
	flag.StringVar(&cfg.Namespace, "namespace", "", "cache namespace passed to database.Connect")
	flag.BoolVar(&cfg.Force, "force", false, "overwrite existing files")
	flag.BoolVar(&cfg.WithWorker, "with-worker", true, "generate worker file")
	flag.BoolVar(&cfg.WithImpex, "with-impex", true, "generate import/export file")
	flag.BoolVar(&cfg.WithDebug, "with-debug", true, "generate debug file")
	flag.Parse()

	if cfg.Package == "" {
		exitf("-pkg is required")
	}
	if cfg.TypeName == "" {
		exitf("-type is required")
	}
	if cfg.TableName == "" {
		cfg.TableName = cfg.TypeName
	}
	if cfg.Namespace == "" {
		cfg.Namespace = "main"
	}

	// Read domain fields from .definition file if it exists
	domainFields, fieldNames, fieldInits, fieldDefs := readDefinitionFile(cfg.OutDir, cfg.TypeName)

	data := templateData{
		PackageName:      cfg.Package,
		TypeName:         cfg.TypeName,
		TableName:        cfg.TableName,
		Namespace:        cfg.Namespace,
		TableVar:         "TableName",
		FieldsVar:        "Fields",
		DomainFields:     domainFields,
		FieldDefinitions: fieldDefs,
	}

	// Add custom functions for template
	customFuncs := template.FuncMap{
		"lowerFirst":       lowerFirst,
		"domainFieldNames": func() string { return fieldNames },
		"domainFieldInits": func() string { return fieldInits },
	}

	if err := os.MkdirAll(cfg.OutDir, 0o755); err != nil {
		exitf("creating out dir: %v", err)
	}

	base := lowerFirst(cfg.TypeName)

	files := []struct {
		tmplName string
		outName  string
		enabled  bool
	}{
		{"model.tmpl", base + "Model.go", true},
		{"db.tmpl", base + "DB.go", true},
		{"cache.tmpl", base + "Cache.go", true},
		{"dao.tmpl", base + ".go", true},
		{"internals.tmpl", base + "Internals.go", true},
		{"helpers.tmpl", base + "Helpers.go", true},
		{"worker.tmpl", base + "Worker.go", cfg.WithWorker},
		{"impex.tmpl", base + "Impex.go", cfg.WithImpex},
		{"debug.tmpl", base + "Debug.go", cfg.WithDebug},
		{"readme.tmpl", "README.md", true},
	}

	for _, f := range files {
		if !f.enabled {
			continue
		}
		outPath := filepath.Join(cfg.OutDir, f.outName)

		if !cfg.Force {
			if _, err := os.Stat(outPath); err == nil {
				exitf("refusing to overwrite existing file: %s (use -force)", outPath)
			}
		}

		if err := renderTemplate(templatesFS, f.tmplName, outPath, data, customFuncs); err != nil {
			exitf("generating %s: %v", outPath, err)
		}
	}
}

func readDefinitionFile(outDir, typeName string) (fields, fieldNames, fieldInits string, fieldDefs []FieldDefinition) {
	defPath := filepath.Join(outDir, typeName+".definition")
	file, err := os.Open(defPath)
	if err != nil {
		// If no definition file exists, return empty strings
		return "", "", "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDomainSection := false
	var fieldLines []string
	var namesList []string
	var initsList []string
	var commentBuffer []string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check for start marker
		if strings.Contains(trimmed, "Domain specific fields, starts") {
			inDomainSection = true
			continue
		}

		// Skip if not in domain section
		if !inDomainSection {
			continue
		}

		// Skip empty lines
		if trimmed == "" {
			fieldLines = append(fieldLines, line)
			commentBuffer = nil // Reset comment buffer on empty line
			continue
		}

		// Collect comments as potential purpose text
		if strings.HasPrefix(trimmed, "//") {
			fieldLines = append(fieldLines, line)
			commentText := strings.TrimPrefix(trimmed, "//")
			commentText = strings.TrimSpace(commentText)
			if commentText != "" {
				commentBuffer = append(commentBuffer, commentText)
			}
			continue
		}

		// Parse field definition
		// Expected format: FieldName Type `tags`
		// Extract field name, type, and tags
		fieldName := ""
		fieldType := ""
		fieldTags := ""

		// Find tags (backtick-enclosed string)
		tagStart := strings.Index(line, "`")
		tagEnd := strings.LastIndex(line, "`")
		if tagStart >= 0 && tagEnd > tagStart {
			fieldTags = line[tagStart+1 : tagEnd]
		}

		// Parse field name and type (before tags)
		fieldDef := line
		if tagStart >= 0 {
			fieldDef = strings.TrimSpace(line[:tagStart])
		}

		parts := strings.Fields(fieldDef)
		if len(parts) >= 2 {
			fieldName = strings.TrimSpace(parts[0])
			// Type could be multiple parts (e.g., "time.Time")
			fieldType = strings.TrimSpace(strings.Join(parts[1:], " "))

			if fieldName != "" && !strings.HasPrefix(fieldName, "//") {
				namesList = append(namesList, fieldName)
				initsList = append(initsList, fmt.Sprintf("\t%s: \"%s\",", fieldName, fieldName))

				// Build purpose from comment buffer
				purpose := strings.Join(commentBuffer, " ")

				// Add to field definitions
				fieldDefs = append(fieldDefs, FieldDefinition{
					Name:    fieldName,
					Type:    fieldType,
					Tags:    fieldTags,
					Purpose: purpose,
				})

				// Reset comment buffer after processing field
				commentBuffer = nil
			}
		}

		fieldLines = append(fieldLines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: error reading definition file: %v\n", err)
		return "", "", "", nil
	}

	// Build the strings
	fields = strings.Join(fieldLines, "\n")

	if len(namesList) > 0 {
		var namesBuilder strings.Builder
		for _, name := range namesList {
			namesBuilder.WriteString(fmt.Sprintf("\t%s entities.Field\n", name))
		}
		fieldNames = namesBuilder.String()
		fieldInits = strings.Join(initsList, "\n")
	}

	return fields, fieldNames, fieldInits, fieldDefs
}

func renderTemplate(fsys fs.FS, tmplName string, outPath string, d templateData, funcs template.FuncMap) error {
	tmplPath := filepath.ToSlash(filepath.Join("templates", tmplName))
	b, err := fs.ReadFile(fsys, tmplPath)
	if err != nil {
		return err
	}

	compiled, err := template.New(tmplName).Funcs(funcs).Parse(string(b))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := compiled.Execute(file, d); err != nil {
		return err
	}
	return file.Close()
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = []rune(strings.ToLower(string(r[0])))[0]
	return string(r)
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}
