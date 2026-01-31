package repositories

import (
	"database/sql"
	"errors"
	"product-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT p.id, p.name, p."desc", p.price, p.stock, p.category_id, c.name as category_name, p.created_at 
			  FROM products p 
			  LEFT JOIN categories c ON p.category_id = c.id`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var categoryID sql.NullInt64
		var categoryName sql.NullString
		err := rows.Scan(&p.ID, &p.Name, &p.Desc, &p.Price, &p.Stock, &categoryID, &categoryName, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			catID := int(categoryID.Int64)
			p.CategoryID = &catID
			if categoryName.Valid {
				p.CategoryName = &categoryName.String
			}
		} else {
			p.CategoryID = nil
			p.CategoryName = nil
		}

		products = append(products, p)
	}
	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, \"desc\", price, stock, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var categoryID interface{}
	if product.CategoryID != nil {
		categoryID = *product.CategoryID
	} else {
		categoryID = nil
	}

	err := repo.db.QueryRow(query, product.Name, product.Desc, product.Price, product.Stock, categoryID).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT p.id, p.name, p."desc", p.price, p.stock, p.category_id, c.name as category_name, p.created_at 
			  FROM products p 
			  LEFT JOIN categories c ON p.category_id = c.id 
			  WHERE p.id = $1`

	var p models.Product
	var categoryID sql.NullInt64
	var categoryName sql.NullString
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Desc, &p.Price, &p.Stock, &categoryID, &categoryName, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	if categoryID.Valid {
		catID := int(categoryID.Int64)
		p.CategoryID = &catID
		if categoryName.Valid {
			p.CategoryName = &categoryName.String
		}
	} else {
		p.CategoryID = nil
		p.CategoryName = nil
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, \"desc\" = $2, price = $3, stock = $4, category_id = $5 WHERE id = $6"

	var categoryID interface{}
	if product.CategoryID != nil {
		categoryID = *product.CategoryID
	} else {
		categoryID = nil
	}

	result, err := repo.db.Exec(query, product.Name, product.Desc, product.Price, product.Stock, categoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return err
}
