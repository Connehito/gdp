package main

import (
	"strings"
	"testing"
	"time"
)

func TestGetNextVersion(t *testing.T) {
	tag, _ := GetNextVersion("v3.1.5")
	expected := "v3.1.6"
	if tag != expected {
		t.Errorf("Output=%q, Expected=%q", tag, expected)
	}
}

func TestGetNextVersion_Semantic(t *testing.T) {
	tag, _ := GetNextVersion("4.2.6")
	expected := "4.2.7"
	if tag != expected {
		t.Errorf("Output=%q, Expected=%q", tag, expected)
	}
}

func TestGetNextVersion_SemanticError(t *testing.T) {
	_, err := GetNextVersion("4.2.semantic")

	expected := "invalid syntax"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.Error(), expected)
	}
}

func TestGetNextVersion_Date(t *testing.T) {
	tag, _ := GetNextVersion("20180525.1")

	const layout = "20060102"
	today := time.Now().Format(layout)
	expected := today + ".1"
	if tag != expected {
		t.Errorf("Output=%q, Expected=%q", tag, expected)
	}
}

func TestGetNextVersion_DatePlus(t *testing.T) {
	const layout = "20060102"
	today := time.Now().Format(layout)
	tag, _ := GetNextVersion(today + ".1")

	expected := today + ".2"
	if tag != expected {
		t.Errorf("Output=%q, Expected=%q", tag, expected)
	}
}

func TestGetNextVersion_DateWithPrefix(t *testing.T) {
	const prefix = "release_"
	tag, _ := GetNextVersion(prefix + "20180525.1")

	const layout = "20060102"
	today := time.Now().Format(layout)
	expected := prefix + today + ".1"
	if tag != expected {
		t.Errorf("Output=%q, Expected=%q", tag, expected)
	}
}

func TestGetNextVersion_DatePlusWithPrefix(t *testing.T) {
	const prefix = "release_"
	const layout = "20060102"
	today := time.Now().Format(layout)
	tag, _ := GetNextVersion(prefix + today + ".1")

	expected := prefix + today + ".2"
	if tag != expected {
		t.Errorf("Output=%q, Expected=%q", tag, expected)
	}
}

func TestGetNextVersion_DateError(t *testing.T) {
	const layout = "20060102"
	today := time.Now().Format(layout)
	_, err := GetNextVersion(today + ".date")

	expected := "invalid syntax"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.Error(), expected)
	}
}

func TestGetNextVersion_EmptyTag(t *testing.T) {
	tag, _ := GetNextVersion("")
	expected := "v1.0.0"
	if tag != expected {
		t.Errorf("Output=%q, Expected=%q", tag, expected)
	}
}

func TestGetReleaseNote(t *testing.T) {
	list := "- itosho: initial commit\n"
	list = list + "- itosho: fix bug"
	note := GetReleaseNote("20180525.1", list)

	expected := "Release 20180525.1\n\n"
	expected = expected + "## 20180525.1\n"
	expected = expected + list

	if note != expected {
		t.Errorf("Output=%q, Expected=%q", note, expected)
	}
}
