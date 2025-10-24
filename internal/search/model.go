package search

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MovieData struct {
	Movies []Movie `json:"movies"`
}
