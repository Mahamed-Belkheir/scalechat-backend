package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/gocql/gocql"
	_ "github.com/gocql/gocql"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	conf := socket_service.GetConfig()

	db := gocql.NewCluster(conf.Clusters...)
	db.Keyspace = "scalechatmessages"
	session, err := db.CreateSession()
	if err != nil {
		log.Fatalf("failed to make session %v", err)
	}
	driver, err := cassandra.WithInstance(session, &cassandra.Config{
		KeyspaceName:    "scalechatmessages",
		MigrationsTable: "go_migrate_table",
	})
	if err != nil {
		log.Fatalf("failed to make instance %v", err)
	}

	p, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("failed to make get path %v", err)
	}
	p = filepath.ToSlash(p)
	p = path.Join(p, "socket_service", "migrations")
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", p),
		"cassandra", driver)
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
