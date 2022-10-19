package nutsdb

import (
	"os"

	"github.com/xujiajun/nutsdb"
	"github.com/ziyht/eden_go/ecache/driver"
)

type cfg struct {
	Dir       string
	InMemory  bool
}

type DB struct {
	db      *nutsdb.DB
	opts    *nutsdb.Options
}

func newDB(cfg *cfg) (*DB, error){

	opts := nutsdb.DefaultOptions
	opts.RWMode      = nutsdb.MMap
	opts.Dir         = cfg.Dir
	opts.SyncEnable  = false
	
	db, err := nutsdb.Open(opts)
	if err != nil {
		return nil, err
	}

	out := &DB{db: db}
	out.opts    = &opts

	return out, nil
}

func (db *DB)TX(tx interface{}) driver.TX {
	return &TX{txn: tx.(*nutsdb.Tx)}
}

func (db *DB)Update(fn func(tx driver.TX) error) error {
	return db.db.Update(func(txn *nutsdb.Tx)error{
		return fn(db.TX(txn))
	})
}

func (db *DB)View(fn func(tx driver.TX) error) error {
	return db.db.View(func(txn *nutsdb.Tx)error{
		return fn(db.TX(txn))
	})
}

func(db *DB)DropPrefix(prefix []byte) (error) {
	return db.db.Update(func(tx *nutsdb.Tx)error{
		return tx.DeleteBucket(nutsdb.DataStructureBPTree, string(prefix))
	})
}

func(db *DB)Truncate() (error) {
	err := db.db.Close()
	if err != nil {
		return err
	}
	err = os.RemoveAll(db.opts.Dir)
	if err != nil {
		return err
	}
	
	db.db, err = nutsdb.Open(*db.opts)

	return err
}

func (db *DB)Close() error{
	return db.db.Close()
}