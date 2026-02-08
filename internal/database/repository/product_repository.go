package repository

import (
	"context"
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}
type Product struct {
	ID         string  `json:"id"`
	CategoryID string  `json:"category_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Picture    string  `json:"picture"`
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) Create(ctx context.Context, categoryID, name string, price float64, picture string) (string, error) {
	var id string
	query := `INSERT INTO products (category_id, name, price, picture) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, categoryID, name, price, picture).Scan(&id)
	return id, err
}
func (r *ProductRepository) GetAll(ctx context.Context) ([]Product, error) {
	query := `SELECT id, category_id, name, price, picture FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Picture); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	var p Product
	query := `SELECT id, category_id, name, price, picture FROM products WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Picture)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id string, categoryID string, name string, price float64, picture string) error {
	query := `UPDATE products SET category_id=$1, name=$2, price=$3, picture=$4 WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, categoryID, name, price, picture, id)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}