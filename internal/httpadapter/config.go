package httpadapter

type Config struct {
	ServeAddress string `yaml:"serve_address"`
	BasePath     string `yaml:"base_path"`
	UseTLS       bool   `yaml:"use_tls"`
	TLSKeyFile   string `yaml:"tls_key_file"`
	TLSCrtFile   string `yaml:"tls_crt_file"`

	AccessTokenCookie  string `yaml:"access_token_cookie"`
	RefreshTokenCookie string `yaml:"refresh_token_cookie"`

	SwaggerAddress string `yaml:"swagger_address"`
}
