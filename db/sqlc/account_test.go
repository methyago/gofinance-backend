package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/methyago/gofinance-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	cat := createRandomCategory(t)
	arg := CreateAccountParams{
		UserID:      cat.UserID,
		CategoryID:  cat.ID,
		Title:       util.RandomString(12),
		Type:        cat.Type,
		Description: util.RandomString(20),
		Value:       10,
		Date:        time.Now(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.UserID, account.UserID)
	require.Equal(t, arg.Title, account.Title)
	require.Equal(t, arg.Type, account.Type)
	require.Equal(t, arg.Description, account.Description)
	require.Equal(t, arg.Value, account.Value)
	require.Equal(t, arg.CategoryID, account.CategoryID)
	require.NotEmpty(t, account.Date)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountById(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.UserID, acc2.UserID)
	require.Equal(t, acc1.Title, acc2.Title)
	require.Equal(t, acc1.Type, acc2.Type)
	require.Equal(t, acc1.Description, acc2.Description)
	require.Equal(t, acc1.Value, acc2.Value)
	require.Equal(t, acc1.CategoryID, acc2.CategoryID)
	require.NotEmpty(t, acc2.Date)
	require.NotEmpty(t, acc2.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	acc := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc.ID)

	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := UpdateAccountsParams{
		ID:          acc1.ID,
		Title:       util.RandomString(12),
		Description: util.RandomString(20),
		Value:       20,
	}

	cat2, err := testQueries.UpdateAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cat2)

	require.Equal(t, arg.ID, cat2.ID)
	require.Equal(t, arg.Title, cat2.Title)
	require.Equal(t, arg.Description, cat2.Description)
	require.Equal(t, arg.Value, cat2.Value)
	require.Equal(t, acc1.Type, cat2.Type)
	require.Equal(t, acc1.Date, cat2.Date)
	require.NotEmpty(t, cat2.CreatedAt)
}

func TestListAccounts(t *testing.T) {
	lastAccount := createRandomAccount(t)

	arg := GetAccountsParams{
		UserID:      lastAccount.UserID,
		Type:        lastAccount.Type,
		Title:       lastAccount.Title,
		Description: lastAccount.Description,
		CategoryID: sql.NullInt32{
			Valid: true,
			Int32: lastAccount.CategoryID,
		},
		Date: sql.NullTime{
			Valid: true,
			Time:  lastAccount.Date,
		},
	}

	accs, err := testQueries.GetAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accs)

	for _, acc := range accs {

		require.Equal(t, lastAccount.ID, acc.ID)
		require.Equal(t, lastAccount.UserID, acc.UserID)
		require.Equal(t, arg.Title, acc.Title)
		require.NotEmpty(t, acc.CategoryTitle)
		require.Equal(t, lastAccount.Description, acc.Description)
		require.Equal(t, lastAccount.Title, acc.Title)
		require.Equal(t, lastAccount.Type, acc.Type)
		require.Equal(t, lastAccount.Value, acc.Value)
		require.NotEmpty(t, acc.Date)
		require.NotEmpty(t, acc.CreatedAt)
	}

}

func TestListGetReports(t *testing.T) {
	lastAccount := createRandomAccount(t)

	arg := GetAccountsReportsParams{
		UserID: lastAccount.UserID,
		Type:   lastAccount.Type,
	}

	total, err := testQueries.GetAccountsReports(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, total)

}

func TestListGetAccountGraph(t *testing.T) {
	lastAccount := createRandomAccount(t)

	arg := GetAccountGraphParams{
		UserID: lastAccount.UserID,
		Type:   lastAccount.Type,
	}

	total, err := testQueries.GetAccountGraph(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, total)

}
