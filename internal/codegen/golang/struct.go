package golang

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type Struct struct {
	Table   *plugin.Identifier
	Name    string
	Fields  []Field
	Comment string
}

func (s Struct) Match(req *plugin.CodeGenRequest, other Struct) bool {
	if len(s.Fields) != len(other.Fields) {
		return false
	}

	for i, f := range s.Fields {
		of := other.Fields[i]

		if !f.Match(of) {
			return false
		}

		if s.Table != nil && !sdk.SameTableName(of.Column.Table, s.Table, req.Catalog.DefaultSchema) {
			return false
		}
	}

	return true
}

type Structs []Struct

// Lookup search for a matching Struct in slice.
//
//   - if found, returns the matching Struct and true
//   - else returns the given Struct and false
func (s Structs) Lookup(req *plugin.CodeGenRequest, other Struct) (Struct, bool) {
	for _, exists := range s {
		if exists.Match(req, other) {
			return exists, true
		}
	}

	return other, false
}

func StructName(name string, settings *plugin.Settings) string {
	if rename := settings.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		if unicode.IsDigit(r) {
			return r
		}
		return rune('_')
	}, name)

	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}

	// If a name has a digit as its first char, prepand an underscore to make it a valid Go name.
	r, _ := utf8.DecodeRuneInString(out)
	if unicode.IsDigit(r) {
		return "_" + out
	} else {
		return out
	}
}
