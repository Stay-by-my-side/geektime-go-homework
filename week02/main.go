package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type tag struct {
	Id    int64  `db:"id"`
	Name  string `db:"name"`
	State string `db:"state"`
}

type product struct {
	Id          int64   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
}

func main() {
	db, err := sqlx.Connect("mysql", "root:111111@tcp(127.0.0.1:3306)/blog_service")
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return
	}
	defer db.Close()

	if err := queryProductRow(db, 1); err != nil {
		fmt.Printf("main.queryRow FATAL: \n%+v\n", err)
	}
	if err := queryRow(db, 1); err != nil {
		fmt.Printf("main.queryRow FATAL: \n%+v\n", err)
	}
}

func queryRow(db *sqlx.DB, id int64) error {
	sqlStr := "SELECT id, name, state FROM blog_tag WHERE id = ?"

	var t tag
	if err := db.Get(&t, sqlStr, id); err != nil {
		return errors.Wrapf(err, "sql: %s", sqlStr)
	}
	fmt.Printf("id:%d, name:%s, age:%s\n", t.Id, t.Name, t.State)
	return nil
}

func queryProductRow(db *sqlx.DB, id int64) error {
	sqlStr := "SELECT id, name, description, price FROM product WHERE id = ?"

	var p product
	if err := db.Get(&p, sqlStr, id); err != nil {
		return errors.Wrapf(err, "sql: %s", sqlStr)
	}
	fmt.Printf("id:%d, name:%s, description:%s, price:%g\n", p.Id, p.Name, p.Description, p.Price)
	return nil
}
