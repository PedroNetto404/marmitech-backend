package pagination

type PagedSlice[T any] struct {
	Meta struct {
		TotalRecords  int `json:"total_records"`
		RecordsLength int `json:"records_length"`
		CurrentPage   int `json:"current_page"`
		TotalPages    int `json:"total_pages"`
	} `json:"meta"`
	Records []T `json:"records"`
}

func New[T any](
	limit int,
	offset int,
	totalRecords int,
	records []T,
) PagedSlice[T] {
	var ps PagedSlice[T]
	ps.Records = records

	ps.Meta.TotalRecords = totalRecords
	ps.Meta.RecordsLength = len(records)

	if limit > 0 {
		ps.Meta.TotalPages = (totalRecords + limit - 1) / limit
		ps.Meta.CurrentPage = (offset / limit) + 1
	} else {
		ps.Meta.TotalPages = 1
		ps.Meta.CurrentPage = 1
	}

	return ps
}

type Selector[T, K any] func(T) K

func Map[T, K any](source PagedSlice[T], selector Selector[T, K]) PagedSlice[K] {
	destinyRecords := make([]K, source.Meta.RecordsLength)
	for index, record := range source.Records {
		destinyRecords[index] = selector(record)
	}

	var destiny PagedSlice[K]

	destiny.Records = destinyRecords
	destiny.Meta = source.Meta

	return destiny
}
