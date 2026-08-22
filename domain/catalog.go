package domain

import "strings"

type Section string

const (
	SectionStories   Section = "stories"
	SectionArtisans  Section = "artisans"
	SectionVisits    Section = "visits"
	SectionNews      Section = "news"
	SectionGuestbook Section = "guestbook"
)

type ContentRecord struct {
	ID, Kind, Title, Summary, Body, ImageURL string
	Published                                bool
	SortOrder                                int
}
type SearchQuery struct {
	Text          string
	Section       Section
	PublishedOnly bool
	Limit         int
}
type SearchResult struct {
	Exhibits []Exhibit
	Artisans []Artisan
	News     []News
	Total    int
}

func NormalizeSearch(q SearchQuery) SearchQuery {
	q.Text = strings.TrimSpace(strings.ToLower(q.Text))
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	return q
}
func MatchText(value, term string) bool {
	if term == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(term))
}
func MatchExhibit(e Exhibit, q SearchQuery) bool {
	if q.PublishedOnly && e.Status != Published {
		return false
	}
	return MatchText(e.Title, q.Text) || MatchText(e.Story, q.Text)
}
func MatchArtisan(a Artisan, q SearchQuery) bool {
	return MatchText(a.Name, q.Text) || MatchText(a.Bio, q.Text) || MatchText(a.Craft, q.Text)
}
func MatchNews(n News, q SearchQuery) bool {
	if q.PublishedOnly && !n.Published {
		return false
	}
	return MatchText(n.Title, q.Text) || MatchText(n.Body, q.Text)
}
func SectionLabel(s Section) string {
	switch s {
	case SectionStories:
		return "展品故事"
	case SectionArtisans:
		return "手艺人"
	case SectionVisits:
		return "参观预约"
	case SectionNews:
		return "活动新闻"
	case SectionGuestbook:
		return "留言板"
	}
	return "展陈"
}
func ValidSection(s Section) bool {
	return s == SectionStories || s == SectionArtisans || s == SectionVisits || s == SectionNews || s == SectionGuestbook
}
func StatusLabel(s ExhibitStatus) string {
	switch s {
	case Draft:
		return "草稿"
	case Submitted:
		return "待发布"
	case Published:
		return "已发布"
	}
	return "未知"
}
func BookingLabel(s BookingStatus) string {
	if s == BookingConfirmed {
		return "已确认"
	}
	return "待处理"
}
func GuestbookLabel(s GuestbookStatus) string {
	if s == GuestbookApproved {
		return "已展示"
	}
	return "待审核"
}
func NewExhibit(id, title, story string) Exhibit {
	return Exhibit{ID: id, Title: strings.TrimSpace(title), Story: strings.TrimSpace(story), Status: Draft}
}
func NewArtisan(id, name, craft, bio string) Artisan {
	return Artisan{ID: id, Name: strings.TrimSpace(name), Craft: strings.TrimSpace(craft), Bio: strings.TrimSpace(bio)}
}
func NewBooking(id, visitor, date string, size int) Booking {
	return Booking{ID: id, VisitorName: strings.TrimSpace(visitor), VisitDate: strings.TrimSpace(date), PartySize: size, Status: BookingPending}
}
func NewGuestbook(id, name, msg string) GuestbookEntry {
	return GuestbookEntry{ID: id, Name: strings.TrimSpace(name), Message: strings.TrimSpace(msg), Status: GuestbookPending}
}
func NewUser(id, email, role, hash string) User {
	return User{ID: id, Email: strings.TrimSpace(email), Role: role, PasswordHash: hash}
}
func NewFavorite(id, user, exhibit, created string) Favorite {
	return Favorite{ID: id, UserID: user, ExhibitID: exhibit, CreatedAt: created}
}
func NewNews(id, title, body string, published bool) News {
	return News{ID: id, Title: strings.TrimSpace(title), Body: strings.TrimSpace(body), Published: published}
}
func CopyExhibit(e Exhibit) Exhibit                 { return e }
func CopyArtisan(a Artisan) Artisan                 { return a }
func CopyBooking(b Booking) Booking                 { return b }
func CopyGuestbook(g GuestbookEntry) GuestbookEntry { return g }
func CopyNews(n News) News                          { return n }
func IsVisibleExhibit(e Exhibit) bool               { return e.Status == Published && e.Title != "" && e.Story != "" }
func IsVisibleArtisan(a Artisan) bool               { return a.Name != "" && a.Bio != "" }
func IsVisibleNews(n News) bool                     { return n.Published }
func IsPendingBooking(b Booking) bool               { return b.Status == BookingPending }
func IsConfirmedBooking(b Booking) bool             { return b.Status == BookingConfirmed }
func IsPendingGuestbook(g GuestbookEntry) bool      { return g.Status == GuestbookPending }
func IsApprovedGuestbook(g GuestbookEntry) bool     { return g.Status == GuestbookApproved }
func CountWords(s string) int                       { return len(strings.Fields(s)) }
func TrimForCard(s string, max int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= max {
		return string(r)
	}
	if max < 4 {
		return string(r[:max])
	}
	return string(r[:max-1]) + "…"
}
func ExhibitCard(e Exhibit) ContentRecord {
	return ContentRecord{ID: e.ID, Kind: "exhibit", Title: e.Title, Summary: TrimForCard(e.Story, 80), Body: e.Story, ImageURL: e.MediaURL, Published: e.Status == Published}
}
func ArtisanCard(a Artisan) ContentRecord {
	return ContentRecord{ID: a.ID, Kind: "artisan", Title: a.Name, Summary: TrimForCard(a.Craft, 80), Body: a.Bio, ImageURL: a.PortraitURL, Published: true}
}
func NewsCard(n News) ContentRecord {
	return ContentRecord{ID: n.ID, Kind: "news", Title: n.Title, Summary: TrimForCard(n.Body, 80), Body: n.Body, Published: n.Published}
}
