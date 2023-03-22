package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

func main() {
	connectionSTR := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		"root", "root", "test_db")

	db, err := sql.Open("postgres", connectionSTR)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	product := NewProduct("Product 1", 10.5)

	err = insertProduct(db, product)
	if err != nil {
		panic(err)
	}

	product.Price = 20.5
	err = updateProduct(db, product)
	if err != nil {
		panic(err)
	}

	prod, err := selectProduct(db, product.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("id: %s, name: %s, price: %f\n\n", prod.ID, prod.Name, prod.Price)

	prods, err := selectAllProduct(db)
	if err != nil {
		panic(err)
	}

	for _, prod := range prods {
		fmt.Printf("id: %s, name: %s, price: %.2f\n", prod.ID, prod.Name, prod.Price)
	}

	err = deleteProduct(db, product.ID)
	if err != nil {
		panic(err)
	}

}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func insertProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("INSERT INTO goexpert.products (id, name, price) VALUES ($1, $2, $3) RETURNING id")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func updateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("UPDATE goexpert.products SET name = $1, price = $2 WHERE id = $3")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func selectProduct(db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("SELECT id, name, price FROM goexpert.products WHERE id = $1")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var product Product
	// err = stmt.QueryRowContext(ctx, id).Scan(&product.ID, &product.Name, &product.Price)
	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func selectAllProduct(db *sql.DB) ([]*Product, error) {
	stmt, err := db.Prepare("SELECT id, name, price FROM goexpert.products")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var products []*Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func deleteProduct(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM goexpert.products WHERE id = $1")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
