package mix_options

//go:generate genflagged -type=MixOptions -outFile=mix_options_flagged.go
type MixOptions struct {
	Flag1  bool
	Field2 int
	Field3 string
	Field4 struct {
		Flag2 bool // won't be generated.
	}
	Flag2 bool
}
