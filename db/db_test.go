package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUndefinedTable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "undefined table", err: &pgconn.PgError{Code: "42P01"}, want: true},
		{name: "other postgres error", err: &pgconn.PgError{Code: "23505"}},
		{name: "ordinary error", err: errors.New("boom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUndefinedTable(tt.err); got != tt.want {
				t.Fatalf("isUndefinedTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
