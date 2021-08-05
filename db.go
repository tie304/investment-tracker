package main

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var Database *pg.DB

func initDB() {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("database_user"),
		Password: os.Getenv("database_password"),
		Database: os.Getenv("database_name"),
	})
	Database = db
	ctx := context.Background()
	err := Database.Ping(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("initalized db")
	createSchema()
	log.Println("Schema created")

}

func createSchema() {
	model := Database.Model(new(Asset))
	err := model.CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		panic(err)
	}

}
