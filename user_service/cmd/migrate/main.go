package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Mahamed-Belkheir/scalechat-backend/user_service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := user_service.GetConfig()
	p, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("failed to make get path %v", err)
	}
	p = filepath.ToSlash(p)
	p = path.Join(p, "user_service", "migrations")
	m, err := migrate.New(
		fmt.Sprintf("file://%s", p),
		config.DBConn)
	if err != nil {
		log.Fatalf("failed to make migrate instance %v", err)
	}
	cmd := os.Args[1]
	switch cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Printf("incorrect command %v", cmd)
	}
	if err != nil {
		log.Print(err)
	}
}
