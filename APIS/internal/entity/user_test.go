package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewUser tests the creation of a new user with valid data.
// It verifies that the user is created successfully, the password is hashed,
// and all fields are populated correctly.
func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "john@example.com", "password123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "password123", user.Password)
}

// TestNewUser_PasswordIsHashed verifies that passwords are properly hashed
// and not stored in plain text. It checks that the hashed password is
// different from the original and longer in length.
func TestNewUser_PasswordIsHashed(t *testing.T) {
	password := "mySecretPassword"
	user, err := NewUser("Jane Doe", "jane@example.com", password)

	assert.Nil(t, err)
	assert.NotEqual(t, password, user.Password, "Password should be hashed, not stored in plain text")
	assert.Greater(t, len(user.Password), len(password), "Hashed password should be longer than original")
}

// TestUser_ValidatePassword tests that the ValidatePassword method
// correctly validates a correct password against the stored hash.
func TestUser_ValidatePassword(t *testing.T) {
	password := "correctPassword"
	user, err := NewUser("Test User", "test@example.com", password)

	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword(password), "Should validate correct password")
}

// TestUser_ValidatePassword_WrongPassword verifies that ValidatePassword
// correctly rejects an incorrect password.
func TestUser_ValidatePassword_WrongPassword(t *testing.T) {
	user, err := NewUser("Test User", "test@example.com", "correctPassword")

	assert.Nil(t, err)
	assert.False(t, user.ValidatePassword("wrongPassword"), "Should reject incorrect password")
}

// TestUser_ValidatePassword_EmptyPassword tests that ValidatePassword
// rejects empty passwords, ensuring security validation.
func TestUser_ValidatePassword_EmptyPassword(t *testing.T) {
	user, err := NewUser("Test User", "test@example.com", "correctPassword")

	assert.Nil(t, err)
	assert.False(t, user.ValidatePassword(""), "Should reject empty password")
}

// TestNewUser_EmptyPassword verifies that even empty passwords are hashed
// rather than stored as empty strings, maintaining consistent security behavior.
func TestNewUser_EmptyPassword(t *testing.T) {
	user, err := NewUser("Test User", "test@example.com", "")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Password, "Even empty password should be hashed")
}

// TestNewUser_SpecialCharactersInPassword tests that passwords containing
// special characters are properly hashed and can be validated correctly.
func TestNewUser_SpecialCharactersInPassword(t *testing.T) {
	password := "p@ssw0rd!#$%^&*()"
	user, err := NewUser("Test User", "test@example.com", password)

	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword(password), "Should handle special characters in password")
}

// TestNewUser_UniqueIDs verifies that each user created receives a unique ID,
// ensuring no ID collisions occur between different users.
func TestNewUser_UniqueIDs(t *testing.T) {
	user1, err1 := NewUser("User One", "user1@example.com", "password1")
	user2, err2 := NewUser("User Two", "user2@example.com", "password2")

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.NotEqual(t, user1.ID, user2.ID, "Each user should have a unique ID")
}
