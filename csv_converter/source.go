package csv_converter

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"os"
)

func CsvFromFile(filePath string) (*csv.Reader, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	return csvFromReader(csvFile)
}

func CsvFromWeb(url string) (*csv.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return csvFromReader(resp.Body)
}

func csvFromReader(src io.Reader) (*csv.Reader, error) {
	var csvRaw bytes.Buffer
	if _, err := io.Copy(&csvRaw, src); err != nil {
		return nil, err
	}
	return csv.NewReader(bytes.NewReader(csvRaw.Bytes())), nil
}
