package main

import (
	"database/sql"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price, category FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price, &p.Category)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2, category=$3 WHERE id=$4",
			p.Name, p.Price, p.Category, p.ID)

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price, category) VALUES($1, $2, $3) RETURNING id",
		p.Name, p.Price, p.Category).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getValuesFromRows(rows *sql.Rows) ([]product, error) {
	products := []product{}
	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name, price, category FROM products LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return getValuesFromRows(rows)
}

func searchProductsByName(db *sql.DB, name string) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name, price, category FROM products WHERE name ILIKE '%'||$1||'%'", name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return getValuesFromRows(rows)
}
