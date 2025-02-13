package database

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CommitOrRollback(tx *sql.Tx, c *fiber.Ctx, err error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("error rolling back transaction: %v\n (original error: %v)", rollbackErr, err)
		}
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			log.Printf("error committing transaction: %v\n (original error: %v)", commitErr, err)
		}
	}
}
