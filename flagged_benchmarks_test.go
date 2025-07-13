package flagged

import (
	"fmt"
	"testing"
)

func helperRunBenchmarkIs[T bitFlags, TP ptrBitFlags[T]](b *testing.B) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testCase struct {
		name     string
		initial  T
		bitIndex int
	}
	tests := []testCase{
		{
			name:     "zero - range start",
			initial:  zero,
			bitIndex: 0,
		},
		{
			name:     "zero - range end",
			initial:  zero,
			bitIndex: size - 1,
		},
		{
			name:     "allset - range start",
			initial:  allset,
			bitIndex: 0,
		},
		{
			name:     "allset - range end",
			initial:  allset,
			bitIndex: size - 1,
		},
	}
	b.Run(fmt.Sprintf("%T", zero), func(b *testing.B) {
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				var f TP = &tt.initial
				for i := 0; i < b.N; i++ {
					f.Is(tt.bitIndex)
				}
			})
		}
	})
}

func BenchmarkBitFlags_Is(b *testing.B) {
	helperRunBenchmarkIs[BitFlags8](b)
	helperRunBenchmarkIs[BitFlags16](b)
	helperRunBenchmarkIs[BitFlags32](b)
	helperRunBenchmarkIs[BitFlags64](b)
}

func helperRunBenchmarkSet[T bitFlags, TP ptrBitFlags[T]](b *testing.B) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testCase struct {
		name     string
		initial  T
		bitIndex int
	}
	tests := []testCase{
		{
			name:     "zero - range start",
			initial:  zero,
			bitIndex: 0,
		},
		{
			name:     "zero - range end",
			initial:  zero,
			bitIndex: size - 1,
		},
		{
			name:     "allset - range start",
			initial:  allset,
			bitIndex: 0,
		},
		{
			name:     "allset - range end",
			initial:  allset,
			bitIndex: size - 1,
		},
	}
	b.Run(fmt.Sprintf("%T", zero), func(b *testing.B) {
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				var f TP = &tt.initial
				for i := 0; i < b.N; i++ {
					f.Set(tt.bitIndex)
				}
			})
		}
	})
}

func BenchmarkBitFlags_Set(b *testing.B) {
	helperRunBenchmarkSet[BitFlags8](b)
	helperRunBenchmarkSet[BitFlags16](b)
	helperRunBenchmarkSet[BitFlags32](b)
	helperRunBenchmarkSet[BitFlags64](b)
}

func helperRunBenchmarkToggle[T bitFlags, TP ptrBitFlags[T]](b *testing.B) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testCase struct {
		name     string
		initial  T
		bitIndex int
	}
	tests := []testCase{
		{
			name:     "zero - range start",
			initial:  zero,
			bitIndex: 0,
		},
		{
			name:     "zero - range end",
			initial:  zero,
			bitIndex: size - 1,
		},
		{
			name:     "allset - range start",
			initial:  allset,
			bitIndex: 0,
		},
		{
			name:     "allset - range end",
			initial:  allset,
			bitIndex: size - 1,
		},
	}
	b.Run(fmt.Sprintf("%T", zero), func(b *testing.B) {
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				var f TP = &tt.initial
				for i := 0; i < b.N; i++ {
					f.Toggle(tt.bitIndex)
				}
			})
		}
	})
}

func BenchmarkBitFlags_Toggle(b *testing.B) {
	helperRunBenchmarkToggle[BitFlags8](b)
	helperRunBenchmarkToggle[BitFlags16](b)
	helperRunBenchmarkToggle[BitFlags32](b)
	helperRunBenchmarkToggle[BitFlags64](b)
}

func helperRunBenchmarkAnyOf[T bitFlags, TP ptrBitFlags[T]](b *testing.B) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testCase struct {
		name     string
		initial  T
		bitIndex []int
	}
	tests := []testCase{
		{
			name:     "zero - empty index",
			initial:  zero,
			bitIndex: []int{},
		},
		{
			name:     "zero - short index",
			initial:  zero,
			bitIndex: []int{0, size - 1, 3, 1},
		},
		{
			name:     "zero - long index",
			initial:  zero,
			bitIndex: []int{0, 3, 1, size - 1, 2, 4, 7, size - 1},
		},
		{
			name:     "allset - empty index",
			initial:  allset,
			bitIndex: []int{},
		},
		{
			name:     "allset - short index",
			initial:  allset,
			bitIndex: []int{0, 3, size - 1, 1},
		},
		{
			name:     "allset - long index",
			initial:  allset,
			bitIndex: []int{size - 1, 3, 6, 2, 0, size - 1, 4, 7},
		},
	}
	b.Run(fmt.Sprintf("%T", zero), func(b *testing.B) {
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				var f TP = &tt.initial
				for i := 0; i < b.N; i++ {
					f.AnyOf(tt.bitIndex...)
				}
			})
		}
	})
}

func BenchmarkBitFlags_AnyOf(b *testing.B) {
	helperRunBenchmarkAnyOf[BitFlags8](b)
	helperRunBenchmarkAnyOf[BitFlags16](b)
	helperRunBenchmarkAnyOf[BitFlags32](b)
	helperRunBenchmarkAnyOf[BitFlags64](b)
}

