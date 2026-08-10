package repository

import (
	"github.com/ostafen/clover/v2"
)

type CloverModuleRepository struct {
	db *clover.DB
}

func NewCloverModuleRepository(db *clover.DB) *CloverModuleRepository {
	return &CloverModuleRepository{db: db}
}
