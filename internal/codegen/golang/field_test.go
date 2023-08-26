package golang

import (
	"testing"
)

func TestField_Match(t *testing.T) {
	t.Parallel()

	field := Field{
		Name: "Name",
		Type: "string",
	}

	tests := []struct {
		name  string
		field Field
		want  bool
	}{
		{
			name: "match",
			field: Field{
				Name: "Name",
				Type: "string",
			},
			want: true,
		},
		{
			name: "name mismatch",
			field: Field{
				Name: "OtherName",
				Type: "string",
			},
		},
		{
			name: "type mismatch",
			field: Field{
				Name: "Name",
				Type: "int",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := field.Match(tt.field); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
