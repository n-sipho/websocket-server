package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	DB = db

	t.Run("successful user retrieval", func(t *testing.T) {
		expectedUID := "test-uid"
		expectedSpotifyID := "spotify-123"

		rows := sqlmock.NewRows([]string{"user_id"}).AddRow(expectedSpotifyID)
		mock.ExpectQuery("SELECT user_id FROM users WHERE uid = ?").
			WithArgs(expectedUID).
			WillReturnRows(rows)

		spotifyID, err := GetUser(expectedUID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSpotifyID, spotifyID)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT user_id FROM users WHERE uid = ?").
			WithArgs("non-existent").
			WillReturnError(sql.ErrNoRows)

		spotifyID, err := GetUser("non-existent")
		assert.Error(t, err)
		assert.Equal(t, "", spotifyID)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT user_id FROM users WHERE uid = ?").
			WithArgs("error-uid").
			WillReturnError(sql.ErrConnDone)

		spotifyID, err := GetUser("error-uid")
		assert.Error(t, err)
		assert.Equal(t, "", spotifyID)
	})
}
