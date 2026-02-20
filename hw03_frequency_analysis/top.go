package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type WordStat struct {
	word  string
	count int
}

func Top10(s string) []string {
	res := make([]string, 0, 10)

	m := make(map[string]int)

	sl := strings.Fields(s)

	for _, v := range sl {
		hasLetterOrDigit := false

		for _, r := range v {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				hasLetterOrDigit = true
				break
			}
		}

		if !hasLetterOrDigit {
			if len(v) > 1 {
				m[v]++
			}
			continue
		}

		word := strings.TrimFunc(v, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if word == "" {
			continue
		}

		isUpper := false
		for _, r := range word {
			if unicode.IsUpper(r) {
				isUpper = true
				break
			}
		}

		if isUpper {
			word = strings.ToLower(word)
		}

		m[word]++
	}

	tmp := make([]WordStat, 0, len(m))

	for k, v := range m {
		tmp = append(tmp, WordStat{word: k, count: v})
	}

	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i].count != tmp[j].count {
			return tmp[i].count > tmp[j].count
		}

		return tmp[i].word < tmp[j].word
	})

	limit := 10

	if limit > len(tmp) {
		limit = len(tmp)
	}

	for _, v := range tmp[:limit] {
		res = append(res, v.word)
	}

	return res
}
