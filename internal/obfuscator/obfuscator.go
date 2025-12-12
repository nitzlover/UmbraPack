package obfuscator

import (
	"strings"

	"github.com/nitzlover/UmbraPack/internal/crypto"
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
	var result strings.Builder

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		if strings.Contains(line, `"`) {
			parts := strings.Split(line, `"`)
			for j := 1; j < len(parts); j += 2 {
				if parts[j] != "" {
					encrypted, err := o.cryptor.Encrypt(parts[j])
					if err != nil {
						return "", err
					}
					parts[j] = encrypted
				}
			}
			line = strings.Join(parts, `"`)
		}
		result.WriteString(line)
	}

	return result.String(), nil
}

func (o *Obfuscator) DeobfuscateStrings(code string) (string, error) {
	lines := strings.Split(code, "\n")
	var result strings.Builder

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		if strings.Contains(line, `"`) {
			parts := strings.Split(line, `"`)
			for j := 1; j < len(parts); j += 2 {
				if parts[j] != "" {
					decrypted, err := o.cryptor.Decrypt(parts[j])
					if err == nil {
						parts[j] = decrypted
					}
				}
			}
			line = strings.Join(parts, `"`)
		}
		result.WriteString(line)
	}

	return result.String(), nil
}
