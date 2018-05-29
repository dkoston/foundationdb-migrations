
package main

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

// Up is executed when this migration is applied
func Up_20180529152851(t fdb.Transactor) error {
	_, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
        tr.Set(fdb.Key("hello2"), []byte("world2"))
        return
    })

    return err
}

// Down is executed when this migration is rolled back
func Down_20180529152851(t fdb.Transactor) error {
	_, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
        tr.Clear(fdb.Key("hello2"))
        return
    })

    return err
}
