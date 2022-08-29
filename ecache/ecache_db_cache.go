package ecache

type DBCache struct {
	db           *db
	dfRegion     *Region
	dfItemRegion *ItemRegion[Item]
}

func newDBCache(db *db)(*DBCache){
	return &DBCache{db: db}
}

func (c *DBCache)DfRegion()(*Region){
	if c.dfRegion == nil {
		c.dfRegion = newRegion(c.db, nil)
	}
	return c.dfRegion
}

func (c *DBCache)NewRegion(keys ...string)(*Region){
	if keys == nil {
		return c.DfRegion()
	}
	return newRegion(c.db, keys)
}

func (c *DBCache)DfItemRegion()(*ItemRegion[Item]){
	if c.dfRegion == nil {
		c.dfItemRegion = newItemRegion[Item](c.db, nil)
	}
	return c.dfItemRegion
}

func (c *DBCache)NewItemRegion(keys ...string)(*ItemRegion[Item]){
	if keys == nil {
		return c.DfItemRegion()
	}
	return newItemRegion[Item](c.db, keys)
}

func (c *DBCache)Truncate() error {
	return c.db.truncate()
}

func (c *DBCache)Close() error {
	return c.db.close()
}