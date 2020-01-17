package storage

import (
	"database/sql"
)

// StringToNull converts a string to nullable string, setting null if the string is empty.
func StringToNull(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{Valid: false}
	} else {
		return sql.NullString{String: s, Valid: true}
	}
}

// Int64ToNull converts an int64 to nullable int64, setting null if the integer value is zero.
func Int64ToNull(n int64) sql.NullInt64 {
	if n == 0 {
		return sql.NullInt64{Valid: false}
	} else {
		return sql.NullInt64{Int64: n, Valid: true}
	}
}

// NullToString converts a nullable string to a standard string, returning an empty string if value is null.
func NullToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	} else {
		return ""
	}
}

// NullToInt64 converts a nullable int to a standard int64, returning 0 if value is null.
func NullToInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	} else {
		return 0
	}
}
