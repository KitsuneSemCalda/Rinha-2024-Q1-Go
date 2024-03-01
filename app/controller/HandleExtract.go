package controller

import (
	"RinhaBackend/app/database"
	"RinhaBackend/app/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func ValidateExtratoResponse(response models.ExtratoResponse) error {
	_, err := time.Parse(time.RFC3339, response.Saldo.DataExtrato)
	if err != nil {
		return errors.New("invalid date format for DataExtrato")
	}

	if response.Saldo.Limite < 0 {
		return errors.New("limit should not be negative")
	}

	for _, transacao := range response.UltimasTransacoes {
		if transacao.Valor < 0 {
			return errors.New("transaction value should not be negative")
		}

		if transacao.Tipo != "c" && transacao.Tipo != "d" {
			return errors.New("invalid transaction type")
		}

		_, err := time.Parse(time.RFC3339, transacao.RealizadaEm.String())
		if err != nil {
			return errors.New("invalid date format for RealizadaEm")
		}
	}

	return nil
}

func HandleExtrato(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Error: Occured an unknown error in convert id: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	clienteChan := make(chan *models.Cliente)
	errChan := make(chan error)

	go func() {
		cliente, err := database.GetCliente(id)
		if err != nil {
			errChan <- err
			return
		}
		clienteChan <- cliente
	}()

	select {
	case err := <-errChan:
		http.Error(w, "Error: Occured an unknown error in get client: "+err.Error(), http.StatusNotFound)
		return
	case cliente := <-clienteChan:
		transacoes := database.GetLast10Transactions(cliente.ID)

		response := models.ExtratoResponse{
			Saldo: models.Saldo{
				Total:       cliente.Saldo,
				DataExtrato: time.Now().Format(time.RFC3339),
				Limite:      cliente.Limite,
			},
			UltimasTransacoes: transacoes,
		}

		err = ValidateExtratoResponse(response)

		if err != nil {
			log.Printf("Error: Occured an unknown error in validate response: %s", err.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		json.NewEncoder(w).Encode(response)
	}
}
