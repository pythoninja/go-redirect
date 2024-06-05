package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func notBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func maxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func minChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func startsWith(value string, characters ...string) bool {
	for _, s := range characters {
		if strings.HasPrefix(value, s) {
			return true
		}
	}

	return false
}
