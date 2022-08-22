package parser

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

var constantDict = map[string]string{
	"<char>":      ".",
	"<space>":     "\\s",
	"<alpha>":     "[a-zA-Z]",
	"<alphanum>":  "\\w",
	"<word>":      "\\w",
	"<num>":       "\\d",
	"<!alpha>":    "[!a-zA-Z]",
	"<!num>":      "[!0-9]",
	"<!alphanum>": "[!0-9a-zA-Z]",
	"<!word>":     "[!\\w]",
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

var literalStringCommand, _ = regexp2.Compile(`^".*"$`, regexp2.RE2)
var rawStringCommand, _ = regexp2.Compile("^`.*?`$", regexp2.RE2)
var amountOfCommand, _ = regexp2.Compile(`(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`, regexp2.RE2)
var rangeOfCommand, _ = regexp2.Compile(`(?P<amountStart>\d+)\s+to\s+(?P<amountEnd>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`, regexp2.RE2)
var atLeastOfCommand, _ = regexp2.Compile(`at\s+least\s+(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`, regexp2.RE2)
var atMostOfCommand, _ = regexp2.Compile(`at\s+most\s+(?P<amount>\d+)\s+of\s+(?P<request>[\w\"\s<>]+)`, regexp2.RE2)
var someOfCommand, _ = regexp2.Compile(`some\s+of\s+(?P<request>[\w\"\s<>]+)`, regexp2.RE2)
var anyOfCommand, _ = regexp2.Compile(`any\s+of\s+(?P<request>[\w\"\s<>]+)`, regexp2.RE2)
var maybeOfCommand, _ = regexp2.Compile(`maybe\s+of\s+(?P<request>[\w\"\s<>]+)`, regexp2.RE2)

// commenting
var multiLineComment, _ = regexp2.Compile(`/\*{2}[\s\S]*?\*/`, regexp2.RE2)
var singleComment, _ = regexp2.Compile(`//[\S\s]*?$`, regexp2.Multiline)

var importRegex, _ = regexp2.Compile(`^#import\s+(?P<filename>[.\w\-\/]+);`, regexp2.RE2)

func doesMatchRegex(r *regexp2.Regexp, str string) (map[string]string, bool) {
	subMatchMap := make(map[string]string)
	if matches, _ := r.MatchString(str); !matches {
		return subMatchMap, false
	}

	match, err := r.FindStringMatch(str)
	if err != nil {
		return subMatchMap, false
	}

	gps := match.Groups()
	for _, g := range gps {
		subMatchMap[g.Name] = strings.TrimSpace(g.String())
	}

	return subMatchMap, true
}

func getWhat(piece string) (string, error) {
	if piece[0] == '"' && piece[len(piece)-1] == '"' {
		return piece[1 : len(piece)-1], nil
	} else if piece[0] == '<' && piece[len(piece)-1] == '>' {
		val, ok := constantDict[piece]
		if !ok {
			return "", errors.New("No constant named " + piece)
		}
		return val, nil
	}
	return piece, nil

}

func normalize(expression string) string {
	if val, ok := constantDict[expression]; ok {
		return val
	}
	return expression
}

func parseLine(line string) (string, error) {
	if val, ok := userVariables[line]; ok {
		return val, nil
	} else if matches, _ := literalStringCommand.MatchString(line); matches {
		return regexp.QuoteMeta(line[1 : len(line)-1]), nil
	} else if matches, _ := rawStringCommand.MatchString(line); matches {
		return line[1 : len(line)-1], nil
	} else if resultMap, ok := doesMatchRegex(rangeOfCommand, line); ok {
		start := resultMap["amountStart"]
		end := resultMap["amountEnd"]
		what, wErr := getWhat(resultMap["request"])
		if wErr != nil {
			return "", wErr
		}
		return "(" + what + ")" + "{" + start + "," + end + "}", nil
	} else if resultMap, ok := doesMatchRegex(atLeastOfCommand, line); ok {
		num := resultMap["amount"]
		what, wErr := getWhat(resultMap["request"])
		if wErr != nil {
			return "", wErr
		}
		return "(" + what + ")" + "{" + num + ",}", nil
	} else if resultMap, ok := doesMatchRegex(atMostOfCommand, line); ok {
		num := resultMap["amount"]
		what, wErr := getWhat(resultMap["request"])
		if wErr != nil {
			return "", wErr
		}
		return "(" + what + ")" + "{," + num + "}", nil
	} else if resultMap, ok := doesMatchRegex(amountOfCommand, line); ok {
		num := resultMap["amount"]
		request := resultMap["request"]
		constant, isAConstant := constantDict[request]
		if isAConstant {
			return constant + "{" + num + "}", nil
		}
		what, wErr := getWhat(resultMap["request"])
		if wErr != nil {
			return "", wErr
		}
		return "(" + what + ")" + "{" + num + "}", nil
	} else if resultMap, ok := doesMatchRegex(someOfCommand, line); ok {
		what, wErr := getWhat(resultMap["request"])
		if wErr != nil {
			return "", wErr
		}
		return normalize(what) + "+", nil
	} else if resultMap, ok := doesMatchRegex(anyOfCommand, line); ok {
		what, wErr := getWhat(resultMap["request"])
		if wErr != nil {
			return "", wErr
		}
		return normalize(what) + "*", nil
	} else if resultMap, ok := doesMatchRegex(maybeOfCommand, line); ok {
		what, wErr := getWhat(resultMap["request"])
		if wErr != nil {
			return "", wErr
		}
		return normalize(what) + "?", nil
	} else {
		return normalize(line), nil
	}
}

func removeCommentWithRegex(r *regexp2.Regexp, data string) (string, error) {
	if matches, _ := r.MatchString(data); !matches {
		return data, nil
	}

	newData := strings.Clone(data)

	match, err := r.FindStringMatch(data)
	if err != nil {
		return "", err
	}

	gps := match.Groups()
	for _, g := range gps {
		newData = strings.ReplaceAll(newData, g.String(), "")
	}

	return newData, nil

}

func removeComments(data string) (string, error) {
	dataSansMlComments, err := removeCommentWithRegex(multiLineComment, data)
	if err != nil {
		return "", err
	}

	return removeCommentWithRegex(singleComment, dataSansMlComments)

}

func parseInteriorOfCommand(interior, joiningChar string) (string, error) {
	pieces := strings.Split(interior, ";")
	var parsedPieces []string
	for _, piece := range pieces {
		if len(piece) == 0 {
			continue
		}
		parsedLine, parseError := parseLine(strings.TrimSpace(piece))
		if parseError != nil {
			return "", parseError
		}
		parsedPieces = append(parsedPieces, parsedLine)
	}
	return strings.Join(parsedPieces, joiningChar), nil
}

func handleCommand(command, interior string) (string, error) {

	if command == "before" {
		interior, err := parseInteriorOfCommand(interior, "")
		if err != nil {
			return "", err
		}
		return "`(?<=" + interior + ")`;", nil
	} else if command == "after" {
		interior, err := parseInteriorOfCommand(interior, "")
		if err != nil {
			return "", err
		}
		return "`(?=" + interior + ")`;", nil
	} else if command == "match" {
		interior, err := parseInteriorOfCommand(interior, "")
		if err != nil {
			return "", err
		}
		return "`(?:" + interior + ")`;", nil
	} else if command == "either" {
		interior, err := parseInteriorOfCommand(interior, "|")
		if err != nil {
			return "", err
		}
		return "`(?:" + interior + ")`;", nil
	} else if strings.HasPrefix(command, "capture as ") {
		variable := strings.TrimSpace(command[11 : len(command)-1])
		interior, err := parseInteriorOfCommand(interior, "")
		if err != nil {
			return "", err
		}
		return "`(?<" + variable + ">" + interior + ")`;", nil
	} else if strings.HasPrefix(command, "let .") {
		variableName := strings.TrimSpace(command[4:])
		value, err := parseInteriorOfCommand(interior, "")
		if err != nil {
			return "", err
		}
		userVariables[variableName] = value
		return "", nil
	} else {
		return "", errors.New("Unknown command: " + command)
	}
}

func simplifyBlocks(data string) (string, error) {
	blockCommand := regexp2.MustCompile(`(?m)^(?P<pattern>[\w\s\.]+\{(.|\n)*?};)`, regexp2.RE2)
	interiorCommand := regexp2.MustCompile(`^(?P<command>[\w\s<>\.]+)\s*\{\n(?P<interior>(.|\s)+?)\}`, regexp2.RE2)

	matches, _ := blockCommand.FindStringMatch(data)
	for matches != nil {
		block := matches.String()
		if propertyMap, ok := doesMatchRegex(interiorCommand, block); ok {
			command := strings.TrimSpace(propertyMap["command"])
			replacementStr, error := handleCommand(command, propertyMap["interior"])
			if error != nil {
				return "", error
			}
			data = strings.ReplaceAll(data, block, replacementStr)
			matches, _ = blockCommand.FindNextMatch(matches)
		} else {
			return "", errors.New("Unable to parse" + block)
		}
	}

	return data, nil
}

func parse(data string) (string, error) {
	dataSansComments, commentError := removeComments(data)
	if commentError != nil {
		return "", commentError
	}
	simplifiedData, error := simplifyBlocks(dataSansComments)
	if error != nil {
		return "", error
	}
	pieces := strings.Split(simplifiedData, ";")
	regOut := ""

	for _, piece := range pieces {
		line, err := parseLine(strings.TrimSpace(piece))
		if err != nil {
			return "", err
		}
		regOut += line
	}

	return regOut, nil
}

func Parse(data string) (string, error) {
	regexStr, parseError := parse(data)
	if parseError != nil {
		return "", parseError
	}
	_, error := regexp2.Compile(regexStr, regexp2.RE2)
	return regexStr, error
}

type stack []string

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}
func (s *stack) push(strs ...string) {
	if len(strs) > 0 {
		*s = append(*s, strs...)
	}

}

