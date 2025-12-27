package bot

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func Connection() (*whatsmeow.Client, error) {
	godotenv.Load()
	dbPath := os.Getenv("DB_PATH")
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", dbPath, dbLog)
	if err != nil {
		return nil, err
	}

	store, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(store, nil)
	return client, nil
}
