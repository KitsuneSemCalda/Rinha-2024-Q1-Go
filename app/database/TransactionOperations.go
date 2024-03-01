package database

import (
	"RinhaBackend/app/models"
	"log"
	"time"
)

func CreateTransaction(clienteID int, t *models.TransacaoRequest) {
	var cliente models.Cliente
	if err := DB.First(&cliente, clienteID).Error; err != nil {
		log.Printf("Error: Cliente with ID %d not found: %v", clienteID, err)
		return
	}

	var transaction models.Transacao
	transaction.ClienteID = clienteID
	transaction.Valor = t.Valor
	transaction.Tipo = t.Tipo
	transaction.Descricao = t.Descricao
	transaction.RealizadaEm = time.Now()

	tx := DB.Begin()
	if err := tx.Create(&transaction).Error; err != nil {
		log.Printf("Error: Failed to create transaction: %v", err)
		tx.Rollback()
		return
	}
	tx.Commit()
}

func GetLast10Transactions(clienteID int) ([]*models.Transacao, error) {
	var transactions []*models.Transacao

	sqlDB, err := DB.DB()

	if err != nil {
		log.Printf("Can't open sql in gorm: %s", err.Error())
		return nil, err
	}

	query := `
	select id, cliente_id, valor, tipo, descricao, realizada_em
	from transacoes
	where cliente_id = $1
	order by realizada_em desc
	limit 10
	`

	rows, err := sqlDB.Query(query, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tr := new(models.Transacao)

		err := rows.Scan(&tr.ID, &tr.ClienteID, &tr.Valor, &tr.Tipo, &tr.Descricao, &tr.RealizadaEm)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tr)
	}

	return transactions, nil
}
