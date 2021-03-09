package converter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/evg1605/csv_arb/csv_source"
)

const (
	columnName        = "name"
	columnDescription = "description"
	columnParameters  = "parameters"
)

var (
	errInvalidCsvStructure = errors.New("invalid csv")
)

type csvIndexes struct {
	name             int
	description      *int
	parameters       *int
	cultures         map[string]int
	countFieldsInRow int
}

type item struct {
	name        string
	description string
	cultures    map[string]string
	parameters  map[string]struct{}
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

	if err := os.RemoveAll(arbFolderPath); err != nil {
		return err
	}

	if err := os.MkdirAll(arbFolderPath, 0777); err != nil {
		return err
	}

	//fieldsIndexes.cultures
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

		name := row[fieldsIndexes.name]
		if _, ok := items[name]; ok {
			return nil, fmt.Errorf("found more than one key with same name %s: %w", name, errInvalidCsvStructure)
		}

		i := &item{
			name:     name,
			cultures: make(map[string]string),
		}

		if fieldsIndexes.description != nil {
			i.description = row[*fieldsIndexes.description]
		}

		if fieldsIndexes.parameters != nil {
			i.parameters = make(map[string]struct{})
			parameters := strings.Split(row[*fieldsIndexes.parameters], ";")
			for _, p := range parameters {
				pName := strings.Trim(p, " ")
				if _, ok := i.parameters[pName]; ok {
					return nil, fmt.Errorf("key %s has more than one parameter with name %s : %w", name, pName, errInvalidCsvStructure)
				}
				i.parameters[pName] = struct{}{}
			}
		}

		for cn, ci := range fieldsIndexes.cultures {
			i.cultures[cn] = row[ci]
		}

		items[name] = i
	}
	return items, nil
}

func getFieldsIndexes(r *csv.Reader, defaultCulture string) (*csvIndexes, error) {
	// read first row and get indexes of name and description fields

	var nameInd, descriptionInd, parametersInd *int

	cultures := make(map[string]int)

	row, err := r.Read()
	if err != nil {
		return nil, err
	}

	m := map[string]**int{
		columnName:        &nameInd,
		columnDescription: &descriptionInd,
		columnParameters:  &parametersInd,
	}

	for i, f := range row {
		if f == "" {
			continue
		}

		ind, ok := m[f]
		if ok {
			if *ind != nil {
				return nil, fmt.Errorf("there should only be one column for the %s: %w", f, errInvalidCsvStructure)
			}
			iTmp := i
			*ind = &iTmp
			continue
		}

		if _, ok := cultures[f]; ok {
			return nil, fmt.Errorf("each culture to be represented by only one column (%s): %w", f, errInvalidCsvStructure)
		}
		cultures[f] = i
	}

	if nameInd == nil {
		return nil, fmt.Errorf("csv must have column for name: %w", errInvalidCsvStructure)
	}

	if len(cultures) == 0 {
		return nil, fmt.Errorf("cultures not found: %w", errInvalidCsvStructure)
	}

	if _, ok := cultures[defaultCulture]; !ok {
		return nil, fmt.Errorf("csv must have column for default culture (%s): %w", defaultCulture, errInvalidCsvStructure)
	}
	return &csvIndexes{
		name:             *nameInd,
		description:      descriptionInd,
		parameters:       parametersInd,
		cultures:         cultures,
		countFieldsInRow: len(row),
	}, nil
}
