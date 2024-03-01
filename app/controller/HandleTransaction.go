package controller

import (
	"RinhaBackend/app/database"
	"RinhaBackend/app/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TransactionTask struct {
	ClienteID int
	Request   models.TransacaoRequest
}

var request models.TransacaoRequest
var transactionQueue = make(chan TransactionTask, 100)

func init() {
	go processTransactionQueue()
}

func processTransactionQueue() {
	for task := range transactionQueue {
		database.CreateTransaction(task.ClienteID, &task.Request)
	}
}

func validateRequest(req models.TransacaoRequest, client *models.Cliente) error {
	if len(req.Descricao) < 1 || len(req.Descricao) > 10 {
		log.Printf("Unconsistent transaction description: size equal %d", len(req.Descricao))
		return errors.New("unconsistent transaction description")
	}

	if req.Tipo == "d" {

		value := client.Saldo - req.Valor

		if value < -client.Limite {
			log.Printf("Unconsistent transaction %d is minor than %d", value, client.Limite)
			return errors.New("unconsistent transaction")
		}
		return nil
	}

	if req.Tipo != "d" && req.Tipo != "c" {
		log.Printf("Unconsistent transaction type: %s", req.Tipo)
		return errors.New("unconsistent transaction type")
	}

	return nil
}

func HandleTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Error: Occured an unknown error decoding request: "+err.Error(), http.StatusUnprocessableEntity)
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
		err = validateRequest(request, cliente)
		if err != nil {
			http.Error(w, "Error: Occured an unknown error in validate request: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}

		transactionTask := TransactionTask{
			ClienteID: id,
			Request:   request,
		}
		transactionQueue <- transactionTask

		var newSaldo int
		if request.Tipo == "c" {
			newSaldo = cliente.Saldo + request.Valor
		}

		if request.Tipo == "d" {
			newSaldo = cliente.Saldo - request.Valor
		}

		response := models.TransacaoResponse{
			Limite: cliente.Limite,
			Saldo:  newSaldo,
		}

		json.NewEncoder(w).Encode(response)
	}
}
