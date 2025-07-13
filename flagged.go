// Package flagged provides a minimal, extensible API for manipulating and
// inspecting compact bitflags, while remaining dependency- and allocation-free.
package flagged

// BitIndex is a marker type denoting that its values should be used
// as bit indexes, passed to the different [BitFlags] methods.
// If the value is outside [BitFlags] range, the methods will panic.
// The allowed range is [0, [BitFlags.Size]-1].
//
// Example:
//
//	 const (
//		myOption1 flagged.BitIndex = iota
//		myOption2
//		myOption3
//	)
type BitIndex = int

// BitFlags8 is a wrapper for uint8 bit flags, carrying 8 flags at max.
type BitFlags8 uint8

// BitFlags16 is a wrapper for uint16 bit flags, carrying 16 flags at max.
type BitFlags16 uint16

// BitFlags32 is a wrapper for uint32 bit flags, carrying 32 flags at max.
type BitFlags32 uint32

// BitFlags64 is a wrapper for uint64 bit flags, carrying 64 flags at max.
type BitFlags64 uint64

// BitFlags represents the set of methods allowed on one of the typed
// bit flags types ([BitFlags8], [BitFlags16], [BitFlags32], [BitFlags64]).
// It makes it easy to write generic flag-aware code, regardless of the
// underlying type or its size.
type BitFlags interface {
	// Is reports whether the bit at index idx is set to true or not.
	// It panics if idx is out of the allowed range [0, Size-1].
	Is(idx BitIndex) (set bool)

	// Set sets the bit at index idx to true, returning its old value.
	// It panics if idx is out of the allowed range [0, Size-1].
	Set(idx BitIndex) (old bool)

	// Reset sets the bit at index idx to false, returning its old value.
	// It panics if idx is out of the allowed range [0, Size-1].
	Reset(idx BitIndex) (old bool)

	// SetTo sets the bit at index idx to new, returning its old value.
	// It panics if idx is out of the allowed range [0, Size-1].
	SetTo(idx BitIndex, new bool) (old bool)

	// Toggle toggles the bit at index idx, returning its new value.
	// It panics if idx is out of the allowed range [0, Size-1].
	Toggle(idx BitIndex) (new bool)

	// SetAll sets all bits to true.
	SetAll()

	// ResetAll sets all bits to false.
	ResetAll()

	// AnySet reports whether any of the bits are set to true.
	AnySet() bool

	// AllSet reports whether all the bits are set to true.
	AllSet() bool

	// AnyOf reports whether any of the bits at indexes idx are set to true.
	// If no indexes are passed, it acts as [BitFlags.AnySet].
	AnyOf(idx ...BitIndex) bool

	// AllOf reports whether all the bits at indexes idx are set to true.
	// If no indexes are passed, it acts as [BitFlags.AllSet].
	AllOf(idx ...BitIndex) bool

	// Size is the number of bits included in this [BitFlags] value.
	// It represents the bit width of the underlying uint.
	// It's one of 8, 16, 32, 64.
	Size() int

	// String returns the binary representation of this [BitFlags] value,
	// formatted like fmt's %b verb, but with leading zeros to preserve
	// the full bit width of the underlying type (Size).
	String() string

	// PrettyString returns a human-readable binary representation of
	// this [BitFlags] value, with 1 represented as 'I' and 0 as 'O', and
	// with '|' delimiter between each bit and '_' delimiter each 8 bits.
	//
	// Example:
	//
	//  String() // "0000010001000100"
	//  PrettyString() // "O|O|O|O|O|I|O|O_O|I|O|O|O|I|O|O"
	PrettyString() string
}

// New is a helper function for creating pointer to one of the BitFlags types.
// It's useful for returning a value that implements the [BitFlags] interface.
func New[T BitFlags8 | BitFlags16 | BitFlags32 | BitFlags64](f T) *T {
	return &f
}

func (f BitFlags8) Is(idx BitIndex) (set bool)               { return is(f, idx) }
func (f *BitFlags8) Set(idx BitIndex) (old bool)             { return set(f, idx, true) }
func (f *BitFlags8) Reset(idx BitIndex) (old bool)           { return set(f, idx, false) }
func (f *BitFlags8) SetTo(idx BitIndex, new bool) (old bool) { return set(f, idx, new) }
func (f *BitFlags8) Toggle(idx BitIndex) (new bool)          { return toggle(f, idx) }
func (f *BitFlags8) SetAll()                                 { setAll(f) }
func (f *BitFlags8) ResetAll()                               { resetAll(f) }
func (f BitFlags8) AnySet() bool                             { return anySet(f) }
func (f BitFlags8) AllSet() bool                             { return allSet(f) }
func (f BitFlags8) AnyOf(idx ...BitIndex) bool               { return anySet(f, idx...) }
func (f BitFlags8) AllOf(idx ...BitIndex) bool               { return allSet(f, idx...) }
func (BitFlags8) Size() int                                  { return 8 }
func (f BitFlags8) String() string                           { return getBinaryString(f) }
func (f BitFlags8) PrettyString() string                     { return getPrettyString(f) }
func (f *BitFlags8) BitFlags() BitFlags                      { return f }

