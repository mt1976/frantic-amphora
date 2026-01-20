package main

import (
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
	PackageName string
	TypeName    string
	TableName   string
	Namespace   string
	TableVar    string
	FieldsVar   string
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
		cfg.Namespace = "cheeseOnToast"
	}

	data := templateData{
		PackageName: cfg.Package,
		TypeName:    cfg.TypeName,
		TableName:   cfg.TableName,
		Namespace:   cfg.Namespace,
		TableVar:    "TableName",
		FieldsVar:   "Fields",
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
		{"new.tmpl", base + "New.go", true},
		{"internals.tmpl", base + "Internals.go", true},
		{"helpers.tmpl", base + "Helpers.go", true},
		{"deprecated.tmpl", base + "Deprecated.go", true},
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

		// For key files where developers often add custom logic (model/new/helpers),
		// never overwrite in-place. Instead rotate the existing file to .OLD and warn.
		if shouldRotateExisting(outPath) {
			if _, err := os.Stat(outPath); err == nil {
				backupPath, backupErr := rotateToOld(outPath)
				if backupErr != nil {
					exitf("failed to rotate existing file %s: %v", outPath, backupErr)
				}
				fmt.Fprintf(os.Stderr, "WARNING: existing file renamed to %s\n", backupPath)
				fmt.Fprintf(os.Stderr, "WARNING: you must manually migrate any custom code from %s into the regenerated file\n", backupPath)
			}
		} else if !cfg.Force {
			if _, err := os.Stat(outPath); err == nil {
				exitf("refusing to overwrite existing file: %s (use -force)", outPath)
			}
		}

		if err := renderTemplate(templatesFS, f.tmplName, outPath, data); err != nil {
			exitf("generating %s: %v", outPath, err)
		}
	}
}

func shouldRotateExisting(outPath string) bool {
	lower := strings.ToLower(outPath)
	return strings.HasSuffix(lower, "model.go") || strings.HasSuffix(lower, "helpers.go") || strings.HasSuffix(lower, "new.go")
}

func rotateToOld(outPath string) (string, error) {
	// Requested behaviour: rename from .go to .OLD.
	backupBase := strings.TrimSuffix(outPath, filepath.Ext(outPath)) + ".OLD"
	backupPath := backupBase
	for i := 1; ; i++ {
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		backupPath = fmt.Sprintf("%s%d", backupBase, i)
	}
	if err := os.Rename(outPath, backupPath); err != nil {
		return "", err
	}
	return backupPath, nil
}

func renderTemplate(fsys fs.FS, tmplName string, outPath string, d templateData) error {
	tmplPath := filepath.ToSlash(filepath.Join("templates", tmplName))
	b, err := fs.ReadFile(fsys, tmplPath)
	if err != nil {
		return err
	}

	funcs := template.FuncMap{
		"lowerFirst": lowerFirst,
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
