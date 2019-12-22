package addressparser

import (
	"fmt"
	"sort"
	"strings"
)

const (
	DefaultHyphens     = '-'
	ProbablyDelimiters = ' '
)

func splitProbablyDelimiters(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case ProbablyDelimiters:
			return true
		default:
			return false
		}
	})
}

func defaultHyphens(r rune) bool {
	switch r {
	case DefaultHyphens:
		return true
	default:
		return false
	}
}

func defaultAddressSplitters(r rune) bool {
	switch r {
	case '.', ',', ';', '/', '\\':
		return true
	default:
		return false
	}
}

type SearchVariant struct {
	addressParts []string
}

func (search *SearchVariant) Parse(address string) []*SearchVariant {
	firstResult := &SearchVariant{}
	firstResult.addressParts = strings.FieldsFunc(address, defaultAddressSplitters)

	firstResult.trim()
	firstResult.splitHyphens()
	firstResult.trimJunkWords()
	firstResult.splitDigitalAndLetter()

	result := append([]*SearchVariant{}, firstResult)
	index := 0

	for index < len(result) {
		for partIndex, part := range result[index].addressParts {
			subStrings := splitProbablyDelimiters(part)
			for subIndex := range subStrings {
				clone := result[index].clone()
				clone.addressParts[partIndex] = asString(subStrings, ProbablyDelimiters, 0, subIndex)
				insert(partIndex + 1, asString(subStrings, ProbablyDelimiters, subIndex, len(subStrings)), &clone.addressParts)
				clone.trimJunkWords()
				clone.splitDigitalAndLetter()
				isDuplicate := false
				for _, resultItem := range result {
					if resultItem.equals(clone) {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					result = append(result, clone)
				}
			}
		}
		index++
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].getRank() < result[j].getRank()
	})
	return result
}

func (search *SearchVariant) splitHyphens() {
	for index := len(search.addressParts) - 1; index >= 0; index-- {
		subStrings := strings.FieldsFunc(search.addressParts[index], defaultHyphens)

		for subIndex := len(subStrings) - 1; subIndex > 0; subIndex-- {
			if isNumeric(subStrings[subIndex]) {
				insert(index + 1, subStrings[subIndex], &search.addressParts)
				search.addressParts[index] = asString(subStrings, DefaultHyphens, 0, subIndex)
			} else {
				break
			}
		}
	}
}

func (search *SearchVariant) trim() {
	trimmed := make([]string, len(search.addressParts))
	for _, part := range search.addressParts {
		trimmed = append(trimmed, strings.Trim(part, " "))
	}
}

func (search *SearchVariant) clone() *SearchVariant {
	addressParts := make([]string, len(search.addressParts))
	copy(addressParts, search.addressParts)
	clone := &SearchVariant{
		addressParts: addressParts,
	}
	return clone
}

func isJunkWords(word string) bool {
	JunkWords := map[string]struct{}{
		"область":    {},
		"обл":        {},
		"город":      {},
		"гор":        {},
		"район":      {},
		"р-н":        {},
		"р-он":       {},
		"улица":      {},
		"ул":         {},
		"дом":        {},
		"квартира":   {},
		"кв":         {},
		"проспект":   {},
		"пр-кт":      {},
		"пр":         {},
		"микрорайон": {},
		"м-р-н":      {},
		"мкр":        {},
		"мкрн":       {},
		"копрус":     {},
		"корп":       {},
		"литера":     {},
		"лит":        {},
		"бульвар":    {},
		"б-р":        {},
		"поселок":    {},
		"посёлок":    {},
		"пос":        {},
		"квартал":    {},
		"кв-л":       {},
		"квл":        {},
		"кварт":      {},
	}
	word = strings.ToLower(word)
	_, ok := JunkWords[word]
	return ok
}

func (search *SearchVariant) trimJunkWords() {
	for index, part := range search.addressParts {
		if isJunkWords(part) {
			search.addressParts = append(search.addressParts[:index], search.addressParts[index+1:]...)
		}
	}
}

func (search *SearchVariant) toString() string {
	var sb strings.Builder
	sep := ""
	for _, part := range search.addressParts {
		_, _ = fmt.Fprintf(&sb, "%s[%s]", sep, part)
		sep = " "
	}
	return sb.String()
}

func (search *SearchVariant) equals(other *SearchVariant) bool {
	if len(search.addressParts) != len(other.addressParts) {
		return false
	}

	for i, v := range search.addressParts {
		if v != other.addressParts[i] {
			return false
		}
	}
	return true
}

func (search *SearchVariant) splitDigitalAndLetter() {
	for index, part := range search.addressParts {
		asRunes := []rune(part)
		beforeLastSymbol := string(asRunes[:len(asRunes)-2])
		symbol := string(asRunes[len(asRunes)-2:])
		if isNumeric(beforeLastSymbol) && isNumeric(symbol) {
			search.addressParts[index] = beforeLastSymbol
			insert(index + 1, symbol, &search.addressParts)
		}
	}
}

func (search *SearchVariant) getRank() float64 {
	var rank float64
	for _, part := range search.addressParts {
		partRank := float64(len(part) * len(part))
		for _, subPart := range splitProbablyDelimiters(part) {
			if isJunkWords(subPart) {
				partRank *= partRank
			}
		}
		rank += partRank
	}
	return rank
}
