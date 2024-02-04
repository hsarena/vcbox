package util

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToF64(i64 []int64, precision float64) []float64 {
	f64 := make([]float64, len(i64))
	var ii int64
	var i int
	for i, ii = range i64 {
		f64[i] = float64(ii) / precision
	}
	return f64
}

func MetricIdToString(m string) string {
	words := strings.Split(m, ".")
	key := "_:"
	for _, word := range words {
		key += cases.Title(language.AmericanEnglish, cases.NoLower).String(word) + " "
	}
	return key
}
