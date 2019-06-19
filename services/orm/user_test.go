package orm

import (
	"github.com/google/uuid"
	"testing"
	"time"
	"weassistant/models"
)

func createTestUser() *models.User {
	now := time.Now()
	return &models.User{
		Nickname:    "test nickname",
		Account:     uuid.New().String(),
		Password:    "test password",
		LastLoginAt: &now,
	}
}

func TestUser(t *testing.T) {
	testService(t, MustNewUserService, createTestUser)
}
