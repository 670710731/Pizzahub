package store

import (
	"context"
	"testing"
)

func TestOpenRejectsInvalidDatabaseURL(t *testing.T) {
	if _, err := Open(context.Background(), "not a database URL"); err == nil {
		t.Fatal("expected invalid database URL to be rejected")
	}
}

func TestInvalidDatabaseValueRecognizesPostgresConstraints(t *testing.T) {
	testCases := []struct {
		code string
		want bool
	}{
		{code: "23503", want: true},
		{code: "23514", want: true},
		{code: "22P02", want: true},
		{code: "22003", want: true},
		{code: "08006", want: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.code, func(t *testing.T) {
			got := isInvalidDatabaseValue(postgresError{code: testCase.code})
			if got != testCase.want {
				t.Fatalf("isInvalidDatabaseValue() = %v, want %v", got, testCase.want)
			}
		})
	}
}

type postgresError struct {
	code string
}

func (e postgresError) Error() string    { return e.code }
func (e postgresError) SQLState() string { return e.code }
