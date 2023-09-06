package database

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	db, err := Connect()
	if err != nil {{
		t.Error("Failed to connect to db")
	}}

	got, err := Search(db, "Product", "prodName", "Nurofen")
	if err != nil {
		t.Errorf("%v", err)
	}
	db.Close()

	want := []Product{}

	fmt.Printf("%v\n", reflect.TypeOf(got))
	fmt.Printf("%v\n", reflect.TypeOf(want))

	if fmt.Sprintf("%v", reflect.TypeOf(got)) != fmt.Sprintf("%v", reflect.TypeOf(want)) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
