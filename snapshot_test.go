package snapshot_test

import (
	"testing"

	"github.com/jhunt/go-snapshot"
)

func TestSimpleStructure(t *testing.T) {
	var o = struct {
		String string
		Bool   bool
		Int    int
	}{
		String: "original value",
		Bool:   true,
		Int:    42,
	}

	ok := func(good bool, got, expect interface{}, message string) {
		if !good {
			t.Errorf("%s failed - got [%v], expected [%v]\n", message, got, expect)
		}
	}

	ok(o.String == "original value", o.String, "original value", "string pre-check")
	ok(o.Bool == true, o.Bool, true, "bool pre-check")
	ok(o.Int == 42, o.Int, 42, "int pre-check")

	ss, err := snapshot.Take(&o)
	ok(err == nil, err, nil, "snapshotting")

	o.String = "updated value"
	o.Bool = false
	o.Int = -556

	ok(o.String == "updated value", o.String, "updated value", "string change post-snap")
	ok(o.Bool == false, o.Bool, false, "bool change post-snap")
	ok(o.Int == -556, o.Int, -556, "int change post-snap")

	err = ss.Revert()
	ok(err == nil, err, nil, "snapshot revert")

	ok(o.String == "original value", o.String, "original value", "string post-revert")
	ok(o.Bool == true, o.Bool, true, "bool post-revert")
	ok(o.Int == 42, o.Int, 42, "int post-revert")
}

func TestNestedStructure(t *testing.T) {
	var o = struct {
		String string
		Struct struct {
			Bool bool
			Int  int
		}
	}{}

	o.String = "original value"
	o.Struct.Bool = true
	o.Struct.Int = 42

	ok := func(good bool, got, expect interface{}, message string) {
		if !good {
			t.Errorf("%s failed - got [%v], expected [%v]\n", message, got, expect)
		}
	}

	ok(o.String == "original value", o.String, "original value", "string pre-check")
	ok(o.Struct.Bool == true, o.Struct.Bool, true, "bool pre-check")
	ok(o.Struct.Int == 42, o.Struct.Int, 42, "int pre-check")

	ss, err := snapshot.Take(&o)
	ok(err == nil, err, nil, "snapshotting")

	o.String = "updated value"
	o.Struct.Bool = false
	o.Struct.Int = -556

	ok(o.String == "updated value", o.String, "updated value", "string change post-snap")
	ok(o.Struct.Bool == false, o.Struct.Bool, false, "bool change post-snap")
	ok(o.Struct.Int == -556, o.Struct.Int, -556, "int change post-snap")

	err = ss.Revert()
	ok(err == nil, err, nil, "snapshot revert")

	ok(o.String == "original value", o.String, "original value", "string post-revert")
	ok(o.Struct.Bool == true, o.Struct.Bool, true, "bool post-revert")
	ok(o.Struct.Int == 42, o.Struct.Int, 42, "int post-revert")
}

func TestSimpleVars(t *testing.T) {
	String := "original value"
	Bool := true
	Int := 42

	ok := func(good bool, got, expect interface{}, message string) {
		if !good {
			t.Errorf("%s failed - got [%v], expected [%v]\n", message, got, expect)
		}
	}

	ok(String == "original value", String, "original value", "string pre-check")
	ok(Bool == true, Bool, true, "bool pre-check")
	ok(Int == 42, Int, 42, "int pre-check")

	ss1, err := snapshot.Take(&String)
	ok(err == nil, err, nil, "snapshotting string")

	ss2, err := snapshot.Take(&Bool)
	ok(err == nil, err, nil, "snapshotting bool")

	ss3, err := snapshot.Take(&Int)
	ok(err == nil, err, nil, "snapshotting int")

	String = "updated value"
	Bool = false
	Int = -556

	ok(String == "updated value", String, "updated value", "string change post-snap")
	ok(Bool == false, Bool, false, "bool change post-snap")
	ok(Int == -556, Int, -556, "int change post-snap")

	err = ss1.Revert()
	ok(err == nil, err, nil, "snapshot revert string")

	err = ss2.Revert()
	ok(err == nil, err, nil, "snapshot revert bool")

	err = ss3.Revert()
	ok(err == nil, err, nil, "snapshot revert int")

	ok(String == "original value", String, "original value", "string post-revert")
	ok(Bool == true, Bool, true, "bool post-revert")
	ok(Int == 42, Int, 42, "int post-revert")
}

