package connect

import "strings"

func IsContextURI(uri string) bool {
	return strings.Contains(uri, ":album:") ||
		strings.Contains(uri, ":playlist:") ||
		strings.Contains(uri, ":show:")
}

func idFromURI(uri string) string {
	parts := strings.Split(uri, ":")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

func typeFromURI(uri string) string {
	parts := strings.Split(uri, ":")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
