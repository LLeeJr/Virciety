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
	Subscribers []string `json:"subscribers"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Location    string   `json:"location"`
	Attendees   []string `json:"attendees"`
}

type EditEventResponse struct {
	Ok   string `json:"ok"`
	Type string `json:"type"`
}

type Event struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	StartDate         string   `json:"startDate"`
	EndDate           string   `json:"endDate"`
	Location          string   `json:"location"`
	Description       string   `json:"description"`
	Subscribers       []string `json:"subscribers"`
	Host              string   `json:"host"`
	Attendees         []string `json:"attendees"`
	CurrentlyAttended *bool    `json:"currentlyAttended"`
}

type GetEventsResponse struct {
	UpcomingEvents []*Event `json:"upcomingEvents"`
	OngoingEvents  []*Event `json:"ongoingEvents"`
	PastEvents     []*Event `json:"pastEvents"`
}

type UserData struct {
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Street      string `json:"street"`
	Housenumber string `json:"housenumber"`
	Postalcode  string `json:"postalcode"`
	City        string `json:"city"`
	Email       string `json:"email"`
}

type UserDataRequest struct {
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Street      string `json:"street"`
	Housenumber string `json:"housenumber"`
	Postalcode  string `json:"postalcode"`
	City        string `json:"city"`
	Email       string `json:"email"`
}