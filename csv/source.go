package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func csvFromFile(logger *logrus.Logger, filePath string) (*csv.Reader, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			logger.Warningf("close csv file error: %v", err)
		}
	}()

	return csvFromReader(csvFile)
}

func csvFromWeb(logger *logrus.Logger, url string) (*csv.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Warningf("close response body error: %v", err)
		}
	}()

	return csvFromReader(resp.Body)
}

func csvFromReader(src io.Reader) (*csv.Reader, error) {
	var csvRaw bytes.Buffer
	if _, err := io.Copy(&csvRaw, src); err != nil {
		return nil, err
	}
	return csv.NewReader(bytes.NewReader(csvRaw.Bytes())), nil
}
