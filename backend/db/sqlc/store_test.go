package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const NUM_OF_GOROUTINES = 10

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	amount := int64(10)
	fmt.Println(">> before:", accountOne.Balance, accountTwo.Balance)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < NUM_OF_GOROUTINES; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: accountOne.ID,
				ToAccountID:   accountTwo.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < NUM_OF_GOROUTINES; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountOne.ID, transfer.FromAccountID)
		require.Equal(t, accountTwo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountOne.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountTwo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountOne.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTwo.ID, toAccount.ID)

		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := accountOne.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - accountTwo.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= NUM_OF_GOROUTINES)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccountOne, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updatedAccountTwo, err := testQueries.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccountOne.Balance, updatedAccountTwo.Balance)
	require.Equal(t, accountOne.Balance-int64(NUM_OF_GOROUTINES)*amount, updatedAccountOne.Balance)
	require.Equal(t, accountTwo.Balance+int64(NUM_OF_GOROUTINES)*amount, updatedAccountTwo.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	amount := int64(10)
	fmt.Println(">> before:", accountOne.Balance, accountTwo.Balance)

	errs := make(chan error)

	for i := 0; i < NUM_OF_GOROUTINES; i++ {
		fromAccountID := accountOne.ID
		toAccountID := accountTwo.ID

		if i%2 == 1 {
			fromAccountID = accountTwo.ID
			toAccountID = accountOne.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < NUM_OF_GOROUTINES; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccountOne, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updatedAccountTwo, err := testQueries.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccountOne.Balance, updatedAccountTwo.Balance)
	require.Equal(t, accountOne.Balance, updatedAccountOne.Balance)
	require.Equal(t, accountTwo.Balance, updatedAccountTwo.Balance)
}
