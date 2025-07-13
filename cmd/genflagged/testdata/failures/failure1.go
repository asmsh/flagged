// Package failures contains all generate attempts here will fail.
package failures

//go:generate genflagged -type=emptyOptions

type emptyOptions struct{}
