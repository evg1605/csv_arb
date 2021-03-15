package main

import (
	"github.com/evg1605/csv_arb/arb"
	"github.com/evg1605/csv_arb/csv"

	"github.com/sirupsen/logrus"
	"github.com/thatisuday/commando"
)

func arb2csv(logger *logrus.Logger, flags map[string]commando.FlagValue) error {
	arbData, err := arb.LoadArb(logger, getStrFromFlag(flags, arbPathFlag), getStrFromFlag(flags, cultureFlag))
	if err != nil {
		return err
	}

	csvParams := csv.Params{
		ColumnName:        getStrFromFlag(flags, colNameFlag),
		ColumnDescription: getStrFromFlag(flags, colDescrFlag),
		ColumnParameters:  getStrFromFlag(flags, colParamsFlag),
		DefaultCulture:    getStrFromFlag(flags, cultureFlag),
	}
	return csv.SaveArb(logger, getStrFromFlag(flags, csvPathFlag), csvParams, arbData)
}
