package controller

import (
	"RinhaBackend/app/database"
	"RinhaBackend/app/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func validateExtratoAndTransactions(c *models.Cliente, t []*models.Transacao) bool {
	if c.Saldo < 0 {
		log.Printf("Error: Invalid Saldo (Cliente ID %d): %d", c.ID, c.Saldo)
		return false
	}

	if c.Limite < 0 {
		log.Printf("Error: Invalid Limite (Cliente ID %d): %d", c.ID, c.Limite)
		return false
	}

	for _, transacao := range t {
		if transacao.Valor < 0 {
			log.Printf("Error: Invalid Valor (Cliente ID %d, Transacao ID %d): %d", c.ID, transacao.ID, transacao.Valor)
			return false
		}

		if transacao.Descricao == "" || (len(transacao.Descricao) < 1 || len(transacao.Descricao) > 10) {
			log.Printf("Error: Invalid Descricao (Cliente ID %d, Transacao ID %d): %s", c.ID, transacao.ID, transacao.Descricao)
			return false
		}

		if transacao.Tipo != "c" && transacao.Tipo != "d" {
			log.Printf("Error: Invalid Tipo (Cliente ID %d, Transacao ID %d): %s", c.ID, transacao.ID, transacao.Tipo)
			return false
		}
	}

	return true
}

func HandleExtrato(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if id < 1 || id > 5 {
		http.Error(w, "Error: User not found", http.StatusNotFound)
		return
	}

	clienteChan := make(chan *models.Cliente, 1)
	transacoesChan := make(chan []*models.Transacao, 1)
	errChan := make(chan error, 1)

	go func() {
		cliente, err := database.GetCliente(ctx, id)
		if err != nil {
			errChan <- err
			return
		}
		clienteChan <- cliente
	}()

	go func() {
		transacoes, err := database.GetLast10Transactions(ctx, id)
		if err != nil {
			errChan <- err
			return
		}
		transacoesChan <- transacoes
	}()

	select {
	case err := <-errChan:
		http.Error(w, "Error: Occured an unknown error: "+err.Error(), http.StatusUnprocessableEntity)
		return
	default:
		cliente := <-clienteChan
		transacoes := <-transacoesChan

		if validateExtratoAndTransactions(cliente, transacoes) {
			response := models.ExtratoResponse{
				Saldo: models.Saldo{
					Total:       cliente.Saldo,
					DataExtrato: time.Now().Format(time.RFC3339),
					Limite:      cliente.Limite,
				},
				UltimasTransacoes: transacoes,
			}

			json.NewEncoder(w).Encode(response)
		}
	}
}
