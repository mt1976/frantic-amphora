// This file contains the boolean entity methods.

package entities

import (
	"strings"
)

// Set sets the StormBool value from a native Go bool.
func (sb *StormBool) Set(b bool) {
	if b {
		sb.Value = constTrue
	} else {
		sb.Value = constFalse
	}
}

// Bool returns the StormBool as a native Go bool.
func (sb *StormBool) Bool() bool {
	return sb.Value == constTrue
}

// String returns the StormBool as a string ("true" or "false").
func (sb *StormBool) String() string {
	return sb.Value
}

// IsTrue returns true if the StormBool represents a true value.
func (sb *StormBool) IsTrue() bool {
	return sb.Bool()
}

// IsFalse returns true if the StormBool represents a false value.
func (sb *StormBool) IsFalse() bool {
	return !sb.Bool()
}

// Set sets the Bool value from a native Go bool.
func (b *Bool) Set(in bool) {
	if in == true {
		b.Value = constTrue
	} else {
		b.Value = constFalse
	}
}

// SetTrue sets the Bool value to true.
func (b *Bool) SetTrue() {
	b.Value = constTrue
}

// SetFalse sets the Bool value to false.
func (b *Bool) SetFalse() {
	b.Value = constFalse
}

// SetFromString sets the Bool value by parsing a string.
// Recognizes "true", "1", "yes", "y", "t" (case-insensitive) as true, everything else as false.
func (b *Bool) SetFromString(in string) {
	in = strings.ToLower(strings.TrimSpace(in))
	if in == "true" || in == "1" || in == "yes" || in == "y" || in == "t" {
		b.Value = constTrue
	} else {
		b.Value = constFalse
	}
}

// Toggle inverts the Bool value (true becomes false, false becomes true).
func (b *Bool) Toggle() {
	if b.Bool() {
		b.Value = constFalse
	} else {
		b.Value = constTrue
	}
}

// HtmlChecked returns "checked" if the Bool is true, otherwise it returns an empty string
func (b *Bool) HtmlChecked() string {
	if b.Bool() {
		return "checked" // Checked
	}
	return "" // Not Checked
}

// HtmlSelected returns "selected" if the Bool is true, otherwise returns an empty string.
// Useful for HTML select elements.
func (b *Bool) HtmlSelected() string {
	if b.Bool() {
		return "selected" // Selected
	}
	return "" // Not Selected
}

// HtmlDisabled returns "disabled" if the Bool is true, otherwise returns an empty string.
// Useful for HTML form elements.
func (b *Bool) HtmlDisabled() string {
	if b.Bool() {
		return "disabled" // Disabled
	}
	return "" // Not Disabled
}

// HtmlReadOnly returns "readonly" if the Bool is true, otherwise returns an empty string.
// Useful for HTML form elements.
func (b *Bool) HtmlReadOnly() string {
	if b.Bool() {
		return "readonly" // ReadOnly
	}
	return "" // Not ReadOnly
}

// Bool converts the Bool to a native Go bool.
// Returns false if the value is empty, otherwise true only if value equals "true".
func (b *Bool) Bool() bool {
	if b.Value == "" {
		// Needs reversing based on current implementation
		return false
	}
	return b.Value == constTrue
}

// Get returns the Bool value as a native Go bool.
func (b *Bool) Get() bool {
	return b.Bool()
}

// String returns the Bool as a string ("true" or "false").
func (b *Bool) String() string {
	return b.Value
}

// IsTrue returns true if the Bool represents a true value.
func (b *Bool) IsTrue() bool {
	return b.Bool()
}

// IsFalse returns true if the Bool represents a false value.
func (b *Bool) IsFalse() bool {
	return !b.Bool()
}
