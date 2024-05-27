package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tomersf/bank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
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
