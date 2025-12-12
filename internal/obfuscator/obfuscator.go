package obfuscator

import (
	"obfuscator/internal/crypto"
	"strings"
)

type Obfuscator struct {
	cryptor *crypto.Cryptor
}

func New(password string) *Obfuscator {
	return &Obfuscator{
		cryptor: crypto.NewCryptor(password),
	}
}

func (o *Obfuscator) ObfuscateStrings(code string) (string, error) {
	lines := strings.Split(code, "\n")
	var result []string

	for _, line := range lines {
		if strings.Contains(line, `"`) {
			parts := strings.Split(line, `"`)
			for i := 1; i < len(parts); i += 2 {
				if parts[i] != "" {
					encrypted, err := o.cryptor.Encrypt(parts[i])
					if err != nil {
						return "", err
					}
					parts[i] = encrypted
				}
			}
			line = strings.Join(parts, `"`)
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n"), nil
}

func (o *Obfuscator) DeobfuscateStrings(code string) (string, error) {
	lines := strings.Split(code, "\n")
	var result []string

	for _, line := range lines {
		if strings.Contains(line, `"`) {
			parts := strings.Split(line, `"`)
			for i := 1; i < len(parts); i += 2 {
				if parts[i] != "" {
					decrypted, err := o.cryptor.Decrypt(parts[i])
					if err == nil {
						parts[i] = decrypted
					}
				}
			}
			line = strings.Join(parts, `"`)
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n"), nil
}
