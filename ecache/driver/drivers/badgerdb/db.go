package badgerdb

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"github.com/ziyht/eden_go/ecache/driver"
)

type cfg struct {
	Dir       string
	InMemory  bool
}

type DB struct {
	db      *badger.DB
	opts    *badger.Options
}

func newDB(cfg *cfg) (*DB, error){

	opts := badger.DefaultOptions(cfg.Dir)
	opts = opts.WithInMemory(cfg.InMemory)
	opts = opts.WithLoggingLevel(badger.WARNING)
	opts = opts.WithCompression(options.Snappy)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	out := &DB{db: db}
	out.opts    = &opts

	return out, nil
}

func (db *DB)TX(tx interface{}) driver.TX {
	return &TX{txn: tx.(*badger.Txn)}
}

func (db *DB)Update(fn func(tx driver.TX) error) error {
	return db.db.Update(func(txn *badger.Txn)error{
		return fn(db.TX(txn))
	})
}

func (db *DB)View(fn func(tx driver.TX) error) error {
	return db.db.View(func(txn *badger.Txn)error{
		return fn(db.TX(txn))
	})
}

func(db *DB)DropPrefix(prefix []byte) (error) {
	return db.db.DropPrefix(prefix)
}

func(db *DB)Truncate() (error) {
	return db.db.DropAll()
}

func (db *DB)Close() error{
	return db.db.Close()
}