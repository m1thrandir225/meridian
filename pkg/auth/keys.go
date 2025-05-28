package auth

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
)

var pasetoKublicKey ed25519.PublicKey

func LoadPublicKeyFromEnv() error {
	hexKey := os.Getenv("INTEGRITY_PASETO_PUBLIC_KEY")
	if hexKey == "" {
	}

	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
	}
	if len(keyBytes) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid ed25519 public key size: expected: %d, got %d", ed25519.PublicKeySize, len(keyBytes))
	}

	pasetoKublicKey = ed25519.PublicKey(keyBytes)
	log.Println("PASETO v4 Public key loaded sucessfully for token verification")
	return nil
}

func GetPublicKey() (ed25519.PublicKey, error) {
	if pasetoKublicKey == nil {
		return nil, errors.New("paseto public key not loaded")
	}
	return pasetoKublicKey, nil
}
