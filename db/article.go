package db

type Article struct {
	Model
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}
