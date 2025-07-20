package options

import "github.com/spf13/pflag"

// SecureOptions contains information for TLS.
type SecureOptions struct {
	KeyFile  string `json:"key-file"               mapstructure:"key-file"`
	CertFile string `json:"cert-file"              mapstructure:"cert-file"`
}

// GetDefaultSecureOptions returns a Secure configuration with default values.
func GetDefaultSecureOptions() *SecureOptions {
	return &SecureOptions{
		KeyFile:  "",
		CertFile: "",
	}
}

// AddFlags adds flags for a specific Server to the specified FlagSet.
func (s *SecureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.KeyFile, "key-file", s.KeyFile,
		"The private key file path.")
	fs.StringVar(&s.CertFile, "cert-file", s.CertFile,
		"The certificate file path.")
}

type ClientSecureOptions struct {
	ServerCertFile string `json:"server-cert-file"              mapstructure:"server-cert-file"`
}

// GetDefaultSecureOptions returns a Secure configuration with default values.
func GetDefaultClientSecureOptions() *ClientSecureOptions {
	return &ClientSecureOptions{
		ServerCertFile: "",
	}
}

// AddFlags adds flags for a specific Server to the specified FlagSet.
func (s *ClientSecureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServerCertFile, "server-cert-file", s.ServerCertFile,
		"The server certificate file path.")
}
