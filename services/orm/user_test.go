package orm_test

import (
	"github.com/google/uuid"
	"testing"
	"time"
	"weassistant/models"
	"weassistant/services/orm"
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
	testService(t, orm.MustNewUserService, createTestUser)
}
