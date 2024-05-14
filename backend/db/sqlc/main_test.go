package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(conn)
	result := m.Run()
	cleanup()
	os.Exit(result)
}

func cleanup() {
	clearAccountsTable()
}

func clearAccountsTable() {
	err := testQueries.DeleteAccounts(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
