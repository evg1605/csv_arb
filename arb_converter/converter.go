package arb_converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	arbExt = ".arb"
)

var (
	ErrArbFile = errors.New("invalid arb file error")
)

func SaveArb(logger *logrus.Logger,
	dataArb *DataArb,
	arbFolderPath,
	arbFileTemplate,
	defaultCulture string) error {
	if err := os.RemoveAll(arbFolderPath); err != nil {
		return err
	}

	if err := os.MkdirAll(arbFolderPath, 0777); err != nil {
		return err
	}

	for _, cn := range dataArb.Cultures {
		m := make(map[string]interface{})

		for name, item := range dataArb.Items {
			m[name] = item.Cultures[cn]
			if cn != defaultCulture {
				continue
			}

			meta := map[string]interface{}{
				"description": item.Description,
			}
			if len(item.Parameters) > 0 {
				placeholders := make(map[string]interface{})
				for pn, _ := range item.Parameters {
					placeholders[pn] = map[string]interface{}{
						"type": "dynamic",
					}
				}
				meta["placeholders"] = placeholders
			}
			m[fmt.Sprintf("@%s", name)] = meta
		}

		buf, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(path.Join(arbFolderPath, strings.ReplaceAll(arbFileTemplate, "{culture}", cn)), buf, 0666); err != nil {
			return err
		}
	}

	return nil
}

func LoadArb(logger *logrus.Logger,
	arbFolderPath,
	defaultCulture string) (*DataArb, error) {

	files, err := ioutil.ReadDir(arbFolderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || strings.ToLower(path.Ext(file.Name())) != arbExt {
			logger.Tracef("skip %s", file.Name())
			continue
		}

		logger.Tracef("process file %s", file.Name())
		fullName := path.Join(arbFolderPath, file.Name())
		rawData, err := ioutil.ReadFile(fullName)
		if err != nil {
			return nil, fmt.Errorf("file read error [%s]: %w", fullName, ErrArbFile)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, fmt.Errorf("file unmarshal error [%s]: %w", fullName, ErrArbFile)
		}
		culture := getCultureFromData(data)
		if culture == "" {
			culture = getCultureFromName(file.Name())
		}
		if culture == "" {
			return nil, fmt.Errorf("can not detect culture for [%s]: %w", fullName, ErrArbFile)
		}
	}

	return nil, err
}

func getCultureFromName(name string) string {
	return ""
}

func getCultureFromData(data map[string]interface{}) string {
	return ""
}
