package models

import "time"

type Transacao struct {
	ID          int
	Valor       int
	Tipo        string
	Descricao   string
	RealizadaEm time.Time
}

func (Transacao) TableName() string {
	return "transacoes"
}
