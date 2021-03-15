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

	csvParams := csv.Params{
		ColumnName:        getStrFromFlag(flags, colNameFlag),
		ColumnDescription: getStrFromFlag(flags, colDescrFlag),
		ColumnParameters:  getStrFromFlag(flags, colParamsFlag),
		DefaultCulture:    getStrFromFlag(flags, cultureFlag),
	}

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
