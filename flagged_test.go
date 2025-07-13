package flagged

import (
	"fmt"
	"testing"
)

type testAPIType struct {
	myOptions BitFlags8
}

func TestAPI(t *testing.T) {
	var testAPI = testAPIType{}

	const bitIndex0 = 0
	if testAPI.myOptions.Set(bitIndex0) {
		t.Errorf("Set() = %v, want = %v", true, false)
	}
	if !testAPI.myOptions.Set(bitIndex0) {
		t.Errorf("Set() = %v, want = %v", false, true)
	}

	if !testAPI.myOptions.Is(bitIndex0) {
		t.Errorf("Is() = %v, want = %v", true, false)
	}

	const bitIndex1 = 1
	if testAPI.myOptions.BitFlags().Set(bitIndex1) {
		t.Errorf("Set() = %v, want = %v", true, false)
	}
	if !testAPI.myOptions.BitFlags().Set(bitIndex1) {
		t.Errorf("Set() = %v, want = %v", false, true)
	}

	if !testAPI.myOptions.BitFlags().Is(bitIndex1) {
		t.Errorf("Is() = %v, want = %v", true, false)
	}
	if !testAPI.myOptions.BitFlags().Is(bitIndex0) {
		t.Errorf("Is() = %v, want = %v", true, false)
	}

	const bitIndex3 = 3
	if testAPI.myOptions.Set(bitIndex3) {
		t.Errorf("Set() = %v, want = %v", true, false)
	}
	if !testAPI.myOptions.Set(bitIndex3) {
		t.Errorf("Set() = %v, want = %v", false, true)
	}

	if !testAPI.myOptions.BitFlags().Is(bitIndex3) {
		t.Errorf("Is() = %v, want = %v", true, false)
	}
	if !testAPI.myOptions.Is(bitIndex1) {
		t.Errorf("Is() = %v, want = %v", true, false)
	}
}

type ptrBitFlags[T bitFlags] interface {
	*T
	BitFlags
}

func helperRunTestIs[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testRun struct {
		bitIndex int
		want     bool
		panics   bool
	}
	type testCase struct {
		name    string
		initial T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero - within range",
			initial: zero,
			runs: []testRun{
				{
					bitIndex: 0,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     false,
					panics:   false,
				},
			},
		},
		{
			name:    "zero - out of range",
			initial: zero,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "allset - within range",
			initial: allset,
			runs: []testRun{
				{
					bitIndex: 0,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     true,
					panics:   false,
				},
			},
		},
		{
			name:    "allset - out of range",
			initial: allset,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					func() {
						defer func() {
							v := recover()
							if v == nil && tr.panics || v != nil && !tr.panics {
								t.Errorf("[%d] Is() panicked = %v, want = %v", ti, v != nil, tr.panics)
							}
						}()

						if got := f.Is(tr.bitIndex); got != tr.want {
							t.Errorf("Is() = %v, want %v", got, tr.want)
						}
					}()
				}
			})
		}
	})
}

func TestTestBitFlags_Is(t *testing.T) {
	helperRunTestIs[BitFlags8](t)
	helperRunTestIs[BitFlags16](t)
	helperRunTestIs[BitFlags32](t)
	helperRunTestIs[BitFlags64](t)
}

func helperRunTestSet[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testRun struct {
		bitIndex int
		want     bool
		panics   bool
	}
	type testCase struct {
		name    string
		initial T
		updated T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero - within range",
			initial: zero,
			updated: zero | T(1) | 1<<(size/2) | 1<<(size-1),
			runs: []testRun{
				{
					bitIndex: 0,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     false,
					panics:   false,
				},
			},
		},
		{
			name:    "zero - out of range",
			initial: zero,
			updated: zero,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "allset - within range",
			initial: allset,
			updated: allset | T(1) | 1<<(size/2) | 1<<(size-1),
			runs: []testRun{
				{
					bitIndex: 0,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     true,
					panics:   false,
				},
			},
		},
		{
			name:    "allset - out of range",
			initial: allset,
			updated: allset,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					func() {
						defer func() {
							v := recover()
							if v == nil && tr.panics || v != nil && !tr.panics {
								t.Errorf("[%d] Set() panicked = %v, want = %v", ti, v != nil, tr.panics)
							}
						}()

						if got := f.Set(tr.bitIndex); got != tr.want {
							t.Errorf("[%d] Set() = %v, want %v", ti, got, tr.want)
						}

						// second Set must return true.
						if got := f.Set(tr.bitIndex); !got {
							t.Errorf("[%d] Set() = %v, want %v", ti, got, true)
						}
					}()
				}

				if tt.initial != tt.updated {
					t.Errorf("Set() updated inital unexpectedly got = %v, want = %v", tt.initial, tt.updated)
				}
			})
		}
	})
}

