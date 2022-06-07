package main

import "testing"

func TestSum(t *testing.T) {
	a := 1
	b := 2
	s := sum(a, b)
	if s != 3 {
		t.Fatalf(`Sum did not equal 3`)
	}
}

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
