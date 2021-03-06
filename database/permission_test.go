package database

import (
	"testing"

	"github.com/CoryEvans2324/eds-enterprise-notes/models"
)

func TestPermissions(t *testing.T) {
	cfg := createConfig()
	CreateDatabaseManager(cfg.Database.DataSourceName())
	ResetDB()

	// create test users
	owner, _ := Mgr.CreateUser("owner", "password")
	other, _ := Mgr.CreateUser("other", "password")

	// create test note
	note := &models.Note{
		Name:    "testing",
		Content: "test",
		Owner:   owner,
		Status:  "In Progress",
		SharedUsers: []models.Permission{
			{
				Permission: "viewer",
				User:       *other,
			},
		},
	}

	note, err := Mgr.CreateNote(note)
	if err != nil {
		t.Errorf("Cannot create note: %v", err)
	}

	// test permissions
	if len(note.SharedUsers) == 0 {
		t.Error("Note SharedUsers list is empty")
	}

	// remove the permission
	err = Mgr.RemovePermission(note.SharedUsers[0])
	checkErrNil(t, err)

	// create
	perm := models.Permission{
		Permission: "editor",
		User:       *note.Owner,
		Note:       *note,
	}
	err = Mgr.CreatePermission(perm)
	checkErrNil(t, err)

	perm.Permission = "viewer"
	// update
	err = Mgr.UpdatePermission(perm)
	checkErrNil(t, err)

}
