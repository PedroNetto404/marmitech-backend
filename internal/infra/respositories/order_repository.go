package respositories

import (
	"database/sql"
	"time"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) ports.IOrderRepository {
	return &orderRepository{db: db}
}

const (
	orderBaseFields = `
		id,
		restaurant_id,
		customer_id,
		status,
		total,
		discount,
		observation,
		deleted_at
	`

	orderDeliveryFields = `
		id,
		order_id,
		address_alias,
		address_street,
		address_number,
		address_complement,
		address_neighborhood,
		address_city,
		address_state,
		address_country,
		address_zip_code,
		address_lat,
		address_lng,
		fee,
		distance,
		average_time_minutes,
		status
	`

	orderItemFields = `
		id,
		order_id,
		product_id,
		product_name,
		product_unit_price,
		quantity,
		observation,
		discount,
		item_total
	`

	orderPaymentFields = `
		id,
		order_id,
		payment_method,
		amount,
		status
	`
)

func (r *orderRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Order], error) {
	offset := (args.Offset - 1) * args.Limit

	// Count total records
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM orders
		WHERE deleted_at IS NULL
	`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get orders
	query := `
		SELECT ` + orderBaseFields + `
		FROM orders
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, args.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []aggregates.Order
	for rows.Next() {
		var order aggregates.Order
		var deletedAt sql.NullTime
		err := rows.Scan(
			&order.Id,
			&order.RestaurantID,
			&order.CustomerID,
			&order.Status,
			&order.Total,
			&order.Discount,
			&order.Observation,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get delivery
		delivery, err := r.getOrderDelivery(order.Id)
		if err != nil {
			return nil, err
		}
		order.Delivery = delivery

		// Get items
		items, err := r.getOrderItems(order.Id)
		if err != nil {
			return nil, err
		}
		order.Items = items

		// Get payments
		payments, err := r.getOrderPayments(order.Id)
		if err != nil {
			return nil, err
		}
		order.Payments = payments

		orders = append(orders, order)
	}

	pagedSlice := types.NewPagedSlice(args.Limit, args.Offset, total, orders)
	return &pagedSlice, nil
}

func (r *orderRepository) FindById(id string) (*aggregates.Order, error) {
	query := `
		SELECT ` + orderBaseFields + `
		FROM orders
		WHERE id = $1 AND deleted_at IS NULL
	`

	var order aggregates.Order
	var deletedAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&order.Id,
		&order.CustomerID,
		&order.RestaurantID,
		&order.Total,
		&order.TotalDiscount,
		&order.Discount,
		&order.Observation,
		&order.CreatedAt,
		&order.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get delivery
	delivery, err := r.getOrderDelivery(order.Id)
	if err != nil {
		return nil, err
	}
	order.Delivery = delivery

	// Get items
	items, err := r.getOrderItems(order.Id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	// Get payments
	payments, err := r.getOrderPayments(order.Id)
	if err != nil {
		return nil, err
	}
	order.Payments = payments

	return &order, nil
}

func (r *orderRepository) Create(order *aggregates.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create order
	_, err = tx.Exec(`
		INSERT INTO orders (
			id,
			restaurant_id,
			customer_id,
			status,
			total,
			discount,
			observation
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, order.Id, order.RestaurantID, order.CustomerID, order.Status, order.Total, order.TotalDiscount, order.Observation)
	if err != nil {
		return err
	}

	// Create delivery
	_, err = tx.Exec(`
		INSERT INTO order_deliveries (
			id,
			order_id,
			address_alias,
			address_street,
			address_number,
			address_complement,
			address_neighborhood,
			address_city,
			address_state,
			address_country,
			address_zip_code,
			address_lat,
			address_lng,
			fee,
			distance,
			average_time_minutes,
			status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`, order.Delivery.Id, order.Id, order.Delivery.Address.Alias, order.Delivery.Address.Street, order.Delivery.Address.Number,
		order.Delivery.Address.Complement, order.Delivery.Address.Neighborhood, order.Delivery.Address.City,
		order.Delivery.Address.State, order.Delivery.Address.Country, order.Delivery.Address.ZipCode,
		order.Delivery.Address.Lat, order.Delivery.Address.Lng, order.Delivery.Fee, order.Delivery.Distance,
		order.Delivery.AverageTimeMinutes, order.Delivery.Status)
	if err != nil {
		return err
	}

	// Create items
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO order_items (
				id,
				order_id,
				product_id,
				product_name,
				product_unit_price,
				quantity,
				observation,
				discount,
				item_total
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, item.Id, order.Id, item.Product.Id, item.ProductName, item.ProductUnitPrice, item.Quantity,
			item.Observation, item.Discount, item.ItemTotal)
		if err != nil {
			return err
		}
	}

	// Create payments
	for _, payment := range order.Payments {
		_, err = tx.Exec(`
			INSERT INTO order_payments (
				id,
				order_id,
				payment_method,
				amount,
				status
			) VALUES ($1, $2, $3, $4, $5)
		`, payment.Id, order.Id, payment.PaymentMethod, payment.Amount, payment.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) Update(order *aggregates.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update order
	query := `
		UPDATE orders
		SET total = $1,
			total_discount = $2,
			has_been_fully_paid_virtual = $3,
			updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
	`

	_, err = tx.Exec(
		query,
		order.Total,
		order.TotalDiscount,
		order.HasBeenFullyPaidVirtual,
		order.UpdatedAt,
		order.Id,
	)
	if err != nil {
		return err
	}

	// Update delivery
	err = r.updateOrderDelivery(tx, order.Id, &order.Delivery)
	if err != nil {
		return err
	}

	// Delete existing items and payments
	err = r.deleteOrderItems(tx, order.Id)
	if err != nil {
		return err
	}

	err = r.deleteOrderPayments(tx, order.Id)
	if err != nil {
		return err
	}

	// Create new items and payments
	for _, item := range order.Items {
		err = r.createOrderItem(tx, order.Id, &item)
		if err != nil {
			return err
		}
	}

	for _, payment := range order.Payments {
		err = r.createOrderPayment(tx, order.Id, &payment)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) Delete(id string) error {
	query := `
		UPDATE orders
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *orderRepository) FindByCustomerId(customerId string) ([]aggregates.Order, error) {
	query := `
		SELECT ` + orderBaseFields + `
		FROM orders
		WHERE customer_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []aggregates.Order
	for rows.Next() {
		var order aggregates.Order
		var deletedAt sql.NullTime
		err := rows.Scan(
			&order.Id,
			&order.CustomerID,
			&order.RestaurantID,
			&order.Total,
			&order.TotalDiscount,
			&order.Discount,
			&order.Observation,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get delivery
		delivery, err := r.getOrderDelivery(order.Id)
		if err != nil {
			return nil, err
		}
		order.Delivery = delivery

		// Get items
		items, err := r.getOrderItems(order.Id)
		if err != nil {
			return nil, err
		}
		order.Items = items

		// Get payments
		payments, err := r.getOrderPayments(order.Id)
		if err != nil {
			return nil, err
		}
		order.Payments = payments

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) FindByRestaurantId(restaurantId string) ([]aggregates.Order, error) {
	query := `
		SELECT ` + orderBaseFields + `
		FROM orders
		WHERE restaurant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []aggregates.Order
	for rows.Next() {
		var order aggregates.Order
		var deletedAt sql.NullTime
		err := rows.Scan(
			&order.Id,
			&order.CustomerID,
			&order.RestaurantID,
			&order.Total,
			&order.TotalDiscount,
			&order.Discount,
			&order.Observation,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get delivery
		delivery, err := r.getOrderDelivery(order.Id)
		if err != nil {
			return nil, err
		}
		order.Delivery = delivery

		// Get items
		items, err := r.getOrderItems(order.Id)
		if err != nil {
			return nil, err
		}
		order.Items = items

		// Get payments
		payments, err := r.getOrderPayments(order.Id)
		if err != nil {
			return nil, err
		}
		order.Payments = payments

		orders = append(orders, order)
	}

	return orders, nil
}

// Helper methods

func (r *orderRepository) getOrderDelivery(orderId string) (aggregates.OrderDelivery, error) {
	query := `
		SELECT ` + orderDeliveryFields + `
		FROM order_deliveries
		WHERE order_id = $1 AND deleted_at IS NULL
	`

	var delivery aggregates.OrderDelivery
	var addressId string
	var deletedAt sql.NullTime
	err := r.db.QueryRow(query, orderId).Scan(
		&delivery.Id,
		&orderId,
		&addressId,
		&delivery.Fee,
		&delivery.Distance,
		&delivery.AverageTimeMinutes,
		&delivery.Status,
		&delivery.DeliveryDate,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return aggregates.OrderDelivery{}, err
	}

	// Get address
	address, err := r.getAddress(addressId)
	if err != nil {
		return aggregates.OrderDelivery{}, err
	}
	delivery.Address = address

	return delivery, nil
}

func (r *orderRepository) getOrderItems(orderId string) ([]aggregates.OrderItem, error) {
	query := `
		SELECT ` + orderItemFields + `
		FROM order_items
		WHERE order_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.Query(query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []aggregates.OrderItem
	for rows.Next() {
		var item aggregates.OrderItem
		var deletedAt sql.NullTime
		err := rows.Scan(
			&item.Id,
			&orderId,
			&item.Product.Id,
			&item.ProductName,
			&item.ProductUnitPrice,
			&item.Quantity,
			&item.Observation,
			&item.Discount,
			&item.ItemTotal,
		)
		if err != nil {
			return nil, err
		}

		// Get product name
		product, err := r.getProduct(item.Product.Id)
		if err != nil {
			return nil, err
		}
		item.Product.Name = product.Name

		items = append(items, item)
	}

	return items, nil
}

func (r *orderRepository) getOrderPayments(orderId string) ([]aggregates.OrderPayment, error) {
	query := `
		SELECT ` + orderPaymentFields + `
		FROM order_payments
		WHERE order_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.Query(query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []aggregates.OrderPayment
	for rows.Next() {
		var payment aggregates.OrderPayment
		var deletedAt sql.NullTime
		err := rows.Scan(
			&payment.Id,
			&orderId,
			&payment.PaymentMethod,
			&payment.Amount,
			&payment.Status,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *orderRepository) createOrderDelivery(tx *sql.Tx, orderId string, delivery *aggregates.OrderDelivery) error {
	query := `
		INSERT INTO order_deliveries (
			id,
			order_id,
			address_id,
			fee,
			distance,
			average_time_minutes,
			status,
			delivery_date,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := tx.Exec(
		query,
		delivery.Id,
		orderId,
		delivery.Address.Id,
		delivery.Fee,
		delivery.Distance,
		delivery.AverageTimeMinutes,
		delivery.Status,
		delivery.DeliveryDate,
		delivery.CreatedAt,
		delivery.UpdatedAt,
	)
	return err
}

func (r *orderRepository) createOrderItem(tx *sql.Tx, orderId string, item *aggregates.OrderItem) error {
	query := `
		INSERT INTO order_items (
			id,
			order_id,
			product_id,
			product_unit_price,
			product_name,
			quantity,
			observation,
			discount,
			item_total,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := tx.Exec(
		query,
		item.Id,
		orderId,
		item.Product.Id,
		item.ProductUnitPrice,
		item.ProductName,
		item.Quantity,
		item.Observation,
		item.Discount,
		item.ItemTotal,
		item.CreatedAt,
		item.UpdatedAt,
	)
	return err
}

func (r *orderRepository) createOrderPayment(tx *sql.Tx, orderId string, payment *aggregates.OrderPayment) error {
	query := `
		INSERT INTO order_payments (
			id,
			order_id,
			payment_method,
			amount,
			status,
			payment_date,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.Exec(
		query,
		payment.Id,
		orderId,
		payment.PaymentMethod,
		payment.Amount,
		payment.Status,
		payment.PaymentDate,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	return err
}

func (r *orderRepository) updateOrderDelivery(tx *sql.Tx, orderId string, delivery *aggregates.OrderDelivery) error {
	query := `
		UPDATE order_deliveries
		SET address_id = $1,
			fee = $2,
			distance = $3,
			average_time_minutes = $4,
			status = $5,
			delivery_date = $6,
			updated_at = $7
		WHERE order_id = $8 AND deleted_at IS NULL
	`

	_, err := tx.Exec(
		query,
		delivery.Address.Id,
		delivery.Fee,
		delivery.Distance,
		delivery.AverageTimeMinutes,
		delivery.Status,
		delivery.DeliveryDate,
		delivery.UpdatedAt,
		orderId,
	)
	return err
}

func (r *orderRepository) deleteOrderItems(tx *sql.Tx, orderId string) error {
	query := `
		UPDATE order_items
		SET deleted_at = $1
		WHERE order_id = $2 AND deleted_at IS NULL
	`

	_, err := tx.Exec(query, time.Now(), orderId)
	return err
}

func (r *orderRepository) deleteOrderPayments(tx *sql.Tx, orderId string) error {
	query := `
		UPDATE order_payments
		SET deleted_at = $1
		WHERE order_id = $2 AND deleted_at IS NULL
	`

	_, err := tx.Exec(query, time.Now(), orderId)
	return err
}

func (r *orderRepository) getAddress(id string) (types.Address, error) {
	query := `
		SELECT id, alias, street, number, complement, neighborhood, city, state, country, zip_code, lat, lng
		FROM addresses
		WHERE id = $1 AND deleted_at IS NULL
	`

	var address types.Address
	err := r.db.QueryRow(query, id).Scan(
		&address.Id,
		&address.Alias,
		&address.Street,
		&address.Number,
		&address.Complement,
		&address.Neighborhood,
		&address.City,
		&address.State,
		&address.Country,
		&address.ZipCode,
		&address.Lat,
		&address.Lng,
	)
	return address, err
}

func (r *orderRepository) getProduct(id string) (aggregates.PartialProduct, error) {
	query := `
		SELECT id, name
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	var product aggregates.PartialProduct
	err := r.db.QueryRow(query, id).Scan(
		&product.Id,
		&product.Name,
	)
	return product, err
}