func TestBitFlags_Set(t *testing.T) {
	helperRunTestSet[BitFlags8](t)
	helperRunTestSet[BitFlags16](t)
	helperRunTestSet[BitFlags32](t)
	helperRunTestSet[BitFlags64](t)
}

func helperRunTestReset[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testRun struct {
		bitIndex int
		want     bool
		panics   bool
	}
	type testCase struct {
		name    string
		initial T
		updated T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero - within range",
			initial: zero,
			updated: zero,
			runs: []testRun{
				{
					bitIndex: 0,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     false,
					panics:   false,
				},
			},
		},
		{
			name:    "zero - out of range",
			initial: zero,
			updated: zero,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "allset - within range",
			initial: allset,
			updated: allset &^ T(1) &^ (1 << (size / 2)) &^ (1 << (size - 1)),
			runs: []testRun{
				{
					bitIndex: 0,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     true,
					panics:   false,
				},
			},
		},
		{
			name:    "allset - out of range",
			initial: allset,
			updated: allset,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					func() {
						defer func() {
							v := recover()
							if v == nil && tr.panics || v != nil && !tr.panics {
								t.Errorf("[%d] Reset() panicked = %v, want = %v", ti, v != nil, tr.panics)
							}
						}()

						if got := f.Reset(tr.bitIndex); got != tr.want {
							t.Errorf("[%d] Reset() = %v, want %v", ti, got, tr.want)
						}

						// second Reset must return false.
						if got := f.Reset(tr.bitIndex); got {
							t.Errorf("[%d] Reset() = %v, want %v", ti, got, false)
						}
					}()
				}

				if tt.initial != tt.updated {
					t.Errorf("Reset() updated inital unexpectedly got = %v, want = %v", tt.initial, tt.updated)
				}
			})
		}
	})
}

func TestBitFlags_Reset(t *testing.T) {
	helperRunTestReset[BitFlags8](t)
	helperRunTestReset[BitFlags16](t)
	helperRunTestReset[BitFlags32](t)
	helperRunTestReset[BitFlags64](t)
}

func helperRunTestToggle[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testRun struct {
		bitIndex int
		want     bool
		panics   bool
	}
	type testCase struct {
		name    string
		initial T
		updated T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero - within range",
			initial: zero,
			updated: zero,
			runs: []testRun{
				{
					bitIndex: 0,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     true,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     true,
					panics:   false,
				},
			},
		},
		{
			name:    "zero - out of range",
			initial: zero,
			updated: zero,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "allset - within range",
			initial: allset,
			updated: allset,
			runs: []testRun{
				{
					bitIndex: 0,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size / 2,
					want:     false,
					panics:   false,
				},
				{
					bitIndex: size - 1,
					want:     false,
					panics:   false,
				},
			},
		},
		{
			name:    "allset - out of range",
			initial: allset,
			updated: allset,
			runs: []testRun{
				{
					bitIndex: -1,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size,
					want:     false,
					panics:   true,
				},
				{
					bitIndex: size * 2,
					want:     false,
					panics:   true,
				},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					func() {
						defer func() {
							v := recover()
							if v == nil && tr.panics || v != nil && !tr.panics {
								t.Errorf("[%d] Toggle() panicked = %v, want = %v", ti, v != nil, tr.panics)
							}
						}()

						if got := f.Toggle(tr.bitIndex); got != tr.want {
							t.Errorf("[%d] Toggle() = %v, want %v", ti, got, tr.want)
						}

						// second Toggle must return opposite of previous.
						if got := f.Toggle(tr.bitIndex); got == tr.want {
							t.Errorf("[%d] Toggle() = %v, want %v", ti, got, !tr.want)
						}
					}()
				}

				if tt.initial != tt.updated {
					t.Errorf("Toggle() updated inital unexpectedly got = %v, want = %v", tt.initial, tt.updated)
				}
			})
		}
	})
}

func TestBitFlags_Toggle(t *testing.T) {
	helperRunTestToggle[BitFlags8](t)
	helperRunTestToggle[BitFlags16](t)
	helperRunTestToggle[BitFlags32](t)
	helperRunTestToggle[BitFlags64](t)
}

