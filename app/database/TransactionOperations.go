package database

import (
	"RinhaBackend/app/models"
	"time"
)

func CreateTransaction(id int, t *models.TransacaoRequest) {
	var transaction models.Transacao

	transaction.ID = id
	transaction.Valor = t.Valor
	transaction.Tipo = t.Tipo
	transaction.Descricao = t.Descricao
	transaction.RealizadaEm = time.Now()

	DB.Create(transaction)
}

func GetLast10Transactions(clienteID int) []models.Transacao {
	var transactions []models.Transacao
	DB.Where("cliente_id = ?", clienteID).Order("realizada_em desc").Limit(10).Find(&transactions)
	return transactions
}
