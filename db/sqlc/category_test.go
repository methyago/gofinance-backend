package db

import (
	"context"
	"testing"

	"github.com/methyago/gofinance-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {
	user1 := createRandomUser(t)
	arg := CreateCategoryParams{
		UserID:      user1.ID,
		Title:       util.RandomString(12),
		Type:        "debit",
		Description: util.RandomString(20),
	}

	category, err := testQueries.CreateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.Equal(t, arg.UserID, category.UserID)
	require.Equal(t, arg.Title, category.Title)
	require.Equal(t, arg.Type, category.Type)
	require.Equal(t, arg.Description, category.Description)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestGetCategoryById(t *testing.T) {
	cat1 := createRandomCategory(t)
	cat2, err := testQueries.GetCategory(context.Background(), cat1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, cat2)

	require.Equal(t, cat1.ID, cat2.ID)
	require.Equal(t, cat1.UserID, cat2.UserID)
	require.Equal(t, cat1.Title, cat2.Title)
	require.Equal(t, cat1.Type, cat2.Type)
	require.Equal(t, cat1.Description, cat2.Description)
	require.NotEmpty(t, cat2.CreatedAt)
}

func TestDeleteCategory(t *testing.T) {
	cat1 := createRandomCategory(t)
	err := testQueries.DeleteCategory(context.Background(), cat1.ID)

	require.NoError(t, err)
}

func TestUpdateCategory(t *testing.T) {
	cat1 := createRandomCategory(t)

	arg := UpdateCategoriesParams{
		ID:          cat1.ID,
		Title:       util.RandomString(12),
		Description: util.RandomString(20),
	}

	cat2, err := testQueries.UpdateCategories(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cat2)

	require.Equal(t, arg.ID, cat2.ID)
	require.Equal(t, arg.Title, cat2.Title)
	require.Equal(t, arg.Description, cat2.Description)
	require.Equal(t, cat1.Type, cat2.Type)
	require.NotEmpty(t, cat2.CreatedAt)
}

func TestListCategories(t *testing.T) {
	lastCategory := createRandomCategory(t)

	arg := GetCategoriesParams{
		UserID:      lastCategory.UserID,
		Type:        lastCategory.Type,
		Title:       lastCategory.Title,
		Description: lastCategory.Description,
	}

	cats, err := testQueries.GetCategories(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cats)

	for _, cat := range cats {

		require.Equal(t, lastCategory.ID, cat.ID)
		require.Equal(t, lastCategory.UserID, cat.UserID)
		require.Equal(t, arg.Title, cat.Title)
		require.Equal(t, lastCategory.Description, cat.Description)
		require.Equal(t, lastCategory.Title, cat.Title)
		require.Equal(t, lastCategory.Type, cat.Type)
		require.NotEmpty(t, cat.CreatedAt)
	}

}
