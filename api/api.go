package api

import (
	"encoding/json"
	"net/http"
	"ruralfolk/domain"
	"ruralfolk/service"
)

type Server struct{ Service *service.Service }

func New(s *service.Service) *Server { return &Server{Service: s} }
func Routes(s *service.Service) http.Handler {
	a := New(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.health)
	mux.HandleFunc("/api/exhibits", a.exhibits)
	mux.HandleFunc("/api/bookings", a.bookings)
	mux.HandleFunc("/api/guestbook", a.guestbook)
	mux.HandleFunc("/api/artisans", a.artisans)
	mux.HandleFunc("/api/news", a.news)
	mux.HandleFunc("/api/search", a.search)
	mux.HandleFunc("/api/admin", a.admin)
	mux.HandleFunc("/api/dashboard", a.dashboard)
	mux.HandleFunc("/api/calendar", a.calendar)
	mux.HandleFunc("/api/featured", a.featured)
	mux.HandleFunc("/api/favorites", a.favorites)
	return mux
}
func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func (a *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}
func (a *Server) exhibits(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := a.Service.ListExhibits()
		if e != nil {
			writeJSON(w, map[string]string{"error": e.Error()}, 500)
			return
		}
		writeJSON(w, v, 200)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var e domain.Exhibit
	if json.NewDecoder(r.Body).Decode(&e) != nil {
		writeJSON(w, map[string]string{"error": "invalid json"}, 400)
		return
	}
	if e.Status == domain.Draft {
		if err := a.Service.CreateDraft(e); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, e, 201)
		return
	}
	if e.Status == domain.Submitted {
		if _, err := a.Service.SubmitExhibit(e.ID); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		v, err := a.Service.PublishExhibit(e.ID)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, 400)
			return
		}
		writeJSON(w, v, 200)
		return
	}
	writeJSON(w, map[string]string{"error": "unsupported state"}, 400)
}
func (a *Server) bookings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := a.Service.ListBookings()
		if e != nil {
			writeJSON(w, map[string]string{"error": e.Error()}, 500)
			return
		}
		writeJSON(w, v, 200)
		return
	}
	var b domain.Booking
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		writeJSON(w, map[string]string{"error": "invalid json"}, 400)
		return
	}
	v, e := a.Service.CreateBooking(b)
	if e != nil {
		writeJSON(w, map[string]string{"error": e.Error()}, 400)
		return
	}
	writeJSON(w, v, 201)
}
func (a *Server) guestbook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := a.Service.ListGuestbook()
		if e != nil {
			writeJSON(w, map[string]string{"error": e.Error()}, 500)
			return
		}
		writeJSON(w, v, 200)
		return
	}
	var g domain.GuestbookEntry
	if json.NewDecoder(r.Body).Decode(&g) != nil {
		writeJSON(w, map[string]string{"error": "invalid json"}, 400)
		return
	}
	v, e := a.Service.SubmitGuestbook(g)
	if e != nil {
		writeJSON(w, map[string]string{"error": e.Error()}, 400)
		return
	}
	writeJSON(w, v, 201)
}
func (a *Server) artisans(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.ListArtisans()
	if e != nil {
		writeJSON(w, map[string]string{"error": e.Error()}, 500)
		return
	}
	writeJSON(w, v, 200)
}
func (a *Server) news(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.ListNews()
	if e != nil {
		writeJSON(w, map[string]string{"error": e.Error()}, 500)
		return
	}
	writeJSON(w, service.FilterNews(v), 200)
}
