package csv_converter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFieldsIndexes(t *testing.T) {
	csvData := fmt.Sprintf(`ru,%s,en,%s,fr,%s`, columnDescription, columnName, columnParameters)
	r := csv.NewReader(bytes.NewReader([]byte(csvData)))
	indexes, err := getFieldsIndexes(r, "en")

	require.NoError(t, err)
	require.NotNil(t, indexes)
	require.Equal(t, 3, indexes.name)
	require.NotNil(t, indexes.description)
	require.Equal(t, 1, *indexes.description)
	require.NotNil(t, indexes.parameters)
	require.Equal(t, 5, *indexes.parameters)
}

func TestGetItems(t *testing.T) {
	csvData := fmt.Sprintf(`item1,descr1,par1;par2,val-ru-1,val-en-1
item2,descr2,para;parb2,val-ru-2,val-en-2`)

	r := csv.NewReader(bytes.NewReader([]byte(csvData)))
	descriptionInd, parametersInd := 1, 2
	indexes := &csvIndexes{
		name:        0,
		description: &descriptionInd,
		parameters:  &parametersInd,
		cultures: map[string]int{
			"ru": 3,
			"en": 4,
		},
		countFieldsInRow: 5,
	}

	items, err := getItems(r, indexes)
	require.NoError(t, err)
	require.NotNil(t, items)
	require.Len(t, items, 2)

	require.Contains(t, items, "item1")
	require.Contains(t, items, "item2")

	require.Contains(t, items["item1"].cultures, "ru")
	require.Contains(t, items["item1"].cultures, "en")

	require.Contains(t, items["item1"].parameters, "par1")
	require.Contains(t, items["item1"].parameters, "par2")
}
