package cmd

import (
	"errors"
	"flag"
	"log"
)

const (
	arb2csvMode convertMode = "arb2csv"
	csv2arbMode convertMode = "csv2arb"
)

type convertMode string

type inputParams struct {
	mode            convertMode
	url             string
	csvPath         string
	arbFolderPath   string
	arbFileTemplate string
	defaultCulture  string
}

func Run() {
	params, err := getParams()
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}
	_ = params
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

	if mode == csv2arbMode {
		return nil, errors.New("csv-path is required for arb2csv mode")
	}

	if *csvPathFlag != "" && *csvUrlFlag != "" {
		return nil, errors.New("you need to specify only one parameter - csv-url or csv-path")
	}

	params := &inputParams{
		mode:            mode,
		url:             *csvUrlFlag,
		csvPath:         *csvPathFlag,
		arbFolderPath:   *arbFolderPathFlag,
		arbFileTemplate: *arbFileTemplateFlag,
		defaultCulture:  *defaultCultureFlag,
	}

	return params, nil
}
