package db
import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
	"fmt"
)

var DB *pgxpool.Pool

func ConnectDB() {
	// dsn := os.Getenv("DATABASE_URL")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbname,
	)
	var err error
	for i := 0; i < 10; i++ {
		DB, err = pgxpool.New(context.Background(), dsn)
		if err == nil {
			err = DB.Ping(context.Background())
			if err == nil {
				break
			}
		}
		log.Printf("DB not ready, retrying...")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Println("Connected to database successfully")
}