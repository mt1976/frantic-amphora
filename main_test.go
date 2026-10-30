package main

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mt1976/frantic-amphora/dao/database"
	"github.com/mt1976/frantic-amphora/dao/entities"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV9"
)

func TestGeneratedEntityTypes(t *testing.T) {
	config, err := os.ReadFile(filepath.Join("data", "config", "common.toml"))
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(t.TempDir())
	for _, directory := range []string{"database", "config", "logs", "dumps", "defaults"} {
		if err := os.MkdirAll(filepath.Join("data", directory), 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join("data", "config", "common.toml"), config, 0600); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	templateStoreV9.Initialise(ctx, false)
	t.Cleanup(func() {
		connections, err := templateStoreV9.GetDatabaseConnections()()
		if err != nil {
			t.Error(err)
			return
		}
		for _, connection := range connections {
			connection.Disconnect()
		}
	})
	records, err := createEntityTypeRecords(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sample := reflect.ValueOf(entityTypeSample()).Elem()
	var tested int
	for _, field := range reflect.VisibleFields(sample.Type()) {
		if !strings.HasPrefix(field.Name, "Entity") {
			continue
		}
		tested++
		t.Run(field.Type.Name(), func(t *testing.T) {
			t.Run("populated", func(t *testing.T) {
				if err := verifyEntityType(entities.Field(field.Name), sample.FieldByIndex(field.Index).Interface(), records[:5]); err != nil {
					t.Fatal(err)
				}
			})
			t.Run("zero", func(t *testing.T) {
				if err := verifyEntityType(entities.Field(field.Name), reflect.Zero(field.Type).Interface(), records[5:]); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
	if tested != 24 {
		t.Fatalf("tested %d entity types, want 24", tested)
	}
	t.Run("all-pages", func(t *testing.T) {
		failures := verifyPages("Entity", "ALL", len(records), 4, 2, recordIDs(records), func(page int) ([]*templateStoreV9.TemplateStoreV9, database.PageInfo, error) {
			return templateStoreV9.GetAllPaged(page, 2)
		})
		if failures != 0 {
			t.Fatalf("unfiltered paging: %d failure(s)", failures)
		}
	})
	t.Run("false-booleans", func(t *testing.T) {
		for _, record := range records[5:] {
			record.EntityBool.Set(false)
			record.EntityStormBool.Set(false)
			if err := record.Update(ctx, "Test false boolean values"); err != nil {
				t.Fatal(err)
			}
		}
		for _, field := range []entities.Field{templateStoreV9.Fields.EntityBool, templateStoreV9.Fields.EntityStormBool} {
			t.Run(field.String(), func(t *testing.T) {
				value := reflect.ValueOf(records[5]).Elem().FieldByName(field.String()).Interface()
				if err := verifyEntityType(field, value, records[5:]); err != nil {
					t.Fatal(err)
				}
			})
		}
	})
}
