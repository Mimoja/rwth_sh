package common

import "strings"

func PrepSubdomain(v string) string {
	v = strings.ToLower(v)

	/* replace all unicode characters that are … problematic

	ZWJ is being removed because there's not really consens
	how this should be encoded. The answer in this issue gives
	a small overview
	https://github.com/mathiasbynens/punycode.js/issues/114#issuecomment-891935957
	*/
	v = strings.ReplaceAll(v, "\u200d", "")

	return v
}

func PrepPath(v string) string {
	v = strings.ToLower(v)

	return v
}

func PrepareURI(subdomain, path string) (string, string) {
	return PrepSubdomain(subdomain), PrepPath(path)
}