func TestSlices(t *testing.T) {
	var Slices = []string{"a", "b", "c"}
	var o = struct {
		Bools []bool
	}{}

	ok := func(good bool, got, expect interface{}, message string) {
		if !good {
			t.Errorf("%s failed - got [%v], expected [%v]\n", message, got, expect)
		}
	}

	o.Bools = make([]bool, 2)
	o.Bools[0] = true
	o.Bools[1] = false

	ok(len(Slices) == 3, len(Slices), 3, "slices pre-check")
	ok(Slices[0] == "a", Slices[0], "a", "slices[0] pre-check")
	ok(Slices[1] == "b", Slices[1], "b", "slices[1] pre-check")
	ok(Slices[2] == "c", Slices[2], "c", "slices[2] pre-check")

	ok(len(o.Bools) == 2, len(o.Bools), 2, "struct slices pre-check")
	ok(o.Bools[0] == true, o.Bools[0], true, "struct slices[0] pre-check")
	ok(o.Bools[1] == false, o.Bools[1], false, "struct slices[1] pre-check")

	ss1, err := snapshot.Take(&Slices)
	ok(err == nil, err, nil, "snapshotting slices")

	ss2, err := snapshot.Take(&o)
	ok(err == nil, err, nil, "snapshotting struct")

	Slices[0] = "ALPHA"
	Slices[1] = "BETA"
	Slices[2] = "GAMMA"

	ok(len(Slices) == 3, len(Slices), 3, "slices pre-check")
	ok(Slices[0] == "ALPHA", Slices[0], "ALPHA", "slices[0] post-snap")
	ok(Slices[1] == "BETA", Slices[1], "BETA", "slices[1] post-snap")
	ok(Slices[2] == "GAMMA", Slices[2], "GAMMA", "slices[1] post-snap")

	o.Bools = []bool{false, false, true}

	ok(len(o.Bools) == 3, len(o.Bools), 3, "struct slices post-snap")
	ok(o.Bools[0] == false, o.Bools[0], false, "struct slices[0] post-snap")
	ok(o.Bools[1] == false, o.Bools[1], false, "struct slices[1] post-snap")
	ok(o.Bools[2] == true, o.Bools[2], true, "struct slices[2] post-snap")

	err = ss1.Revert()
	ok(err == nil, err, nil, "snapshot revert slices")

	err = ss2.Revert()
	ok(err == nil, err, nil, "snapshot revert struct")

	ok(len(Slices) == 3, len(Slices), 3, "slices post-revert")
	ok(Slices[0] == "a", Slices[0], "a", "slices[0] post-revert")
	ok(Slices[1] == "b", Slices[1], "b", "slices[1] post-revert")
	ok(Slices[2] == "c", Slices[2], "c", "slices[2] post-revert")

	ok(len(o.Bools) == 2, len(o.Bools), 2, "struct slices post-revert")
	ok(o.Bools[0] == true, o.Bools[0], true, "struct slices[0] post-revert")
	ok(o.Bools[1] == false, o.Bools[1], false, "struct slices[1] post-revert")
}

func TestKitchenSink(t *testing.T) {
	var o = struct {
		Bool       bool
		Int        int
		Int8       int8
		Int16      int16
		Int32      int32
		Int64      int64
		Uint       uint
		Uint8      uint8
		Uint16     uint16
		Uint32     uint32
		Uint64     uint64
		Uintptr    uintptr
		Float32    float32
		Float64    float64
		Complex64  complex64
		Complex128 complex128
		Array      [5]string
		Chan       chan int
		Func       func() error
		Interface  interface{} // fails
		Map        map[string]string
		Ptr        *string
		Slice      []int
		String     string
		Struct     struct {
			Member string
		}
		//UnsafePointer

		private string
	}{}

	ok := func(good bool, got, expect interface{}, message string) {
		if !good {
			t.Errorf("%s failed - got [%v], expected [%v]\n", message, got, expect)
		}
	}

	ss, err := snapshot.Take(&o)
	ok(err == nil, err, nil, "snapshotting")

	err = ss.Revert()
	ok(err == nil, err, nil, "reverting")
}
