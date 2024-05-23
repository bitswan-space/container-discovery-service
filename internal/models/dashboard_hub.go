package models

type DashboardEntryListResponse struct {
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []*DashboardEntry `json:"results"`
}

type DashboardEntry struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}
