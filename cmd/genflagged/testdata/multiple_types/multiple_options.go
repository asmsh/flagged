package multiple_types

//go:generate genflagged -type=options,MaxOptions -size=32 -outType=OptionsBitFlags,_ -outFile=multiple_options_flagged.go

type options struct {
	Flag0 bool
	Flag1 bool
	Flag2 bool
	Flag3 bool
	Flag4 bool
	Flag5 bool
}

type MaxOptions struct {
	Flag0  bool
	Flag1  bool
	Flag2  bool
	Flag3  bool
	Flag4  bool
	Flag5  bool
	Flag6  bool
	Flag7  bool
	Flag8  bool
	Flag9  bool
	Flag10 bool
	Flag11 bool
	Flag12 bool
	Flag13 bool
	Flag14 bool
	Flag15 bool
	Flag16 bool
	Flag17 bool
	Flag18 bool
	Flag19 bool
	Flag20 bool
	Flag21 bool
	Flag22 bool
	Flag23 bool
	Flag24 bool
	Flag25 bool
}
