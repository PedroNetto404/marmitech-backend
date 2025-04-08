package respositories

import (
	"database/sql"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type restaurantRepository struct {
	db *database.Db
}

func NewRestaurantRepository(db *database.Db) *restaurantRepository {
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
			address_id,
			show_cnpj_in_receipt,
			delivery_enabled,
			delivery_fee_per_km,
			delivery_minimum_order_value,
			delivery_max_radius_km,
			delivery_average_time_minutes,
			ecommerce_enabled,
			ecommerce_minimum_order_value,
			ecommerce_delivery_fee_per_km,
			ecommerce_delivery_max_radius_km,
			ecommerce_delivery_average_time_minutes,
			customer_post_paid_orders_enabled,
			customer_post_paid_orders_minimum_order_value,
			customer_post_paid_orders_average_time_minutes,
			customer_post_paid_orders_delivery_fee_per_km,
			customer_post_paid_orders_delivery_max_radius_km,
			customer_post_paid_orders_delivery_average_time_minutes,
			created_at,
			updated_at,
			active
		FROM restaurants
		WHERE deleted_at IS NULL AND active = 1`

	countQuery := `SELECT COUNT(*) FROM restaurants WHERE deleted_at IS NULL AND active = 1`
	query, count, params := r.db.ConstructFindQuery(baseQuery, countQuery, args)

	rows, err := r.db.Instance.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := make([]aggregates.Restaurant, 0, 10)
	for rows.Next() {
		var restaurant aggregates.Restaurant
		err := rows.Scan(
			&restaurant.Id,
			&restaurant.TradeName,
			&restaurant.LegalName,
			&restaurant.CNPJ,
			&restaurant.ContactPhone,
			&restaurant.WhatsAppPhone,
			&restaurant.Email,
			&restaurant.Slug,
			&restaurant.Address.Id,
			&restaurant.Settings.ShowCnpjInReceipt,
			&restaurant.Settings.Delivery.Enabled,
			&restaurant.Settings.Delivery.FeePerKm,
			&restaurant.Settings.Delivery.MinimumOrderValue,
			&restaurant.Settings.Delivery.MaxRadiusDelivery,
			&restaurant.Settings.Delivery.AverageTimeMinutes,
			&restaurant.Settings.Ecommerce.Enabled,
			&restaurant.Settings.Ecommerce.MinimumOrderValue,
			&restaurant.Settings.CustomerPostPaidOrders.Enabled,
			&restaurant.Settings.CustomerPostPaidOrders.Acquired,
			&restaurant.Settings.CustomerPostPaidOrders.AcquiredAt,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
			&restaurant.Active,
		)
		if err != nil {
			return nil, err
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
			address_id,
			show_cnpj_in_receipt,
			delivery_enabled,
			delivery_fee_per_km,
			delivery_minimum_order_value,
			delivery_max_radius_km,
			delivery_average_time_minutes,
			ecommerce_enabled,
			ecommerce_minimum_order_value,
			ecommerce_delivery_fee_per_km,
			ecommerce_delivery_max_radius_km,
			ecommerce_delivery_average_time_minutes,
			customer_post_paid_orders_enabled,
			customer_post_paid_orders_minimum_order_value,
			customer_post_paid_orders_average_time_minutes,
			customer_post_paid_orders_delivery_fee_per_km,
			customer_post_paid_orders_delivery_max_radius_km,
			customer_post_paid_orders_delivery_average_time_minutes,
			logo_url,
			banner_url,
			created_at,
			updated_at,
			active
		FROM restaurants
		WHERE deleted_at IS NULL AND id = ?`

	var r aggregates.Restaurant
	err := repo.db.Instance.QueryRow(query, id).Scan(
		&r.Id,
		&r.TradeName,
		&r.LegalName,
		&r.CNPJ,
		&r.ContactPhone,
		&r.WhatsAppPhone,
		&r.Email,
		&r.Slug,
		&r.AddressId,
		&r.Settings.ShowCnpjInReceipt,
		&r.Settings.Delivery.Enabled,
		&r.Settings.Delivery.FeePerKm,
		&r.Settings.Delivery.MinimumOrderValue,
		&r.Settings.Delivery.MaxRadiusKm,
		&r.Settings.Delivery.AverageTimeMinutes,
		&r.Settings.Ecommerce.Enabled,
		&r.Settings.Ecommerce.MinimumOrderValue,
		&r.Settings.Ecommerce.DeliveryFeePerKm,
		&r.Settings.Ecommerce.DeliveryMaxRadiusKm,
		&r.Settings.Ecommerce.DeliveryAverageTimeMinutes,
		&r.Settings.CustomerPostPaidOrders.Enabled,
		&r.Settings.CustomerPostPaidOrders.MinimumOrderValue,
		&r.Settings.CustomerPostPaidOrders.AverageTimeMinutes,
		&r.Settings.CustomerPostPaidOrders.DeliveryFeePerKm,
		&r.Settings.CustomerPostPaidOrders.DeliveryMaxRadiusKm,
		&r.Settings.CustomerPostPaidOrders.DeliveryAverageTimeMinutes,
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

func (r *restaurantRepository) Create(record *aggregates.Restaurant) error {
	query := `
		INSERT INTO restaurants (
			id, trade_name, legal_name, cnpj, contact_phone, whatsapp_phone, email, slug,
			address_id, show_cnpj_in_receipt, delivery_enabled, delivery_fee_per_km, delivery_minimum_order_value,
			delivery_max_radius_km, delivery_average_time_minutes, ecommerce_enabled, ecommerce_minimum_order_value,
			ecommerce_delivery_fee_per_km, ecommerce_delivery_max_radius_km, ecommerce_delivery_average_time_minutes,
			customer_post_paid_orders_enabled, customer_post_paid_orders_minimum_order_value, customer_post_paid_orders_average_time_minutes,
			customer_post_paid_orders_delivery_fee_per_km, customer_post_paid_orders_delivery_max_radius_km, customer_post_paid_orders_delivery_average_time_minutes,
			logo_url, banner_url, created_at, updated_at, active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Instance.Exec(query,
		record.Id, record.TradeName, record.LegalName, record.CNPJ,
		record.ContactPhone, record.WhatsAppPhone, record.Email, record.Slug,
		record.AddressId, record.Settings.ShowCnpjInReceipt, record.Settings.Delivery.Enabled,
		record.Settings.Delivery.FeePerKm, record.Settings.Delivery.MinimumOrderValue,
		record.Settings.Delivery.MaxRadiusKm, record.Settings.Delivery.AverageTimeMinutes,
		record.Settings.Ecommerce.Enabled, record.Settings.Ecommerce.MinimumOrderValue,
		record.Settings.Ecommerce.DeliveryFeePerKm, record.Settings.Ecommerce.DeliveryMaxRadiusKm,
		record.Settings.Ecommerce.DeliveryAverageTimeMinutes, record.Settings.CustomerPostPaidOrders.Enabled,
		record.Settings.CustomerPostPaidOrders.MinimumOrderValue, record.Settings.CustomerPostPaidOrders.AverageTimeMinutes,
		record.Settings.CustomerPostPaidOrders.DeliveryFeePerKm, record.Settings.CustomerPostPaidOrders.DeliveryMaxRadiusKm,
		record.Settings.CustomerPostPaidOrders.DeliveryAverageTimeMinutes, record.LogoUrl, record.BannerUrl,
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
			address_id = ?, 
			show_cnpj_in_receipt = ?, 
			delivery_enabled = ?, 
			delivery_fee_per_km = ?, 
			delivery_minimum_order_value = ?,
			delivery_max_radius_km = ?, 
			delivery_average_time_minutes = ?,
			ecommerce_enabled = ?, 
			ecommerce_minimum_order_value = ?, 
			ecommerce_delivery_fee_per_km = ?, 
			ecommerce_delivery_max_radius_km = ?, 
			ecommerce_delivery_average_time_minutes = ?,
			customer_post_paid_orders_enabled = ?, 
			customer_post_paid_orders_minimum_order_value = ?, 
			customer_post_paid_orders_average_time_minutes = ?, 
			customer_post_paid_orders_delivery_fee_per_km = ?, 
			customer_post_paid_orders_delivery_max_radius_km = ?, 
			customer_post_paid_orders_delivery_average_time_minutes = ?,
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
		record.AddressId,
		record.Settings.ShowCnpjInReceipt,
		record.Settings.Delivery.Enabled,
		record.Settings.Delivery.FeePerKm,
		record.Settings.Delivery.MinimumOrderValue,
		record.Settings.Delivery.MaxRadiusKm,
		record.Settings.Delivery.AverageTimeMinutes,
		record.Settings.Ecommerce.Enabled,
		record.Settings.Ecommerce.MinimumOrderValue,
		record.Settings.Ecommerce.DeliveryFeePerKm,
		record.Settings.Ecommerce.DeliveryMaxRadiusKm,
		record.Settings.Ecommerce.DeliveryAverageTimeMinutes,
		record.Settings.CustomerPostPaidOrders.Enabled,
		record.Settings.CustomerPostPaidOrders.MinimumOrderValue,
		record.Settings.CustomerPostPaidOrders.AverageTimeMinutes,
		record.Settings.CustomerPostPaidOrders.DeliveryFeePerKm,
		record.Settings.CustomerPostPaidOrders.DeliveryMaxRadiusKm,
		record.Settings.CustomerPostPaidOrders.DeliveryAverageTimeMinutes,
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
		address_id, 
		show_cnpj_in_receipt,
		delivery_enabled, 
		delivery_fee_per_km, 
		delivery_minimum_order_value, 
		delivery_max_radius_km, 
		delivery_average_time_minutes,
		ecommerce_enabled, 
		ecommerce_minimum_order_value, 
		ecommerce_delivery_fee_per_km, 
		ecommerce_delivery_max_radius_km, 
		ecommerce_delivery_average_time_minutes,
		customer_post_paid_orders_enabled, 
		customer_post_paid_orders_minimum_order_value, 
		customer_post_paid_orders_average_time_minutes, 
		customer_post_paid_orders_delivery_fee_per_km, 
		customer_post_paid_orders_delivery_max_radius_km, 
		customer_post_paid_orders_delivery_average_time_minutes,
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
		&r.AddressId,
		&r.Settings.ShowCnpjInReceipt,
		&r.Settings.Delivery.Enabled,
		&r.Settings.Delivery.FeePerKm,
		&r.Settings.Delivery.MinimumOrderValue,
		&r.Settings.Delivery.MaxRadiusKm,
		&r.Settings.Delivery.AverageTimeMinutes,
		&r.Settings.Ecommerce.Enabled,
		&r.Settings.Ecommerce.MinimumOrderValue,
		&r.Settings.Ecommerce.DeliveryFeePerKm,
		&r.Settings.Ecommerce.DeliveryMaxRadiusKm,
		&r.Settings.Ecommerce.DeliveryAverageTimeMinutes,
		&r.Settings.CustomerPostPaidOrders.Enabled,
		&r.Settings.CustomerPostPaidOrders.MinimumOrderValue,
		&r.Settings.CustomerPostPaidOrders.AverageTimeMinutes,
		&r.Settings.CustomerPostPaidOrders.DeliveryFeePerKm,
		&r.Settings.CustomerPostPaidOrders.DeliveryMaxRadiusKm,
		&r.Settings.CustomerPostPaidOrders.DeliveryAverageTimeMinutes,
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
