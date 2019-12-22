package addressparser

import (
	"strconv"
	"strings"
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func asString(list []string, delimiter rune, startInc, stopExc int) string{
	if (startInc >= len(list)) || (startInc >= stopExc) {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(list[startInc])
	for index := startInc + 1; index < stopExc; index++ {
		sb.WriteRune(delimiter)
		sb.WriteString(list[index])
	}

	return sb.String()
}

func insert(position int, element string, list *[]string) {
	if position > len(*list) {
		position = len(*list)
	}
	*list = append(
		(*list)[:position],
		append(
			[]string{element},
			(*list)[position:]...
			)...
		)
}
