package respositories

import (
	"database/sql"
	"encoding/json"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type userRepository struct {
	db *database.Db
}

func NewUserRepository(db *database.Db) *userRepository {
	return &userRepository{
		db: db,
	}
}

const (
	userBaseFields = `
		u.id,
		u.email,
		u.pwd_hash,
		u.active,
		u.role,
		u.permissions,
		u.created_at,
		u.updated_at,
		u.deleted_at`
)

func (r *userRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.User], error) {
	baseQuery := `
		SELECT 
			` + userBaseFields + `
		FROM users u
		WHERE u.deleted_at IS NULL`

	countQuery := `
		SELECT COUNT(*)
		FROM users u
		WHERE u.deleted_at IS NULL`

	query, count, params := r.db.ConstructFindQuery(baseQuery, countQuery, args)

	rows, err := r.db.Instance.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]aggregates.User, 0, 10)
	for rows.Next() {
		var user aggregates.User
		var permissionsJSON []byte
		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.PwdHash,
			&user.Active,
			&user.Role,
			&permissionsJSON,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(permissionsJSON, &user.Permissions); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, count, users)
	return &pagedSlice, nil
}

func (r *userRepository) FindById(id string) (*aggregates.User, error) {
	query := `
		SELECT 
			` + userBaseFields + `
		FROM users u
		WHERE u.deleted_at IS NULL AND u.id = ?`

	var user aggregates.User
	var permissionsJSON []byte
	err := r.db.Instance.QueryRow(query, id).Scan(
		&user.Id,
		&user.Email,
		&user.PwdHash,
		&user.Active,
		&user.Role,
		&permissionsJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(permissionsJSON, &user.Permissions); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user *aggregates.User) error {
	permissionsJSON, err := json.Marshal(user.Permissions)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (
			id, email, pwd_hash, active, role, permissions, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = r.db.Instance.Exec(
		query,
		user.Id,
		user.Email,
		user.PwdHash,
		user.Active,
		user.Role,
		permissionsJSON,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *userRepository) Update(user *aggregates.User) error {
	permissionsJSON, err := json.Marshal(user.Permissions)
	if err != nil {
		return err
	}

	query := `
		UPDATE users SET
			email = ?,
			pwd_hash = ?,
			active = ?,
			role = ?,
			permissions = ?,
			updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err = r.db.Instance.Exec(
		query,
		user.Email,
		user.PwdHash,
		user.Active,
		user.Role,
		permissionsJSON,
		user.UpdatedAt,
		user.Id,
	)
	return err
}

func (r *userRepository) Delete(id string) error {
	query := `
		UPDATE users 
		SET deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Instance.Exec(query, id)
	return err
}

func (r *userRepository) Exists(email string) (bool, error) {
	const sql = `
		SELECT COUNT(1) 
		FROM users 
		WHERE email = ? 
		AND deleted_at IS NULL`
	var count int
	if err := r.db.Instance.QueryRow(sql, email).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) FindByEmail(email string) (*aggregates.User, error) {
	query := `
		SELECT 
			` + userBaseFields + `
		FROM users u
		WHERE u.deleted_at IS NULL AND u.email = ?`

	var user aggregates.User
	var permissionsJSON []byte
	err := r.db.Instance.QueryRow(query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PwdHash,
		&user.Active,
		&user.Role,
		&permissionsJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(permissionsJSON, &user.Permissions); err != nil {
		return nil, err
	}

	return &user, nil
}
