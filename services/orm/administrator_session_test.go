package services

import (
	"testing"
	"weassistant/models"

	"github.com/google/uuid"
)

func createTestAdministratorSession() *models.AdministratorSession {
	return &models.AdministratorSession{
		AdministratorID: 1,
		Effective:       true,
		Token:           uuid.New().String(),
	}
}

func TestAdministratorSession(t *testing.T) {
	administratorSessionService := MustNewAdministratorSessionService(config.GetMainDB(), config.GetMainRedis())
	testService(t, administratorSessionService, createTestAdministratorSession)
}

// TestAdministratorSessionValid 测试session在增删改查中是否如预期的可验证、不可被验证
func TestAdministratorSessionValid(t *testing.T) {
	administratorSessionService := MustNewAdministratorSessionService(config.GetMainDB(), config.GetMainRedis())
	// 测试路径：Create->Save(!Effective)
	administratorSession := createTestAdministratorSession()
	effective, err := administratorSessionService.ValidSessionToken(administratorSession.AdministratorID, administratorSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("administrator session %s effective before created", administratorSession.Token)
	}
	err = administratorSessionService.Save(administratorSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = administratorSessionService.ValidSessionToken(administratorSession.AdministratorID, administratorSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if !effective {
		t.Fatalf("administrator session %s not effective after created", administratorSession.Token)
	}
	administratorSession.Effective = false
	err = administratorSessionService.Save(administratorSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = administratorSessionService.ValidSessionToken(administratorSession.AdministratorID, administratorSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("administrator session %s effective after disabled", administratorSession.Token)
	}
	// 测试路径：Create->Delete
	administratorSession = createTestAdministratorSession()
	err = administratorSessionService.Save(administratorSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = administratorSessionService.ValidSessionToken(administratorSession.AdministratorID, administratorSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if !effective {
		t.Fatalf("administrator session %s not effective after created", administratorSession.Token)
	}
	err = administratorSessionService.Delete(administratorSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = administratorSessionService.ValidSessionToken(administratorSession.AdministratorID, administratorSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("administrator session %s effective after disabled", administratorSession.Token)
	}
	// 测试路径：Create->DeleteByWhereOptions
	administratorSession = createTestAdministratorSession()
	err = administratorSessionService.Save(administratorSession)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = administratorSessionService.ValidSessionToken(administratorSession.AdministratorID, administratorSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if !effective {
		t.Fatalf("administrator session %s not effective after created", administratorSession.Token)
	}
	whereOptions := []OrmWhereOption{
		OrmWhereOption{Query: "id = ?", Item: []interface{}{administratorSession.ID}},
	}
	err = administratorSessionService.DeleteByWhereOptions(whereOptions)
	if err != nil {
		t.Fatal(err)
	}
	effective, err = administratorSessionService.ValidSessionToken(administratorSession.AdministratorID, administratorSession.Token)
	if err != nil {
		t.Fatal(err)
	}
	if effective {
		t.Fatalf("administrator session %s effective after disabled", administratorSession.Token)
	}
}
