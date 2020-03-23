package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
	"unicode/utf8"
)

func Top10(text string) []string {
	// Place your code here
	if text == "" {
		return nil
	}
	if !utf8.ValidString(text) {
		return nil
	}

	var (
		wordCounter = map[string]int{}
		uniq        []string
		top         = 10
		result      = make([]string, 0, top)
	)

	//re := regexp.MustCompile(`\p{L}+-\p{L}+|\p{L}+| - |\t- |\n- `)
	//words := re.FindAllString(text, -1)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return strings.ContainsRune("0123456789\t\n\r ", r)
	})

	if len(words) == 0 {
		return nil
	}

	for i, v := range words {
		v = strings.TrimSpace(v)
		words[i] = v
		// selecting only uniq
		if _, ok := wordCounter[v]; !ok {
			uniq = append(uniq, v)
		}
		// counting words
		wordCounter[v]++
	}

	sort.SliceStable(uniq, func(i, j int) bool {
		return wordCounter[uniq[i]] > wordCounter[uniq[j]]
	})

	if len(uniq) < top {
		top = len(uniq)
	}

	result = append(result, uniq[:top]...)

	return result
}
