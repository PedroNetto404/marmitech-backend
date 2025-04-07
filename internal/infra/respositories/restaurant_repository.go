package respositories

import (
	"database/sql"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type restaurantRepository struct {
	db *database.Database
}

func NewRestaurantRepository(db *database.Database) *restaurantRepository {
	return &restaurantRepository{
		db: db,
	}
}

func (r *restaurantRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Restaurant], error) {
	baseQuery := `
		SELECT 
			id,
			trade_name,
			legal_name,
			cnpj,
			contact_phone,
			whatsapp_phone,
			email,
			slug,
			address_street,
			address_number,
			address_complement,
			address_neighborhood,
			address_city,
			address_state,
			address_zip_code,
			address_country,
			address_lat,
			address_lng,
			logo_url,
			banner_url,
			show_cnpj_in_receipt,
			delivery_enabled,
			delivery_fee_per_km,
			delivery_minimum_order_value,
			delivery_max_radius_delivery,
			delivery_average_time_minutes,
			ecommerce_minimum_order_value,
			ecommerce_acquired,
			ecommerce_acquired_at,
			ecommerce_online,
			customer_post_paid_orders_acquired,
			customer_post_paid_orders_acquired_at,
			customer_post_paid_orders_enabled,
			created_at,
			updated_at
		FROM restaurants
		WHERE deleted_at IS NULL AND active = 1`

	countQuery := `SELECT COUNT(*) FROM restaurants WHERE deleted_at IS NULL AND active = 1`
	query, count, params := r.db.BuildFindQuery(baseQuery, countQuery, args)

	rows, err := r.db.Instance.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := make([]aggregates.Restaurant, 0, 10)
	for rows.Next() {
		var restaurant aggregates.Restaurant
		var addressComplement sql.NullString
		err := rows.Scan(
			&restaurant.Id,
			&restaurant.TradeName,
			&restaurant.LegalName,
			&restaurant.CNPJ,
			&restaurant.ContactPhone,
			&restaurant.WhatsAppPhone,
			&restaurant.Email,
			&restaurant.Slug,
			&restaurant.Address.Street,
			&restaurant.Address.Number,
			&addressComplement,
			&restaurant.Address.Neighborhood,
			&restaurant.Address.City,
			&restaurant.Address.State,
			&restaurant.Address.ZipCode,
			&restaurant.Address.Country,
			&restaurant.Address.Lat,
			&restaurant.Address.Lng,
			&restaurant.LogoUrl,
			&restaurant.BannerUrl,
			&restaurant.Settings.ShowCnpjInReceipt,
			&restaurant.Settings.Delivery.Enabled,
			&restaurant.Settings.Delivery.FeePerKm,
			&restaurant.Settings.Delivery.MinimumOrderValue,
			&restaurant.Settings.Delivery.MaxRadiusDelivery,
			&restaurant.Settings.Delivery.AverageTimeMinutes,
			&restaurant.Settings.Ecommerce.MinimumOrderValue,
			&restaurant.Settings.Ecommerce.Acquired,
			&restaurant.Settings.Ecommerce.AcquiredAt,
			&restaurant.Settings.Ecommerce.Online,
			&restaurant.Settings.CustomerPostPaidOrders.Acquired,
			&restaurant.Settings.CustomerPostPaidOrders.AcquiredAt,
			&restaurant.Settings.CustomerPostPaidOrders.Enabled,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if addressComplement.Valid {
			restaurant.Address.Complement = addressComplement.String
		} else {
			restaurant.Address.Complement = ""
		}
		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, count, restaurants)
	return &pagedSlice, nil
}

func (repo *restaurantRepository) FindById(id string) (*aggregates.Restaurant, error) {
	query := `
		SELECT 
			id,
			trade_name,
			legal_name,
			cnpj,
			contact_phone,
			whatsapp_phone,
			email,
			slug,
			address_street,
			address_number,
			address_complement,
			address_neighborhood,
			address_city,
			address_state,
			address_zip_code,
			address_country,
			address_lat,
			address_lng,
			show_cnpj_in_receipt,
			delivery_enabled,
			delivery_fee_per_km,
			delivery_minimum_order_value,
			delivery_max_radius_delivery,
			delivery_average_time_minutes,
			ecommerce_minimum_order_value,
			ecommerce_acquired,
			ecommerce_acquired_at,
			ecommerce_online,
			customer_post_paid_orders_acquired,
			customer_post_paid_orders_acquired_at,
			customer_post_paid_orders_enabled,
			logo_url,
			banner_url,
			created_at,
			updated_at,
			active
		FROM restaurants
		WHERE deleted_at IS NULL AND id = ?`

	var r aggregates.Restaurant
	var addressComplement sql.NullString
	err := repo.db.Instance.QueryRow(query, id).Scan(
		&r.Id,
		&r.TradeName,
		&r.LegalName,
		&r.CNPJ,
		&r.ContactPhone,
		&r.WhatsAppPhone,
		&r.Email,
		&r.Slug,
		&r.Address.Street,
		&r.Address.Number,
		&addressComplement,
		&r.Address.Neighborhood,
		&r.Address.City,
		&r.Address.State,
		&r.Address.ZipCode,
		&r.Address.Country,
		&r.Address.Lat,
		&r.Address.Lng,
		&r.Settings.ShowCnpjInReceipt,
		&r.Settings.Delivery.Enabled,
		&r.Settings.Delivery.FeePerKm,
		&r.Settings.Delivery.MinimumOrderValue,
		&r.Settings.Delivery.MaxRadiusDelivery,
		&r.Settings.Delivery.AverageTimeMinutes,
		&r.Settings.Ecommerce.MinimumOrderValue,
		&r.Settings.Ecommerce.Acquired,
		&r.Settings.Ecommerce.AcquiredAt,
		&r.Settings.Ecommerce.Online,
		&r.Settings.CustomerPostPaidOrders.Acquired,
		&r.Settings.CustomerPostPaidOrders.AcquiredAt,
		&r.Settings.CustomerPostPaidOrders.Enabled,
		&r.LogoUrl,
		&r.BannerUrl,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.Active,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if addressComplement.Valid {
		r.Address.Complement = addressComplement.String
	} else {
		r.Address.Complement = ""
	}

	return &r, nil
}

func (r *restaurantRepository) Create(record *aggregates.Restaurant) error {
	query := `
		INSERT INTO restaurants (
			id, trade_name, legal_name, cnpj, contact_phone, whatsapp_phone, email, slug,
			address_street, address_number, address_complement, address_neighborhood,
			address_city, address_state, address_zip_code, address_country, address_lat, address_lng,
			show_cnpj_in_receipt, delivery_enabled, delivery_fee_per_km, delivery_minimum_order_value,
			delivery_max_radius_delivery, delivery_average_time_minutes,
			ecommerce_minimum_order_value, ecommerce_acquired, ecommerce_acquired_at, ecommerce_online,
			customer_post_paid_orders_acquired, customer_post_paid_orders_acquired_at, customer_post_paid_orders_enabled,
			logo_url, banner_url, created_at, updated_at, active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Instance.Exec(query,
		record.Id, record.TradeName, record.LegalName, record.CNPJ,
		record.ContactPhone, record.WhatsAppPhone, record.Email, record.Slug,
		record.Address.Street, record.Address.Number, record.Address.Complement, record.Address.Neighborhood,
		record.Address.City, record.Address.State, record.Address.ZipCode, record.Address.Country,
		record.Address.Lat, record.Address.Lng,
		record.Settings.ShowCnpjInReceipt,
		record.Settings.Delivery.Enabled,
		record.Settings.Delivery.FeePerKm,
		record.Settings.Delivery.MinimumOrderValue,
		record.Settings.Delivery.MaxRadiusDelivery,
		record.Settings.Delivery.AverageTimeMinutes,
		record.Settings.Ecommerce.MinimumOrderValue,
		record.Settings.Ecommerce.Acquired,
		record.Settings.Ecommerce.AcquiredAt,
		record.Settings.Ecommerce.Online,
		record.Settings.CustomerPostPaidOrders.Acquired,
		record.Settings.CustomerPostPaidOrders.AcquiredAt,
		record.Settings.CustomerPostPaidOrders.Enabled,
		record.LogoUrl, record.BannerUrl,
		record.CreatedAt, record.UpdatedAt, record.Active,
	)
	return err
}

func (r *restaurantRepository) Update(record *aggregates.Restaurant) error {
	query := `
		UPDATE restaurants SET
			trade_name = ?, 
			legal_name = ?, 
			cnpj = ?, 
			contact_phone = ?, 
			whatsapp_phone = ?, 
			email = ?, 
			slug = ?,
			address_street = ?, 
			address_number = ?, 
			address_complement = ?, 
			address_neighborhood = ?,
			address_city = ?, 
			address_state = ?, 
			address_zip_code = ?, 
			address_country = ?, 
			address_lat = ?, 
			address_lng = ?,
			show_cnpj_in_receipt = ?, 
			delivery_enabled = ?, 
			delivery_fee_per_km = ?, 
			delivery_minimum_order_value = ?,
			delivery_max_radius_delivery = ?, 
			delivery_average_time_minutes = ?,
			ecommerce_minimum_order_value = ?, 
			ecommerce_acquired = ?, 
			ecommerce_acquired_at = ?, 
			ecommerce_online = ?,
			customer_post_paid_orders_acquired = ?, 
			customer_post_paid_orders_acquired_at = ?, 
			customer_post_paid_orders_enabled = ?,
			logo_url = ?, 
			banner_url = ?, 
			updated_at = ?, 
			active = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.Instance.Exec(
		query,
		record.TradeName,
		record.LegalName,
		record.CNPJ,
		record.ContactPhone,
		record.WhatsAppPhone,
		record.Email,
		record.Slug,
		record.Address.Street,
		record.Address.Number,
		record.Address.Complement,
		record.Address.Neighborhood,
		record.Address.City,
		record.Address.State,
		record.Address.ZipCode,
		record.Address.Country,
		record.Address.Lat,
		record.Address.Lng,
		record.Settings.ShowCnpjInReceipt,
		record.Settings.Delivery.Enabled,
		record.Settings.Delivery.FeePerKm,
		record.Settings.Delivery.MinimumOrderValue,
		record.Settings.Delivery.MaxRadiusDelivery,
		record.Settings.Delivery.AverageTimeMinutes,
		record.Settings.Ecommerce.MinimumOrderValue,
		record.Settings.Ecommerce.Acquired,
		record.Settings.Ecommerce.AcquiredAt,
		record.Settings.Ecommerce.Online,
		record.Settings.CustomerPostPaidOrders.Acquired,
		record.Settings.CustomerPostPaidOrders.AcquiredAt,
		record.Settings.CustomerPostPaidOrders.Enabled,
		record.LogoUrl,
		record.BannerUrl,
		record.UpdatedAt,
		record.Active,
		record.Id,
	)
	return err
}

func (r *restaurantRepository) Delete(id string) error {
	query := `
	UPDATE restaurants 
	SET deleted_at = NOW() 
	WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Instance.Exec(query, id)
	return err
}

func (repo *restaurantRepository) FindByDocument(cnpj string) (*aggregates.Restaurant, error) {
	query := `
	SELECT 
		id, 
		trade_name, 
		legal_name, 
		cnpj, 
		contact_phone, 
		whatsapp_phone, 
		email, 
		slug,
		address_street, 
		address_number, 
		address_complement, 
		address_neighborhood,
		address_city, 
		address_state, 
		address_zip_code, 
		address_country, 
		address_lat, 
		address_lng,
		show_cnpj_in_receipt,
		delivery_enabled, 
		delivery_fee_per_km, 
		delivery_minimum_order_value, 
		delivery_max_radius_delivery, 
		delivery_average_time_minutes,
		ecommerce_minimum_order_value, 
		ecommerce_acquired, 
		ecommerce_acquired_at, 
		ecommerce_online,
		customer_post_paid_orders_acquired, 
		customer_post_paid_orders_acquired_at, 
		customer_post_paid_orders_enabled,
		logo_url, 
		banner_url,
		created_at, 
		updated_at, 
		active
	FROM restaurants 
	WHERE deleted_at IS NULL AND cnpj = ?`

	var r aggregates.Restaurant
	err := repo.db.Instance.QueryRow(query, cnpj).Scan(
		&r.Id,
		&r.TradeName,
		&r.LegalName,
		&r.CNPJ,
		&r.ContactPhone,
		&r.WhatsAppPhone,
		&r.Email,
		&r.Slug,
		&r.Address.Street,
		&r.Address.Number,
		&r.Address.Complement,
		&r.Address.Neighborhood,
		&r.Address.City,
		&r.Address.State,
		&r.Address.ZipCode,
		&r.Address.Country,
		&r.Address.Lat,
		&r.Address.Lng,
		&r.Settings.ShowCnpjInReceipt,
		&r.Settings.Delivery.Enabled,
		&r.Settings.Delivery.FeePerKm,
		&r.Settings.Delivery.MinimumOrderValue,
		&r.Settings.Delivery.MaxRadiusDelivery,
		&r.Settings.Delivery.AverageTimeMinutes,
		&r.Settings.Ecommerce.MinimumOrderValue,
		&r.Settings.Ecommerce.Acquired,
		&r.Settings.Ecommerce.AcquiredAt,
		&r.Settings.Ecommerce.Online,
		&r.Settings.CustomerPostPaidOrders.Acquired,
		&r.Settings.CustomerPostPaidOrders.AcquiredAt,
		&r.Settings.CustomerPostPaidOrders.Enabled,
		&r.LogoUrl,
		&r.BannerUrl,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.Active,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}
