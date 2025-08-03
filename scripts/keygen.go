package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	GeneratePasetoSecurityKeys()
}

// Generates PASETO Environment Keys
func GeneratePasetoSecurityKeys() {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	privateKeyHex := hex.EncodeToString(privateKey)
	publicKeyHex := hex.EncodeToString(publicKey)

	fmt.Println("🔑 PASETO V4 Keys Generated Successfully!")
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("PRIVATE KEY (64 bytes / 128 hex chars):\n%s\n\n", privateKeyHex)
	fmt.Printf("PUBLIC KEY (32 bytes / 64 hex chars):\n%s\n\n", publicKeyHex)
	fmt.Println("----------------------------------------------------------------")
}
