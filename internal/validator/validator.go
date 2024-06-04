package validator

import (
	"net/url"
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

func (v *Validator) notBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func (v *Validator) maxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func (v *Validator) minChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
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

func (v *Validator) matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) startsWith(value string, characters ...string) bool {
	for _, s := range characters {
		if strings.HasPrefix(value, s) {
			return true
		}
	}

	return false
}

func ValidateURL(v *Validator, u *url.URL) {
	if u == nil {
		v.AddError("url", "cannot be null")
		return
	}

	v.check(u.Scheme == "http" || u.Scheme == "https", "protocol", "must be http or https")
	v.check(v.notBlank(u.Host), "hostname", "must not be empty")
}

func ValidateAlias(v *Validator, alias string) {
	re := regexp.MustCompile("^[0-9A-Za-z_-]+$")

	v.check(v.notBlank(alias), "alias", "must not be empty")
	v.check(v.maxChars(alias, 15), "alias", "must be less then 16 symbols")
	v.check(v.minChars(alias, 3), "alias", "must be more then 2 symbols")
	v.check(v.matches(alias, re), "alias", "must contain alphabetical (both uppercase and lowercase) "+
		"characters (A-Z and a-z), underscore, and dash symbols")
	v.check(!v.startsWith(alias, "-", "_"), "alias", "must start with alphabetical (both uppercase and lowercase "+
		"characters (A-Z and a-z)")
}
