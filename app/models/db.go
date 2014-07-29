package models

import (
    "github.com/revel/revel"
    _ "github.com/mattn/go-sqlite3"
    "github.com/coopernurse/gorp"
    "github.com/revel/revel/modules/db/app"
)

var (
    DbMap *gorp.DbMap // このデータベースマッパーからSQLを流す
)

type Transactional struct {
    *revel.Controller
    Txn *gorp.Transaction
}

func init() {
    revel.OnAppStart(InitDB)
    revel.InterceptMethod((*Transactional).Begin, revel.BEFORE)
    revel.InterceptMethod((*Transactional).Commit, revel.AFTER)
    revel.InterceptMethod((*Transactional).Rollback, revel.PANIC)
}

func InitDB() {
    db.Init()
    DbMap = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

    // ここで好きにテーブルを定義する
    t := DbMap.AddTable(User{}).SetKeys(true, "Id")
    t.ColMap("Name").MaxSize = 20

    DbMap.CreateTables()
}

func (c *Transactional) Begin() revel.Result {
    c.Txn, _ = DbMap.Begin()
    return nil
}

func (c *Transactional) Commit() revel.Result {
    _ = c.Txn.Commit()
    c.Txn = nil
    return nil
}

func (c *Transactional) Rollback() revel.Result {
    _ = c.Txn.Rollback()
    c.Txn = nil
    return nil
}