func helperRunBenchmarkAllOf[T bitFlags, TP ptrBitFlags[T]](b *testing.B) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testCase struct {
		name     string
		initial  T
		bitIndex []int
	}
	tests := []testCase{
		{
			name:     "zero - empty index",
			initial:  zero,
			bitIndex: []int{},
		},
		{
			name:     "zero - short index",
			initial:  zero,
			bitIndex: []int{0, size - 1, 3, 1},
		},
		{
			name:     "zero - long index",
			initial:  zero,
			bitIndex: []int{0, 3, 1, size - 1, 2, 4, 7, size - 1},
		},
		{
			name:     "allset - empty index",
			initial:  allset,
			bitIndex: []int{},
		},
		{
			name:     "allset - short index",
			initial:  allset,
			bitIndex: []int{0, 3, size - 1, 1},
		},
		{
			name:     "allset - long index",
			initial:  allset,
			bitIndex: []int{size - 1, 3, 6, 2, 0, size - 1, 4, 7},
		},
	}
	b.Run(fmt.Sprintf("%T", zero), func(b *testing.B) {
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				var f TP = &tt.initial
				for i := 0; i < b.N; i++ {
					f.AllOf(tt.bitIndex...)
				}
			})
		}
	})
}

func BenchmarkBitFlags_AllOf(b *testing.B) {
	helperRunBenchmarkAllOf[BitFlags8](b)
	helperRunBenchmarkAllOf[BitFlags16](b)
	helperRunBenchmarkAllOf[BitFlags32](b)
	helperRunBenchmarkAllOf[BitFlags64](b)
}

func helperRunBenchmarkString[T bitFlags, TP ptrBitFlags[T]](b *testing.B) {
	var (
		zero   T
		allset = ^zero
	)
	type testCase struct {
		name    string
		initial T
	}
	tests := []testCase{
		{
			name:    "zero",
			initial: zero,
		},
		{
			name:    "allset",
			initial: allset,
		},
	}
	b.Run(fmt.Sprintf("%T", zero), func(b *testing.B) {
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				var f TP = &tt.initial
				for i := 0; i < b.N; i++ {
					f.String()
				}
			})
		}
	})
}

func BenchmarkBitFlags_String(b *testing.B) {
	helperRunBenchmarkString[BitFlags8](b)
	helperRunBenchmarkString[BitFlags16](b)
	helperRunBenchmarkString[BitFlags32](b)
	helperRunBenchmarkString[BitFlags64](b)
}

func helperRunBenchmarkPrettyString[T bitFlags, TP ptrBitFlags[T]](b *testing.B) {
	var (
		zero   T
		allset = ^zero
	)
	type testCase struct {
		name    string
		initial T
	}
	tests := []testCase{
		{
			name:    "zero",
			initial: zero,
		},
		{
			name:    "allset",
			initial: allset,
		},
	}
	b.Run(fmt.Sprintf("%T", zero), func(b *testing.B) {
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				var f TP = &tt.initial
				for i := 0; i < b.N; i++ {
					f.PrettyString()
				}
			})
		}
	})
}

func BenchmarkBitFlags_PrettyString(b *testing.B) {
	helperRunBenchmarkPrettyString[BitFlags8](b)
	helperRunBenchmarkPrettyString[BitFlags16](b)
	helperRunBenchmarkPrettyString[BitFlags32](b)
	helperRunBenchmarkPrettyString[BitFlags64](b)
}

func BenchmarkBitFlags_BitFlags(b *testing.B) {
	b.Run("BitFlags8", func(b *testing.B) {
		type testCase struct {
			name    string
			initial BitFlags8
		}
		tests := []testCase{
			{
				name:    "zero",
				initial: BitFlags8(0),
			},
			{
				name:    "allset",
				initial: ^BitFlags8(0),
			},
		}
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					tt.initial.BitFlags()
				}
			})
		}
	})
	b.Run("BitFlags16", func(b *testing.B) {
		type testCase struct {
			name    string
			initial BitFlags16
		}
		tests := []testCase{
			{
				name:    "zero",
				initial: BitFlags16(0),
			},
			{
				name:    "allset",
				initial: ^BitFlags16(0),
			},
		}
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					tt.initial.BitFlags()
				}
			})
		}
	})
	b.Run("BitFlags32", func(b *testing.B) {
		type testCase struct {
			name    string
			initial BitFlags32
		}
		tests := []testCase{
			{
				name:    "zero",
				initial: BitFlags32(0),
			},
			{
				name:    "allset",
				initial: ^BitFlags32(0),
			},
		}
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					tt.initial.BitFlags()
				}
			})
		}
	})
	b.Run("BitFlags64", func(b *testing.B) {
		type testCase struct {
			name    string
			initial BitFlags64
		}
		tests := []testCase{
			{
				name:    "zero",
				initial: BitFlags64(0),
			},
			{
				name:    "allset",
				initial: ^BitFlags64(0),
			},
		}
		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					tt.initial.BitFlags()
				}
			})
		}
	})
}
