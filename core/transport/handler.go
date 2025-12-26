package transport

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"time"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return
	}

	err := tlsConn.Handshake()
	if err != nil {
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return
	}

	clientCert := state.PeerCertificates[0]
	fingerprint := CertFingerprint(clientCert)

	log.Println("[TRANSPORT] Cliente conectado | FP:", fingerprint)

	reader := bufio.NewReader(tlsConn)

	for {
		tlsConn.SetReadDeadline(time.Now().Add(60 * time.Second))

		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("[TRANSPORT] Cliente desconectado | FP:", fingerprint)
			return
		}

		// Base do protocolo (ainda neutra)
		log.Println("[TRANSPORT] RX:", msg)

		_, err = tlsConn.Write([]byte("ACK\n"))
		if err != nil {
			return
		}
	}
}
