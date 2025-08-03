package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	// ed25519.GenerateKey uses crypto/rand to generate a cryptographically secure key pair.
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// The ed25519.PrivateKey is 64 bytes (32-byte seed followed by 32-byte public key).
	// The ed25519.PublicKey is just the 32-byte public key portion.

	// Encode keys to hex strings for easy use in environment variables.
	privateKeyHex := hex.EncodeToString(privateKey)
	publicKeyHex := hex.EncodeToString(publicKey)

	fmt.Println("🔑 PASETO V4 Keys Generated Successfully!")
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("PRIVATE KEY (64 bytes / 128 hex chars):\n%s\n\n", privateKeyHex)
	fmt.Printf("PUBLIC KEY (32 bytes / 64 hex chars):\n%s\n\n", publicKeyHex)
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("✅ Action Required:")
	fmt.Println("1. Set this PRIVATE KEY for the 'IDENTITY_PASETO_PRIVATE_KEY' environment variable in your Identity Service.")
	fmt.Println("2. Set this PUBLIC KEY for the 'INTEGRITY_PASETO_PUBLIC_KEY' environment variable in all services that need to verify tokens (including the Identity Service itself, other services, and your API Gateway).")
}
