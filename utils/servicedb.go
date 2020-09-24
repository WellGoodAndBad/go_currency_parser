package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

func InsertData(curdata *[]DataCur, dataParse string, connStr string) {

	//connStr := "user=<your_user> dbname=<database_name> password=<your_password> host=<host_address> sslmode=disable"
	dataParse = strings.Replace(dataParse,"-","_", -1)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Print(err)
	}
	// create schema for partitions
	createSchema := "create schema if not exists currency_partitions;"
	db.Query(createSchema)
	// create main table for inherits partitions
	createTable := "CREATE TABLE if not exists public.currency (currency_code text  primary key,currency_name text, units_per_usd float, usd_per_unit float);"
	db.Query(createTable)
	// create partition about date
	createPart :=fmt.Sprintf("CREATE TABLE IF NOT EXISTS currency_partitions.currency_%s () INHERITS (public.currency);", dataParse)
	db.Query(createPart)
	// my batch insert
	BulkInsert(*curdata, db, dataParse)
	defer db.Close()
}

func BulkInsert(unsavedRows []DataCur, db *sql.DB, dataParse string)  {
	valueStrings := make([]string, 0, len(unsavedRows))
	for _, curData := range unsavedRows {
		valueInsert := fmt.Sprintf("('%s', '%s', %s, %s)", curData.CurrencyCode, strings.Replace(curData.CurrencyName,"'","`", -1), curData.UnitsPerUSD, curData.USDPerUnit)
		valueStrings = append(valueStrings, valueInsert)
	}
	txn, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	trnt := fmt.Sprintf("TRUNCATE currency_partitions.currency_%s", dataParse)
	_, err = txn.Exec(trnt)
	if err != nil {
		txn.Rollback()
		fmt.Println(err)
	}

	stmt := fmt.Sprintf("INSERT INTO currency_partitions.currency_%s (currency_code, currency_name, units_per_usd, usd_per_unit) VALUES %s", dataParse, strings.Join(valueStrings, ","))
	_, err = txn.Exec(stmt)
	if err != nil {
		txn.Rollback()
		fmt.Println(err)
	}
	err = txn.Commit()
	if err != nil {
		fmt.Println(err)
	}


}