func (f BitFlags16) Is(idx BitIndex) (set bool)               { return is(f, idx) }
func (f *BitFlags16) Set(idx BitIndex) (old bool)             { return set(f, idx, true) }
func (f *BitFlags16) Reset(idx BitIndex) (old bool)           { return set(f, idx, false) }
func (f *BitFlags16) SetTo(idx BitIndex, new bool) (old bool) { return set(f, idx, new) }
func (f *BitFlags16) Toggle(idx BitIndex) (new bool)          { return toggle(f, idx) }
func (f *BitFlags16) SetAll()                                 { setAll(f) }
func (f *BitFlags16) ResetAll()                               { resetAll(f) }
func (f BitFlags16) AnySet() bool                             { return anySet(f) }
func (f BitFlags16) AllSet() bool                             { return allSet(f) }
func (f BitFlags16) AnyOf(idx ...BitIndex) bool               { return anySet(f, idx...) }
func (f BitFlags16) AllOf(idx ...BitIndex) bool               { return allSet(f, idx...) }
func (BitFlags16) Size() int                                  { return 16 }
func (f BitFlags16) String() string                           { return getBinaryString(f) }
func (f BitFlags16) PrettyString() string                     { return getPrettyString(f) }
func (f *BitFlags16) BitFlags() BitFlags                      { return f }

func (f BitFlags32) Is(idx BitIndex) (set bool)               { return is(f, idx) }
func (f *BitFlags32) Set(idx BitIndex) (old bool)             { return set(f, idx, true) }
func (f *BitFlags32) Reset(idx BitIndex) (old bool)           { return set(f, idx, false) }
func (f *BitFlags32) SetTo(idx BitIndex, new bool) (old bool) { return set(f, idx, new) }
func (f *BitFlags32) Toggle(idx BitIndex) (new bool)          { return toggle(f, idx) }
func (f *BitFlags32) SetAll()                                 { setAll(f) }
func (f *BitFlags32) ResetAll()                               { resetAll(f) }
func (f BitFlags32) AnySet() bool                             { return anySet(f) }
func (f BitFlags32) AllSet() bool                             { return allSet(f) }
func (f BitFlags32) AnyOf(idx ...BitIndex) bool               { return anySet(f, idx...) }
func (f BitFlags32) AllOf(idx ...BitIndex) bool               { return allSet(f, idx...) }
func (BitFlags32) Size() int                                  { return 32 }
func (f BitFlags32) String() string                           { return getBinaryString(f) }
func (f BitFlags32) PrettyString() string                     { return getPrettyString(f) }
func (f *BitFlags32) BitFlags() BitFlags                      { return f }

func (f BitFlags64) Is(idx BitIndex) (set bool)               { return is(f, idx) }
func (f *BitFlags64) Set(idx BitIndex) (old bool)             { return set(f, idx, true) }
func (f *BitFlags64) Reset(idx BitIndex) (old bool)           { return set(f, idx, false) }
func (f *BitFlags64) SetTo(idx BitIndex, new bool) (old bool) { return set(f, idx, new) }
func (f *BitFlags64) Toggle(idx BitIndex) (new bool)          { return toggle(f, idx) }
func (f *BitFlags64) SetAll()                                 { setAll(f) }
func (f *BitFlags64) ResetAll()                               { resetAll(f) }
func (f BitFlags64) AnySet() bool                             { return anySet(f) }
func (f BitFlags64) AllSet() bool                             { return allSet(f) }
func (f BitFlags64) AnyOf(idx ...BitIndex) bool               { return anySet(f, idx...) }
func (f BitFlags64) AllOf(idx ...BitIndex) bool               { return allSet(f, idx...) }
func (BitFlags64) Size() int                                  { return 64 }
func (f BitFlags64) String() string                           { return getBinaryString(f) }
func (f BitFlags64) PrettyString() string                     { return getPrettyString(f) }
func (f *BitFlags64) BitFlags() BitFlags                      { return f }

type bitFlags interface {
	BitFlags8 | BitFlags16 | BitFlags32 | BitFlags64
	Size() int
}

