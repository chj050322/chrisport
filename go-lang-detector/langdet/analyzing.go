package langdet

import (
	"bytes"
	"sort"
	"strings"
	"unicode/utf8"
)

// maxSampleSize represents the maximum number of tokens per sample, low number can
// cause bad accuracy, but better performance.
// -1 for no maximum
var maxSampleSize = 100000

func Analyze(text, name string) Language {
	return AnalyzeWithNDepth(text, name, DEFAULT_NDEPTH)
}

// Analyze creates the language profile from a given Text and returns it in a Language struct.
func AnalyzeWithNDepth(text, name string, nDepth int) Language {
	theMap := CreateOccurenceMap(text, nDepth)
	//fmt.Printf("theMap:%+v\n", theMap)
	ranked := CreateRankLookupMap(theMap)
	//fmt.Printf("ranked:%+v\n", ranked)
	return Language{Name: name, Profile: ranked}
}

// creates the map [token] rank from a map [token] occurrence
func CreateRankLookupMap(input map[string]int) map[string]int {
	tokens := make([]Token, len(input))
	counter := 0
	for k, v := range input {
		tokens[counter] = Token{Key: k, Occurrence: v}
		counter++
	}
	sort.Sort(ByOccurrence(tokens))
	result := make(map[string]int)
	length := len(tokens)
	locMaxL := maxSampleSize
	if locMaxL < 0 {
		locMaxL = length
	}
	for i := length - 1; i >= 0 && i > length-locMaxL; i-- {
		result[tokens[i].Key] = length - i

		//fmt.Println("tokens[i].Key", tokens[i].Key, result[tokens[i].Key])
	}
	return result
}

// CreateOccurenceMap creates a map[token]occurrence from a given text and up to a given gram depth
// gramDepth=1 means only 1-letter tokens are created, gramDepth=2 means 1- and 2-letters token are created, etc.
func CreateOccurenceMap(text string, gramDepth int) map[string]int {
	text = cleanText(text)
	tokens := strings.Split(text, " ")
	result := make(map[string]int)

	//fmt.Println("len(tokens)", len(tokens))
	for _, token := range tokens {
		analyseToken(result, token, gramDepth)
	}
	return result
}

// analyseToken analyses a token to a certain gramDepth and stores the result in resultMap
func analyseToken(resultMap map[string]int, token string, gramDepth int) {
	if len(token) == 0 {
		return
	}
	for i := 1; i <= gramDepth+1; i++ {
		generateNthGrams(resultMap, token, i)
	}
}

// generateNthGrams creates n-gram tokens from the input string and
// adds the mapping from token to its number of occurrences to the resultMap
func generateNthGrams(resultMap map[string]int, text string, n int) {
	//padding := createPadding(n - 1)
	padding := " "

	if CalByteOffsetToRunOffset(text, 1) != 1 {
		padding = ""
	}

	text = padding + text + padding
	upperBound := utf8.RuneCountInString(text) - (n - 1)
	for p := 0; p < upperBound; p++ {

		p1 := CalByteOffsetToRunOffset(text, p)
		pn := CalByteOffsetToRunOffset(text, p+n)

		if p1 == pn || pn > len(text) || p1 >= len(text) {
			break
		}

		currentToken := text[p1:pn]

		if currentToken == " " {
			continue
		}

		resultMap[currentToken] += 1

		//fmt.Println("currentToken", currentToken, resultMap[currentToken], text)

	}
}

// createPadding surrounds text with a padding
func createPadding(length int) string {
	var buffer bytes.Buffer
	padding := "_"
	for i := 0; i < length; i++ {
		buffer.WriteString(padding)
	}
	return buffer.String()
}

// cleanText removes newlines, special characters and numbers from a input text
func cleanText(text string) string {
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, ",", " ", -1)
	text = strings.Replace(text, "#", " ", -1)
	text = strings.Replace(text, "/", " ", -1)
	text = strings.Replace(text, "\\", " ", -1)
	text = strings.Replace(text, ".", " ", -1)
	text = strings.Replace(text, "!", " ", -1)
	text = strings.Replace(text, "?", " ", -1)
	text = strings.Replace(text, ":", " ", -1)
	text = strings.Replace(text, ";", " ", -1)
	text = strings.Replace(text, "-", " ", -1)
	text = strings.Replace(text, "'", " ", -1)
	text = strings.Replace(text, "\"", " ", -1)
	text = strings.Replace(text, "_", " ", -1)
	text = strings.Replace(text, "*", " ", -1)
	text = strings.Replace(text, "1", "", -1)
	text = strings.Replace(text, "2", "", -1)
	text = strings.Replace(text, "3", "", -1)
	text = strings.Replace(text, "4", "", -1)
	text = strings.Replace(text, "5", "", -1)
	text = strings.Replace(text, "6", "", -1)
	text = strings.Replace(text, "7", "", -1)
	text = strings.Replace(text, "8", "", -1)
	text = strings.Replace(text, "9", "", -1)
	text = strings.Replace(text, "0", "", -1)
	text = strings.Replace(text, "  ", " ", -1)
	return text
}

func CalByteOffsetToRunOffset(in string, offset int) int {
	if offset < 0 || utf8.RuneCountInString(in) <= offset {
		return len(in)
	}

	if offset == 0 {
		return 0
	}
	for i := offset; i < len(in); i++ {
		sub := in[:i]
		err, _ := utf8.DecodeLastRuneInString(sub)
		if err != utf8.RuneError {
			if utf8.RuneCountInString(sub) == offset {
				return i
			}
		}
	}

	return len(in)
}
