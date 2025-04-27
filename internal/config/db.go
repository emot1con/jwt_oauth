package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func Connect() *sql.DB {
	ticker := 1

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_Name")

	logrus.Info("connecting to postgres")

	for {
		DSN := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbUser, dbPassword, dbName)
		db, err := sql.Open("pgx", DSN)

		if ticker >= 5 {
			logrus.Fatalf("failed connect to postgres: %v", err)
			return nil
		}

		if err != nil {
			logrus.Infof("retrying connect to postgres (%d/5) attempt", ticker)
			time.Sleep(2 * time.Second)
			ticker++
			continue
		}

		if err := db.Ping(); err != nil {
			logrus.Fatalf("failed ping db: %v", err)
		}

		db.SetMaxOpenConns(20)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(1 * time.Hour)
		db.SetConnMaxIdleTime(15 * time.Minute)

		return db
	}
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "default",
		DB:       1,
	})
}
