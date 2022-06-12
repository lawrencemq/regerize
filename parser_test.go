package main

import (
	"testing"
)

func ensureEquals(t *testing.T, expected, result string) {
	if expected != result {
		t.Fatalf("'%s' != '%s'", expected, result)
	}
}
func TestBatman(t *testing.T) {
	expected := "(na){16}\\sbatman"
	input := `16 of "na";
<space>;
batman;`
	result := parse(input)
	ensureEquals(t, expected, result)

}

func TestSomeOf(t *testing.T) {
	input := "some of <alpha>"
	expected := "[a-zA-Z]+"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestStart(t *testing.T) {
	input := "<start>;hello world;"
	expected := "^hello world"
	result := parse(input)
	ensureEquals(t, expected, result)
}
func TestEnd(t *testing.T) {
	input := "world;<end>;"
	expected := "world$"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestMaybe(t *testing.T) {
	input := "maybe of <space>;"
	expected := "\\s?"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestAny(t *testing.T) {
	input := "any of <alpha>;"
	expected := "[a-zA-Z]*"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestRange(t *testing.T) {
	input := "5 to 9 of \"hello\""
	expected := "(hello){5,9}"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestAtLeast(t *testing.T) {
	input := "at least 5 of \"ducks\""
	expected := "(ducks){5,}"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestAtMost(t *testing.T) {
	input := "at most 5 of \"ducks\""
	expected := "(ducks){,5}"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestBeforeBlock(t *testing.T) {
	input := `before {
"www.";
};
some of <word>;
".com"`
	expected := "(?<=www\\.)\\w+\\.com"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestAfterBlock(t *testing.T) {
	input := `"www.google.com";
after {
"/";
some of <word>;
};`
	expected := "www\\.google\\.com(?=/\\w+)"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestRaw(t *testing.T) {
	input := "`abc123(?=something)`;"
	expected := "abc123(?=something)"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestCap(t *testing.T) {
	input := `capture as shiba {
"inu";
<space>;
at least 5 of "wow";
};`
	expected := "(?<shib>inu\\s(wow){5,})"
	result := parse(input)
	ensureEquals(t, expected, result)
}
