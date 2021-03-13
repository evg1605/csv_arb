package main

import (
	"errors"
	"flag"
	"fmt"
	"path"
	"runtime"

	"github.com/evg1605/csv_arb/arb_converter"
	"github.com/evg1605/csv_arb/common"
	"github.com/evg1605/csv_arb/csv_converter"
	"github.com/sirupsen/logrus"
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
	csvParams       csv_converter.CsvParams
	arbFolderPath   string
	arbFileTemplate string
	defaultCulture  string
	logLevel        logrus.Level
}

func main() {
	logger := createLogger(logrus.ErrorLevel)

	params, err := getParams()
	if err != nil {
		flag.PrintDefaults()
		logger.Fatal(err)
	}
	logger.SetLevel(params.logLevel)

	logger.Traceln("input params:")
	logger.Tracef("mode: %s", params.mode)
	logger.Tracef("logLevel: %s", params.logLevel)
	logger.Tracef("defaultCulture: %s", params.defaultCulture)
	logger.Tracef("csvUrl: %s", params.csvUrl)
	logger.Tracef("csvColumnName: %s", params.csvParams.ColumnName)
	logger.Tracef("csvColumnDescription: %s", params.csvParams.ColumnDescription)
	logger.Tracef("csvColumnParameters: %s", params.csvParams.ColumnParameters)
	logger.Tracef("csvPath: %s", params.csvPath)
	logger.Tracef("arbFolderPath: %s", params.arbFolderPath)
	logger.Tracef("arbFileTemplate: %s", params.arbFileTemplate)
	logger.Traceln()

	switch params.mode {
	case csv2arbMode:
		logger.Tracef("convert csv to arb in folder %s", params.arbFolderPath)
		if err := convertCsvToArb(logger, params); err != nil {
			logger.Fatal(err)
		}
	default:
		logger.Fatal(fmt.Errorf("unsupported mode %s", params.mode))
	}
	logger.Traceln("converted")
}

func convertCsvToArb(logger *logrus.Logger, params *inputParams) error {
	var dataArb *common.DataArb
	if params.csvUrl != "" {
		logger.Tracef("load arb data from csv web source %s", params.csvUrl)
		d, err := csv_converter.LoadArbFromWeb(logger, params.csvUrl, params.csvParams)
		if err != nil {
			return err
		}
		logger.Traceln("arb data loaded")
		dataArb = d
	} else if params.csvPath != "" {
		logger.Tracef("load arb data from csv file source %s", params.csvUrl)
		d, err := csv_converter.LoadArbFromFile(logger, params.csvPath, params.csvParams)
		if err != nil {
			return err
		}
		logger.Traceln("arb data loaded")
		dataArb = d
	} else {
		return errors.New("need to pass csv-url or csv-path")
	}

	logger.Tracef("save arb data to %s", params.arbFolderPath)
	if err := arb_converter.SaveArb(logger, dataArb, params.arbFolderPath, params.arbFileTemplate, params.defaultCulture); err != nil {
		return err
	}
	logger.Traceln("arb data saved")
	return nil
}

func getParams() (*inputParams, error) {
	modeFlag := flag.String("mode", "", "mode of conversion (arb2csv or csv2arb)")
	csvUrlFlag := flag.String("csv-url", "", "url of csv file")
	csvPathFlag := flag.String("csv-path", "", "csv file path")
	csvColNameFlag := flag.String("csv-col-name", "name", "name column name in csv table")
	csvColDescrFlag := flag.String("csv-col-descr", "description", "name column description in csv table")
	csvColParamsFlag := flag.String("csv-col-params", "parameters", "name column parameters in csv table")
	arbFolderPathFlag := flag.String("arb-path", "", "arb folder path (folder contains arb files - one for every culture)")
	arbFileTemplateFlag := flag.String("arb-template", "app_{culture}.arb", "arb file template")
	defaultCultureFlag := flag.String("default-culture", "en", "default culture")
	logLevelFlag := flag.String("log", "error", "log level (trace, debug, info, warning, error, fatal, panic)")

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

	logLevel, err := logrus.ParseLevel(*logLevelFlag)
	if err != nil {
		return nil, err
	}

	params := &inputParams{
		mode:    mode,
		csvUrl:  *csvUrlFlag,
		csvPath: *csvPathFlag,
		csvParams: csv_converter.CsvParams{
			ColumnName:        *csvColNameFlag,
			ColumnDescription: *csvColDescrFlag,
			ColumnParameters:  *csvColParamsFlag,
			DefaultCulture:    *defaultCultureFlag},
		arbFolderPath:   *arbFolderPathFlag,
		arbFileTemplate: *arbFileTemplateFlag,
		defaultCulture:  *defaultCultureFlag,
		logLevel:        logLevel,
	}

	return params, nil
}

func createLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		PadLevelText: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
	logger.SetLevel(level)
	return logger
}
