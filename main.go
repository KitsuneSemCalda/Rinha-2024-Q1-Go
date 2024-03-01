package main

import (
	"RinhaBackend/app/controller"
	"RinhaBackend/app/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	port := utils.GetPort()

	r := mux.NewRouter()

	r.HandleFunc("/clientes/{id}/transacoes", controller.HandleTransaction).Methods("POST")
	r.HandleFunc("/clientes/{id}/extrato", controller.HandleExtrato).Methods("GET")

	http.ListenAndServe(port, r)
}
