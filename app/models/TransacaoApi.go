package models

type TransacaoRequest struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacaoResponse struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}
