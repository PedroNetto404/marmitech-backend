package types

type Filter map[string]any

type FindArgs struct {
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	SortBy  string `json:"sort_by"`
	SortAsc bool   `json:"sort_asc"`
	Filter  Filter `json:"filter"`
}
