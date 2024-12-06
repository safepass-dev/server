package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func main() {
	curve := elliptic.P256()

	// ECDSA özel anahtarını oluşturun
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	code, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		panic(err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: code,
	})

	// PEM formatındaki özel anahtarı yazdır
	fmt.Println("PEM Formatındaki Özel Anahtar:")
	fmt.Println(string(privPEM))
}