func helperRunTestSetAll[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
	)
	type testRun struct{}
	type testCase struct {
		name    string
		initial T
		updated T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero",
			initial: zero,
			updated: allset,
			runs: []testRun{
				{},
				{},
			},
		},
		{
			name:    "allset",
			initial: allset,
			updated: allset,
			runs: []testRun{
				{},
				{},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for range tt.runs {
					f.SetAll()
				}

				if tt.initial != tt.updated {
					t.Errorf("SetAll() updated inital unexpectedly got = %v, want = %v", tt.initial, tt.updated)
				}
			})
		}
	})
}

func TestBitFlags_SetAll(t *testing.T) {
	helperRunTestSetAll[BitFlags8](t)
	helperRunTestSetAll[BitFlags16](t)
	helperRunTestSetAll[BitFlags32](t)
	helperRunTestSetAll[BitFlags64](t)
}

func helperRunTestResetAll[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
	)
	type testRun struct{}
	type testCase struct {
		name    string
		initial T
		updated T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero",
			initial: zero,
			updated: zero,
			runs: []testRun{
				{},
				{},
			},
		},
		{
			name:    "allset",
			initial: allset,
			updated: zero,
			runs: []testRun{
				{},
				{},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for range tt.runs {
					f.ResetAll()
				}

				if tt.initial != tt.updated {
					t.Errorf("ResetAll() updated inital unexpectedly got = %v, want = %v", tt.initial, tt.updated)
				}
			})
		}
	})
}

func TestBitFlags_ResetAll(t *testing.T) {
	helperRunTestResetAll[BitFlags8](t)
	helperRunTestResetAll[BitFlags16](t)
	helperRunTestResetAll[BitFlags32](t)
	helperRunTestResetAll[BitFlags64](t)
}

func helperRunTestAnySet[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
	)
	type testRun struct {
		want bool
	}
	type testCase struct {
		name    string
		initial T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero",
			initial: zero,
			runs: []testRun{
				{want: false},
			},
		},
		{
			name:    "allset",
			initial: allset,
			runs: []testRun{
				{want: true},
			},
		},
		{
			name:    "partial",
			initial: zero | T(1)<<1,
			runs: []testRun{
				{want: true},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					if got := f.AnySet(); got != tr.want {
						t.Errorf("[%d] AnySet() = %v, want = %v", ti, got, tr.want)
					}
				}
			})
		}
	})
}

func TestBitFlags_AnySet(t *testing.T) {
	helperRunTestAnySet[BitFlags8](t)
	helperRunTestAnySet[BitFlags16](t)
	helperRunTestAnySet[BitFlags32](t)
	helperRunTestAnySet[BitFlags64](t)
}

func helperRunTestAllSet[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
	)
	type testRun struct {
		want bool
	}
	type testCase struct {
		name    string
		initial T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero",
			initial: zero,
			runs: []testRun{
				{want: false},
			},
		},
		{
			name:    "allset",
			initial: allset,
			runs: []testRun{
				{want: true},
			},
		},
		{
			name:    "partial",
			initial: zero | T(1)<<1,
			runs: []testRun{
				{want: false},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					if got := f.AllSet(); got != tr.want {
						t.Errorf("[%d] AllSet() = %v, want = %v", ti, got, tr.want)
					}
				}
			})
		}
	})
}

func TestBitFlags_AllSet(t *testing.T) {
	helperRunTestAllSet[BitFlags8](t)
	helperRunTestAllSet[BitFlags16](t)
	helperRunTestAllSet[BitFlags32](t)
	helperRunTestAllSet[BitFlags64](t)
}

