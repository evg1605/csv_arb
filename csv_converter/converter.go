package csv_converter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/evg1605/csv_arb/common"
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

func LoadArbFromWeb(csvUrl, defaultCulture string) (*common.DataArb, error) {
	r, err := csvFromWeb(csvUrl)
	if err != nil {
		return nil, err
	}

	return loadArb(r, defaultCulture)
}

func LoadArbFromFile(csvPath, defaultCulture string) (*common.DataArb, error) {
	r, err := csvFromFile(csvPath)
	if err != nil {
		return nil, err
	}

	return loadArb(r, defaultCulture)
}

func loadArb(r *csv.Reader, defaultCulture string) (*common.DataArb, error) {
	fieldsIndexes, err := getFieldsIndexes(r, defaultCulture)
	if err != nil {
		return nil, err
	}

	items, err := getItems(r, fieldsIndexes)
	if err != nil {
		return nil, err
	}

	var cultures []string
	for cn, _ := range fieldsIndexes.cultures {
		cultures = append(cultures, cn)
	}

	return &common.DataArb{
		Cultures: cultures,
		Items:    items,
	}, nil
}

func getItems(r *csv.Reader, fieldsIndexes *csvIndexes) (map[string]*common.ItemArb, error) {
	items := make(map[string]*common.ItemArb)

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
			return nil, fmt.Errorf("found more than one key with same Name %s: %w", name, errInvalidCsvStructure)
		}

		i := &common.ItemArb{
			Cultures: make(map[string]string),
		}

		if fieldsIndexes.description != nil {
			i.Description = row[*fieldsIndexes.description]
		}

		if fieldsIndexes.parameters != nil {
			parameters := make(map[string]struct{})
			parametersRaw := strings.Split(row[*fieldsIndexes.parameters], ";")
			for _, p := range parametersRaw {
				pName := strings.Trim(p, " ")
				if pName == "" {
					continue
				}
				if _, ok := parameters[pName]; ok {
					return nil, fmt.Errorf("key %s has more than one parameter with Name %s : %w", name, pName, errInvalidCsvStructure)
				}
				parameters[pName] = struct{}{}
			}
			if len(parameters) > 0 {
				i.Parameters = parameters
			}
		}

		for cn, ci := range fieldsIndexes.cultures {
			i.Cultures[cn] = row[ci]
		}

		items[name] = i
	}
	return items, nil
}

func getFieldsIndexes(r *csv.Reader, defaultCulture string) (*csvIndexes, error) {
	// read first row and get indexes of Name and Description fields

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
		return nil, fmt.Errorf("csv must have column for Name: %w", errInvalidCsvStructure)
	}

	if len(cultures) == 0 {
		return nil, fmt.Errorf("Cultures not found: %w", errInvalidCsvStructure)
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
