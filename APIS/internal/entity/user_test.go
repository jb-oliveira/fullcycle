package entity

import (
	"testing"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
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

// TestNewUser_EmptyPassword verifies that empty passwords are rejected
// by validation, ensuring security requirements.
func TestNewUser_EmptyPassword(t *testing.T) {
	user, err := NewUser("Test User", "test@example.com", "")

	assert.ErrorIs(t, err, ErrPasswordRequired)
	assert.Nil(t, user)
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

// TestNewUser_ValidatesFields tests the validation logic in NewUser.
// It uses table-driven tests to verify that validation errors are returned
// for invalid inputs (empty name, empty email, empty password, fields too long).
func TestNewUser_ValidatesFields(t *testing.T) {
	name256 := string(make([]byte, 256))
	for i := range name256 {
		name256 = name256[:i] + "A" + name256[i+1:]
	}

	email256 := string(make([]byte, 256))
	for i := range email256 {
		email256 = email256[:i] + "B" + email256[i+1:]
	}

	password256 := string(make([]byte, 256))
	for i := range password256 {
		password256 = password256[:i] + "C" + password256[i+1:]
	}

	tests := []struct {
		name        string
		userName    string
		email       string
		password    string
		expectError error
	}{
		{
			name:        "valid user",
			userName:    "John Doe",
			email:       "john@example.com",
			password:    "password123",
			expectError: nil,
		},
		{
			name:        "empty name",
			userName:    "",
			email:       "john@example.com",
			password:    "password123",
			expectError: ErrNameRequired,
		},
		{
			name:        "name too long",
			userName:    name256,
			email:       "john@example.com",
			password:    "password123",
			expectError: ErrNameTooLong,
		},
		{
			name:        "empty email",
			userName:    "John Doe",
			email:       "",
			password:    "password123",
			expectError: ErrEmailRequired,
		},
		{
			name:        "email too long",
			userName:    "John Doe",
			email:       email256,
			password:    "password123",
			expectError: ErrEmailTooLong,
		},
		{
			name:        "empty password",
			userName:    "John Doe",
			email:       "john@example.com",
			password:    "",
			expectError: ErrPasswordRequired,
		},
		{
			name:        "password too long",
			userName:    "John Doe",
			email:       "john@example.com",
			password:    password256,
			expectError: ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.userName, tt.email, tt.password)

			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
				assert.Nil(t, user)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userName, user.Name)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

// TestUser_Validate tests the Validate method directly on User instances.
// It verifies that validation rules are correctly enforced for various
// user configurations.
func TestUser_Validate(t *testing.T) {
	name256 := string(make([]byte, 256))
	for i := range name256 {
		name256 = name256[:i] + "X" + name256[i+1:]
	}

	email256 := string(make([]byte, 256))
	for i := range email256 {
		email256 = email256[:i] + "Y" + email256[i+1:]
	}

	password256 := string(make([]byte, 256))
	for i := range password256 {
		password256 = password256[:i] + "Z" + password256[i+1:]
	}

	tests := []struct {
		name        string
		user        *User
		expectError error
	}{
		{
			name: "valid user",
			user: &User{
				ID:       entity.NewID(),
				Name:     "Valid User",
				Email:    "valid@example.com",
				Password: "hashedpassword",
			},
			expectError: nil,
		},
		{
			name: "empty name",
			user: &User{
				ID:       entity.NewID(),
				Name:     "",
				Email:    "user@example.com",
				Password: "hashedpassword",
			},
			expectError: ErrNameRequired,
		},
		{
			name: "name too long",
			user: &User{
				ID:       entity.NewID(),
				Name:     name256,
				Email:    "user@example.com",
				Password: "hashedpassword",
			},
			expectError: ErrNameTooLong,
		},
		{
			name: "empty email",
			user: &User{
				ID:       entity.NewID(),
				Name:     "User Name",
				Email:    "",
				Password: "hashedpassword",
			},
			expectError: ErrEmailRequired,
		},
		{
			name: "email too long",
			user: &User{
				ID:       entity.NewID(),
				Name:     "User Name",
				Email:    email256,
				Password: "hashedpassword",
			},
			expectError: ErrEmailTooLong,
		},
		{
			name: "empty password",
			user: &User{
				ID:       entity.NewID(),
				Name:     "User Name",
				Email:    "user@example.com",
				Password: "",
			},
			expectError: ErrPasswordRequired,
		},
		{
			name: "password too long",
			user: &User{
				ID:       entity.NewID(),
				Name:     "User Name",
				Email:    "user@example.com",
				Password: password256,
			},
			expectError: ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// TestNewUser_WithNameAt255Characters tests that user names
// with exactly 255 characters are accepted as valid.
func TestNewUser_WithNameAt255Characters(t *testing.T) {
	name255 := string(make([]byte, 255))
	for i := range name255 {
		name255 = name255[:i] + "A" + name255[i+1:]
	}

	user, err := NewUser(name255, "user@example.com", "password123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 255, len(user.Name))
	assert.Equal(t, name255, user.Name)
}

// TestNewUser_WithEmailAt255Characters tests that user emails
// with exactly 255 characters are accepted as valid.
func TestNewUser_WithEmailAt255Characters(t *testing.T) {
	email255 := string(make([]byte, 255))
	for i := range email255 {
		email255 = email255[:i] + "B" + email255[i+1:]
	}

	user, err := NewUser("John Doe", email255, "password123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 255, len(user.Email))
	assert.Equal(t, email255, user.Email)
}
