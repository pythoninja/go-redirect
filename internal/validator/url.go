package validator

import (
	"net/url"
	"regexp"
)

func ValidateURL(v *Validator, u *url.URL) {
	if u == nil {
		v.AddError("url", "cannot be null")

		return
	}

	v.check(u.Scheme == "http" || u.Scheme == "https", "protocol", "must be http or https")
	v.check(notBlank(u.Host), "hostname", "must not be empty")
}

func ValidateAlias(v *Validator, alias string) {
	re := regexp.MustCompile("^[0-9A-Za-z_-]+$")

	v.check(notBlank(alias), "alias", "must not be empty")
	v.check(maxChars(alias, 15), "alias", "must be less then 16 symbols")
	v.check(minChars(alias, 3), "alias", "must be more then 2 symbols")
	v.check(matches(alias, re), "alias", "must contain alphabetical (both uppercase and lowercase) "+
		"characters (A-Z and a-z), underscore, and dash symbols")
	v.check(!startsWith(alias, "-", "_"), "alias", "must start with alphabetical or "+
		"numerical characters (A-Z, a-z, 0-9)")
}
