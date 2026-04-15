package pg_crud

import (
	"regexp"
	"strings"
)

var (
	matchFirstCapRe = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCapRe   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// PascalToSnake преобразования названия из PascalCase в snake_case
func PascalToSnake(str string) string {
	snake := matchFirstCapRe.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCapRe.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// SliceToMap сформировать карту элементов слайса для удобного обращения values[key] в условиях
func SliceToMap[T comparable](source []T) map[T]bool {
	res := make(map[T]bool, len(source))
	for _, s := range source {
		res[s] = true
	}
	return res
}

// Ptr получить ссылку на аргумент, удобно для получения ссылок на константы
func Ptr[T any](source T) *T {
	return &source
}
