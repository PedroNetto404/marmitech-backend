package respositories

import (
	"database/sql"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type categoryRepository struct {
	database *database.Db
}

func NewCategoryRepository(db *database.Db) ports.ICategoryRepository {
	return &categoryRepository{
		database: db,
	}
}

const (
	categoryBaseFields = `
		c.id,
		c.name,
		c.picture_url,
		c.priority,
		c.active,
		c.restaurant_id`
)

func (r *categoryRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Category], error) {
	baseQuery := `
		SELECT 
			` + categoryBaseFields + `
		FROM categories c
		WHERE c.deleted_at IS NULL`

	countQuery := `
		SELECT COUNT(*)
		FROM categories c
		WHERE c.deleted_at IS NULL`

	query, count, params := r.database.ConstructFindQuery(baseQuery, countQuery, args)

	rows, err := r.database.Instance.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]aggregates.Category, 0, 10)
	for rows.Next() {
		var category aggregates.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.PictureUrl,
			&category.Priority,
			&category.Active,
			&category.Restaurant.Id,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, count, categories)
	return &pagedSlice, nil
}

func (r *categoryRepository) FindById(id string) (*aggregates.Category, error) {
	query := `
		SELECT 
			` + categoryBaseFields + `
		FROM categories c
		WHERE c.deleted_at IS NULL AND c.id = ?`

	var category aggregates.Category
	err := r.database.Instance.QueryRow(query, id).Scan(
		&category.Id,
		&category.Name,
		&category.PictureUrl,
		&category.Priority,
		&category.Active,
		&category.Restaurant.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Create(category *aggregates.Category) error {
	query := `
		INSERT INTO categories (
			id, name, picture_url, priority, active, restaurant_id
		) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.database.Instance.Exec(
		query,
		category.Id,
		category.Name,
		category.PictureUrl,
		category.Priority,
		category.Active,
		category.Restaurant.Id,
	)
	return err
}

func (r *categoryRepository) Update(category *aggregates.Category) error {
	query := `
		UPDATE categories SET
			name = ?,
			picture_url = ?,
			priority = ?,
			active = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.database.Instance.Exec(
		query,
		category.Name,
		category.PictureUrl,
		category.Priority,
		category.Active,
		category.Id,
	)
	return err
}

func (r *categoryRepository) Delete(id string) error {
	query := `
		UPDATE categories 
		SET deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL`
	_, err := r.database.Instance.Exec(query, id)
	return err
}

func (r *categoryRepository) FindByRestaurantId(restaurantId string) ([]aggregates.Category, error) {
	query := `
		SELECT 
			` + categoryBaseFields + `
		FROM categories c
		WHERE c.deleted_at IS NULL 
		AND c.restaurant_id = ?
		ORDER BY c.priority ASC`

	rows, err := r.database.Instance.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]aggregates.Category, 0, 10)
	for rows.Next() {
		var category aggregates.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.PictureUrl,
			&category.Priority,
			&category.Active,
			&category.Restaurant.Id,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Exists(name, restaurantId string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM categories c
		WHERE c.deleted_at IS NULL 
		AND c.name = ? 
		AND c.restaurant_id = ?`

	var count int
	err := r.database.Instance.QueryRow(query, name, restaurantId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
