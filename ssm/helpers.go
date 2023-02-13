package ssm

import (
	"strings"
)

// buildKeyAliasPath generates alias string.
// It cuts `/` from begining of user provided string if it exists
func buildKeyAliasPath(kmsKeyAlias string) string {
	return "alias/" + removeSlashPrefix(kmsKeyAlias)
}

func removeSlashPrefix(s string) string {
	return strings.TrimPrefix(s, "/")
}
