package database

import (
	"RinhaBackend/app/models"
	"errors"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var Cache = cache.New(5*time.Minute, 10*time.Minute)

func GetCliente(id int) (*models.Cliente, error) {
	cliente, found := getClienteInCache(id)

	if found {
		return cliente, nil
	}

	cliente, found = getClienteInDB(id)

	if !found {
		return nil, gorm.ErrRecordNotFound
	}

	Cache.Set(strconv.Itoa(id), cliente, cache.DefaultExpiration)
	return cliente, nil
}

func UpdateCliente(id int, updatedCliente models.Cliente) error {
	if !updateClienteInCache(id, updatedCliente) {
		return errors.New("failed to update cache")
	}

	if !updateClienteInDB(id, updatedCliente) {
		return errors.New("failed to update database")
	}

	return nil
}

// This private function garrants syncronization in cache and database

func updateClienteInDB(id int, updatedCliente models.Cliente) bool {
	var cliente models.Cliente

	tx := DB.Begin()

	err := tx.Where("id = $1", id).First(&cliente).Error

	if err != nil {
		tx.Rollback()
		return false
	}

	if cliente.Version != updatedCliente.Version {
		tx.Rollback()
		return false
	}

	updatedCliente.Version++

	err = tx.Model(&cliente).Updates(updatedCliente).Error

	if err != nil {
		tx.Rollback()
		return false
	}

	tx.Commit()
	return true
}

func updateClienteInCache(id int, updatedCliente models.Cliente) bool {
	_, err := GetCliente(id)

	if err != nil {
		return false
	}

	Cache.Set(strconv.Itoa(id), updatedCliente, cache.DefaultExpiration)
	return true
}

func getClienteInDB(id int) (*models.Cliente, bool) {
	var cliente models.Cliente

	err := DB.Where("id = $1", id).First(&cliente).Error

	if err != nil {
		return nil, false
	}

	return &cliente, true
}

func getClienteInCache(id int) (*models.Cliente, bool) {
	cliente, found := Cache.Get(strconv.Itoa(id))

	if !found {
		return nil, false
	}

	return cliente.(*models.Cliente), true
}
