package respositories

import (
	"database/sql"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type customerRepository struct {
	database *database.Db
}

func NewCustomerRepository(db *database.Db) ports.ICustomerRepository {
	return &customerRepository{
		database: db,
	}
}

const (
	customerBaseFields = `
		c.id,
		c.first_name,
		c.last_name,
		c.contact_email,
		c.contact_phone,
		c.deleted_at`
)

func (r *customerRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Customer], error) {
	baseQuery := `
		SELECT 
			` + customerBaseFields + `,
			` + AddressFields + `
		FROM customers c
		JOIN addresses a ON c.address_id = a.id
		WHERE c.deleted_at IS NULL`

	countQuery := `
		SELECT COUNT(*)
		FROM customers c
		WHERE c.deleted_at IS NULL`

	query, count, params := r.database.ConstructFindQuery(baseQuery, countQuery, args)

	rows, err := r.database.Instance.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]aggregates.Customer, 0, 10)
	for rows.Next() {
		var customer aggregates.Customer
		var deletedAt sql.NullTime
		err := rows.Scan(
			&customer.Id,
			&customer.FirstName,
			&customer.LastName,
			&customer.ContactEmail,
			&customer.ContactPhone,
			&deletedAt,
			&customer.Address.Id,
			&customer.Address.Alias,
			&customer.Address.Street,
			&customer.Address.Number,
			&customer.Address.Complement,
			&customer.Address.Neighborhood,
			&customer.Address.City,
			&customer.Address.State,
			&customer.Address.Country,
			&customer.Address.ZipCode,
			&customer.Address.Lat,
			&customer.Address.Lng,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, count, customers)
	return &pagedSlice, nil
}

func (r *customerRepository) FindById(id string) (*aggregates.Customer, error) {
	query := `
		SELECT 
			` + customerBaseFields + `,
			` + AddressFields + `
		FROM customers c
		JOIN addresses a ON c.address_id = a.id
		WHERE c.deleted_at IS NULL AND c.id = ?`

	var customer aggregates.Customer
	var deletedAt sql.NullTime
	err := r.database.Instance.QueryRow(query, id).Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.LastName,
		&customer.ContactEmail,
		&customer.ContactPhone,
		&deletedAt,
		&customer.Address.Id,
		&customer.Address.Alias,
		&customer.Address.Street,
		&customer.Address.Number,
		&customer.Address.Complement,
		&customer.Address.Neighborhood,
		&customer.Address.City,
		&customer.Address.State,
		&customer.Address.Country,
		&customer.Address.ZipCode,
		&customer.Address.Lat,
		&customer.Address.Lng,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) Create(customer *aggregates.Customer) error {
	tx, err := r.database.Instance.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First create the address
	addressQuery := `
		INSERT INTO addresses (
			id, alias, street, number, complement, neighborhood,
			city, state, country, zip_code, lat, lng
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(
		addressQuery,
		customer.Address.Id,
		customer.Address.Alias,
		customer.Address.Street,
		customer.Address.Number,
		customer.Address.Complement,
		customer.Address.Neighborhood,
		customer.Address.City,
		customer.Address.State,
		customer.Address.Country,
		customer.Address.ZipCode,
		customer.Address.Lat,
		customer.Address.Lng,
	)
	if err != nil {
		return err
	}

	// Then create the customer
	customerQuery := `
		INSERT INTO customers (
			id, first_name, last_name, contact_email, contact_phone,
			address_id
		) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(
		customerQuery,
		customer.Id,
		customer.FirstName,
		customer.LastName,
		customer.ContactEmail,
		customer.ContactPhone,
		customer.Address.Id,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *customerRepository) Update(customer *aggregates.Customer) error {
	tx, err := r.database.Instance.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First update the address
	addressQuery := `
		UPDATE addresses SET
			alias = ?,
			street = ?,
			number = ?,
			complement = ?,
			neighborhood = ?,
			city = ?,
			state = ?,
			country = ?,
			zip_code = ?,
			lat = ?,
			lng = ?
		WHERE id = ?`

	_, err = tx.Exec(
		addressQuery,
		customer.Address.Alias,
		customer.Address.Street,
		customer.Address.Number,
		customer.Address.Complement,
		customer.Address.Neighborhood,
		customer.Address.City,
		customer.Address.State,
		customer.Address.Country,
		customer.Address.ZipCode,
		customer.Address.Lat,
		customer.Address.Lng,
		customer.Address.Id,
	)
	if err != nil {
		return err
	}

	// Then update the customer
	customerQuery := `
		UPDATE customers SET
			first_name = ?,
			last_name = ?,
			contact_email = ?,
			contact_phone = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err = tx.Exec(
		customerQuery,
		customer.FirstName,
		customer.LastName,
		customer.ContactEmail,
		customer.ContactPhone,
		customer.Id,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *customerRepository) Delete(id string) error {
	query := `
		UPDATE customers 
		SET deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL`
	_, err := r.database.Instance.Exec(query, id)
	return err
}

func (r *customerRepository) FindByEmail(email string) (*aggregates.Customer, error) {
	query := `
		SELECT 
			` + customerBaseFields + `,
			` + AddressFields + `
		FROM customers c
		JOIN addresses a ON c.address_id = a.id
		WHERE c.deleted_at IS NULL AND c.contact_email = ?`

	var customer aggregates.Customer
	var deletedAt sql.NullTime
	err := r.database.Instance.QueryRow(query, email).Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.LastName,
		&customer.ContactEmail,
		&customer.ContactPhone,
		&deletedAt,
		&customer.Address.Id,
		&customer.Address.Alias,
		&customer.Address.Street,
		&customer.Address.Number,
		&customer.Address.Complement,
		&customer.Address.Neighborhood,
		&customer.Address.City,
		&customer.Address.State,
		&customer.Address.Country,
		&customer.Address.ZipCode,
		&customer.Address.Lat,
		&customer.Address.Lng,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) Exists(email string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM customers c
		WHERE c.deleted_at IS NULL 
		AND c.contact_email = ?`

	var count int
	err := r.database.Instance.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
