package main

import "testing"

func ensureEquals(t *testing.T, expected, result string) {
	if expected != result {
		t.Fatalf("'%s' != '%s'", expected, result)
	}
}
func TestBatman(t *testing.T) {
	expected := "na{16}\\sbatman"
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
