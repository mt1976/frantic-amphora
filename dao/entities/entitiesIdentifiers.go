// This file contains the identifier entity methods (Field, Table, KEY, URI, ACTION).

package entities

import (
	"fmt"
)

// String returns the Field as a string.
func (f Field) String() string {
	return string(f)
}

// String returns the Table name as a string.
func (t *Table) String() string {
	return fmt.Sprintf("%v", *t)
}

// String returns the URI as a string.
func (u *URI) String() string {
	return string(*u)
}

// Set sets the URI value from a string.
func (u *URI) Set(in string) {
	*u = URI(in)
}

// Get returns the URI as a string.
func (u *URI) Get() string {
	return u.String()
}

// String returns the primary key as a string.
func (p *KEY) String() string {
	return string(*p)
}

// Set sets the primary key value from a string.
func (p *KEY) Set(in string) {
	*p = KEY(in)
}

// Get returns the primary key as a string.
func (p *KEY) Get() string {
	return p.String()
}

// Key returns the primary key wrapped in curly braces.
// Example: if the key is "123", this returns "{123}".
func (p *KEY) Key() string {
	return "{" + p.String() + "}"
}

// String returns the action as a string.
func (a *ACTION) String() string {
	return string(*a)
}

// Set sets the action value from a string.
func (a *ACTION) Set(in string) {
	*a = ACTION(in)
}

// Get returns the action as a string.
func (a *ACTION) Get() string {
	return a.String()
}
