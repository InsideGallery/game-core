package engine

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestNewAttributes(t *testing.T) {
	a := NewAttributes()
	if a == nil {
		t.Fatal("NewAttributes returned nil")
	}
	if a.values == nil {
		t.Fatal("NewAttributes values map is nil")
	}
}

func TestSetAndGetAttribute(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("key1", "value1")

	v, exists := a.GetAttribute("key1")
	testutils.Equal(t, exists, true)
	testutils.Equal(t, v, "value1")
}

func TestGetAttributeNonExistent(t *testing.T) {
	a := NewAttributes()

	v, exists := a.GetAttribute("missing")
	testutils.Equal(t, exists, false)
	testutils.Equal(t, v, nil)
}

func TestRemoveAttribute(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("key1", "value1")

	a.RemoveAttribute("key1")

	_, exists := a.GetAttribute("key1")
	testutils.Equal(t, exists, false)
}

func TestRemoveAttributeNonExistent(t *testing.T) {
	a := NewAttributes()
	// Should not panic
	a.RemoveAttribute("missing")
}

// Note: UpdateAttribute has the signature (name, f func(interface{}) interface{}),
// meaning both parameters are function-typed. Since function values are unhashable
// in Go, they cannot be used as map keys. This makes UpdateAttribute unusable at
// runtime -- any call will panic. We skip testing it as it is a broken API.

func TestSetAttributeOverwrite(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("key", "first")
	a.SetAttribute("key", "second")

	v, exists := a.GetAttribute("key")
	testutils.Equal(t, exists, true)
	testutils.Equal(t, v, "second")
}

func TestGetUint32(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("u32", uint32(42))

	result := a.GetUint32("u32")
	testutils.Equal(t, result, uint32(42))
}

func TestGetUint32NonExistent(t *testing.T) {
	a := NewAttributes()

	result := a.GetUint32("missing")
	testutils.Equal(t, result, uint32(0))
}

func TestGetUint32TypeMismatchPanics(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("wrong", "not a uint32")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for type mismatch, got none")
		}
	}()

	a.GetUint32("wrong")
}

func TestGetUint64(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("u64", uint64(999))

	result := a.GetUint64("u64")
	testutils.Equal(t, result, uint64(999))
}

func TestGetUint64NonExistent(t *testing.T) {
	a := NewAttributes()

	result := a.GetUint64("missing")
	testutils.Equal(t, result, uint64(0))
}

func TestGetUint64TypeMismatchPanics(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("wrong", "not a uint64")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for type mismatch, got none")
		}
	}()

	a.GetUint64("wrong")
}

func TestGetFloat64(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("f64", float64(3.14))

	result := a.GetFloat64("f64")
	testutils.Equal(t, result, float64(3.14))
}

func TestGetFloat64NonExistent(t *testing.T) {
	a := NewAttributes()

	result := a.GetFloat64("missing")
	testutils.Equal(t, result, float64(0))
}

func TestGetFloat64TypeMismatchPanics(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("wrong", "not a float64")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for type mismatch, got none")
		}
	}()

	a.GetFloat64("wrong")
}

func TestGetInt(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("i", 77)

	result := a.GetInt("i")
	testutils.Equal(t, result, 77)
}

func TestGetIntNonExistent(t *testing.T) {
	a := NewAttributes()

	result := a.GetInt("missing")
	testutils.Equal(t, result, 0)
}

func TestGetIntTypeMismatchPanics(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("wrong", "not an int")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for type mismatch, got none")
		}
	}()

	a.GetInt("wrong")
}

func TestGetUint8(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("u8", uint8(255))

	result := a.GetUint8("u8")
	testutils.Equal(t, result, uint8(255))
}

func TestGetUint8NonExistent(t *testing.T) {
	a := NewAttributes()

	result := a.GetUint8("missing")
	testutils.Equal(t, result, uint8(0))
}

func TestGetUint8TypeMismatchPanics(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("wrong", "not a uint8")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for type mismatch, got none")
		}
	}()

	a.GetUint8("wrong")
}

func TestGetBool(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("b", true)

	result := a.GetBool("b")
	testutils.Equal(t, result, true)
}

func TestGetBoolFalse(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("b", false)

	result := a.GetBool("b")
	testutils.Equal(t, result, false)
}

func TestGetBoolNonExistent(t *testing.T) {
	a := NewAttributes()

	result := a.GetBool("missing")
	testutils.Equal(t, result, false)
}

func TestGetBoolTypeMismatchPanics(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("wrong", "not a bool")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for type mismatch, got none")
		}
	}()

	a.GetBool("wrong")
}

func TestGetString(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("s", "hello")

	result := a.GetString("s")
	testutils.Equal(t, result, "hello")
}

func TestGetStringNonExistent(t *testing.T) {
	a := NewAttributes()

	result := a.GetString("missing")
	testutils.Equal(t, result, "")
}

func TestGetStringTypeMismatchPanics(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("wrong", 12345)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for type mismatch, got none")
		}
	}()

	a.GetString("wrong")
}

func TestTypedGettersTableDriven(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("u32", uint32(1))
	a.SetAttribute("u64", uint64(2))
	a.SetAttribute("f64", float64(3.5))
	a.SetAttribute("int", 4)
	a.SetAttribute("u8", uint8(5))
	a.SetAttribute("bool", true)
	a.SetAttribute("str", "test")

	t.Run("GetUint32", func(t *testing.T) {
		testutils.Equal(t, a.GetUint32("u32"), uint32(1))
	})
	t.Run("GetUint64", func(t *testing.T) {
		testutils.Equal(t, a.GetUint64("u64"), uint64(2))
	})
	t.Run("GetFloat64", func(t *testing.T) {
		testutils.Equal(t, a.GetFloat64("f64"), float64(3.5))
	})
	t.Run("GetInt", func(t *testing.T) {
		testutils.Equal(t, a.GetInt("int"), 4)
	})
	t.Run("GetUint8", func(t *testing.T) {
		testutils.Equal(t, a.GetUint8("u8"), uint8(5))
	})
	t.Run("GetBool", func(t *testing.T) {
		testutils.Equal(t, a.GetBool("bool"), true)
	})
	t.Run("GetString", func(t *testing.T) {
		testutils.Equal(t, a.GetString("str"), "test")
	})
}

func TestMultipleAttributes(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute("a", 1)
	a.SetAttribute("b", 2)
	a.SetAttribute("c", 3)

	testutils.Equal(t, a.GetInt("a"), 1)
	testutils.Equal(t, a.GetInt("b"), 2)
	testutils.Equal(t, a.GetInt("c"), 3)
}

func TestAttributeWithDifferentKeyTypes(t *testing.T) {
	a := NewAttributes()
	a.SetAttribute(1, "int key")
	a.SetAttribute("str", "string key")
	a.SetAttribute(true, "bool key")

	v1, e1 := a.GetAttribute(1)
	testutils.Equal(t, e1, true)
	testutils.Equal(t, v1, "int key")

	v2, e2 := a.GetAttribute("str")
	testutils.Equal(t, e2, true)
	testutils.Equal(t, v2, "string key")

	v3, e3 := a.GetAttribute(true)
	testutils.Equal(t, e3, true)
	testutils.Equal(t, v3, "bool key")
}
