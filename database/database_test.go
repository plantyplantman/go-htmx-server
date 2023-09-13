package database

import (
	"fmt"
	// "reflect"
	"testing"
)

func TestGetStoreId(t *testing.T) {
	db, e := Connect()
	if e != nil {
		panic(e)
	}

	id,e := GetStoreId(db, PETRIE)
	if e != nil {
		panic(e)
	}
	if id != 1 {
		t.Fatalf("\nPETRIE ID OF 1 != %v", id)
	}
}

func TestSearch(t *testing.T) {
	db, err := Connect()
	if err != nil {{
		t.Error("Failed to connect to db")
	}}

	got, err := SearchProductNames(db, `A%`)
	if err != nil {
		t.Errorf("%v", err)
	}
	db.Close()

	if len(got) < 1 {
		t.Fatal("NO PRODUCTS FOUND")
	}

	fmt.Printf("\n%+v", got)
}
