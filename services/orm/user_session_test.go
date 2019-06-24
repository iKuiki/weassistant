package orm_test

import (
	"testing"
	"weassistant/models"
	"weassistant/services/orm"

	"github.com/google/uuid"
)

func createTestUserSession() *models.UserSession {
	return &models.UserSession{
		UserID:    1,
		Effective: true,
		Token:     uuid.New().String(),
	}
}

func TestUserSession(t *testing.T) {
	userSessionService := orm.MustNewUserSessionService(config.GetMainDB(), config.GetMainRedis())
	testService(t, userSessionService, createTestUserSession)
}

// TestUserSessionValid 测试session在增删改查中是否如预期的可验证、不可被验证
func TestUserSessionValid(t *testing.T) {
	userSessionService := orm.MustNewUserSessionService(config.GetMainDB(), config.GetMainRedis())
	// 测试路径：Create->Save(!Effective)
	userSession := createTestUserSession()
	effective, err := userSessionService.ValidSessionToken(userSession.UserID, userSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("user session %s effective before created", userSession.Token)
	}
	err = userSessionService.Save(userSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = userSessionService.ValidSessionToken(userSession.UserID, userSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if !effective {
		t.Fatalf("user session %s not effective after created", userSession.Token)
	}
	userSession.Effective = false
	err = userSessionService.Save(userSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = userSessionService.ValidSessionToken(userSession.UserID, userSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("user session %s effective after disabled", userSession.Token)
	}
	// 测试路径：Create->Delete
	userSession = createTestUserSession()
	err = userSessionService.Save(userSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = userSessionService.ValidSessionToken(userSession.UserID, userSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if !effective {
		t.Fatalf("user session %s not effective after created", userSession.Token)
	}
	err = userSessionService.Delete(userSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = userSessionService.ValidSessionToken(userSession.UserID, userSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("user session %s effective after disabled", userSession.Token)
	}
	// 测试路径：Create->DeleteByWhereOptions
	userSession = createTestUserSession()
	err = userSessionService.Save(userSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = userSessionService.ValidSessionToken(userSession.UserID, userSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if !effective {
		t.Fatalf("user session %s not effective after created", userSession.Token)
	}
	whereOptions := []orm.WhereOption{
		orm.WhereOption{Query: "id = ?", Item: []interface{}{userSession.ID}},
	}
	err = userSessionService.DeleteByWhereOptions(whereOptions)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = userSessionService.ValidSessionToken(userSession.UserID, userSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("user session %s effective after disabled", userSession.Token)
	}
}
