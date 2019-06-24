package orm_test

import (
	"github.com/google/uuid"
	"testing"
	"time"
	"weassistant/models"
	"weassistant/services/orm"
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
	testService(t, orm.MustNewAdministratorService, createTestAdministrator)
}
