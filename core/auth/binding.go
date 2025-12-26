package auth

import "errors"

// Garante que o JWT pertence ao certificado usado na conexão
func EnforceBinding(jwtCertFP, connCertFP string) error {
	if jwtCertFP != connCertFP {
		return errors.New("jwt not bound to this certificate")
	}
	return nil
}
