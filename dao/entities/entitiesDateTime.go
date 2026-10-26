// This file contains the Date, Time and DateTime entity methods.
//
// Values are stored as time.Time values in UTC. Each type offers three ways to set and get a value:
//   - Set / Get             work with UTC (Set converts any time.Time to UTC).
//   - SetLocal / GetLocal   work with the machine's local timezone.
//   - SetIn / GetIn         work with any IANA timezone name, e.g. "America/New_York".
//
// For DateTime and Time, SetLocal and SetIn treat the wall-clock fields of the supplied
// time.Time as being in that zone, ignoring the time.Time's own location. Time retains the
// reference date as well as the time of day so timezone conversions remain unambiguous.
//
// A Date has no time of day, so it is never shifted between zones: SetLocal and SetIn store the
// calendar date the given instant falls on in that zone, and GetLocal and GetIn return midnight
// of the stored date in that zone, so a date set in a zone reads back unchanged in that zone.

package entities

import (
	"fmt"
	"time"
)

// loadZone resolves an IANA timezone name. An empty name or "UTC" returns time.UTC and
// "Local" returns time.Local.
func loadZone(zone string) (*time.Location, error) {
	if zone == "" {
		return time.UTC, nil
	}
	if zone == "Local" {
		return time.Local, nil
	}
	return time.LoadLocation(zone)
}

// wallClockIn re-interprets the wall-clock fields of t as being in loc.
func wallClockIn(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}

// ---------------------------------------------------------------------------
// DateTime
// ---------------------------------------------------------------------------

// Set stores t converted to UTC and returns the updated DateTime.
func (d *DateTime) Set(t time.Time) DateTime {
	d.Value = t.UTC()
	return *d
}

// SetNow stores the current time and returns the updated DateTime.
func (d *DateTime) SetNow() DateTime {
	return d.Set(time.Now())
}

// SetLocal treats t's wall-clock fields as local time, stores them as UTC and returns the updated DateTime.
func (d *DateTime) SetLocal(t time.Time) DateTime {
	return d.Set(wallClockIn(t, time.Local))
}

// SetIn treats t's wall-clock fields as being in the named timezone, stores them as UTC and
// returns the updated DateTime. The value is unchanged if the zone is invalid.
func (d *DateTime) SetIn(zone string, t time.Time) (DateTime, error) {
	loc, err := loadZone(zone)
	if err != nil {
		return *d, fmt.Errorf("invalid timezone %q: %w", zone, err)
	}
	return d.Set(wallClockIn(t, loc)), nil
}

// Get returns the stored value in UTC. Returns the zero time if empty.
func (d *DateTime) Get() time.Time {
	if d.Value.IsZero() {
		return time.Time{}
	}
	return d.Value.UTC()
}

// GetLocal returns the stored value in the machine's local timezone.
func (d *DateTime) GetLocal() time.Time {
	return d.Get().In(time.Local)
}

// GetIn returns the stored value in the named timezone.
func (d *DateTime) GetIn(zone string) (time.Time, error) {
	loc, err := loadZone(zone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone %q: %w", zone, err)
	}
	return d.Get().In(loc), nil
}

// String returns the stored UTC value.
func (d *DateTime) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Value.UTC().Format(time.RFC3339Nano)
}

// IsZero returns true if no value has been set.
func (d *DateTime) IsZero() bool {
	return d.Value.IsZero()
}

// Equals returns true if both represent the same instant.
func (d *DateTime) Equals(other DateTime) bool {
	return d.Get().Equal(other.Get())
}

// Before returns true if this DateTime is earlier than other.
func (d *DateTime) Before(other DateTime) bool {
	return d.Get().Before(other.Get())
}

// After returns true if this DateTime is later than other.
func (d *DateTime) After(other DateTime) bool {
	return d.Get().After(other.Get())
}

// ---------------------------------------------------------------------------
// Date
// ---------------------------------------------------------------------------

// Set stores the UTC calendar date of the instant t and returns the updated Date.
func (d *Date) Set(t time.Time) Date {
	utc := t.UTC()
	d.Value = time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	return *d
}

// SetNow stores today's UTC date and returns the updated Date.
func (d *Date) SetNow() Date {
	return d.Set(time.Now())
}

// SetLocal stores the calendar date that the instant t falls on in the machine's local
// timezone and returns the updated Date.
func (d *Date) SetLocal(t time.Time) Date {
	d.setDate(t.In(time.Local))
	return *d
}

