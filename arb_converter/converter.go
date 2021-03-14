package arb_converter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

func SaveArb(logger *logrus.Logger, dataArb *DataArb,
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
