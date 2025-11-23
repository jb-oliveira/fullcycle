package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("João Silva", "joao@exemplo.com", "senha123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "João Silva", user.Name)
	assert.Equal(t, "joao@exemplo.com", user.Email)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "senha123", user.Password)
}

func TestNewUser_PasswordIsHashed(t *testing.T) {
	password := "minhaSenhaSecreta"
	user, err := NewUser("Maria Santos", "maria@exemplo.com", password)

	assert.Nil(t, err)
	assert.NotEqual(t, password, user.Password, "Password should be hashed, not stored in plain text")
	assert.Greater(t, len(user.Password), len(password), "Hashed password should be longer than original")
}

func TestUser_ValidatePassword(t *testing.T) {
	password := "senhaCorreta"
	user, err := NewUser("Usuário Teste", "teste@exemplo.com", password)

	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword(password), "Should validate correct password")
}

func TestUser_ValidatePassword_WrongPassword(t *testing.T) {
	user, err := NewUser("Usuário Teste", "teste@exemplo.com", "senhaCorreta")

	assert.Nil(t, err)
	assert.False(t, user.ValidatePassword("senhaErrada"), "Should reject incorrect password")
}

func TestUser_ValidatePassword_EmptyPassword(t *testing.T) {
	user, err := NewUser("Usuário Teste", "teste@exemplo.com", "senhaCorreta")

	assert.Nil(t, err)
	assert.False(t, user.ValidatePassword(""), "Should reject empty password")
}

func TestNewUser_EmptyPassword(t *testing.T) {
	user, err := NewUser("Usuário Teste", "teste@exemplo.com", "")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Password, "Even empty password should be hashed")
}

func TestNewUser_SpecialCharactersInPassword(t *testing.T) {
	password := "s3nh@!#$%^&*()"
	user, err := NewUser("Usuário Teste", "teste@exemplo.com", password)

	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword(password), "Should handle special characters in password")
}

func TestNewUser_UniqueIDs(t *testing.T) {
	user1, err1 := NewUser("Usuário Um", "usuario1@exemplo.com", "senha1")
	user2, err2 := NewUser("Usuário Dois", "usuario2@exemplo.com", "senha2")

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.NotEqual(t, user1.ID, user2.ID, "Each user should have a unique ID")
}
