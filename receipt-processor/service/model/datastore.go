package model

import (
	"sync"
)

type ReceiptDB struct {
	Store map[string]*Receipt
	sync.RWMutex
}

func (db *ReceiptDB) Create(r *Receipt) (id string, err error) {
	// yes, generate an id & proceed
	if id, err := idFactory(); err != nil {
		return "", ErrInternalServer(err.Error())
	} else {
		db.Lock()
		defer db.Unlock()
		db.Store[id] = r
		return id, nil
	}
}

func (db *ReceiptDB) Set(idSet string, r *Receipt) (id string, err error) {
	if idSet == "" {
		return "", ErrInternalServer("No Receipt ID was provided to Set")
	} else {
		id = idSet
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
