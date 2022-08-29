package pebble

import (
	"github.com/cockroachdb/pebble"
)

type cfg struct {
	Dir       string
	InMemory  bool
}

type DB struct {
	db *pebble.DB
}

func newDB(cfg *cfg) (*DB, error){

	opt := pebble.Options{}

	db, err := pebble.Open(cfg.Dir, &opt)
	if err != nil {
		return nil, err
	}

	// b := db.NewBatch()
	
	// db.Set()
	// db.DeleteRange()
	// db.Get()
	

	return &DB{db: db}, nil
}


func(db *DB)Set(k []byte, v []byte) error {

	return db.db.Set(k, v, nil)

}
