package respositories

import (
	"database/sql"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type dishRepository struct {
	db *database.Db
}

func NewDishRepository(db *database.Db) *dishRepository {
	return &dishRepository{
		db: db,
	}
}

const (
	dishBaseFields = `
		d.id,
		d.name,
		d.type,
		d.picture_url,
		d.restaurant_id`
)

func (r *dishRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Dish], error) {
	baseQuery := `
		SELECT 
			` + dishBaseFields + `
		FROM dishes d
		WHERE d.deleted_at IS NULL`

	countQuery := `
		SELECT COUNT(*)
		FROM dishes d
		WHERE d.deleted_at IS NULL`

	query, count, params := r.db.ConstructFindQuery(baseQuery, countQuery, args)

	rows, err := r.db.Instance.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dishes := make([]aggregates.Dish, 0, 10)
	for rows.Next() {
		var dish aggregates.Dish
		err := rows.Scan(
			&dish.Id,
			&dish.Name,
			&dish.Type,
			&dish.PictureUrl,
			&dish.Restaurant.Id,
		)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, count, dishes)
	return &pagedSlice, nil
}

func (r *dishRepository) FindById(id string) (*aggregates.Dish, error) {
	query := `
		SELECT 
			` + dishBaseFields + `
		FROM dishes d
		WHERE d.deleted_at IS NULL AND d.id = ?`

	var dish aggregates.Dish
	err := r.db.Instance.QueryRow(query, id).Scan(
		&dish.Id,
		&dish.Name,
		&dish.Type,
		&dish.PictureUrl,
		&dish.Restaurant.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &dish, nil
}

func (r *dishRepository) Create(dish *aggregates.Dish) error {
	query := `
		INSERT INTO dishes (
			id, name, type, picture_url, restaurant_id
		) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Instance.Exec(
		query,
		dish.Id,
		dish.Name,
		dish.Type,
		dish.PictureUrl,
		dish.Restaurant.Id,
	)
	return err
}

func (r *dishRepository) Update(dish *aggregates.Dish) error {
	query := `
		UPDATE dishes SET
			name = ?,
			type = ?,
			picture_url = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.Instance.Exec(
		query,
		dish.Name,
		dish.Type,
		dish.PictureUrl,
		dish.Id,
	)
	return err
}

func (r *dishRepository) Delete(id string) error {
	query := `
		UPDATE dishes 
		SET deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Instance.Exec(query, id)
	return err
}

func (r *dishRepository) FindByRestaurantId(restaurantId string) ([]aggregates.Dish, error) {
	query := `
		SELECT 
			` + dishBaseFields + `
		FROM dishes d
		WHERE d.deleted_at IS NULL AND d.restaurant_id = ?`

	rows, err := r.db.Instance.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dishes := make([]aggregates.Dish, 0, 10)
	for rows.Next() {
		var dish aggregates.Dish
		err := rows.Scan(
			&dish.Id,
			&dish.Name,
			&dish.Type,
			&dish.PictureUrl,
			&dish.Restaurant.Id,
		)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return dishes, nil
}

func (r *dishRepository) Exists(name string, restaurantId string) (bool, error) {
	// Implementation of the Exists method
	return false, nil
}
