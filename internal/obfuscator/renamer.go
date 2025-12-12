package obfuscator

import (
	"strings"
)

func RenameVariables(code string) string {
	replacements := map[string]string{
		"username": "a1",
		"password": "b2",
		"data":     "c3",
		"result":   "d4",
		"temp":     "e5",
		"value":    "f6",
		"config":   "g7",
		"handler":  "h8",
		"response": "i9",
		"request":  "j10",
	}

	for old, new := range replacements {
		code = strings.ReplaceAll(code, old, new)
	}

	return code
}

func RemoveComments(code string) string {
	lines := strings.Split(code, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "//") {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
