package homekit

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"strconv"
)

type Store struct {
	db *buntdb.DB
}

func NewStore() Store {
	db, _ := buntdb.Open("db")
	return Store{db}
}

func (s Store) SetId(id uint64, itemName string) {
	s.db.Update(func(tx *buntdb.Tx) (err error) {
		tx.Set(itemName, fmt.Sprintf(`%d`, id), nil)
		return
	})
}

func (s Store) GetId(itemName string) (ok bool, id uint64) {
	s.db.View(func(tx *buntdb.Tx) (err error) {
		res, err := tx.Get(itemName)
		ok = err == nil
		id, _ = strconv.ParseUint(res, 10, 64)
		return
	})
	return
}
