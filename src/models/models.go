package models

import (
    "database/sql"
	"fmt"
)

var Db *sql.DB

func Ping() (error) {
    fmt.Println("... verifying database connection")
    if err := Db.Ping(); err != nil {
        return err
    }
    return nil
}