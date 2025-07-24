package main

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationAreaResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

type pokemon struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	BaseXp int    `json:"base_experience"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		}
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type encounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type locationEncounters struct {
	Encounters []encounter `json:"pokemon_encounters"`
}

type config struct {
	next string
	prev string
}
