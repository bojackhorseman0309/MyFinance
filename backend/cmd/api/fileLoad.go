package main

import (
	"backend/internal/data"
	"backend/internal/helpers"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
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

	app.logger.Info(fmt.Sprintf("Loading file name %s", header.Filename))

	_, err = io.Copy(&buf, file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	latestTransaction, err := app.models.WalletTransactions.GetLatestTransaction()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.logger.Warn("No latest transaction could be found")
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.logger.Info(fmt.Sprintf("Latest transaction: %v", latestTransaction))

	filterTransactionsAfter := time.Time{}
	if latestTransaction != nil {
		filterTransactionsAfter = latestTransaction.CreatedAt
	}

	transactions, err := parseWalletExport(buf, filterTransactionsAfter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, transaction := range transactions {
		err = app.models.WalletTransactions.Insert(&transaction)
		if err != nil {
			app.logger.Warn(fmt.Sprintf("Error inserting transaction: %s", err))
		}
	}

	buf.Reset()

	err = app.writeJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func parseWalletExport(export bytes.Buffer, latestDate time.Time) ([]data.WalletTransaction, error) {
	var transactionModelList []data.WalletTransaction
	csvReader := csv.NewReader(&export)
	csvReader.Comma = '|'

	for {
		transaction, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		createdAtParsed, err := helpers.ParseDateTime(transaction[0])
		if err != nil {
			return nil, err
		}

		if !latestDate.IsZero() &&
			(createdAtParsed.Before(latestDate) || createdAtParsed.Equal(latestDate)) {
			continue
		}

		dueDateParsed := time.Time{}
		if transaction[13] != "" {
			dueDateParsed, err = helpers.ParseDateTime(transaction[13])
			if err != nil {
				return nil, err
			}
		}

		parsedAmount := 0.0
		if transaction[4] != "" {
			parsedAmount, err = helpers.ConvertStringFloatToFloat(transaction[4])
			if err != nil {
				return nil, err
			}
		}

		parsedTransferAmount := 0.0
		if transaction[7] != "" {
			parsedTransferAmount, err = helpers.ConvertStringFloatToFloat(transaction[7])
			if err != nil {
				return nil, err
			}
		}

		parsedReceiveAmount := 0.0
		if transaction[10] != "" {
			parsedReceiveAmount, err = helpers.ConvertStringFloatToFloat(transaction[10])
			if err != nil {
				return nil, err
			}
		}

		transactionModel := data.WalletTransaction{
			CreatedAt:        createdAtParsed,
			Title:            transaction[1],
			Category:         transaction[2],
			Account:          transaction[3],
			Amount:           parsedAmount,
			Currency:         transaction[5],
			TransactionType:  transaction[6],
			TransferAmount:   parsedTransferAmount,
			TransferCurrency: transaction[8],
			ToAccount:        transaction[9],
			ReceiveAmount:    parsedReceiveAmount,
			ReceiveCurrency:  transaction[11],
			Description:      transaction[12],
			DueDate:          dueDateParsed,
			WalletId:         transaction[14],
		}

		transactionModelList = append(transactionModelList, transactionModel)
	}

	return transactionModelList, nil
}
