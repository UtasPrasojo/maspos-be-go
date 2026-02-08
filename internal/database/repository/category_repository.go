package repository

import (
    "context"
    "database/sql"
)

type Category struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type CategoryRepository struct {
    db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
    return &CategoryRepository{db}
}

func (r *CategoryRepository) Create(ctx context.Context, name string) (string, error) {
    var id string
    query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
    err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
    return id, err
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]Category, error) {
    query := `SELECT id, name FROM categories`
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []Category
    for rows.Next() {
        var c Category
        if err := rows.Scan(&c.ID, &c.Name); err != nil {
            return nil, err
		}
        categories = append(categories, c)
    }
    return categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*Category, error) {
    var c Category
    query := `SELECT id, name FROM categories WHERE id = $1`
    err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name)
    return &c, err
}

func (r *CategoryRepository) Update(ctx context.Context, id string, name string) error {
    query := `UPDATE categories SET name = $1 WHERE id = $2`
    _, err := r.db.ExecContext(ctx, query, name, id)
    return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
    query := `DELETE FROM categories WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}