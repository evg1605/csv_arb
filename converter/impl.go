package converter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/evg1605/csv_arb/csv_source"
)

const (
	columnName        = "name"
	columnDescription = "description"
)

var (
	errInvalidCsvStructure = errors.New("invalid csv")
)

type csvIndexes struct {
	name             int
	description      *int
	cultures         []*cultureInd
	countFieldsInRow int
}

type cultureInd struct {
	name string
	ind  int
}

type item struct {
	name        string
	description string
	cultures    map[string]string
}

func WebCsvToArb(csvUrl,
	arbFolderPath,
	arbFileTemplate,
	defaultCulture string) error {
	r, err := csv_source.CsvFromWeb(csvUrl)
	if err != nil {
		return err
	}

	return csvToArb(r,
		arbFolderPath,
		arbFileTemplate,
		defaultCulture)
}

func FileCsvToArb(csvPath,
	arbFolderPath,
	arbFileTemplate,
	defaultCulture string) error {
	r, err := csv_source.CsvFromFile(csvPath)
	if err != nil {
		return err
	}

	return csvToArb(r,
		arbFolderPath,
		arbFileTemplate,
		defaultCulture)
}

func csvToArb(r *csv.Reader, arbFolderPath,
	arbFileTemplate,
	defaultCulture string) error {

	fieldsIndexes, err := getFieldsIndexes(r, defaultCulture)
	if err != nil {
		return err
	}

	items, err := getItems(r, fieldsIndexes)
	if err != nil {
		return err
	}

	_ = items
	return nil
}

func getItems(r *csv.Reader, fieldsIndexes *csvIndexes) (map[string]*item, error) {
	items := make(map[string]*item)

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(row) != fieldsIndexes.countFieldsInRow {
			return nil, fmt.Errorf("invalid row with fields count %v, but expect %v: %w", len(row), fieldsIndexes.countFieldsInRow, errInvalidCsvStructure)
		}

		i := &item{
			name:     row[fieldsIndexes.name],
			cultures: make(map[string]string),
		}
		if fieldsIndexes.description != nil {
			i.description = row[*fieldsIndexes.description]
		}
		for _, ci := range fieldsIndexes.cultures {
			i.cultures[ci.name] = row[ci.ind]
		}
	}
	return items, nil
}

func getFieldsIndexes(r *csv.Reader, defaultCulture string) (*csvIndexes, error) {
	// read first row and get indexes of name and description fields

	var culturesInd []*cultureInd
	var nameInd, descriptionInd *int

	allCultures := make(map[string]struct{})

	row, err := r.Read()
	if err != nil {
		return nil, err
	}

	for i, f := range row {
		if f == "" {
			continue
		}
		iTmp := i
		switch f {
		case columnName:
			if nameInd != nil {
				return nil, fmt.Errorf("there should only be one column for the name: %w", errInvalidCsvStructure)
			}
			nameInd = &iTmp

		case columnDescription:
			if descriptionInd != nil {
				return nil, fmt.Errorf("there should only be one column for the description: %w", errInvalidCsvStructure)

			}
			descriptionInd = &iTmp
		default:
			if _, ok := allCultures[f]; ok {
				return nil, fmt.Errorf("each culture to be represented by only one column (%s): %w", f, errInvalidCsvStructure)
			}
			allCultures[f] = struct{}{}
			culturesInd = append(culturesInd, &cultureInd{
				name: f,
				ind:  i,
			})
		}
	}

	if nameInd == nil {
		return nil, fmt.Errorf("csv must have column for name: %w", errInvalidCsvStructure)
	}

	if len(allCultures) == 0 {
		return nil, fmt.Errorf("cultures not found: %w", errInvalidCsvStructure)
	}

	if _, ok := allCultures[defaultCulture]; !ok {
		return nil, fmt.Errorf("csv must have column for default culture (%s): %w", defaultCulture, errInvalidCsvStructure)
	}
	return &csvIndexes{
		name:             *nameInd,
		description:      descriptionInd,
		cultures:         culturesInd,
		countFieldsInRow: len(row),
	}, nil
}
