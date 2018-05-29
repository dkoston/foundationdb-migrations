package main

import (
    "github.com/apple/foundationdb/bindings/go/src/fdb"
)


func Up_20130106222315(t fdb.Transactor) error {
    _, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
        tr.Set(fdb.Key("hello"), []byte("world"))
        return
    })

    return err
}

func Down_20130106222315(t fdb.Transactor) error {
    _, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
        tr.Clear(fdb.Key("hello"))
        return
    })

    return err
}
