package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connection() (*sql.DB, error) {
	url := "root:root@/"
	db, erro := sql.Open("mysql", url)

	db.Exec("create database if not exists books")
	db.Exec("use books")
	db.Exec(`create table if not exists users (
		id integer auto_increment,
		name VARCHAR(80) NOT NULL,
		email VARCHAR(80) NOT NULL,
		PRIMARY KEY (id)
	)`)

	if erro != nil {
		return nil, erro
	}
	if erro = db.Ping(); erro != nil {
		return nil, erro
	}
	return db, nil
}