// SetIn stores the calendar date that the instant t falls on in the named timezone and
// returns the updated Date. The value is unchanged if the zone is invalid.
func (d *Date) SetIn(zone string, t time.Time) (Date, error) {
	loc, err := loadZone(zone)
	if err != nil {
		return *d, fmt.Errorf("invalid timezone %q: %w", zone, err)
	}
	d.setDate(t.In(loc))
	return *d, nil
}

func (d *Date) setDate(t time.Time) {
	d.Value = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// Get returns the stored date as midnight UTC. Returns the zero time if empty.
func (d *Date) Get() time.Time {
	if d.Value.IsZero() {
		return time.Time{}
	}
	return time.Date(d.Value.UTC().Year(), d.Value.UTC().Month(), d.Value.UTC().Day(), 0, 0, 0, 0, time.UTC)
}

// GetLocal returns the stored date as midnight in the machine's local timezone.
func (d *Date) GetLocal() time.Time {
	date := d.Get()
	if date.IsZero() {
		return date
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}

// GetIn returns the stored date as midnight in the named timezone.
func (d *Date) GetIn(zone string) (time.Time, error) {
	loc, err := loadZone(zone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone %q: %w", zone, err)
	}
	date := d.Get()
	if date.IsZero() {
		return date, nil
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), nil
}

// String returns the stored UTC date.
func (d *Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Get().Format("2006-01-02")
}

// IsZero returns true if no value has been set.
func (d *Date) IsZero() bool {
	return d.Value.IsZero()
}

// Equals returns true if both represent the same date.
func (d *Date) Equals(other Date) bool {
	return d.Get().Equal(other.Get())
}

// Before returns true if this Date is earlier than other.
func (d *Date) Before(other Date) bool {
	return d.Get().Before(other.Get())
}

// After returns true if this Date is later than other.
func (d *Date) After(other Date) bool {
	return d.Get().After(other.Get())
}

// ---------------------------------------------------------------------------
// Time
// ---------------------------------------------------------------------------

// Set stores the UTC time of day of t and returns the updated Time.
func (t *Time) Set(in time.Time) Time {
	t.Value = in.UTC()
	return *t
}

// SetNow stores the current UTC time of day and returns the updated Time.
func (t *Time) SetNow() Time {
	return t.Set(time.Now())
}

// SetLocal treats in's wall-clock fields as local time, stores the result in UTC and returns
// the updated Time.
func (t *Time) SetLocal(in time.Time) Time {
	return t.Set(wallClockIn(in, time.Local))
}

// SetIn treats in's time of day as being in the named timezone, stores it as UTC and returns
// the updated Time. The value is unchanged if the zone is invalid.
func (t *Time) SetIn(zone string, in time.Time) (Time, error) {
	loc, err := loadZone(zone)
	if err != nil {
		return *t, fmt.Errorf("invalid timezone %q: %w", zone, err)
	}
	return t.Set(wallClockIn(in, loc)), nil
}

// Get returns the stored time in UTC, including its reference date.
func (t *Time) Get() time.Time {
	if t.Value.IsZero() {
		return time.Time{}
	}
	return t.Value.UTC()
}

// GetLocal returns the stored time in the machine's local timezone.
func (t *Time) GetLocal() time.Time {
	return t.Get().In(time.Local)
}

// GetIn returns the stored time in the named timezone.
func (t *Time) GetIn(zone string) (time.Time, error) {
	loc, err := loadZone(zone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone %q: %w", zone, err)
	}
	return t.Get().In(loc), nil
}

// String returns the stored UTC time of day.
func (t *Time) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Value.UTC().Format("15:04:05.999999999")
}

// IsZero returns true if no value has been set.
func (t *Time) IsZero() bool {
	return t.Value.IsZero()
}

// Equals returns true if both represent the same UTC time of day.
func (t *Time) Equals(other Time) bool {
	return compareTimeOfDay(t.Get(), other.Get()) == 0
}

// Before returns true if this Time is earlier in the day than other.
func (t *Time) Before(other Time) bool {
	return compareTimeOfDay(t.Get(), other.Get()) < 0
}

// After returns true if this Time is later in the day than other.
func (t *Time) After(other Time) bool {
	return compareTimeOfDay(t.Get(), other.Get()) > 0
}

func compareTimeOfDay(a, b time.Time) int {
	a, b = a.UTC(), b.UTC()
	for _, pair := range [][2]int{
		{a.Hour(), b.Hour()},
		{a.Minute(), b.Minute()},
		{a.Second(), b.Second()},
		{a.Nanosecond(), b.Nanosecond()},
	} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	return 0
}
