package main

import (
	"os"
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
	"<hex>":       "[0-9a-fA-F]",
}

var userVariables = map[string]string{}

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

// other useful regexes
var importRegex, _ = regexp.Compile(`^#import\s+(?P<filename>[.\w\-\/]+);`)

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
	} else if piece[0] == '<' && piece[len(piece)-1] == '>' {
		val, ok := constantDict[piece]
		if !ok {
			panic("No constant named " + piece)
		}
		return val
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
	if val, ok := userVariables[line]; ok {
		return val
	} else if literalStringCommand.MatchString(line) {
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
		request := resultMap["request"]
		constant, isAConstant := constantDict[request]
		if isAConstant {
			return constant + "{" + num + "}"
		}
		what := getWhat(request)
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
	} else if strings.HasPrefix(command, "let .") {
		variableName := strings.TrimSpace(command[4:])
		value := parseInteriorOfCommand(interior, "")
		userVariables[variableName] = value
		return ""
	} else {
		panic("Unknown command: " + command)
	}
}

func simplifyBlocks(data string) string {
	blockCommand := regexp.MustCompile(`(?m)^(?P<pattern>[\w\s\.]+\{(.|\n)*?};)`)
	interiorCommand := regexp.MustCompile(`^(?P<command>[\w\s<>\.]+)\s*\{\n(?P<interior>(.|\s)+?)\}`)

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

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}
func (s *Stack) Push(strs ...string) {
	if len(strs) > 0 {
		*s = append(*s, strs...)
	}

}

func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}
func (s *Stack) JoinInOrder() string {
	dst := make([]string, len(*s))
	copy(dst, *s)

	// reversing slick
	for i, j := 0, len(dst)-1; i < j; i, j = i+1, j-1 {
		dst[i], dst[j] = dst[j], dst[i]
	}
	return strings.Join(dst, "\n")
}

func readInFile(filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic("Unable to read " + filename + ". " + err.Error())
	}
	return string(dat)
}

func findImports(contents string) []string {

	filenamesFound := []string{}
	if propertyMap, ok := doesMatchRegex(importRegex, contents); ok {
		filename := strings.TrimSpace(propertyMap["filename"]) + ".rgr"
		filenamesFound = append(filenamesFound, filename)
	}
	return filenamesFound
}

func removeImportsFromContents(contents string) string {

	return importRegex.ReplaceAllString(contents, "")

}

func parseFile(filename string) string {

	setOfFilesRead := map[string]bool{}
	contents := &Stack{}
	fileStack := &Stack{}
	fileStack.Push(filename)

	for !fileStack.IsEmpty() {
		toRead, _ := fileStack.Pop()
		if _, exists := setOfFilesRead[toRead]; exists {
			panic("Double import of file " + toRead)
		}

		setOfFilesRead[toRead] = true
		fileContents := readInFile(toRead)
		fileStack.Push(findImports(fileContents)...)
		contents.Push(removeImportsFromContents(fileContents))
	}

	allContents := contents.JoinInOrder()

	return parse(allContents)

}