func helperRunTestAnyOf[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testRun struct {
		bitIndex []int
		want     bool
		panics   bool
	}
	type testCase struct {
		name    string
		initial T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero - within range",
			initial: zero,
			runs: []testRun{
				{
					bitIndex: []int{},
					want:     false,
					panics:   false,
				},
				{
					bitIndex: []int{0, 3, 1},
					want:     false,
					panics:   false,
				},
			},
		},
		{
			name:    "zero - out of range",
			initial: zero,
			runs: []testRun{
				{
					bitIndex: []int{-1},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{0, 3, size},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{size, 7, 6},
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "allset - within range",
			initial: allset,
			runs: []testRun{
				{
					bitIndex: []int{},
					want:     true,
					panics:   false,
				},
				{
					bitIndex: []int{0, 3, 1},
					want:     true,
					panics:   false,
				},
			},
		},
		{
			name:    "allset - out of range",
			initial: allset,
			runs: []testRun{
				{
					bitIndex: []int{-size},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{size, 3, 1},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{5, size, 0},
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "partial",
			initial: zero | T(1)<<1,
			runs: []testRun{
				{
					bitIndex: []int{},
					want:     true,
					panics:   false,
				},
				{
					bitIndex: []int{3, 1},
					want:     true,
					panics:   false,
				},
				{
					bitIndex: []int{3, 5, 7},
					want:     false,
					panics:   false,
				},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					func() {
						defer func() {
							v := recover()
							if v == nil && tr.panics || v != nil && !tr.panics {
								t.Errorf("[%d] AnyOf() panicked = %v, want = %v", ti, v != nil, tr.panics)
							}
						}()

						if got := f.AnyOf(tr.bitIndex...); got != tr.want {
							t.Errorf("[%d] AnyOf() = %v, want = %v", ti, got, tr.want)
						}
					}()
				}
			})
		}
	})
}

func TestBitFlags_AnyOf(t *testing.T) {
	helperRunTestAnyOf[BitFlags8](t)
	helperRunTestAnyOf[BitFlags16](t)
	helperRunTestAnyOf[BitFlags32](t)
	helperRunTestAnyOf[BitFlags64](t)
}

func helperRunTestAllOf[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
		size   = zero.Size()
	)
	type testRun struct {
		bitIndex []int
		want     bool
		panics   bool
	}
	type testCase struct {
		name    string
		initial T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero - within range",
			initial: zero,
			runs: []testRun{
				{
					bitIndex: []int{},
					want:     false,
					panics:   false,
				},
				{
					bitIndex: []int{0, 3, 1},
					want:     false,
					panics:   false,
				},
			},
		},
		{
			name:    "zero - out of range",
			initial: zero,
			runs: []testRun{
				{
					bitIndex: []int{-1},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{0, 3, size},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{size, 4},
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "allset - within range",
			initial: allset,
			runs: []testRun{
				{
					bitIndex: []int{},
					want:     true,
					panics:   false,
				},
				{
					bitIndex: []int{0, 3, 1},
					want:     true,
					panics:   false,
				},
			},
		},
		{
			name:    "allset - out of range",
			initial: allset,
			runs: []testRun{
				{
					bitIndex: []int{-size},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{size, 3, 1},
					want:     false,
					panics:   true,
				},
				{
					bitIndex: []int{7, size, 0},
					want:     false,
					panics:   true,
				},
			},
		},
		{
			name:    "partial",
			initial: zero | T(1)<<1 | T(1)<<3,
			runs: []testRun{
				{
					bitIndex: []int{},
					want:     false,
					panics:   false,
				},
				{
					bitIndex: []int{3, 1},
					want:     true,
					panics:   false,
				},
				{
					bitIndex: []int{3, 5, 7},
					want:     false,
					panics:   false,
				},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti, tr := range tt.runs {
					func() {
						defer func() {
							v := recover()
							if v == nil && tr.panics || v != nil && !tr.panics {
								t.Errorf("[%d] AllOf() panicked = %v, want = %v", ti, v != nil, tr.panics)
							}
						}()

						if got := f.AllOf(tr.bitIndex...); got != tr.want {
							t.Errorf("[%d] AllOf() = %v, want = %v", ti, got, tr.want)
						}
					}()
				}
			})
		}
	})
}

func TestBitFlags_AllOf(t *testing.T) {
	helperRunTestAllOf[BitFlags8](t)
	helperRunTestAllOf[BitFlags16](t)
	helperRunTestAllOf[BitFlags32](t)
	helperRunTestAllOf[BitFlags64](t)
}

func helperRunTestString[T bitFlags, TP ptrBitFlags[T]](t *testing.T) {
	var (
		zero   T
		allset = ^zero
	)
	type testRun struct{}
	type testCase struct {
		name    string
		initial T
		runs    []testRun
	}
	tests := []testCase{
		{
			name:    "zero",
			initial: zero,
			runs: []testRun{
				{},
			},
		},
		{
			name:    "allset",
			initial: allset,
			runs: []testRun{
				{},
			},
		},
		{
			name:    "partial",
			initial: zero | T(1)<<1 | T(1)<<6,
			runs: []testRun{
				{},
			},
		},
	}
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				for ti := range tt.runs {
					got := f.String()

					format := fmt.Sprintf("%%0%db", f.Size())
					want := fmt.Sprintf(format, tt.initial)

					if got != want {
						t.Errorf("[%d] String() = %v, want = %v", ti, got, want)
					}
				}
			})
		}
	})
}

