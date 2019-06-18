package services

import (
	"github.com/google/uuid"
	"testing"
	"time"
	"weassistant/models"
)

func createTestAdministrator() *models.Administrator {
	now := time.Now()
	return &models.Administrator{
		Name:        "test name",
		Account:     uuid.New().String(),
		Password:    "test password",
		LastLoginAt: &now,
	}
}

func TestAdministrator(t *testing.T) {
	testService(t, MustNewAdministratorService, createTestAdministrator)
}
