package options

import "github.com/spf13/pflag"

//ServerOptions contains information for a client service.
type ServerOptions struct {
	BindEndpoint      string `json:"bind-endpoint"           mapstructure:"bind-endpoint"`
	TokenParserPlugin string `json:"token-parser-plugin" mapstructure:"token-parser-plugin"`
	TokenParserKey    string `json:"token-parser-key"    mapstructure:"token-parser-key"`
}

// GetDefaultServerOptions returns a server configuration with default values.
func GetDefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		BindEndpoint:      "[::]:1220",
		TokenParserPlugin: "Cleartext",
		TokenParserKey:    "",
	}
}

// AddFlags adds flags for a specific Server to the specified FlagSet.
func (s *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindEndpoint, "bind-endpoint", s.BindEndpoint,
		"The socket that the server side endpoint listen on")
	fs.StringVar(&s.TokenParserPlugin, "token-parser-plugin", s.TokenParserPlugin,
		"The token parser plugin.")
	fs.StringVar(&s.TokenParserKey, "token-parser-key", s.TokenParserKey,
		"An argument to be passed to the token parse plugin on instantiation.")
}
