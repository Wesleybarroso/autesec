package transport

type TLSConfig struct {
	Address     string
	CertFile    string
	KeyFile     string
	CAFile      string
	RequireMTLS bool
}
