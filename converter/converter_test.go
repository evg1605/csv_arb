package converter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFieldsIndexes(t *testing.T) {
	csvData := fmt.Sprintf(`ru,%s,en,%s,fr`, columnDescription, columnName)
	r := csv.NewReader(bytes.NewReader([]byte(csvData)))
	indexes, err := getFieldsIndexes(r, "en")

	require.NoError(t, err)
	require.NotNil(t, indexes)
	require.Equal(t, 3, indexes.name)
	require.NotNil(t, indexes.description)
	require.Equal(t, 1, *indexes.description)
}
