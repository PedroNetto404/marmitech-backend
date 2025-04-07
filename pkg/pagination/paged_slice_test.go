package pagination_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PedroNetto404/marmitech-backend/pkg/pagination"
)

func TestPagedSliceMap(t *testing.T) {
	// arrange
	records := []int{1, 2, 3}
	limit := 10
	offset := 0
	totalRecords := 23

	pagedSlice := pagination.New(limit, offset, totalRecords, records)

	// act
	result := pagination.Map(pagedSlice, func(num int) string {
		return strconv.Itoa(num)
	})

	//assert
	totalPagesExpeced := 3
	currentPageExpected := 1
	recordsLengthExpected := 3
	recordsExpected := []string{"1", "2", "3"}

	assert := assert.New(t)

	assert.Equal(recordsLengthExpected, result.Meta.RecordsLength, "deveria ter 3 registros")
	assert.Equal(totalRecords, result.Meta.TotalRecords, "deveria ter 23 registros totais")
	assert.Equal(totalPagesExpeced, result.Meta.TotalPages, "deveria ter 3 páginas (ceil de 23/10)")
	assert.Equal(currentPageExpected, result.Meta.CurrentPage, "primeira página com offset 0")
	assert.Equal(recordsExpected, result.Records, "registros devem ser transformados em string")
}
