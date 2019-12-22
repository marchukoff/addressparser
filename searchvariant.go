package addressparser

import (
	"fmt"
	"sort"
	"strconv"
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

func Parse(address string) []*SearchVariant {
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
				clone.addressParts[partIndex] = fmt.Sprintf("%s0%s", ProbablyDelimiters, subIndex)
				clone.addressParts = append(
					clone.addressParts[:partIndex+1],
					fmt.Sprintf("%s%d%d", ProbablyDelimiters, subIndex, len(subStrings)))
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

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (sv *SearchVariant) splitHyphens() {
	for index := len(sv.addressParts) - 1; index >= 0; index-- {
		subStrings := strings.FieldsFunc(sv.addressParts[index], defaultHyphens)

		for subIndex := len(subStrings) - 1; subIndex > 0; subIndex-- {
			if isNumeric(subStrings[subIndex]) {
				sv.addressParts = append(
					sv.addressParts[:index+1],
					append(
						[]string{subStrings[subIndex]}, sv.addressParts[index+1:]...)...,
				)
				sv.addressParts[index] = fmt.Sprintf("%s0%d", DefaultHyphens, subIndex)
			} else {
				break
			}
		}
	}
}

func (sv *SearchVariant) trim() {
	trimmed := make([]string, len(sv.addressParts))
	for _, part := range sv.addressParts {
		trimmed = append(trimmed, strings.Trim(part, " "))
	}
}

func (sv *SearchVariant) clone() *SearchVariant {
	addressParts := make([]string, len(sv.addressParts))
	copy(addressParts, sv.addressParts)
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

func (sv *SearchVariant) trimJunkWords() {
	for index, part := range sv.addressParts {
		if isJunkWords(part) {
			sv.addressParts = append(sv.addressParts[:index], sv.addressParts[index+1:]...)
		}
	}
}

func (sv *SearchVariant) toString() string {
	var sb strings.Builder
	sep := ""
	for _, part := range sv.addressParts {
		_, _ = fmt.Fprintf(&sb, "%s[%s]", sep, part)
		sep = " "
	}
	return sb.String()
}

func (sv *SearchVariant) equals(other *SearchVariant) bool {
	if len(sv.addressParts) != len(other.addressParts) {
		return false
	}

	for i, v := range sv.addressParts {
		if v != other.addressParts[i] {
			return false
		}
	}
	return true
}

func (sv *SearchVariant) splitDigitalAndLetter() {
	for index, part := range sv.addressParts {
		asRunes := []rune(part)
		beforeLastSymbol := string(asRunes[:len(asRunes)-2])
		symbol := string(asRunes[len(asRunes)-2:])
		if isNumeric(beforeLastSymbol) && isNumeric(symbol) {
			sv.addressParts[index] = beforeLastSymbol
			sv.addressParts = append(
				sv.addressParts[:index+1],
				append(
					[]string{symbol}, sv.addressParts[index+1:]...)...,
			) // addressParts.Insert(index+1, symbol)
		}
	}
}

func (sv *SearchVariant) getRank() float64 {
	var rank float64
	for _, part := range sv.addressParts {
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

func main() {
	return
}
