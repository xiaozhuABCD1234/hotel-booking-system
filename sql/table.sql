CREATE TYPE user_role AS ENUM ('customer', 'vip', 'hotel_manager', 'admin');
-- 用户有customer(普通用户)、vip（vip用户）、hotel_manager（酒店管理员）、admin(平台管理员) 
CREATE TABLE user_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    -- postgresql 18 集成 uuidv7
    username TEXT NOT NULL,
    password VARCHAR(255) NOT NULL,
    -- 密码保存hash
    role user_role NOT NULL DEFAULT 'customer',
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    status SMALLINT NOT NULL DEFAULT 1
);

CREATE TABLE region_1718 (
    id SERIAL PRIMARY KEY,
    region_name TEXT NOT NULL,
    parents_id INT REFERENCES region_1718 (id)
);

CREATE TABLE hotel_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    hotel_name TEXT NOT NULL,
    region_id INT REFERENCES region_1718 (id),
    address TEXT NOT NULL,
    telephone VARCHAR(20) NOT NULL,
    star_level SMALLINT CHECK(star_level BETWEEN 1 AND 5),
    description TEXT,
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    status SMALLINT NOT NULL DEFAULT 1
);