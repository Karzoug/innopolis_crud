package domain

type Ingredient struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

type Recipe struct {
	ID          string       `json:"id"`
	AuthorID    string       `json:"author_id"`
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Temperature int          `json:"temperature"`
}
