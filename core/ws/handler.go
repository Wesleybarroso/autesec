package ws

import (
	"log"
	"net/http"
	"time"

	"secure-core/core/auth"
	"secure-core/core/session"

	"github.com/gorilla/websocket"
)

var (
	jwtManager   *auth.JWTManager
	sessionMgr   *session.Manager
)

func Init(jwt *auth.JWTManager, sm *session.Manager) {
	jwtManager = jwt
	sessionMgr = sm
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	sessionID := r.URL.Query().Get("session_id")

	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userID, certFP, err := jwtManager.Validate(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	sess, err := sessionMgr.Resume(r.Context(), sessionID, certFP)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx := &WSContext{
		UserID:  userID,
		Session: sess,
	}

	log.Println("[WS] conectado | user:", ctx.UserID)

	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[WS] desconectado | user:", ctx.UserID)
			return
		}

		if err := HandleBridge(ctx, msg); err != nil {
			return
		}

		conn.WriteMessage(websocket.TextMessage, []byte("ACK"))
	}
}
