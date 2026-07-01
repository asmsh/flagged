package tested_options

//go:generate genflagged -type=Options -tests -outFile=tested_options_flagged.go
type Options struct {
	Flag0 bool
	Flag1 bool
	Flag2 bool
}
