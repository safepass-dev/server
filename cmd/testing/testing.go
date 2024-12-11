package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/database"
	"github.com/safepass/server/internal/logging"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/internal/services"
	"github.com/safepass/server/pkg/dotenv"
	"github.com/supabase-community/supabase-go"
)

func generatePrivateKey() {
	curve := elliptic.P256()

	// ECDSA özel anahtarını oluşturun
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

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

func vaults(client *supabase.Client, appConfig config.Config, logger *logging.Logger) {
	vaultRepository := repositories.NewVaultRepository(client, logger)
	vaultServices := services.NewVaultServices(vaultRepository, &appConfig)

	vault, err := vaultServices.GetVaultByUserID("10")
	if err != nil {
		fmt.Println(err.Description)
		return
	}

	fmt.Println(vault.User)
}

func main() {
	dotenv.LoadEnv(".env")

	var appConfig config.Config
	config.LoadConfig(&appConfig)

	context, err := database.NewAppContextDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := context.GetSupabaseClient()

	logger, err := logging.NewLogger(logging.INFO, "log.txt")
	if err != nil {
		return
	}

	vaults(client, appConfig, logger)
}
