package arb_converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	arbExt     = ".arb"
	localeAttr = "@@locale"
	metaPrefix = "@"
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

	cultures := make(map[string]string)
	arbItems := make(map[string]*ItemArb)

	for _, file := range files {
		if file.IsDir() || strings.ToLower(path.Ext(file.Name())) != arbExt {
			logger.Tracef("skip %s", file.Name())
			continue
		}

		logger.Tracef("process file %s", file.Name())
		rawData, err := ioutil.ReadFile(path.Join(arbFolderPath, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("file read error [%s]: %w", file.Name(), ErrArbFile)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, fmt.Errorf("file unmarshal error [%s]: %w", file.Name(), ErrArbFile)
		}
		culture := getStrByKey(localeAttr, data)
		if culture == "" {
			culture = getCultureFromFileName(file.Name())
		}
		if culture == "" {
			return nil, fmt.Errorf("can not detect culture for [%s]: %w", file.Name(), ErrArbFile)
		}
		culture = strings.ToLower(culture)

		if f, ok := cultures[culture]; ok {
			return nil, fmt.Errorf("same cultures in [%s] and [%s]: %w", f, file.Name(), ErrArbFile)
		}
		cultures[culture] = file.Name()
		processCulture(culture, culture == strings.ToLower(defaultCulture), data, arbItems)
	}

	dataArb := &DataArb{
		Cultures: nil,
		Items:    arbItems,
	}
	for c := range cultures {
		dataArb.Cultures = append(dataArb.Cultures, c)
	}
	return dataArb, err
}

func processCulture(culture string,
	isDefaultCulture bool,
	data map[string]interface{},
	arbItems map[string]*ItemArb) {
	for k := range data {
		if strings.HasPrefix(k, metaPrefix) {
			continue
		}
		item, ok := arbItems[k]
		if !ok {
			item = &ItemArb{
				Description: "",
				Cultures:    make(map[string]string),
			}
			arbItems[k] = item
		}
		item.Cultures[culture] = getStrByKey(k, data)

		if !isDefaultCulture {
			continue
		}
		setMeta(k, item, data)
	}
}

func setMeta(itemName string, item *ItemArb, data map[string]interface{}) {
	meta := getMapByKey(metaPrefix+itemName, data)
	item.Description = getStrByKey("description", meta)
	placeholders := getMapByKey("placeholders", meta)
	if len(placeholders) > 0 {
		item.Parameters = make(map[string]struct{})
		for k := range placeholders {
			item.Parameters[k] = struct{}{}
		}
	}
}

func getCultureFromFileName(name string) string {
	nameWithoutExt := name[:len(name)-len(filepath.Ext(name))]
	parts := strings.Split(nameWithoutExt, "_")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func getStrByKey(k string, m map[string]interface{}) string {
	v, ok := m[k]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func getMapByKey(k string, m map[string]interface{}) map[string]interface{} {
	v, ok := m[k]
	if !ok {
		return nil
	}
	res, _ := v.(map[string]interface{})
	return res
}
