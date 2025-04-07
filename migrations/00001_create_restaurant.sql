CREATE TABLE restaurants (
    id CHAR(36) PRIMARY KEY,
    trade_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    cnpj CHAR(14) NOT NULL UNIQUE,
    contact_phone VARCHAR(20),
    whatsapp_phone VARCHAR(20),
    email VARCHAR(255),
    slug VARCHAR(255) UNIQUE,
    address_street VARCHAR(255),
    address_number VARCHAR(20),
    address_complement VARCHAR(255),
    address_neighborhood VARCHAR(255),
    address_city VARCHAR(100),
    address_state VARCHAR(100),
    address_zip_code VARCHAR(20),
    address_country VARCHAR(100),
    address_lat DOUBLE,
    address_lng DOUBLE,
    show_cnpj_in_receipt BOOLEAN,
    delivery_enabled BOOLEAN,
    delivery_fee_per_km INT,
    delivery_minimum_order_value INT,
    delivery_max_radius_delivery INT,
    delivery_average_time_minutes INT,
    ecommerce_minimum_order_value INT,
    ecommerce_acquired BOOLEAN,
    ecommerce_acquired_at TIMESTAMP NULL,
    ecommerce_online BOOLEAN,
    customer_post_paid_orders_acquired BOOLEAN,
    customer_post_paid_orders_acquired_at TIMESTAMP NULL,
    customer_post_paid_orders_enabled BOOLEAN,
    logo_url TEXT,
    banner_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE INDEX idx_restaurants_slug ON restaurants(slug);
CREATE INDEX idx_restaurants_active ON restaurants(active);
CREATE INDEX idx_restaurants_created_at ON restaurants(created_at);
