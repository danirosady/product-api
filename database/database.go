package database

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string, dbPort string) (*sql.DB, error) {

	// Handle host:port
	if dbPort != "" {
		re := regexp.MustCompile(`(@[^:]+):\d+`)
		connectionString = re.ReplaceAllString(connectionString, "${1}:"+dbPort)

		re2 := regexp.MustCompile(`(://[^:]+):\d+`)
		connectionString = re2.ReplaceAllString(connectionString, "${1}:"+dbPort)
	} else {
		connectionString = strings.Replace(connectionString, ":6543", ":5432", 1)
	}

	// Simple Protocol
	params := "prefer_simple_protocol=true"

	if !strings.Contains(connectionString, "prefer_simple_protocol") {
		if strings.Contains(connectionString, "?") {
			connectionString += "&" + params
		} else {
			connectionString += "?" + params
		}
	}

	log.Println("Connecting to DB...")

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	log.Println("Database connected successfully")
	return db, nil
}
