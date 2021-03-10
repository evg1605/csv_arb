package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/evg1605/csv_arb/arb_converter"
	"github.com/evg1605/csv_arb/common"
	"github.com/evg1605/csv_arb/csv_converter"
)

const (
	arb2csvMode convertMode = "arb2csv"
	csv2arbMode convertMode = "csv2arb"
)

type convertMode string

type inputParams struct {
	mode            convertMode
	csvUrl          string
	csvPath         string
	arbFolderPath   string
	arbFileTemplate string
	defaultCulture  string
}

func main() {
	params, err := getParams()
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	switch params.mode {
	case csv2arbMode:
		if err := convertCsvToArb(params); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(fmt.Errorf("unsupported mode %s", params.mode))
	}
}

func convertCsvToArb(params *inputParams) error {
	var dataArb *common.DataArb
	if params.csvUrl != "" {
		d, err := csv_converter.LoadArbFromWeb(params.csvUrl, params.defaultCulture)
		if err != nil {
			return err
		}
		dataArb = d
	} else if params.csvPath != "" {
		d, err := csv_converter.LoadArbFromFile(params.csvPath, params.defaultCulture)
		if err != nil {
			return err
		}
		dataArb = d
	} else {
		return errors.New("need to pass csv-url or csv-path")
	}
	return arb_converter.SaveArb(dataArb, params.arbFolderPath, params.arbFileTemplate, params.defaultCulture)
}

func getParams() (*inputParams, error) {
	modeFlag := flag.String("mode", "", "mode of conversion (arb2csv or csv2arb)")
	csvUrlFlag := flag.String("csv-url", "", "url of csv file")
	csvPathFlag := flag.String("csv-path", "", "csv file path")
	arbFolderPathFlag := flag.String("arb-path", "", "arb folder path (folder contains arb files - one for every culture)")
	arbFileTemplateFlag := flag.String("arb-template", "app_{culture}.arb", "arb file template")
	defaultCultureFlag := flag.String("default-culture", "en", "default culture")

	flag.Parse()
	if !flag.Parsed() {
		return nil, errors.New("invalid parameters")
	}

	mode := convertMode(*modeFlag)
	if mode != arb2csvMode && mode != csv2arbMode {
		return nil, errors.New("invalid mode")
	}

	if (*csvPathFlag != "" && *csvUrlFlag != "") || (*csvPathFlag == "" && *csvUrlFlag == "") {
		return nil, errors.New("you need to specify one parameter - csv-url or csv-path")
	}

	params := &inputParams{
		mode:            mode,
		csvUrl:          *csvUrlFlag,
		csvPath:         *csvPathFlag,
		arbFolderPath:   *arbFolderPathFlag,
		arbFileTemplate: *arbFileTemplateFlag,
		defaultCulture:  *defaultCultureFlag,
	}

	return params, nil
}
