package models

type Cliente struct {
	ID      int
	Saldo   int
	Limite  int
	Version int
}

func (Cliente) TableName() string {
	return "clientes"
}
