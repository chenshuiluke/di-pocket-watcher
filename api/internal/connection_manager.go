package connection_manager

import (
	"context"
	"fmt"
	"os"

	"github.com/chenshuiluke/di-pocket-watcher/api/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	Conn    *pgxpool.Pool
	Queries *db.Queries
}

var Mgr Manager

func init() {
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, url)
	Mgr.Conn = conn
	Mgr.Queries = db.New(Mgr.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}
