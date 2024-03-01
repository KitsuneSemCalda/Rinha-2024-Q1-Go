package models

type Saldo struct {
	Total       int    `json:"total"`
	DataExtrato string `json:"data_extrato"`
	Limite      int    `json:"limite"`
}
