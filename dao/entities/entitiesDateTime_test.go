package entities

import (
	"testing"
	"time"
)

func TestDateTimeStoresUTC(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")
	var d DateTime
	input := time.Date(2026, 1, 15, 9, 30, 0, 123, ny)
	d.Set(input)
	if !d.Value.Equal(input.UTC()) || d.Value.Location() != time.UTC {
		t.Fatalf("got %v", d.Value)
	}
	if d.String() != "2026-01-15T14:30:00.000000123Z" {
		t.Fatalf("string got %q", d.String())
	}

	got, err := d.GetIn("America/New_York")
	if err != nil || got.Hour() != 9 {
		t.Fatalf("GetIn: %v %v", got, err)
	}
	if got, _ := d.GetIn("Asia/Tokyo"); got.Hour() != 23 {
		t.Fatalf("Tokyo hour %d", got.Hour())
	}
}

func TestDateTimeSetIn(t *testing.T) {
	var d DateTime
	if _, err := d.SetIn("Asia/Tokyo", time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	expected := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	if !d.Value.Equal(expected) || d.Value.Location() != time.UTC {
		t.Fatalf("got %v", d.Value)
	}
	if _, err := d.SetIn("Not/AZone", time.Now()); err == nil {
		t.Fatal("expected error for bad zone")
	}
	if !d.Value.Equal(expected) {
		t.Fatal("value changed on error")
	}
}

func TestDateTimeLocalRoundTrip(t *testing.T) {
	var d DateTime
	in := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	d.SetLocal(in)
	got := d.GetLocal()
	if got.Hour() != 12 || got.Day() != 1 {
		t.Fatalf("round trip: %v", got)
	}
}

func TestDateSetIn(t *testing.T) {
	var d Date
	// 20:00 UTC on 9 March is already 10 March in Tokyo.
	if _, err := d.SetIn("Asia/Tokyo", time.Date(2026, 3, 9, 20, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if !d.Value.Equal(time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("got %q", d.Value)
	}
	got, _ := d.GetIn("Asia/Tokyo")
	if got.Day() != 10 || got.Hour() != 0 || got.Location().String() != "Asia/Tokyo" {
		t.Fatalf("Tokyo date %v", got)
	}
}

func TestTimeSetIn(t *testing.T) {
	var tm Time
	if _, err := tm.SetIn("Asia/Tokyo", time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	expected := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if !tm.Value.Equal(expected) || tm.Value.Location() != time.UTC {
		t.Fatalf("got %v", tm.Value)
	}
	got, _ := tm.GetIn("Asia/Tokyo")
	if got.Hour() != 9 {
		t.Fatalf("Tokyo hour %d", got.Hour())
	}
	if !tm.Before(Time{Value: time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)}) {
		t.Fatal("Before")
	}
	if !tm.Equals(Time{Value: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)}) {
		t.Fatal("Equals should compare the UTC time of day rather than the reference date")
	}
}

func TestEmptyValues(t *testing.T) {
	var d DateTime
	var dt Date
	var tm Time
	if !d.Get().IsZero() || !dt.Get().IsZero() || !tm.Get().IsZero() || !d.IsZero() || !dt.IsZero() || !tm.IsZero() {
		t.Fatal("empty values should yield zero time")
	}
}
