package ws

import (
	"errors"
)

func HandleBridge(ctx *WSContext, payload []byte) error {
	// Aqui é o ponto neutro onde:
	// - eventos do site entram
	// - comandos são roteados
	// - integrações futuras entram (WhatsApp, rastreamento, etc)

	if ctx.Session == nil {
		return errors.New("invalid session")
	}

	// Por enquanto é neutro
	return nil
}
