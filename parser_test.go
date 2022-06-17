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

func TestCapture(t *testing.T) {
	input := `capture as shiba {
"inu";
<space>;
at least 5 of "wow";
};`
	expected := "(?<shib>inu\\s(wow){5,})"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestMatch(t *testing.T) {
	input := `match {
		16 of "na";
		<space>;
		"BATMAN";
		};`
	expected := "(?:(na){16}\\sBATMAN)"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestEither(t *testing.T) {
	input := `either {
		"corgi";
		"shiba";
		"husky";
		};`
	expected := "(?:corgi|shiba|husky)"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestVariable(t *testing.T) {
	input := `let .uuidv4 {
8 of <hex>;
"-";
4 of <hex>;
"-";
4 of <hex>: 
"-";
4 of <hex>: 
"-";
12 of <hex>: 
};
"user id:";
.uuidv4;`

	expected := "user id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}[0-9a-fA-F]{4}[0-9a-fA-F]{12}"
	result := parse(input)
	ensureEquals(t, expected, result)
}

func TestImport(t *testing.T) {
	startingFile := "test_data/import1.rgr"

	expected := "(?:(ring ){8},\\sba(na){2}\\sphone)"
	result := parseFile(startingFile)
	ensureEquals(t, expected, result)
}
