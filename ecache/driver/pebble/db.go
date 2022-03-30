package pebble

import (
	"github.com/cockroachdb/pebble"
)

type cfg struct {
	Dir  *string
}

type DB struct {
	db *pebble.DB
}

func newDB(cfg *cfg) (*DB, error){

	opt := pebble.Options{}

	db, err := pebble.Open(*cfg.Dir, &opt)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}


func(db *DB)Set(k []byte, v []byte) error {

	return db.db.Set(k, v, nil)

}
