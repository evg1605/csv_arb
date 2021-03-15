package main

import (
	"strings"

	"github.com/evg1605/csv_arb/arb"
	"github.com/evg1605/csv_arb/csv"
	"github.com/sirupsen/logrus"
	"github.com/thatisuday/commando"
)

func csv2arb(logger *logrus.Logger, flags map[string]commando.FlagValue) error {
	src := getStrFromFlag(flags, csvPathFlag)

	csvParams := csv.Params{}
	csvParams.ColumnName, _ = flags[colNameFlag].GetString()
	csvParams.ColumnDescription, _ = flags[colDescrFlag].GetString()
	csvParams.ColumnParameters, _ = flags[colParamsFlag].GetString()
	csvParams.DefaultCulture, _ = flags[cultureFlag].GetString()

	var arbData *arb.Data
	var arbDataErr error
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		arbData, arbDataErr = csv.LoadArbFromWeb(logger, src, csvParams)
	} else {
		arbData, arbDataErr = csv.LoadArbFromFile(logger, src, csvParams)
	}
	if arbDataErr != nil {
		return arbDataErr
	}

	return arb.SaveArb(logger,
		arbData,
		getStrFromFlag(flags, arbPathFlag),
		getStrFromFlag(flags, arbTemplateFlag),
		csvParams.DefaultCulture)
}
