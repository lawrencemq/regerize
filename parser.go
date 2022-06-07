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

var amountOfCommand, _ = regexp.Compile(`(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`)
var rangeCommand, _ = regexp.Compile(`(?P<amountStart>\d+)\s+to\s+(?P<amountEnd>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`)
var atLeastCommand, _ = regexp.Compile(`at\s+least\s+(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`)
var someOfCommand, _ = regexp.Compile(`some\s+of\s+(?P<request>[\w\"\s<>]+)`)
var anyOfCommand, _ = regexp.Compile(`any\s+of\s+(?P<request>[\w\"\s<>]+)`)
var maybeCommand, _ = regexp.Compile(`maybe\s+of\s+(?P<request>[\w\"\s<>]+)`)

func doesMatchRegex(r *regexp.Regexp, str string) (map[string]string, bool) {
	subMatchMap := make(map[string]string)
	if !r.MatchString(str) {
		return subMatchMap, false
	}

	match := r.FindStringSubmatch(str)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = strings.TrimSpace(match[i])
		}
	}

	return subMatchMap, true
}

func getWhat(piece string) string {
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
		if resultMap, ok := doesMatchRegex(rangeCommand, piece); ok {
			start := resultMap["amountStart"]
			end := resultMap["amountEnd"]
			what := getWhat(resultMap["request"])
			regOut += "(" + what + ")" + "{" + start + "," + end + "}"
		} else if resultMap, ok := doesMatchRegex(atLeastCommand, piece); ok {
			num := resultMap["amount"]
			what := getWhat(resultMap["request"])
			regOut += "(" + what + ")" + "{" + num + ",}"
		} else if resultMap, ok := doesMatchRegex(amountOfCommand, piece); ok {
			num := resultMap["amount"]
			what := getWhat(resultMap["request"])
			regOut += "(" + what + ")" + "{" + num + "}"
		} else if resultMap, ok := doesMatchRegex(someOfCommand, piece); ok {
			what := getWhat(resultMap["request"])
			regOut += normalize(what) + "+"
		} else if resultMap, ok := doesMatchRegex(anyOfCommand, piece); ok {
			what := getWhat(resultMap["request"])
			regOut += normalize(what) + "*"
		} else if resultMap, ok := doesMatchRegex(maybeCommand, piece); ok {
			what := getWhat(resultMap["request"])
			regOut += normalize(what) + "?"
		} else if canNormalize(piece) {
			regOut += normalize(piece)
		} else {
			regOut += piece
		}
	}

	return regOut
}
