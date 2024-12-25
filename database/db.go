package database

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
    "log"
    "os"
)

var DbPool *pgxpool.Pool // Exported variable

// InitializeDB initializes the connection pool
func InitializeDB() {
    config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("Unable to parse DATABASE_URL:", err)
    }

    pool, err := pgxpool.ConnectConfig(context.Background(), config)
    if err != nil {
        log.Fatal("Unable to connect to database:", err)
    }

    DbPool = pool
    log.Println("Database connection pool initialized.")
}

// GetDB retrieves the connection pool
func GetDB() *pgxpool.Pool {
    return DbPool
}

// CloseDB closes the database connection pool
func CloseDB() {
    DbPool.Close()
}
