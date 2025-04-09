package respositories

import (
	"database/sql"
	"encoding/json"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type productRepository struct {
	database *database.Db
}

func NewProductRepository(db *database.Db) ports.IProductRepository {
	return &productRepository{
		database: db,
	}
}

const (
	productBaseFields = `
		p.id,
		p.name,
		p.description,
		p.sales_price,
		p.cost_price,
		p.picture_url,
		p.dish_type_map,
		p.active,
		p.category_id,
		p.restaurant_id`
)

func (r *productRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Product], error) {
	baseQuery := `
		SELECT 
			` + productBaseFields + `
		FROM products p
		WHERE p.deleted_at IS NULL`

	countQuery := `
		SELECT COUNT(*)
		FROM products p
		WHERE p.deleted_at IS NULL`

	query, count, params := r.database.ConstructFindQuery(baseQuery, countQuery, args)

	rows, err := r.database.Instance.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]aggregates.Product, 0, 10)
	for rows.Next() {
		var product aggregates.Product
		var dishTypeMapJSON []byte
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.SalesPrice,
			&product.CostPrice,
			&product.PictureUrl,
			&dishTypeMapJSON,
			&product.Active,
			&product.Category.Id,
			&product.Restaurant.Id,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(dishTypeMapJSON, &product.DishTypeMap); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, count, products)
	return &pagedSlice, nil
}

func (r *productRepository) FindById(id string) (*aggregates.Product, error) {
	query := `
		SELECT 
			` + productBaseFields + `
		FROM products p
		WHERE p.deleted_at IS NULL AND p.id = ?`

	var product aggregates.Product
	var dishTypeMapJSON []byte
	err := r.database.Instance.QueryRow(query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.SalesPrice,
		&product.CostPrice,
		&product.PictureUrl,
		&dishTypeMapJSON,
		&product.Active,
		&product.Category.Id,
		&product.Restaurant.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(dishTypeMapJSON, &product.DishTypeMap); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Create(product *aggregates.Product) error {
	dishTypeMapJSON, err := json.Marshal(product.DishTypeMap)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO products (
			id, name, description, sales_price, cost_price, picture_url,
			dish_type_map, active, category_id, restaurant_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = r.database.Instance.Exec(
		query,
		product.Id,
		product.Name,
		product.Description,
		product.SalesPrice,
		product.CostPrice,
		product.PictureUrl,
		dishTypeMapJSON,
		product.Active,
		product.Category.Id,
		product.Restaurant.Id,
	)
	return err
}

func (r *productRepository) Update(product *aggregates.Product) error {
	dishTypeMapJSON, err := json.Marshal(product.DishTypeMap)
	if err != nil {
		return err
	}

	query := `
		UPDATE products SET
			name = ?,
			description = ?,
			sales_price = ?,
			cost_price = ?,
			picture_url = ?,
			dish_type_map = ?,
			active = ?,
			category_id = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err = r.database.Instance.Exec(
		query,
		product.Name,
		product.Description,
		product.SalesPrice,
		product.CostPrice,
		product.PictureUrl,
		dishTypeMapJSON,
		product.Active,
		product.Category.Id,
		product.Id,
	)
	return err
}

func (r *productRepository) Delete(id string) error {
	query := `
		UPDATE products 
		SET deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL`
	_, err := r.database.Instance.Exec(query, id)
	return err
}

func (r *productRepository) Exists(name, restaurantId string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM products p
		WHERE p.deleted_at IS NULL 
		AND p.name = ? 
		AND p.restaurant_id = ?`

	var count int
	err := r.database.Instance.QueryRow(query, name, restaurantId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
