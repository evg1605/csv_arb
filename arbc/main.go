package main

import (
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/thatisuday/commando"
)

const (
	csvPathFlag     = "csv-path"
	arbTemplateFlag = "arb-template"
	srcFlag         = "src"
	colNameFlag     = "col-name"
	arbPathFlag     = "arb-path"
	colDescrFlag    = "col-descr"
	colParamsFlag   = "col-params"
	cultureFlag     = "culture"
	logLevelFlag    = "log-level"
)

var AppVersion = "develop"

func main() {
	r := commando.
		SetExecutableName("arbc").
		SetVersion(AppVersion).
		SetDescription("Convertor from csv to arb and from arb to csv.")

	var rootCmd *commando.Command
	rootCmd = commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			r.PrintHelp(rootCmd)
		})

	var csv2arbCmd *commando.Command
	csv2arbCmd = commando.
		Register("csv2arb").
		SetDescription("convert csv to arb").
		SetShortDescription("convert csv to arb").
		AddFlag(srcFlag, "url or path of csv file", commando.String, "").
		AddFlag(arbTemplateFlag, "arb file template", commando.String, "app_{culture}.arb").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			baseAction(r, csv2arbCmd, flags, csv2arb)
		})
	addCommonFlags(csv2arbCmd)

	var arb2csvCmd *commando.Command
	arb2csvCmd = commando.
		Register("arb2csv").
		SetDescription("convert arb to csv").
		SetShortDescription("convert arb to csv").
		AddFlag(csvPathFlag, "path to csv file", commando.String, "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			baseAction(r, arb2csvCmd, flags, csv2arb)
		})
	addCommonFlags(arb2csvCmd)

	commando.Parse(nil)
}

func addCommonFlags(c *commando.Command) *commando.Command {
	c.
		AddFlag(arbPathFlag, "arb folder path (folder contains arb files - one for every culture)", commando.String, "").
		AddFlag(colNameFlag, "name column name in csv table", commando.String, "name").
		AddFlag(colDescrFlag, "name column name in csv table", commando.String, "description").
		AddFlag(colParamsFlag, "name column name in csv table", commando.String, "parameters").
		AddFlag(cultureFlag, "default culture", commando.String, "en").
		AddFlag(logLevelFlag, "log level (trace, debug, info, warning, error, fatal, panic)", commando.String, "error")
	return c
}

func baseAction(r *commando.CommandRegistry, c *commando.Command, flags map[string]commando.FlagValue, action func(*logrus.Logger, map[string]commando.FlagValue) error) {
	logger, err := createLogger(flags["log-level"])
	if err != nil {
		r.PrintHelp(c)
		log.Fatal(err)
	}

	logger.Traceln("flags:")
	for k, v := range flags {
		logger.Tracef("%s: %v", k, v.Value)
	}

	if err := action(logger, flags); err != nil {
		logger.Fatal(err)
	}
	logger.Traceln("success!!!")
}

func createLogger(levelFlag commando.FlagValue) (*logrus.Logger, error) {
	level, err := levelFlag.GetString()
	if err != nil {
		return nil, err
	}

	lLvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		PadLevelText: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
	logger.SetLevel(lLvl)
	return logger, nil
}
