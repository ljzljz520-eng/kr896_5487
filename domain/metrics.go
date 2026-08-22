package domain

type Metrics struct{ Exhibits, Published, Artisans, Bookings, Confirmed, Messages, Approved, News, Users, Favorites int }

func (m Metrics) CompletionRate() float64 {
	if m.Exhibits == 0 {
		return 0
	}
	return float64(m.Published) / float64(m.Exhibits)
}
func (m Metrics) BookingConfirmationRate() float64 {
	if m.Bookings == 0 {
		return 0
	}
	return float64(m.Confirmed) / float64(m.Bookings)
}
func (m Metrics) ModerationRate() float64 {
	if m.Messages == 0 {
		return 0
	}
	return float64(m.Approved) / float64(m.Messages)
}
func (m Metrics) HasContent() bool  { return m.Exhibits > 0 || m.Artisans > 0 || m.News > 0 }
func (m Metrics) HasVisitors() bool { return m.Bookings > 0 || m.Messages > 0 || m.Users > 0 }
func (m Metrics) Add(other Metrics) Metrics {
	return Metrics{Exhibits: m.Exhibits + other.Exhibits, Published: m.Published + other.Published, Artisans: m.Artisans + other.Artisans, Bookings: m.Bookings + other.Bookings, Confirmed: m.Confirmed + other.Confirmed, Messages: m.Messages + other.Messages, Approved: m.Approved + other.Approved, News: m.News + other.News, Users: m.Users + other.Users, Favorites: m.Favorites + other.Favorites}
}
func (m Metrics) Empty() bool {
	return m.Exhibits == 0 && m.Artisans == 0 && m.Bookings == 0 && m.Messages == 0 && m.News == 0 && m.Users == 0 && m.Favorites == 0
}
func (m Metrics) Score() int {
	score := 0
	if m.Published > 0 {
		score++
	}
	if m.Artisans > 0 {
		score++
	}
	if m.Confirmed > 0 {
		score++
	}
	if m.Approved > 0 {
		score++
	}
	if m.News > 0 {
		score++
	}
	if m.Favorites > 0 {
		score++
	}
	return score
}
func (m Metrics) Labels() []string {
	labels := []string{}
	if m.Exhibits > 0 {
		labels = append(labels, "展品")
	}
	if m.Artisans > 0 {
		labels = append(labels, "手艺人")
	}
	if m.Bookings > 0 {
		labels = append(labels, "预约")
	}
	if m.Messages > 0 {
		labels = append(labels, "留言")
	}
	if m.News > 0 {
		labels = append(labels, "新闻")
	}
	if m.Favorites > 0 {
		labels = append(labels, "收藏")
	}
	return labels
}
func (m Metrics) PublicCount() int { return m.Published + m.Artisans + m.Approved + m.News }
func (m Metrics) PrivateCount() int {
	return m.Exhibits - m.Published + (m.Bookings - m.Confirmed) + (m.Messages - m.Approved) + m.Users
}
func (m Metrics) Normalize() Metrics {
	if m.Exhibits < 0 {
		m.Exhibits = 0
	}
	if m.Published < 0 {
		m.Published = 0
	}
	if m.Artisans < 0 {
		m.Artisans = 0
	}
	if m.Bookings < 0 {
		m.Bookings = 0
	}
	if m.Confirmed < 0 {
		m.Confirmed = 0
	}
	if m.Messages < 0 {
		m.Messages = 0
	}
	if m.Approved < 0 {
		m.Approved = 0
	}
	if m.News < 0 {
		m.News = 0
	}
	if m.Users < 0 {
		m.Users = 0
	}
	if m.Favorites < 0 {
		m.Favorites = 0
	}
	return m
}
func (m Metrics) AtCapacity() bool       { return m.Confirmed >= 50 }
func (m Metrics) NeedsModeration() bool  { return m.Messages > m.Approved }
func (m Metrics) NeedsPublication() bool { return m.Exhibits > m.Published }
