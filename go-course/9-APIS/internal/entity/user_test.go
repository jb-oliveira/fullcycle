package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "john@example.com", "password123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "password123", user.Password) // Password should be hashed
}

func TestNewUser_PasswordIsHashed(t *testing.T) {
	password := "mySecretPassword"
	user, err := NewUser("Jane Doe", "jane@example.com", password)

	assert.Nil(t, err)
	assert.NotEqual(t, password, user.Password, "Password should be hashed, not stored in plain text")
	assert.Greater(t, len(user.Password), len(password), "Hashed password should be longer than original")
}

func TestUser_ValidatePassword(t *testing.T) {
	password := "correctPassword"
	user, err := NewUser("Test User", "test@example.com", password)

	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword(password), "Should validate correct password")
}

func TestUser_ValidatePassword_WrongPassword(t *testing.T) {
	user, err := NewUser("Test User", "test@example.com", "correctPassword")

	assert.Nil(t, err)
	assert.False(t, user.ValidatePassword("wrongPassword"), "Should reject incorrect password")
}

func TestUser_ValidatePassword_EmptyPassword(t *testing.T) {
	user, err := NewUser("Test User", "test@example.com", "correctPassword")

	assert.Nil(t, err)
	assert.False(t, user.ValidatePassword(""), "Should reject empty password")
}

func TestNewUser_EmptyPassword(t *testing.T) {
	user, err := NewUser("Test User", "test@example.com", "")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Password, "Even empty password should be hashed")
}

func TestNewUser_SpecialCharactersInPassword(t *testing.T) {
	password := "p@ssw0rd!#$%^&*()"
	user, err := NewUser("Test User", "test@example.com", password)

	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword(password), "Should handle special characters in password")
}

func TestNewUser_UniqueIDs(t *testing.T) {
	user1, err1 := NewUser("User One", "user1@example.com", "password1")
	user2, err2 := NewUser("User Two", "user2@example.com", "password2")

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.NotEqual(t, user1.ID, user2.ID, "Each user should have a unique ID")
}
