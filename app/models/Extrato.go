package models

type ExtratoResponse struct {
	Saldo             Saldo       `json:"saldo"`
	UltimasTransacoes []*Transacao `json:"ultimas_transacoes"`
}
