package database

import (
	"RinhaBackend/app/models"
	"context"
	"log"
	"sync"
	"time"
)

var mu sync.Mutex

func validateExtratoAndTransactions(cliente *models.Cliente, transacoes []*models.Transacao) bool {
	if cliente.Saldo < 0 {
		log.Printf("Error: Invalid Saldo (Cliente ID %d): %d", cliente.ID, cliente.Saldo)
		return false
	}

	if cliente.Limite < 0 {
		log.Printf("Error: Invalid Limite (Cliente ID %d): %d", cliente.ID, cliente.Limite)
		return false
	}

	for _, transacao := range transacoes {
		if transacao.Valor < 0 {
			log.Printf("Error: Invalid Valor (Cliente ID %d, Transacao ID %d): %d", cliente.ID, transacao.ID, transacao.Valor)
			return false
		}

		if len(transacao.Descricao) < 1 || len(transacao.Descricao) > 10 {
			log.Printf("Error: Invalid Descricao (Cliente ID %d, Transacao ID %d): %s", cliente.ID, transacao.ID, transacao.Descricao)
			return false
		}

		if transacao.Tipo != "c" && transacao.Tipo != "d" {
			log.Printf("Error: Invalid Tipo (Cliente ID %d, Transacao ID %d): %s", cliente.ID, transacao.ID, transacao.Tipo)
			return false
		}

		// Additional check: Ensure the transaction value does not exceed the client's limit
		if transacao.Tipo == "d" && transacao.Valor > cliente.Limite {
			log.Printf("Error: Transaction value exceeds limit (Cliente ID %d, Transacao ID %d): %d", cliente.ID, transacao.ID, transacao.Valor)
			return false
		}
	}

	return true
}

func CreateTransaction(clienteID int, t *models.TransacaoRequest) {
	var cliente models.Cliente
	if err := DB.First(&cliente, clienteID).Error; err != nil {
		log.Printf("Error: Cliente with ID %d not found: %v", clienteID, err)
		return
	}

	if !validateExtratoAndTransactions(&cliente, []*models.Transacao{{Valor: t.Valor, Tipo: t.Tipo, Descricao: t.Descricao}}) {
		log.Println("Error: Invalid transaction data")
		return
	}

	var transaction models.Transacao
	transaction.ClienteID = clienteID
	transaction.Valor = t.Valor
	transaction.Tipo = t.Tipo
	transaction.Descricao = t.Descricao
	transaction.RealizadaEm = time.Now()

	go func() {
		tx := DB.Begin()
		if err := tx.Create(&transaction).Error; err != nil {
			log.Printf("Error: Failed to create transaction: %v", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
}

func GetLast10Transactions(ctx context.Context, clienteID int) ([]*models.Transacao, error) {
	var transactions []*models.Transacao

	result := DB.Where("cliente_id = ?", clienteID).Order("realizada_em desc").Limit(10).Find(&transactions)

	if result.Error != nil {
		log.Printf("Error: Failed to fetch transactions: %v", result.Error)
	}

	return transactions, nil
}
