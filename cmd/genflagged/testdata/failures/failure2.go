package failures

//go:generate genflagged -type=nonBoolOptions

type nonBoolOptions struct {
	field1 int
}
