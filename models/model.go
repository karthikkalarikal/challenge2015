package models

type Config struct {
	Actor1     string
	Actor2     string
	Separation uint
}

type MovieConn struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type Actor struct {
	URL    string      `json:"url"`
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Movies []MovieConn `json:"movies"`
}

type Cast struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Crew struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type Movie struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Name string `json:"name"`
	Cast []Cast `json:"cast"`
	Crew []Crew `json:"crew"`
}
