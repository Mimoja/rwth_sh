package common

import "strings"

func PrepSubdomain(v string) string {
	v = strings.ToLower(v)

	// replace all unicode characters that are … problematic
	// looking at you ZWJ (domain encoding of ZWJ is problematic at best)
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