func (s *stack) pop() (string, bool) {
	if s.isEmpty() {
		return "", false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}
func (s *stack) JoinInOrder() string {
	dst := make([]string, len(*s))
	copy(dst, *s)

	// reversing slick
	for i, j := 0, len(dst)-1; i < j; i, j = i+1, j-1 {
		dst[i], dst[j] = dst[j], dst[i]
	}
	return strings.Join(dst, "\n")
}

func readInFile(filename string) (string, error) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(dat), nil
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

	// todo check for error
	result, _ := importRegex.Replace(contents, "", -1, -1)
	return result

}

func ParseFile(filename string) (string, error) {

	setOfFilesRead := map[string]bool{}
	contents := &stack{}
	fileStack := &stack{}
	fileStack.push(filename)

	for !fileStack.isEmpty() {
		toRead, _ := fileStack.pop()
		if _, exists := setOfFilesRead[toRead]; exists {
			return "", errors.New("Double import of file " + toRead)
		}

		setOfFilesRead[toRead] = true
		fileContents, error := readInFile(toRead)
		if error != nil {
			return "", error
		}
		fileStack.push(findImports(fileContents)...)
		contents.push(removeImportsFromContents(fileContents))
	}

	allContents := contents.JoinInOrder()

	regexStr, parseError := parse(allContents)
	if parseError != nil {
		return "", parseError
	}
	_, error := regexp2.Compile(regexStr, regexp2.RE2)
	return regexStr, error

}
