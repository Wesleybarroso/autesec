package transport

import (
	"bufio"
	"context"
	"crypto/tls"
	"log"
	"net"
	"time"

	"autesec/core/session"
)

type Handler struct {
	SessionManager *session.Manager
}

func NewHandler(sm *session.Manager) *Handler {
	return &Handler{
		SessionManager: sm,
	}
}

func (h *Handler) HandleConnection(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return
	}

	// Handshake TLS (Core Shield)
	if err := tlsConn.Handshake(); err != nil {
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return
	}

	// Core Identity (fingerprint)
	clientCert := state.PeerCertificates[0]
	fingerprint := CertFingerprint(clientCert)

	ctx := context.Background()

	// Core Session + Core Bind
	sess, err := h.SessionManager.Create(ctx, fingerprint, fingerprint)
	if err != nil {
		log.Println("[SESSION] erro ao criar sessão:", err)
		return
	}

	log.Println("[TRANSPORT] Cliente conectado | FP:", fingerprint, "| Session:", sess.SessionID)

	reader := bufio.NewReader(tlsConn)

	for {
		tlsConn.SetReadDeadline(time.Now().Add(60 * time.Second))

		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("[TRANSPORT] Cliente desconectado | FP:", fingerprint)
			_ = h.SessionManager.Close(ctx, sess.SessionID)
			return
		}

		// 🔹 Aqui entra o protocolo da aplicação
		log.Println("[TRANSPORT] RX:", msg)

		// Atualiza last_seen da sessão
		_, _ = h.SessionManager.Resume(ctx, sess.SessionID, fingerprint)

		// Resposta simples
		if _, err := tlsConn.Write([]byte("ACK\n")); err != nil {
			return
		}
	}
}
