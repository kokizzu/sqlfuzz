package main

import (
	"log"
	"time"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/connector"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	"github.com/PumpkinSeed/sqlfuzz/pkg/fuzzer"
	"github.com/brianvoe/gofakeit/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	f := flags.Get()
	gofakeit.Seed(0)
	driver := drivers.New(f.Driver)
	db := connector.Connection(driver)
	defer db.Close()

	var tables []string
	if f.Table == "" {
		var err error
		tables, err = driver.ShowTables(db)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tables = []string{f.Table}
	}
	for _, table := range tables {
		f.Table = table
		fields, err := driver.DescribeFields(f.Table, db)
		if err != nil {
			log.Fatal(err.Error())
		}
		t := time.Now()
		if err := fuzzer.Run(db, fields, f); err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("Fuzzing %s table taken: %v \n", table, time.Since(t))
	}
}
