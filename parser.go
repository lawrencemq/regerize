package main

import (
	"regexp"
	"strings"
)

var constantDict = map[string]string{
	"<space>":     "\\s",
	"<alpha>":     "[a-zA-Z]",
	"<alphanum>":  "\\w",
	"<word>":      "\\w",
	"<num>":       "\\d",
	"<!alpha>":    "",
	"<!num>":      "",
	"<!alphanum>": "",
	"<!word>":     "",
	"<start>":     "^",
	"<end>":       "$",
}

var ofCommand, _ = regexp.Compile("\\d+\\s+of[\\w\\\"\\s]+")
var someOfCommand, _ = regexp.Compile("some\\s+of[\\w\\\"\\s]+")

func sum(a, b int) int {
	return a + b
}

func getWhat(line string) string {
	piece := strings.TrimSpace(line[strings.LastIndex(line, " "):])
	if piece[0] == '"' && piece[len(piece)-1] == '"' {
		return piece[1 : len(piece)-1]
	}
	return piece

}

func canNormalize(expression string) bool {
	_, ok := constantDict[expression]
	return ok
}
func normalize(expression string) string {
	if val, ok := constantDict[expression]; ok {
		return val
	}
	return expression
}

func parse(data string) string {
	pieces := strings.Split(data, ";")
	regOut := ""

	for _, piece := range pieces {
		piece = strings.TrimSpace(piece)
		if ofCommand.MatchString(piece) {
			num := piece[:strings.Index(piece, " ")]
			what := getWhat(piece)
			regOut += what + "{" + num + "}"
		} else if someOfCommand.MatchString(piece) {
			what := getWhat(piece)
			regOut += normalize(what) + "+"
		} else if canNormalize(piece) {
			regOut += normalize(piece)
		} else {
			regOut += piece
		}
	}

	return regOut
}
