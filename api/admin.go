package api

import (
	"encoding/json"
	"net/http"
	"ruralfolk/domain"
	"ruralfolk/service"
)

type adminPayload struct {
	ID, Title, Story, Name, Message, Kind, Email, Role string
	Body                                               string
	Published, Confirmed, Approve                      bool
	PartySize                                          int
	VisitDate                                          string
}

func (a *Server) admin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var p adminPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSON(w, map[string]string{"error": "invalid json"}, 400)
		return
	}
	session, err := service.AuthenticateAdmin(domain.NewUser("admin", p.Email, p.Role, "fixture"))
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 403)
		return
	}
	switch p.Kind {
	case "exhibit":
		e := domain.NewExhibit(p.ID, p.Title, p.Story)
		e.MediaURL = p.Body
		if err = a.Service.EditExhibit(session, e); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, e, 200)
	case "artisan":
		a1 := domain.NewArtisan(p.ID, p.Name, p.Role, p.Body)
		if err = a.Service.EditArtisan(session, a1); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, a1, 200)
	case "news":
		n := domain.NewNews(p.ID, p.Title, p.Body, p.Published)
		if err = a.Service.EditNews(session, n); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, n, 200)
	case "guestbook":
		g, err := a.Service.ModerateGuestbook(session, p.ID, p.Approve)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, g, 200)
	default:
		writeJSON(w, map[string]string{"error": "unknown admin kind"}, 400)
	}
}
func (a *Server) search(w http.ResponseWriter, r *http.Request) {
	q := domain.SearchQuery{Text: r.URL.Query().Get("q"), Section: domain.Section(r.URL.Query().Get("section")), PublishedOnly: r.URL.Query().Get("published") != "false"}
	v, err := a.Service.Search(q)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 500)
		return
	}
	writeJSON(w, v, 200)
}
func (a *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	v, err := a.Service.Dashboard()
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 500)
		return
	}
	writeJSON(w, v, 200)
}
func (a *Server) calendar(w http.ResponseWriter, r *http.Request) {
	v, err := a.Service.Calendar()
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 500)
		return
	}
	writeJSON(w, v, 200)
}
func parseBool(value string, fallback bool) bool {
	if value == "" {
		return fallback
	}
	return value == "1" || value == "true" || value == "yes"
}
func (a *Server) featured(w http.ResponseWriter, r *http.Request) {
	items, err := a.Service.Featured(6)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 500)
		return
	}
	writeJSON(w, items, 200)
}
func (a *Server) favorites(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("user")
	if r.Method == http.MethodGet {
		items, err := a.Service.ListFavorites(uid)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, items, 200)
		return
	}
	var p adminPayload
	if json.NewDecoder(r.Body).Decode(&p) != nil {
		writeJSON(w, map[string]string{"error": "invalid json"}, 400)
		return
	}
	if r.Method == http.MethodPost {
		if err := a.Service.EnsureUserFavorite(uid, p.ID); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, map[string]string{"status": "saved"}, 201)
		return
	}
	if err := a.Service.RemoveFavorite(uid, p.ID); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	writeJSON(w, map[string]string{"status": "removed"}, 200)
}
