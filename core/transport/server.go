package transport

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
)

type TCPServer struct {
	cfg TLSConfig
}

func NewTCPServer(cfg TLSConfig) *TCPServer {
	return &TCPServer{cfg: cfg}
}

func (s *TCPServer) Start() error {
	tlsConfig, err := s.buildTLSConfig()
	if err != nil {
		return err
	}

	listener, err := tls.Listen("tcp", s.cfg.Address, tlsConfig)
	if err != nil {
		return err
	}

	log.Println("[TRANSPORT] TCP seguro iniciado em", s.cfg.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go HandleConnection(conn)
	}
}

func (s *TCPServer) buildTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(s.cfg.CertFile, s.cfg.KeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	if s.cfg.RequireMTLS {
		caCert, err := ioutil.ReadFile(s.cfg.CAFile)
		if err != nil {
			return nil, err
		}

		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)

		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = caPool
	}

	return tlsConfig, nil
}
