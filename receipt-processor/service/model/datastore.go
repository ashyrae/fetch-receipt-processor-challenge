package model

import (
	"sync"
)

type ReceiptDB struct {
	Store map[string]*Receipt
	sync.RWMutex
}

func (db *ReceiptDB) Set(r *Receipt) (id string, err error) {
	if id, err := idFactory(); err != nil {
		return "", ErrInternalServer(err.Error())
	} else {
		db.Lock()
		defer db.Unlock()
		db.Store[id] = r
		return id, nil
	}
}

func (db *ReceiptDB) Get(id string) (receipt *Receipt, err error) {
	db.RLock()
	defer db.RUnlock()
	if r, exists := db.Store[id]; !exists {
		return &Receipt{}, ErrNotFound("Receipt was not found: " + id)
	} else {
		return r, nil
	}
}
