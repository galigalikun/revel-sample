package controllers

import (
    "github.com/revel/revel"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/coopernurse/gorp"
    "myapp/app/models" // revel new APP_NAME の APP_NAME
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
    driver, _ := revel.Config.String("db.driver")
    spec, _ := revel.Config.String("db.spec")
    db, err := sql.Open(driver, spec) // "sqlite3", "./app.db")
    if err != nil {
        panic(err.Error())
    }
    DbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

    // ここで好きにテーブルを定義する
    t := DbMap.AddTable(models.User{}).SetKeys(true, "Id")
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

