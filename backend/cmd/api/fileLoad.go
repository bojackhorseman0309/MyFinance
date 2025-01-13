package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

func (app *application) processIvyWalletTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var buf bytes.Buffer

	file, header, err := r.FormFile("file")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	app.logger.Info("File name %s\n", header.Filename)

	_, err = io.Copy(&buf, file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = parseWalletExport(buf)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	buf.Reset()

	err = app.writeJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func parseWalletExport(export bytes.Buffer) error {
	csvReader := csv.NewReader(&export)
	csvReader.Comma = '|'

	for {
		transaction, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Println(transaction)
	}

	return nil
}
