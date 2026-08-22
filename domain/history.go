package domain

import "sort"

type HistoryEvent struct{ ID, EntityID, EntityType, From, To, Actor, At, Note string }
type Timeline struct{ Events []HistoryEvent }

func (t *Timeline) Append(event HistoryEvent) error {
	if event.ID == "" || event.EntityID == "" || event.EntityType == "" {
		return errHistoryIdentity
	}
	t.Events = append(t.Events, event)
	return nil
}
func (t Timeline) ForEntity(id string) []HistoryEvent {
	out := []HistoryEvent{}
	for _, e := range t.Events {
		if e.EntityID == id {
			out = append(out, e)
		}
	}
	return out
}
func (t Timeline) Latest(id string) (HistoryEvent, bool) {
	events := t.ForEntity(id)
	if len(events) == 0 {
		return HistoryEvent{}, false
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].At < events[j].At })
	return events[len(events)-1], true
}
func (t Timeline) HasTransition(id, from, to string) bool {
	for _, e := range t.Events {
		if e.EntityID == id && e.From == from && e.To == to {
			return true
		}
	}
	return false
}
func (t Timeline) Actors(id string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, e := range t.Events {
		if e.EntityID == id && !seen[e.Actor] {
			seen[e.Actor] = true
			out = append(out, e.Actor)
		}
	}
	return out
}
func (t Timeline) Count(id string) int { return len(t.ForEntity(id)) }
func (t Timeline) FilterType(kind string) []HistoryEvent {
	out := []HistoryEvent{}
	for _, e := range t.Events {
		if e.EntityType == kind {
			out = append(out, e)
		}
	}
	return out
}
func (t Timeline) ContainsNote(id, note string) bool {
	for _, e := range t.Events {
		if e.EntityID == id && e.Note == note {
			return true
		}
	}
	return false
}
func EventForExhibit(id, from, to, actor, at, note string) HistoryEvent {
	return HistoryEvent{ID: id + ":" + to, EntityID: id, EntityType: "exhibit", From: from, To: to, Actor: actor, At: at, Note: note}
}
func EventForBooking(id, from, to, actor, at, note string) HistoryEvent {
	return HistoryEvent{ID: id + ":" + to, EntityID: id, EntityType: "booking", From: from, To: to, Actor: actor, At: at, Note: note}
}
func EventForGuestbook(id, from, to, actor, at, note string) HistoryEvent {
	return HistoryEvent{ID: id + ":" + to, EntityID: id, EntityType: "guestbook", From: from, To: to, Actor: actor, At: at, Note: note}
}
func (t Timeline) Clone() Timeline { return Timeline{Events: append([]HistoryEvent(nil), t.Events...)} }
func (t *Timeline) Clear()         { t.Events = nil }

var errHistoryIdentity = historyError("history identity required")

type historyError string

func (e historyError) Error() string { return string(e) }
