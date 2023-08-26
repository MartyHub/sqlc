package golang

import (
	"reflect"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func TestStruct_Match(t *testing.T) {
	t.Parallel()

	req := &plugin.CodeGenRequest{Catalog: &plugin.Catalog{}}
	tableId := &plugin.Identifier{Name: "Table"}
	tableStruct := Struct{
		Table: tableId,
		Name:  "Name",
		Fields: []Field{
			{Name: "Field", Type: "string"},
		},
	}

	tests := []struct {
		name  string
		str   Struct
		other Struct
		want  bool
	}{
		{
			name: "match",
			str:  tableStruct,
			other: Struct{
				Table: tableId,
				Name:  "Name",
				Fields: []Field{
					{Name: "Field", Type: "string", Column: &plugin.Column{Table: tableId}},
				},
			},
			want: true,
		},
		{
			name: "table mismatch",
			str:  tableStruct,
			other: Struct{
				Table: tableId,
				Name:  "Name",
				Fields: []Field{
					{Name: "Field", Type: "string", Column: &plugin.Column{Table: &plugin.Identifier{Name: "OtherTable"}}},
				},
			},
		},
		{
			name: "other table nil",
			str:  tableStruct,
			other: Struct{
				Table: tableId,
				Name:  "Name",
				Fields: []Field{
					{Name: "Field", Type: "string", Column: &plugin.Column{}},
				},
			},
		},
		{
			name: "field count mismatch",
			str:  tableStruct,
			other: Struct{
				Table: tableId,
				Name:  "Name",
				Fields: []Field{
					{Name: "Field1", Type: "string"},
					{Name: "Field2", Type: "string"},
				},
			},
		},
		{
			name: "field mismatch",
			str:  tableStruct,
			other: Struct{
				Table: tableId,
				Name:  "Name",
				Fields: []Field{
					{Name: "OtherField", Type: "string", Column: &plugin.Column{Table: tableId}},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.str.Match(req, tt.other); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructs_Lookup(t *testing.T) {
	t.Parallel()

	req := &plugin.CodeGenRequest{Catalog: &plugin.Catalog{}}
	str := Struct{
		Fields: []Field{
			{Name: "Field", Type: "string"},
		},
	}
	other := Struct{
		Fields: []Field{
			{Name: "OtherField", Type: "string"},
		},
		Comment: "OtherStruct",
	}
	structs := Structs{str}

	tests := []struct {
		name      string
		other     Struct
		want      Struct
		wantFound bool
	}{
		{
			name: "found",
			other: Struct{
				Fields: []Field{
					{Name: "Field", Type: "string"},
				},
				Comment: "Matching Struct",
			},
			want:      str,
			wantFound: true,
		},
		{
			name:  "not found",
			other: other,
			want:  other,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := structs.Lookup(req, tt.other)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lookup() got = %v, want %v", got, tt.want)
			}

			if found != tt.wantFound {
				t.Errorf("Lookup() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}
