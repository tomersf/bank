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
var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(testDB)
	result := m.Run()
	cleanup()
	os.Exit(result)
}

func cleanup() {
	clearTransfersTable()
	clearEntriesTable()
	clearAccountsTable()
}

func clearTransfersTable() {
	err := testQueries.DeleteTransfers(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func clearEntriesTable() {
	err := testQueries.DeleteEntries(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func clearAccountsTable() {
	err := testQueries.DeleteAccounts(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
