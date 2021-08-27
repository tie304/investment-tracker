package database

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var Database *pg.DB

func InitDB() {
	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
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
