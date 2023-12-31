package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
)

const (
	PETRIE   string = "PETRIE"
	BUNDA           = "BUNDA"
	FRANKLIN        = "FRANKLIN"
	CON             = "CON"
)

type Store struct {
	Name string
	ID int
}

type Product struct {
	Id int
	Sku        uint64
	ProdName   string
	Soh        int
	Price      float64
	PromoPrice float64
}

type StoreStock struct {
	ID int
	StoreID int
	ProductID int
	Soh int
	UnitCost float64
	LastOrdered sql.NullTime
}

func Connect() (*sql.DB, error) {
	var dbUrl = "libsql://equal-gargoyle-plantyplantman.turso.io?authToken=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2OTM4ODIwNDQsImlkIjoiOTAxODk2MGEtNDk5Yi0xMWVlLWE1YTQtM2UxNDk0NmVhNTY0In0.0fUZjHUY-hvj1VFwrqzsN5cD2m6zyqu6c7Q-NSiGc_B45-aZMyb4oiscg32xE2oaCrMPQr8DK83bULfE4AZMAA"
	return sql.Open("libsql", dbUrl)
}

func Seed(db *sql.DB) (sql.Result, error) {
	sql := `
-- SQLite3 Database seed file

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Creating the 'Product' table
CREATE TABLE Product (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sku INTEGER UNIQUE NOT NULL,
    prodName TEXT NOT NULL,
    price REAL NOT NULL DEFAULT 0,
		promoPrice REAL NOT NULL DEFAULT 0
);

-- Creating the 'Store' table
CREATE TABLE Store (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Inserting predefined stores
INSERT INTO Store (name) VALUES ('PETRIE'), ('BUNDA'), ('FRANKLIN'), ('CON');

-- Creating the 'StoreStock' table to map Products to Stores and their stock on hand
CREATE TABLE StoreStock (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    store_id INTEGER REFERENCES Store(id),
    product_id INTEGER REFERENCES Product(id),
    soh INTEGER NOT NULL DEFAULT 0,
    unitCost REAL NOT NULL DEFAULT 0,
    lastOrdered DATETIME,
    UNIQUE(store_id, product_id) -- Ensuring that there's only one row for each store-product combination
);
`
	return db.Exec(sql)
}

func GetStoreId(db *sql.DB, store string) (int64, error) {
	// fmt.Printf("getting store id for %v", store)

	q := `SELECT id FROM Store WHERE name = ?;`
	rows, err := db.Query(q, store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query db: %s", err)
		return 0, err
	}
	defer rows.Close()
	var id int64
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
			return 0, err
		}
	}
	return id, nil
}

func UpsertProduct(db *sql.DB, product Product) (sql.Result, error) {
	q := `
INSERT INTO Product (sku, prodName, price, promoPrice)
VALUES (?, ?, ?, ?)
ON CONFLICT(sku) DO UPDATE SET
	prodName = excluded.prodName,
	price = excluded.price,
	promoPrice = excluded.promoPrice
`

	if product.ProdName == "" || product.Sku == 0 {
		return nil, fmt.Errorf("missing required fields")
	}

	return db.Exec(q, product.Sku, product.ProdName, product.Price, product.PromoPrice)
}

func GetAllProducts(db *sql.DB, page int, pageSize int) ([]Product, error) {
	q := `SELECT sku, prodName, price, promoPrice FROM Product LIMIT ? OFFSET ?`

  offset := (page - 1) * pageSize
  rows, err := db.Query(q, pageSize, offset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query db: %s", err)
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.Sku, &product.ProdName, &product.Price, &product.PromoPrice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func GetProductFromSku(db *sql.DB, sku int64) (Product, error) {
  q := `SELECT sku, prodName, price, promoPrice FROM Product WHERE sku = ?`
  rows, err := db.Query(q, sku)
  if err != nil {
    return Product{}, err 
  }
  defer rows.Close()
  var product Product
	for rows.Next() {
		err := rows.Scan(&product.Sku, &product.ProdName, &product.Price, &product.PromoPrice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
			return Product{}, err
		}
	}

	return product, nil
}


func SearchProductNames(db *sql.DB, searchQuery string) ([]Product, error) {
	q := `SELECT * FROM Product WHERE prodName LIKE ?;`
	// q := fmt.Sprintf("SELECT * FROM %v WHERE %v LIKE %%v;%", tableName, searchQuery, fieldName)
	rows, err := db.Query(q, searchQuery)
	if err != nil {
		return []Product{}, err
	}

	defer rows.Close()
	retv := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.Id, &product.Sku, &product.ProdName,
	&product.Price, &product.PromoPrice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
			return []Product{}, err
		}
		retv = append(retv, product)
	}
	return retv, nil
}

