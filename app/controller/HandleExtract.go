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
		http.Error(w, "Error: Occured an unknown error in get client: "+err.Error(),
			http.StatusNotFound)
		return
	case cliente := <-clienteChan:
		transacoes, err := database.GetLast10Transactions(id)

		if err != nil {
			log.Printf("Occured an Unknown error in GetLast10Transactions: %s",
				err.Error())
		}

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
