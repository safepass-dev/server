package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/database"
	"github.com/safepass/server/internal/logging"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/internal/services"
	"github.com/safepass/server/pkg/dotenv"
	"github.com/safepass/server/pkg/dtos/user"
	"github.com/supabase-community/supabase-go"
)

func login(client *supabase.Client, appConfig config.Config) {
	userRepository := repositories.NewUserRepository(client)

	userServices := services.NewUserServices(userRepository)
	authServices := services.NewAuthServices(userServices, &appConfig)

	userRequest := &user.LoginRequest{
		Email:              "test58588@gmail.com",
		MasterPasswordHash: base64.StdEncoding.EncodeToString([]byte("testing")),
	}
	token, merror := authServices.Login(userRequest)
	if merror != nil {
		fmt.Println(merror)
		return
	}

	fmt.Println(token)
}

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

func vaults(client *supabase.Client, logger *logging.Logger) {
	vaultRepository := repositories.NewVaultRepository(client, logger)

	vaultRepository.GetVault("1")
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

	generatePrivateKey()
	login(client, appConfig)
	vaults(client, logger)
}
