package models

type Config struct {
	Actor1     string
	Actor2     string
	Separation uint
}

type Movie struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type Actor struct {
	URL    string  `json:"url"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies"`
}
