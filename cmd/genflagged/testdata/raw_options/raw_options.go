package raw_options

//go:generate genflagged -type=rawOptions -raw -size 64 -outFile=raw_options_flagged.go
type rawOptions struct {
	Flag0 bool
	Flag1 bool
	Flag2 bool
}
