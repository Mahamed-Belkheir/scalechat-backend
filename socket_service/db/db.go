package db

import (
	"log"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/gocql/gocql"
)

func GetConnection(conf service.Config) *gocql.Session {
	cluster := gocql.NewCluster(conf.Clusters...)
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "scalechatmessages"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("fatal: could not connect to db: %v", err)
	}
	return session
}
