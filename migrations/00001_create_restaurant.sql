CREATE TABLE restaurants (
    id CHAR(36) PRIMARY KEY,
    trade_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    cnpj CHAR(14) NOT NULL UNIQUE,
    contact_phone VARCHAR(20),
    whatsapp_phone VARCHAR(20),
    email VARCHAR(255),
    accepted_payment_methods JSON NOT NULL,
    slug VARCHAR(255) UNIQUE,
    address_id CHAR(36) NOT NULL,
    show_cnpj_in_receipt BOOLEAN,
    delivery_enabled BOOLEAN,
    delivery_fee_per_km INT,
    delivery_minimum_order_value INT,
    delivery_max_radius_km INT,
    delivery_average_time_minutes INT,
    ecommerce_enabled BOOLEAN,
    ecommerce_minimum_order_value INT,
    ecommerce_delivery_fee_per_km INT,
    ecommerce_delivery_max_radius_km INT,
    ecommerce_delivery_average_time_minutes INT,
    customer_post_paid_orders_enabled BOOLEAN,
    customer_post_paid_orders_minimum_order_value INT,
    customer_post_paid_orders_average_time_minutes INT,
    customer_post_paid_orders_delivery_fee_per_km INT,
    customer_post_paid_orders_delivery_max_radius_km INT,
    customer_post_paid_orders_delivery_average_time_minutes INT,
    logo_url VARCHAR(255),
    banner_url VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at DATETIME NULL,
    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_restaurants_slug ON restaurants(slug); 
CREATE INDEX idx_restaurants_active ON restaurants(active);
CREATE INDEX idx_restaurants_created_at ON restaurants(created_at);

CREATE TABLE restaurant_pix_keys(
    id CHAR(36) PRIMARY KEY,
    restaurant_id CHAR(36) NOT NULL,
    pix_key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE ON UPDATE CASCADE
);