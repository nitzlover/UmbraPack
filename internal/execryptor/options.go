package execryptor

type Metadata struct {
	CompanyName     string
	FileDescription string
	ProductName     string
	FileVersion     string
	ProductVersion  string
}

type BuildOptions struct {
	Metadata           Metadata
	IconPath           string
	EnableObfuscation  bool
	KeepOriginalBinary bool
}
