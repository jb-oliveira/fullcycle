package database

import (
	"testing"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	// Use auto migration
	err = db.AutoMigrate(&entity.User{})
	require.NoError(t, err)

	return db
}

func TestNewUser(t *testing.T) {
	db := setupTestDB(t)
	userDB := NewUser(db)

	assert.NotNil(t, userDB)
	assert.NotNil(t, userDB.db)
}

func TestUser_Create(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		db := setupTestDB(t)
		userDB := NewUser(db)

		user, err := entity.NewUser("John Doe", "john@example.com", "password123")
		require.NoError(t, err)

		err = userDB.Create(user)
		assert.NoError(t, err)

		// Verify user was created
		var count int64
		db.Model(&entity.User{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("should return error for duplicate email", func(t *testing.T) {
		db := setupTestDB(t)
		userDB := NewUser(db)

		user1, err := entity.NewUser("John Doe", "john@example.com", "password123")
		require.NoError(t, err)

		err = userDB.Create(user1)
		require.NoError(t, err)

		user2, err := entity.NewUser("Jane Doe", "john@example.com", "password456")
		require.NoError(t, err)

		err = userDB.Create(user2)
		assert.Error(t, err)
	})
}

func TestUser_FindByEmail(t *testing.T) {
	t.Run("should find user by email successfully", func(t *testing.T) {
		db := setupTestDB(t)
		userDB := NewUser(db)

		user, err := entity.NewUser("John Doe", "john@example.com", "password123")
		require.NoError(t, err)

		err = userDB.Create(user)
		require.NoError(t, err)

		foundUser, err := userDB.FindByEmail("john@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		db := setupTestDB(t)
		userDB := NewUser(db)

		foundUser, err := userDB.FindByEmail("nonexistent@example.com")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NotNil(t, foundUser)
	})

	t.Run("should find correct user when multiple users exist", func(t *testing.T) {
		db := setupTestDB(t)
		userDB := NewUser(db)

		user1, err := entity.NewUser("John Doe", "john@example.com", "password123")
		require.NoError(t, err)
		err = userDB.Create(user1)
		require.NoError(t, err)

		user2, err := entity.NewUser("Jane Doe", "jane@example.com", "password456")
		require.NoError(t, err)
		err = userDB.Create(user2)
		require.NoError(t, err)

		foundUser, err := userDB.FindByEmail("jane@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user2.ID, foundUser.ID)
		assert.Equal(t, user2.Name, foundUser.Name)
		assert.Equal(t, user2.Email, foundUser.Email)
	})
}
