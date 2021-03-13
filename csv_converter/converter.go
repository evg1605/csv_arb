package csv_converter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/evg1605/csv_arb/common"
	"github.com/sirupsen/logrus"
)

type CsvParams struct {
	ColumnName        string
	ColumnDescription string
	ColumnParameters  string
	DefaultCulture    string
}

var (
	errInvalidCsvParams    = errors.New("invalid csv params")
	errInvalidCsvStructure = errors.New("invalid csv")
)

type csvIndexes struct {
	name             int
	description      *int
	parameters       *int
	cultures         map[string]int
	countFieldsInRow int
}

func LoadArbFromWeb(logger *logrus.Logger, csvUrl string, csvParams CsvParams) (*common.DataArb, error) {
	logger.Tracef("download csv from url %s", csvUrl)
	r, err := csvFromWeb(logger, csvUrl)
	if err != nil {
		return nil, err
	}
	logger.Traceln("csv downloaded")

	return convertCsvToArb(logger, r, csvParams)
}

func LoadArbFromFile(logger *logrus.Logger, csvPath string, csvParams CsvParams) (*common.DataArb, error) {
	logger.Tracef("load csv from file %s", csvPath)
	r, err := csvFromFile(logger, csvPath)
	if err != nil {
		return nil, err
	}
	logger.Traceln("csv loaded")

	return convertCsvToArb(logger, r, csvParams)
}

func checkCsvParams(csvParams CsvParams) error {
	if csvParams.DefaultCulture == "" {
		return fmt.Errorf("invalid DefaultCulture: %w", errInvalidCsvParams)
	}
	if csvParams.ColumnName == "" {
		return fmt.Errorf("invalid ColumnName: %w", errInvalidCsvParams)
	}
	return nil
}

func convertCsvToArb(logger *logrus.Logger, r *csv.Reader, csvParams CsvParams) (*common.DataArb, error) {
	logger.Traceln("convert csv to arb")
	if err := checkCsvParams(csvParams); err != nil {
		return nil, err
	}
	fieldsIndexes, err := getFieldsIndexes(logger, r, csvParams)
	if err != nil {
		return nil, err
	}

	items, err := getArbItems(logger, r, fieldsIndexes)
	if err != nil {
		return nil, err
	}

	var cultures []string
	for cn := range fieldsIndexes.cultures {
		cultures = append(cultures, cn)
	}

	logger.Traceln("csv to arb converted")
	return &common.DataArb{
		Cultures: cultures,
		Items:    items,
	}, nil
}

func getArbItems(logger *logrus.Logger, r *csv.Reader, fieldsIndexes *csvIndexes) (map[string]*common.ItemArb, error) {
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

func getFieldsIndexes(logger *logrus.Logger, r *csv.Reader, csvParams CsvParams) (*csvIndexes, error) {
	// read first row and get indexes of Name and Description fields

	var nameInd, descriptionInd, parametersInd *int

	cultures := make(map[string]int)

	row, err := r.Read()
	if err != nil {
		return nil, err
	}

	m := map[string]**int{
		csvParams.ColumnName: &nameInd,
	}
	if csvParams.ColumnDescription != "" {
		m[csvParams.ColumnDescription] = &descriptionInd
	}
	if csvParams.ColumnParameters != "" {
		m[csvParams.ColumnParameters] = &parametersInd
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

	if _, ok := cultures[csvParams.DefaultCulture]; !ok {
		return nil, fmt.Errorf("csv must have column for default culture (%s): %w", csvParams.DefaultCulture, errInvalidCsvStructure)
	}
	return &csvIndexes{
		name:             *nameInd,
		description:      descriptionInd,
		parameters:       parametersInd,
		cultures:         cultures,
		countFieldsInRow: len(row),
	}, nil
}
