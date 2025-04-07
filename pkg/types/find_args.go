package types

type Filter map[string]any

type FindArgs struct {
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	SortBy  string `json:"sort_by"`
	SortAsc bool   `json:"sort_asc"`
	Filter  Filter `json:"filter"`
}

func NewDefaultFindArgs() FindArgs {
	return FindArgs{
		Limit:   99999999999,
		Offset:  0,
		SortBy:  "id",
		SortAsc: true,
	}
}