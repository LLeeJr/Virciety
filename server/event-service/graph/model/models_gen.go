// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateEventRequest struct {
	Title       string `json:"title"`
	Host        string `json:"host"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Location    string `json:"location"`
}

type CreateEventResponse struct {
	Event *Event `json:"event"`
	Type  string `json:"type"`
}

type EditEventRequest struct {
	EventID     string   `json:"eventID"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Location    string   `json:"location"`
}

type Event struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
	Host        string   `json:"host"`
}

type GetEventsResponse struct {
	UpcomingEvents []*Event `json:"upcomingEvents"`
	OngoingEvents  []*Event `json:"ongoingEvents"`
	PastEvents     []*Event `json:"pastEvents"`
}
