package main

import (
	"regexp"
	"strings"
)

var constantDict = map[string]string{
	"<char>":      ".",
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
	"<newline>":   "\n",
	"<tab>":       "\t",
	"<return>":    "\r",
	"<null>":      "\\0",
	"<feed>":      "\\f",
	"<vertical>":  "\\v",
}

// todo add in ranges, "a to z; 1-3 of a to z"
var literalStringCommand, _ = regexp.Compile(`^".*"$`)
var rawStringCommand, _ = regexp.Compile("^`.*?`$")
var amountOfCommand, _ = regexp.Compile(`(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`)
var rangeOfCommand, _ = regexp.Compile(`(?P<amountStart>\d+)\s+to\s+(?P<amountEnd>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`)
var atLeastOfCommand, _ = regexp.Compile(`at\s+least\s+(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`)
var atMostOfCommand, _ = regexp.Compile(`at\s+most\s+(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`)
var someOfCommand, _ = regexp.Compile(`some\s+of\s+(?P<request>[\w\"\s<>]+)`)
var anyOfCommand, _ = regexp.Compile(`any\s+of\s+(?P<request>[\w\"\s<>]+)`)
var maybeOfCommand, _ = regexp.Compile(`maybe\s+of\s+(?P<request>[\w\"\s<>]+)`)

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

func normalize(expression string) string {
	if val, ok := constantDict[expression]; ok {
		return val
	}
	return expression
}

func parseLine(line string) string {
	if literalStringCommand.MatchString(line) {
		return regexp.QuoteMeta(line[1 : len(line)-1])

	} else if rawStringCommand.MatchString(line) {
		return line[1 : len(line)-1]
	} else if resultMap, ok := doesMatchRegex(rangeOfCommand, line); ok {
		start := resultMap["amountStart"]
		end := resultMap["amountEnd"]
		what := getWhat(resultMap["request"])
		return "(" + what + ")" + "{" + start + "," + end + "}"
	} else if resultMap, ok := doesMatchRegex(atLeastOfCommand, line); ok {
		num := resultMap["amount"]
		what := getWhat(resultMap["request"])
		return "(" + what + ")" + "{" + num + ",}"
	} else if resultMap, ok := doesMatchRegex(atMostOfCommand, line); ok {
		num := resultMap["amount"]
		what := getWhat(resultMap["request"])
		return "(" + what + ")" + "{," + num + "}"
	} else if resultMap, ok := doesMatchRegex(amountOfCommand, line); ok {
		num := resultMap["amount"]
		what := getWhat(resultMap["request"])
		return "(" + what + ")" + "{" + num + "}"
	} else if resultMap, ok := doesMatchRegex(someOfCommand, line); ok {
		what := getWhat(resultMap["request"])
		return normalize(what) + "+"
	} else if resultMap, ok := doesMatchRegex(anyOfCommand, line); ok {
		what := getWhat(resultMap["request"])
		return normalize(what) + "*"
	} else if resultMap, ok := doesMatchRegex(maybeOfCommand, line); ok {
		what := getWhat(resultMap["request"])
		return normalize(what) + "?"
	} else {
		return normalize(line)
	}
}

func parseInteriorOfCommand(interior, joiningChar string) string {
	pieces := strings.Split(interior, ";")
	var parsedPieces []string
	for _, piece := range pieces {
		if len(piece) == 0 {
			continue
		}
		parsedLine := parseLine(strings.TrimSpace(piece))
		parsedPieces = append(parsedPieces, parsedLine)
	}
	return strings.Join(parsedPieces, joiningChar)
}

func handleCommand(command, interior string) string {

	if command == "before" {
		return "`(?<=" + parseInteriorOfCommand(interior, "") + ")`;"
	} else if command == "after" {
		return "`(?=" + parseInteriorOfCommand(interior, "") + ")`;"
	} else if command == "match" {
		return "`(?:" + parseInteriorOfCommand(interior, "") + ")`;"
	} else if command == "either" {
		return "`(?:" + parseInteriorOfCommand(interior, "|") + ")`;"
	} else if strings.HasPrefix(command, "capture as ") {
		variable := strings.TrimSpace(command[11 : len(command)-1])
		return "`(?<" + variable + ">" + parseInteriorOfCommand(interior, "") + ")`;"
	} else {
		panic("Unknown command: " + command)
	}
}

func simplifyBlocks(data string) string {
	blockCommand := regexp.MustCompile(`(?P<pattern>[\w\s]+\{(.|\n)*?};)`)
	interiorCommand := regexp.MustCompile(`(?P<command>[\w\s<>]+)\s+\{\n(?P<interior>(.|\s)+?)\};`)

	blocks := blockCommand.FindAllString(data, -1)
	for _, block := range blocks {
		if propertyMap, ok := doesMatchRegex(interiorCommand, block); ok {
			command := strings.TrimSpace(propertyMap["command"])
			replacementStr := handleCommand(command, propertyMap["interior"])
			data = strings.ReplaceAll(data, block, replacementStr)
		} else {
			panic("Unable to parse" + block)
		}
	}

	return data
}

func parse(data string) string {
	simplifiedData := simplifyBlocks(data)
	pieces := strings.Split(simplifiedData, ";")
	regOut := ""

	for _, piece := range pieces {
		regOut += parseLine(strings.TrimSpace(piece))
	}

	return regOut
}
