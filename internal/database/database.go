package database

import (
	"os"

	"github.com/safepass/server/pkg/dotenv"
	"github.com/supabase-community/supabase-go"
)

type Database interface {
	// GetSupabaseClient returns the Supabase client
	GetSupabaseClient() *supabase.Client
}

// AppContextDB is the struct that holds the Supabase client
type AppContextDB struct {
	SupabaseClient *supabase.Client

	Database
}

// NewAppContextDB creates a new AppContextDB
func NewAppContextDB() (appContextDB *AppContextDB, err error) {
	err = dotenv.LoadEnv(".env")
	if err != nil {
		return
	}

	restUrl := os.Getenv("SUPABASE_REST_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")

	client, err := supabase.NewClient(restUrl, apiKey, &supabase.ClientOptions{})
	if err != nil {
		return
	}

	appContextDB = &AppContextDB{
		SupabaseClient: client,
	}

	return
}

// GetSupabaseClient returns the Supabase client
func (a *AppContextDB) GetSupabaseClient() *supabase.Client {
	return a.SupabaseClient
}
