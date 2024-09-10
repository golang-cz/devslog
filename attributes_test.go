package devslog

import (
	"log/slog"
	"testing"
)

func Test_Attributes(t *testing.T) {
	test_AttributesLen(t)
	test_AttributesSwap(t)
	test_AttributesLess(t)
	test_AttributesLessGroupTrue(t)
	test_AttributesLessGroupFalse(t)
	test_AttributesPadding(t)
}

func test_AttributesLen(t *testing.T) {
	someValue := slog.StringValue("value")
	attrs := attributes{
		slog.Attr{Key: "key1", Value: someValue},
		slog.Attr{Key: "key2", Value: someValue},
		slog.Attr{Key: "key3", Value: someValue},
	}

	expectedLen := 3
	actualLen := attrs.Len()

	if actualLen != expectedLen {
		t.Errorf("Expected length: %d, but got: %d", expectedLen, actualLen)
	}
}

func test_AttributesSwap(t *testing.T) {
	attr1 := slog.Attr{Key: "key1", Value: slog.StringValue("value")}
	attr2 := slog.Attr{Key: "key2", Value: slog.StringValue("value")}
	attrs := attributes{
		attr1,
		attr2,
	}

	attrs.Swap(0, 1)

	if attrs[0].Key != attr2.Key || attrs[1].Key != attr1.Key {
		t.Error("attributes were not swapped correctly")
	}
}

func test_AttributesLess(t *testing.T) {
	someValue := slog.StringValue("value")
	attrs := attributes{
		slog.Attr{Key: "key1", Value: someValue},
		slog.Attr{Key: "key2", Value: someValue},
	}

	less := attrs.Less(0, 1)

	if !less {
		t.Error("Expected the first attribute to be less than the second")
	}
}

func test_AttributesLessGroupTrue(t *testing.T) {
	attrs := attributes{
		slog.String("key1", "someValue"),
		slog.Group("key2", slog.String("someString", "someValue")),
	}

	less := attrs.Less(0, 1)

	if !less {
		t.Error("Expected the first attribute to be less than the second")
	}
}

func test_AttributesLessGroupFalse(t *testing.T) {
	attrs := attributes{
		slog.Group("key1", slog.String("someString", "someValue")),
		slog.String("key2", "someValue"),
	}

	less := attrs.Less(0, 1)

	if less {
		t.Error("Expected the first attribute to be less than the second")
	}
}

func test_AttributesPadding(t *testing.T) {
	someValue := slog.StringValue("value")
	attrs := attributes{
		slog.Attr{Key: "key1", Value: someValue},
		slog.Attr{Key: "key2", Value: someValue},
	}

	h := NewHandler(nil, nil)
	padding := attrs.padding(fgMagenta, h.cs)

	expectedPadding := 13
	if padding != expectedPadding {
		t.Errorf("Expected padding: %d, but got: %d", expectedPadding, padding)
	}
}
