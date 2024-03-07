package database

import (
	"RinhaBackend/app/models"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var Cache = cache.New(5*time.Minute, 10*time.Minute)

func GetCliente(ctx context.Context, id int) (*models.Cliente, error) {
	clienteChan := make(chan *models.Cliente, 1)
	errChan := make(chan error, 1)

	go func() {
		cliente, found := getClienteInCache(id)
		if found {
			clienteChan <- cliente
			return
		}

		cliente, found = getClienteInDB(ctx, id)
		if !found {
			errChan <- gorm.ErrRecordNotFound
			return
		}

		Cache.Set(strconv.Itoa(id), cliente, cache.DefaultExpiration)
		clienteChan <- cliente
	}()

	select {
	case err := <-errChan:
		return nil, err
	case cliente := <-clienteChan:
		return cliente, nil
	}
}

func UpdateCliente(ctx context.Context, id int, updatedCliente models.Cliente) error {
	errChan := make(chan error, 1)

	go func() {
		if !updateClienteInCache(id, updatedCliente) {
			errChan <- errors.New("failed to update cache")
			return
		}

		if !updateClienteInDB(ctx, id, updatedCliente) {
			errChan <- errors.New("failed to update database")
			return
		}

		errChan <- nil
	}()

	return <-errChan
}

// This private function garrants syncronization in cache and database

func updateClienteInDB(ctx context.Context, id int, updatedCliente models.Cliente) bool {
	var cliente models.Cliente

	tx := DB.WithContext(ctx).Begin()

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
	_, err := GetCliente(context.Background(), id)

	if err != nil {
		return false
	}

	Cache.Set(strconv.Itoa(id), updatedCliente, cache.DefaultExpiration)
	return true
}

func getClienteInDB(ctx context.Context, id int) (*models.Cliente, bool) {
	var cliente models.Cliente

	err := DB.WithContext(ctx).Where("id = $1", id).First(&cliente).Error

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
