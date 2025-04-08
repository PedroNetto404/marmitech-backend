package respositories

import (
	"fmt"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

const projection = `
	id,
	email,
	pwd_hash,
	active,
	role,
	permissions
`

type userRepository struct {
	database *database.Db
}

func NewUserRepository(database *database.Db) *userRepository {
	return &userRepository{
		database: database,
	}
}

func (r *userRepository) Create(user *aggregates.User) error {
	const sql = `
		INSERT INTO users (
			id, email, pwd_hash, active, role, permissions, 
			created_at, updated_at, deleted_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.database.Instance.Exec(
		sql, 
		user.Id, user.Email, user.PwdHash, user.Active, 
		user.Role, user.Permissions, user.CreatedAt, 
		user.UpdatedAt, user.DeletedAt,
	)
	return err
}

func (r *userRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.User], error) {
	baseQuery := fmt.Sprintf(`SELECT %s FROM users`, projection)
	countQuery := `SELECT COUNT(1) FROM users`

	finalQuery, totalCount, params := r.database.ConstructFindQuery(baseQuery, countQuery, args)

	rows, err := r.database.Instance.Query(finalQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]aggregates.User, 0)
	for rows.Next() {
		var user aggregates.User
		if err := rows.Scan(
			&user.Id, &user.Email, &user.PwdHash, 
			&user.Active, &user.Role, &user.Permissions,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, totalCount, users)
	return &pagedSlice, nil
}

func (r *userRepository) FindById(id string) (*aggregates.User, error) {
	sql := fmt.Sprintf(`SELECT %s FROM users WHERE id = ?`, projection)
	row := r.database.Instance.QueryRow(sql, id)

	var user aggregates.User
	if err := row.Scan(
		&user.Id, &user.Email, &user.PwdHash, 
		&user.Active, &user.Role, &user.Permissions,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(user *aggregates.User) error {
	const sql = `
		UPDATE users SET 
			email = ?, pwd_hash = ?, active = ?, role = ?, permissions = ?, 
			updated_at = ?, deleted_at = ?
		WHERE id = ?
	`

	_, err := r.database.Instance.Exec(
		sql, 
		user.Email, user.PwdHash, user.Active, user.Role, user.Permissions, 
		user.UpdatedAt, user.DeletedAt, user.Id,
	)
	return err
}

func (r *userRepository) Delete(id string) error {
	const sql = `
		DELETE FROM users WHERE id = ?
	`

	_, err := r.database.Instance.Exec(sql, id)
	return err
}

func (r *userRepository) Exists(email string) (bool, error) {
	const sql = `SELECT COUNT(1) FROM users WHERE email = ?`
	var count int
	if err := r.database.Instance.QueryRow(sql, email).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) FindByEmail(email string) (*aggregates.User, error) {
	sql := fmt.Sprintf(`SELECT %s FROM users WHERE email = ?`, projection)
	row := r.database.Instance.QueryRow(sql, email)

	var user aggregates.User
	if err := row.Scan(
		&user.Id, &user.Email, &user.PwdHash, 
		&user.Active, &user.Role, &user.Permissions,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
