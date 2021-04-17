package db

import (
	"database/sql"
	"log"

	service "github.com/Mahamed-Belkheir/scalechat-backend/chat_service"
	_ "github.com/lib/pq"
)

func GetConnection(config service.Config) *sql.DB {
	conn, err := sql.Open(config.DB, config.DBConn)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

type statementPreparer struct {
	err  error
	conn *sql.DB
}

func prep(conn *sql.DB) statementPreparer {
	return statementPreparer{
		conn: conn,
	}
}

func (q statementPreparer) prepare(statement string) *sql.Stmt {
	if q.err != nil {
		return nil
	}
	stmt, err := q.conn.Prepare(statement)
	if err != nil {
		q.err = err
		return nil
	}
	return stmt
}
