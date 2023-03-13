package testUtils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host                  = "localhost"
	port                  = 5432
	user                  = "postgres"
	password              = "SnkPekgXpSob73ANUKSTLNMgAJIWr3Pp"
	dbname                = "orchestrator"
	GETDATA               = "GetData"
	UPDATE_OR_DELETE_DATA = "UpdateOrDeleteData"
	INSERT_DATA           = "InsertData"
)

func ConnectToDB() *sql.DB {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConn)
	CheckError(err)
	return db
}

func UpdateDeleteData(db *sql.DB, deleteOrUpdateStmt string) {
	_, e := db.Exec(deleteOrUpdateStmt)
	CheckError(e)
}

func GetData(getQuery string) *sql.Rows {
	db := ConnectToDB()
	rows, err := db.Query(getQuery)
	CheckError(err)
	defer db.Close()
	return rows
}

func InsertData(db *sql.DB, insertQuery string) {
	_, e := db.Exec(insertQuery)
	CheckError(e)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
