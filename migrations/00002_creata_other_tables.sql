CREATE TABLE addresses(
    id CHAR(36) PRIMARY KEY,
    alias VARCHAR(255),
    street VARCHAR(255) NOT NULL,
    number VARCHAR(255) NOT NULL,
    complement VARCHAR(255),
    neighborhood VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(255),
    zip_code VARCHAR(255),
    country VARCHAR(255),
    lat FLOAT,
    lng FLOAT
);

CREATE TABLE users(
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    pwd_hash VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    role VARCHAR(255) NOT NULL,
    permissions JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

CREATE TABLE dishes(
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    picture_url VARCHAR(255),
    restaurant_id CHAR(36) NOT NULL,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE menus(
    id CHAR(36) PRIMARY KEY,
    picture_url VARCHAR(255),
    restaurant_id CHAR(36) NOT NULL,
    offer_date DATE NOT NULL,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE menu_items(
    id CHAR(36) PRIMARY KEY,
    menu_id CHAR(36) NOT NULL,
    dish_id CHAR(36) NOT NULL,
    additional_price DECIMAL(10, 2) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    can_be_used_as_additional BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE categories(
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    picture_url VARCHAR(255),
    restaurant_id CHAR(36) NOT NULL,
    priority INT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE products(
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    sales_price DECIMAL(10, 2) NOT NULL,
    cost_price DECIMAL(10, 2) NOT NULL,
    picture_url VARCHAR(255),
    dish_type_map JSON NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    category_id CHAR(36) NOT NULL,
    restaurant_id CHAR(36) NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE customers(
    id CHAR(36) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    contact_phone VARCHAR(255),
    address_id CHAR(36) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE orders(
    id CHAR(36) PRIMARY KEY,
    customer_id CHAR(36) NOT NULL,
    restaurant_id CHAR(36) NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    total_discount DECIMAL(10, 2) NOT NULL,
    has_been_fully_paid_virtual BOOLEAN GENERATED ALWAYS AS ((total - total_discount) = (SELECT SUM(amount) FROM order_payments WHERE order_id = id)) STORED,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (delivery_address_id) REFERENCES addresses(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE order_deliveries(
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    address_id CHAR(36) NOT NULL,
    fee DECIMAL(10, 2) NOT NULL,
    distance DECIMAL(10, 2) NOT NULL,
    average_time_minutes INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    delivery_date DATETIME NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (delivery_address_id) REFERENCES addresses(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_order_deliveries_order_id ON order_deliveries(order_id);

CREATE TABLE order_items(
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    product_unit_price DECIMAL(10, 2) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    observation VARCHAR(255),
    discount DECIMAL(10, 2) NOT NULL,
    item_total DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE lunchbox_order_items(
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    allowed_protein_count INT NOT NULL,
    allowed_side_count INT NOT NULL,
    allowed_accompaniment_count INT NOT NULL,
    wants_flatware BOOLEAN NOT NULL,
    observation VARCHAR(255),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE lunchbox_selected_menu_items(
    id CHAR(36) PRIMARY KEY,
    lunchbox_order_item_id CHAR(36) NOT NULL,
    menu_item_id CHAR(36) NOT NULL,
    quantity INT NOT NULL,
    observation VARCHAR(255),
    item_total DECIMAL(10, 2) NOT NULL,
    is_additional BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (lunchbox_order_item_id) REFERENCES lunchbox_order_items(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE order_payments(
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    payment_method VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(255) NOT NULL,
    payment_date DATETIME NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE
);
