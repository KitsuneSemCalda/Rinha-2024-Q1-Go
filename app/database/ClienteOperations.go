package database

import (
	"RinhaBackend/app/models"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache = cache.New(5*time.Second, 10*time.Second)

func GetCliente(ctx context.Context, id int) (*models.Cliente, error) {
	cliente, found := getClienteInCache(id)
	if found {
		return cliente, nil
	}

	cliente, err := getClienteInDB(ctx, id)
	if err != nil {
		return nil, err
	}

	Cache.Set(strconv.Itoa(id), cliente, cache.DefaultExpiration)
	return cliente, nil
}

func UpdateCliente(ctx context.Context, id int, updatedCliente models.Cliente) error {
	err := updateClienteInDB(ctx, id, updatedCliente)
	if err != nil {
		return err
	}

	updateClienteInCache(id, updatedCliente)
	return nil
}

func updateClienteInDB(ctx context.Context, id int, updatedCliente models.Cliente) error {
	var cliente models.Cliente

	tx := DB.WithContext(ctx).Begin()

	err := tx.Where("id = ?", id).First(&cliente).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if cliente.Version != updatedCliente.Version {
		tx.Rollback()
		return errors.New("version mismatch")
	}

	updatedCliente.Version++

	err = tx.Model(&cliente).Updates(updatedCliente).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func updateClienteInCache(id int, updatedCliente models.Cliente) {
	Cache.Set(strconv.Itoa(id), updatedCliente, cache.DefaultExpiration)
}

func getClienteInDB(ctx context.Context, id int) (*models.Cliente, error) {
	var cliente models.Cliente

	err := DB.WithContext(ctx).Where("id = ?", id).First(&cliente).Error
	if err != nil {
		return nil, err
	}

	return &cliente, nil
}

func getClienteInCache(id int) (*models.Cliente, bool) {
	cliente, found := Cache.Get(strconv.Itoa(id))
	if !found {
		return nil, false
	}

	return cliente.(*models.Cliente), true
}
