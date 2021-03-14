package main

import (
	"strings"

	"github.com/evg1605/csv_arb/arb_converter"
	"github.com/evg1605/csv_arb/common"
	"github.com/evg1605/csv_arb/csv_converter"
	"github.com/sirupsen/logrus"
	"github.com/thatisuday/commando"
)

func csv2arb(logger *logrus.Logger, flags map[string]commando.FlagValue) error {
	src, _ := flags[srcFlag].GetString()

	csvParams := csv_converter.CsvParams{}
	csvParams.ColumnName, _ = flags[colNameFlag].GetString()
	csvParams.ColumnDescription, _ = flags[colDescrFlag].GetString()
	csvParams.ColumnParameters, _ = flags[colParamsFlag].GetString()
	csvParams.DefaultCulture, _ = flags[cultureFlag].GetString()

	var arbData *common.DataArb
	var arbDataErr error
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		arbData, arbDataErr = csv_converter.LoadArbFromWeb(logger, src, csvParams)
	} else {
		arbData, arbDataErr = csv_converter.LoadArbFromFile(logger, src, csvParams)
	}
	if arbDataErr != nil {
		return arbDataErr
	}

	arbFolderPath, _ := flags[arbPathFlag].GetString()
	arbFileTemplate, _ := flags[arbTemplateFlag].GetString()
	return arb_converter.SaveArb(logger, arbData, arbFolderPath, arbFileTemplate, csvParams.DefaultCulture)
}
