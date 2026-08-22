package service

import (
	"ruralfolk/domain"
	"sort"
)

func SortExhibitsByTitle(items []domain.Exhibit) []domain.Exhibit {
	out := append([]domain.Exhibit(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}
func PublishedOnly(items []domain.Exhibit) []domain.Exhibit {
	out := []domain.Exhibit{}
	for _, e := range items {
		if e.Status == domain.Published {
			out = append(out, e)
		}
	}
	return out
}
func FilterNews(items []domain.News) []domain.News {
	out := []domain.News{}
	for _, n := range items {
		if n.Published {
			out = append(out, n)
		}
	}
	return out
}
func SummarizeBooking(b domain.Booking) string    { return b.VisitorName + " @ " + b.VisitDate }
func DisplayMessage(g domain.GuestbookEntry) bool { return g.Status == domain.GuestbookApproved }
func FavoriteIDs(items []domain.Favorite) []string {
	out := make([]string, 0, len(items))
	for _, f := range items {
		out = append(out, f.ExhibitID)
	}
	return out
}
