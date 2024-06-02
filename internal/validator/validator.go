package validator

import (
	"net/url"
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

func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func (v *Validator) MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func (v *Validator) MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func ValidateURL(v *Validator, u *url.URL) {
	if u == nil {
		v.AddError("url", "cannot be null")
		return
	}

	v.Check(u.Scheme == "http" || u.Scheme == "https", "protocol", "must be http or https")
	v.Check(v.NotBlank(u.Host), "hostname", "must be not empty")
}

func ValidateAlias(v *Validator, alias string) {
	v.Check(v.NotBlank(alias), "alias", "must be not empty")
	v.Check(v.MaxChars(alias, 15), "alias", "must be less then 16 symbols")
	v.Check(v.MinChars(alias, 3), "alias", "must be more then 2 symbols")
}