func validateBitIndex(size int, idx BitIndex) {
	if idx < 0 || idx >= size {
		// print a helpful panic message without using fmt or strconv.
		strLen := 30 // of "index -00 out of range [0..00]"
		panicStr := make(stringBuilder, 0, strLen)
		panicStr.WriteString("index ")

		// only print the idx if it's between -nSmalls and nSmalls.
		if -nSmalls < idx && idx < nSmalls { // 2-digit number
			if idx < 0 {
				idx = -idx
				panicStr.WriteByte('-')
			}
			panicStr.WriteString(small(idx))
			panicStr.WriteByte(' ')
		}

		panicStr.WriteString("out of range [0..")
		panicStr.WriteString(sizeIndexString(size))
		panicStr.WriteString("]")
		panic(panicStr.String())
	}
}

func is[T bitFlags](f T, idx BitIndex) (set bool) {
	validateBitIndex(f.Size(), idx)
	return (f & (1 << idx)) != 0
}

func set[T bitFlags](f *T, idx BitIndex, new bool) (old bool) {
	validateBitIndex((*f).Size(), idx)
	old = (*f & (1 << idx)) != 0
	if new {
		*f |= 1 << idx
	} else {
		*f &^= 1 << idx
	}
	return
}

func toggle[T bitFlags](f *T, idx BitIndex) (new bool) {
	validateBitIndex((*f).Size(), idx)
	*f ^= 1 << idx
	return (*f & (1 << idx)) != 0
}

func setAll[T bitFlags](f *T) {
	var all = ^T(0)
	*f = all
}

func resetAll[T bitFlags](f *T) {
	*f = 0
}

func anySet[T bitFlags](f T, idx ...BitIndex) bool {
	if len(idx) == 0 {
		return f != T(0)
	}
	size := f.Size()
	foundSet := false
	for _, bi := range idx {
		validateBitIndex(size, bi)
		if (f & (1 << bi)) != 0 {
			foundSet = true
		}
	}
	return foundSet
}

func allSet[T bitFlags](f T, idx ...BitIndex) bool {
	if len(idx) == 0 {
		return f == ^T(0)
	}
	size := f.Size()
	foundUnset := true
	for _, bi := range idx {
		validateBitIndex(size, bi)
		if (f & (1 << bi)) == 0 {
			foundUnset = false
		}
	}
	return foundUnset
}

func getBinaryString[T bitFlags](f T) string {
	size := f.Size()
	str := make(stringBuilder, 0, size)
	for i := range size {
		if (f & (1 << (size - i - 1))) != 0 {
			str.WriteByte('1')
		} else {
			str.WriteByte('0')
		}
	}
	return str.String()
}

// getPrettyString prints f like "O|I|O|O|O|I|O|O_O|I|O|O|O|I|O|O"
func getPrettyString[T bitFlags](f T) string {
	size := f.Size()
	str := make(stringBuilder, 0, size+(size-1)+(size/8-1))
	for i := range size {
		if (f & (1 << (size - i - 1))) != 0 {
			if i == size-1 {
				str.WriteString("I")
			} else if (i+1)%8 == 0 && i != 0 {
				str.WriteString("I_")
			} else {
				str.WriteString("I|")
			}
		} else {
			if i == size-1 {
				str.WriteString("O")
			} else if (i+1)%8 == 0 && i != 0 {
				str.WriteString("O_")
			} else {
				str.WriteString("O|")
			}
		}
	}
	return str.String()
}

// stringBuilder is a simplified version of [strings.Builder],
// but without depending on the strings package, and without
// using the unsafe package.
// the result is 1 extra allocation, for avoiding importing
// the strings package.
type stringBuilder []byte

func (sb *stringBuilder) WriteByte(b byte) {
	*sb = append(*sb, b)
}
func (sb *stringBuilder) WriteString(s string) {
	*sb = append(*sb, s...)
}
func (sb *stringBuilder) String() string {
	return string(*sb)
}

// sizeIndexString returns size-1 as a string.
func sizeIndexString(size int) string {
	switch size {
	case 8:
		return "7"
	case 16:
		return "15"
	case 32:
		return "31"
	default:
		return "63"
	}
}

// small returns the string for an i with 0 <= i < nSmalls.
// copied from strconv.AppendUint implementation.
func small(i int) string {
	if i < 10 {
		return digits[i : i+1]
	}
	return smallsString[i*2 : i*2+2]
}

// copied from strconv.AppendUint implementation.
const nSmalls = 100

// copied from strconv.AppendUint implementation.
const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

// copied from strconv.AppendUint implementation.
const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
