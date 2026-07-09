-- 酒店预订管理系统 — 数据库定义
-- PostgreSQL 18
-- 自定义类型
CREATE TYPE user_role AS ENUM ('customer', 'vip', 'hotel_manager', 'admin');
CREATE TYPE order_status AS ENUM (
    'pending',
    'booked',
    'checked_in',
    'completed',
    'cancelled'
);
CREATE TYPE education_level AS ENUM ('小学', '初中', '高中', '中专', '大专', '本科', '硕士', '博士', '其他');
-- 表定义
-- VIP 等级定义
CREATE TABLE vip_level_1718 (
    level SMALLINT PRIMARY KEY CHECK (level >= 0),
    level_name TEXT NOT NULL,
    min_points INT NOT NULL CHECK (min_points >= 0),
    discount_rate DECIMAL(3, 2) NOT NULL CHECK (
        discount_rate BETWEEN 0 AND 1
    )
);
-- 用户
CREATE TABLE user_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    username TEXT NOT NULL,
    password VARCHAR(255) NOT NULL,
    -- bcrypt 哈希
    phone VARCHAR(20),
    email TEXT,
    real_name TEXT,
    id_card VARCHAR(18) CHECK (
        id_card ~ '^\d{17}[\dXx]$'
        OR id_card IS NULL
    ),
    role user_role NOT NULL DEFAULT 'customer',
    points INT NOT NULL DEFAULT 0 CHECK (points >= 0),
    vip_level SMALLINT NOT NULL DEFAULT 0 REFERENCES vip_level_1718 (level),
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    status SMALLINT NOT NULL DEFAULT 1
);
-- 地区（省/市/区 三级层次）
CREATE TABLE region_1718 (
    id SERIAL PRIMARY KEY,
    region_name TEXT NOT NULL,
    parents_id INT REFERENCES region_1718 (id)
);
-- 酒店
CREATE TABLE hotel_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    hotel_name TEXT NOT NULL,
    region_id INT NOT NULL REFERENCES region_1718 (id),
    address TEXT NOT NULL,
    telephone VARCHAR(20) NOT NULL,
    star_level SMALLINT CHECK (
        star_level BETWEEN 1 AND 5
    ),
    description TEXT,
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    status SMALLINT NOT NULL DEFAULT 1
);
-- 酒店图片
CREATE TABLE hotel_image_1718 (
    hotel_id UUID NOT NULL REFERENCES hotel_1718 (id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    PRIMARY KEY (hotel_id, image_url)
);
-- 客房（房型定义）
CREATE TABLE room_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    hotel_id UUID NOT NULL REFERENCES hotel_1718 (id) ON DELETE CASCADE,
    type_name TEXT NOT NULL,
    total_quantity INT NOT NULL CHECK (total_quantity > 0),
    available_quantity INT NOT NULL CHECK (available_quantity >= 0),
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    weekend_price DECIMAL(10, 2) CHECK (weekend_price > 0),
    description TEXT,
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    status SMALLINT NOT NULL DEFAULT 1,
    CHECK (available_quantity <= total_quantity)
);
-- 客房图片
CREATE TABLE room_image_1718 (
    room_id UUID NOT NULL REFERENCES room_1718 (id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    PRIMARY KEY (room_id, image_url)
);
-- 客房设施
CREATE TABLE room_facility_1718 (
    room_id UUID NOT NULL REFERENCES room_1718 (id) ON DELETE CASCADE,
    facility_name TEXT NOT NULL,
    PRIMARY KEY (room_id, facility_name)
);
-- 人员（入住人身份信息）
CREATE TABLE person_1718 (
    id_card CHAR(18) PRIMARY KEY CHECK (id_card ~ '^\d{17}[\dXx]$'),
    name TEXT NOT NULL,
    phone VARCHAR(20) CHECK (phone ~ '^\+?[0-9\-]+$'),
    occupation TEXT,
    education education_level,
    income numrange
);
-- 收入范围 GiST 索引（支持范围查询）
CREATE INDEX idx_person_1718_income ON person_1718 USING GIST (income);
-- 订单
CREATE TABLE order_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    user_id UUID NOT NULL REFERENCES user_1718 (id),
    -- 下订单的用户不是入住人
    room_id UUID NOT NULL REFERENCES room_1718 (id),
    quantity INT NOT NULL CHECK (quantity > 0),
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL CHECK (total_price >= 0),
    -- discount DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (discount >= 0),
    actual_price DECIMAL(10, 2) NOT NULL CHECK (actual_price >= 0),
    status order_status NOT NULL DEFAULT 'pending',
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    CHECK (check_out_date > check_in_date)
);
-- 入住人关联（订单 ↔ 入住人，多对多）
CREATE TABLE order_guest_1718 (
    order_id UUID NOT NULL REFERENCES order_1718 (id) ON DELETE CASCADE,
    id_card CHAR(18) NOT NULL REFERENCES person_1718 (id_card),
    PRIMARY KEY (order_id, id_card)
);
-- 评价
CREATE TABLE review_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    user_id UUID NOT NULL REFERENCES user_1718 (id),
    order_id UUID NOT NULL REFERENCES order_1718 (id),
    hotel_id UUID NOT NULL REFERENCES hotel_1718 (id),
    rating SMALLINT NOT NULL CHECK (
        rating BETWEEN 1 AND 5
    ),
    content TEXT,
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    UNIQUE (user_id, order_id)
);
-- JWT Token 黑名单
CREATE TABLE jwt_blacklist_1718 (
    jti TEXT PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL
);
-- 索引
-- 用户
CREATE INDEX idx_user_phone_1718 ON user_1718 (phone);
CREATE INDEX idx_user_role_1718 ON user_1718 (role);
CREATE INDEX idx_user_points_1718 ON user_1718 (points DESC);
-- 地区
CREATE INDEX idx_region_parents_1718 ON region_1718 (parents_id);
-- 酒店
CREATE INDEX idx_hotel_region_1718 ON hotel_1718 (region_id);
CREATE INDEX idx_hotel_name_1718 ON hotel_1718 (hotel_name);
CREATE INDEX idx_hotel_star_1718 ON hotel_1718 (star_level);
-- 客房
CREATE INDEX idx_room_hotel_1718 ON room_1718 (hotel_id);
CREATE INDEX idx_room_type_1718 ON room_1718 (type_name);
CREATE INDEX idx_room_price_1718 ON room_1718 (price);
CREATE INDEX idx_room_available_1718 ON room_1718 (available_quantity);
-- 订单
CREATE INDEX idx_order_user_1718 ON order_1718 (user_id);
CREATE INDEX idx_order_room_1718 ON order_1718 (room_id);
CREATE INDEX idx_order_status_1718 ON order_1718 (status);
CREATE INDEX idx_order_dates_1718 ON order_1718 (check_in_date, check_out_date);
CREATE INDEX idx_order_create_1718 ON order_1718 (create_at);
-- 评价
CREATE INDEX idx_review_hotel_1718 ON review_1718 (hotel_id);
CREATE INDEX idx_review_user_1718 ON review_1718 (user_id);
CREATE INDEX idx_review_rating_1718 ON review_1718 (rating);
-- 入住人
CREATE INDEX idx_order_guest_id_card_1718 ON order_guest_1718 (id_card);
-- JWT
CREATE INDEX idx_blacklist_expires_1718 ON jwt_blacklist_1718 (expires_at);