func TestBitFlags_String(t *testing.T) {
	helperRunTestString[BitFlags8](t)
	helperRunTestString[BitFlags16](t)
	helperRunTestString[BitFlags32](t)
	helperRunTestString[BitFlags64](t)
}

type prettyStringTest[T bitFlags] struct {
	name    string
	initial T
	want    string
}

func helperRunTestPrettyString[T bitFlags, TP ptrBitFlags[T]](
	t *testing.T,
	tests []prettyStringTest[T],
) {
	var zero T
	t.Run(fmt.Sprintf("%T", zero), func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var f TP = &tt.initial

				got := f.PrettyString()

				if got != tt.want {
					t.Errorf("PrettyString() = %v, want = %v", got, tt.want)
				}
			})
		}
	})
}

func TestBitFlags_PrettyString(t *testing.T) {
	helperRunTestPrettyString(
		t,
		[]prettyStringTest[BitFlags8]{
			{
				name:    "zero",
				initial: BitFlags8(0),
				want:    "O|O|O|O|O|O|O|O",
			},
			{
				name:    "allset",
				initial: ^BitFlags8(0),
				want:    "I|I|I|I|I|I|I|I",
			},
			{
				name:    "partial",
				initial: BitFlags8(0) | BitFlags8(1)<<1 | BitFlags8(1)<<6,
				want:    "O|I|O|O|O|O|I|O",
			},
		},
	)
	helperRunTestPrettyString(
		t,
		[]prettyStringTest[BitFlags16]{
			{
				name:    "zero",
				initial: BitFlags16(0),
				want:    "O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O",
			},
			{
				name:    "allset",
				initial: ^BitFlags16(0),
				want:    "I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I",
			},
			{
				name:    "partial",
				initial: BitFlags16(0) | BitFlags16(1)<<1 | BitFlags16(1)<<6,
				want:    "O|O|O|O|O|O|O|O_O|I|O|O|O|O|I|O",
			},
		},
	)
	helperRunTestPrettyString(
		t,
		[]prettyStringTest[BitFlags32]{
			{
				name:    "zero",
				initial: BitFlags32(0),
				want:    "O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O",
			},
			{
				name:    "allset",
				initial: ^BitFlags32(0),
				want:    "I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I",
			},
			{
				name:    "partial",
				initial: BitFlags32(0) | BitFlags32(1)<<1 | BitFlags32(1)<<31,
				want:    "I|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|I|O",
			},
		},
	)
	helperRunTestPrettyString(
		t,
		[]prettyStringTest[BitFlags64]{
			{
				name:    "zero",
				initial: BitFlags64(0),
				want:    "O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O",
			},
			{
				name:    "allset",
				initial: ^BitFlags64(0),
				want:    "I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I_I|I|I|I|I|I|I|I",
			},
			{
				name:    "partial",
				initial: BitFlags64(0) | BitFlags64(1)<<1 | BitFlags64(1)<<63,
				want:    "I|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|O|O_O|O|O|O|O|O|I|O",
			},
		},
	)
}

func Test_validateBitIndex_panic(t *testing.T) {
	tests := []struct {
		name   string
		size   int
		idx    BitIndex
		panicV any
	}{
		{
			name:   "no panic",
			size:   8,
			idx:    7,
			panicV: nil,
		},
		{
			name:   "positive panic - small idx",
			size:   16,
			idx:    16,
			panicV: "index 16 out of range [0..15]",
		},
		{
			name:   "negative panic - small idx",
			size:   64,
			idx:    -99,
			panicV: "index -99 out of range [0..63]",
		},
		{
			name:   "positive panic - big idx",
			size:   32,
			idx:    100,
			panicV: "index out of range [0..31]",
		},
		{
			name:   "negative panic - big idx",
			size:   64,
			idx:    -9999,
			panicV: "index out of range [0..63]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				v := recover()
				if v != tt.panicV {
					t.Errorf("got panicV: %v; want: %v", v, tt.panicV)
				}
			}()

			validateBitIndex(tt.size, tt.idx)
		})
	}
}
