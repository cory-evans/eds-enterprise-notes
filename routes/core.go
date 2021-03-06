package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/CoryEvans2324/eds-enterprise-notes/database"
	"github.com/CoryEvans2324/eds-enterprise-notes/middleware"
	"github.com/CoryEvans2324/eds-enterprise-notes/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	jwtUser := middleware.GetUser(r)
	if jwtUser == nil {
		tmpl, _ := template.ParseFiles("web/index.html", "web/base.layout.html")
		tmpl.Execute(w, nil)
		return
	}

	tmpl, _ := template.ParseFiles("web/index-with-user.html", "web/base.layout.html")

	// find all notes relevant to the current user.
	user, _ := database.Mgr.GetUserByID(jwtUser.UserID)
	owned, err := database.Mgr.GetNotesByOwner(user)
	if err != nil {
		tmpl.Execute(w, struct{ User *models.User }{User: user})
		log.Println("owned: ", err)
		return
	}
	delegated, err := database.Mgr.GetNotesByDelegatedUser(user)
	if err != nil {
		tmpl.Execute(w, struct{ User *models.User }{User: user})
		log.Println("delegated: ", err)
		return
	}
	shared, err := database.Mgr.GetNotesSharedWith(user)
	if err != nil {
		tmpl.Execute(w, struct{ User *models.User }{User: user})
		log.Println("shared: ", err)
		return
	}

	tmpl.Execute(
		w,
		struct {
			User           *models.User
			OwnedNotes     []models.Note
			DelegatedNotes []models.Note
			SharedNotes    []models.Note
		}{
			User:           user,
			OwnedNotes:     owned,
			DelegatedNotes: delegated,
			SharedNotes:    shared,
		},
	)
}
