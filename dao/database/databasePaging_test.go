package database

import (
	"path/filepath"
	"testing"

	"github.com/asdine/storm/v3"
	"github.com/mt1976/frantic-amphora/dao/entities"
)

type pagingTestRecord struct {
	ID     int `storm:"id,increment"`
	Status string
}

func newPagingTestDB(t *testing.T, records, open int) *DB {
	t.Helper()
	conn, err := storm.Open(filepath.Join(t.TempDir(), "paging.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	for i := 0; i < records; i++ {
		status := "CLOSED"
		if i < open {
			status = "OPEN"
		}
		if err := conn.Save(&pagingTestRecord{Status: status}); err != nil {
			t.Fatalf("save: %v", err)
		}
	}
	return &DB{connection: conn, Name: "paging"}
}

func TestTotalPages(t *testing.T) {
	cases := []struct{ total, size, want int }{
		{0, 10, 1}, {1, 10, 1}, {10, 10, 1}, {11, 10, 2}, {95, 10, 10}, {5, 0, 1},
	}
	for _, c := range cases {
		if got := TotalPages(c.total, c.size); got != c.want {
			t.Errorf("TotalPages(%d,%d)=%d want %d", c.total, c.size, got, c.want)
		}
	}
}

func TestGetAllPagedTyped(t *testing.T) {
	db := newPagingTestDB(t, 23, 0)

	recs, info, err := GetAllPagedTyped[pagingTestRecord](db, 2, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 10 || recs[0].ID != 11 || info.TotalPages != 3 || info.TotalRecords != 23 || info.String() != "Page 2 of 3" {
		t.Fatalf("unexpected page 2: len=%d first=%d info=%+v", len(recs), recs[0].ID, info)
	}
	if !info.HasPrevious() || !info.HasNext() {
		t.Fatalf("expected previous and next on page 2: %+v", info)
	}

	recs, info, err = GetAllPagedTyped[pagingTestRecord](db, 3, 10)
	if err != nil || len(recs) != 3 || info.HasNext() {
		t.Fatalf("unexpected last page: len=%d info=%+v err=%v", len(recs), info, err)
	}

	recs, _, err = GetAllPagedTyped[pagingTestRecord](db, 9, 10)
	if err != nil || len(recs) != 0 {
		t.Fatalf("expected empty page beyond range: len=%d err=%v", len(recs), err)
	}

	recs, info, err = GetAllPagedTyped[pagingTestRecord](db, 0, 10)
	if err != nil || info.Page != 1 || recs[0].ID != 1 {
		t.Fatalf("expected page 0 to normalise to 1: info=%+v err=%v", info, err)
	}

	if _, _, err = GetAllPagedTyped[pagingTestRecord](db, 1, 0); err == nil {
		t.Fatal("expected error for page size 0")
	}
}

func TestGetAllPagedWhereTyped(t *testing.T) {
	db := newPagingTestDB(t, 30, 12)

	recs, info, err := GetAllPagedWhereTyped[pagingTestRecord](db, entities.Field("Status"), "OPEN", 3, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 2 || info.TotalRecords != 12 || info.TotalPages != 3 || recs[0].ID != 11 {
		t.Fatalf("unexpected filtered page: len=%d info=%+v", len(recs), info)
	}
}

func TestDBGetAllPaged(t *testing.T) {
	db := newPagingTestDB(t, 23, 12)

	var page []pagingTestRecord
	recs, info, err := db.GetAllPaged(&page, 3, 10)
	if err != nil || len(recs) != 3 || len(page) != 3 || info.TotalPages != 3 {
		t.Fatalf("unexpected: len=%d info=%+v err=%v", len(recs), info, err)
	}

	var filtered []pagingTestRecord
	recs, info, err = db.GetAllPagedWhere(entities.Field("Status"), "OPEN", &filtered, 2, 10)
	if err != nil || len(recs) != 2 || info.TotalRecords != 12 || info.TotalPages != 2 {
		t.Fatalf("unexpected filtered: len=%d info=%+v err=%v", len(recs), info, err)
	}
}
