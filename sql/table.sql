-- ============================================================
-- 酒店预订管理系统 — 数据库完整定义
-- 上海电力大学 数据库原理课程设计 | 选题四：酒店预订管理系统
-- PostgreSQL 18（uuidv7 为内置函数）
-- ============================================================
-- ============================================================
-- 1. 自定义类型
-- ============================================================
CREATE TYPE user_role AS ENUM ('customer', 'vip', 'hotel_manager', 'admin');
CREATE TYPE order_status AS ENUM ('pending', 'booked', 'checked_in', 'completed', 'cancelled');
-- ============================================================
-- 2. 表定义
-- ============================================================
-- 2.1 VIP 等级定义表
CREATE TABLE vip_level_1718 (
    level SMALLINT PRIMARY KEY CHECK (level >= 0),
    level_name TEXT NOT NULL,
    min_points INT NOT NULL CHECK (min_points >= 0),
    discount_rate DECIMAL(3, 2) NOT NULL CHECK (
        discount_rate BETWEEN 0 AND 1
    )
);


-- 2.2 用户表
CREATE TABLE user_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    username TEXT NOT NULL,
    password VARCHAR(255) NOT NULL,
    -- 密码哈希值
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

-- 2.3 地区表（省/市/区 三级层次）
CREATE TABLE region_1718 (
    id SERIAL PRIMARY KEY,
    region_name TEXT NOT NULL,
    parents_id INT REFERENCES region_1718 (id)
);
-- 2.4 酒店表
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
-- 2.4 酒店图片表
CREATE TABLE hotel_image_1718 (
    hotel_id UUID NOT NULL REFERENCES hotel_1718 (id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    PRIMARY KEY (hotel_id, image_url)
);
-- 2.5 客房表（房型定义）
CREATE TABLE room_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    hotel_id UUID NOT NULL REFERENCES hotel_1718 (id) ON DELETE CASCADE,
    type_name TEXT NOT NULL,
    -- 房型名称：标准间、大床房、豪华套房等
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
-- 2.6 客房图片表
CREATE TABLE room_image_1718 (
    room_id UUID NOT NULL REFERENCES room_1718 (id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    PRIMARY KEY (room_id, image_url)
);
-- 2.7 客房设施表
CREATE TABLE room_facility_1718 (
    room_id UUID NOT NULL REFERENCES room_1718 (id) ON DELETE CASCADE,
    facility_name TEXT NOT NULL,
    -- WiFi、空调、电视等
    PRIMARY KEY (room_id, facility_name)
);

-- 2.8 人员表（入住人身份信息）
--     性别、年龄可通过身份证号推导，参见视图 view_person_info
CREATE TABLE person_1718 (
    id_card VARCHAR(18) PRIMARY KEY CHECK (id_card ~ '^\d{17}[\dXx]$'),
    name TEXT NOT NULL,
    phone VARCHAR(20)
);

-- 2.10 订单表
CREATE TABLE order_1718 (
    id UUID PRIMARY KEY DEFAULT uuidv7 (),
    user_id UUID NOT NULL REFERENCES user_1718 (id),
    room_id UUID NOT NULL REFERENCES room_1718 (id),
    quantity INT NOT NULL CHECK (quantity > 0),
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL CHECK (total_price >= 0),
    discount DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (discount >= 0),
    actual_price DECIMAL(10, 2) NOT NULL CHECK (actual_price >= 0),
    status order_status NOT NULL DEFAULT 'pending',
    create_at TIMESTAMPTZ DEFAULT now (),
    update_at TIMESTAMPTZ DEFAULT now (),
    CHECK (check_out_date > check_in_date)
);
-- 2.11 入住人员关联表（订单 ↔ 入住人，多对多）
CREATE TABLE order_guest_1718 (
    order_id UUID NOT NULL REFERENCES order_1718 (id) ON DELETE CASCADE,
    id_card VARCHAR(18) NOT NULL REFERENCES person_1718 (id_card),
    PRIMARY KEY (order_id, id_card)
);
-- 2.12 评价表
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
    UNIQUE (user_id, order_id) -- 每笔订单每个用户只能评价一次
);

-- ============================================================
-- 3. 索引
-- ============================================================
-- 用户表
CREATE INDEX idx_user_phone ON user_1718 (phone);
CREATE INDEX idx_user_role ON user_1718 (role);
CREATE INDEX idx_user_points ON user_1718 (points DESC);
-- 地区表
CREATE INDEX idx_region_parents ON region_1718 (parents_id);
-- 酒店表
CREATE INDEX idx_hotel_region ON hotel_1718 (region_id);
CREATE INDEX idx_hotel_name ON hotel_1718 (hotel_name);
CREATE INDEX idx_hotel_star ON hotel_1718 (star_level);
-- 客房表
CREATE INDEX idx_room_hotel ON room_1718 (hotel_id);
CREATE INDEX idx_room_type ON room_1718 (type_name);
CREATE INDEX idx_room_price ON room_1718 (price);
CREATE INDEX idx_room_available ON room_1718 (available_quantity);
-- 订单表
CREATE INDEX idx_order_user ON order_1718 (user_id);
CREATE INDEX idx_order_room ON order_1718 (room_id);
CREATE INDEX idx_order_status ON order_1718 (status);
CREATE INDEX idx_order_dates ON order_1718 (check_in_date, check_out_date);
CREATE INDEX idx_order_create ON order_1718 (create_at);
-- 评价表
CREATE INDEX idx_review_hotel ON review_1718 (hotel_id);
CREATE INDEX idx_review_user ON review_1718 (user_id);
CREATE INDEX idx_review_rating ON review_1718 (rating);
-- 入住人员
CREATE INDEX idx_order_guest_id_card ON order_guest_1718 (id_card);
-- ============================================================
-- 4. 视图
-- ============================================================
-- 4.1 按城市 → 区域 → 酒店查看所有客房详细信息
CREATE VIEW view_room_details AS
SELECT r_province.region_name AS province,
    r_city.region_name AS city,
    r_district.region_name AS district,
    h.id AS hotel_id,
    h.hotel_name,
    h.address,
    h.star_level,
    h.telephone AS hotel_telephone,
    rm.id AS room_id,
    rm.type_name,
    rm.total_quantity,
    rm.available_quantity,
    rm.price,
    rm.weekend_price,
    rm.description AS room_description,
    STRING_AGG (
        DISTINCT rf.facility_name,
        ', '
        ORDER BY rf.facility_name
    ) AS facilities,
    rev_stats.avg_rating,
    rev_stats.review_count
FROM hotel_1718 h
    JOIN region_1718 r_district ON h.region_id = r_district.id
    LEFT JOIN region_1718 r_city ON r_district.parents_id = r_city.id
    LEFT JOIN region_1718 r_province ON r_city.parents_id = r_province.id
    JOIN room_1718 rm ON rm.hotel_id = h.id
    AND rm.status = 1
    LEFT JOIN room_facility_1718 rf ON rf.room_id = rm.id
    LEFT JOIN (
        SELECT hotel_id,
            AVG(rating)::DECIMAL(3, 2) AS avg_rating,
            COUNT(*) AS review_count
        FROM review_1718
        GROUP BY hotel_id
    ) rev_stats ON rev_stats.hotel_id = h.id
WHERE h.status = 1
GROUP BY r_province.region_name,
    r_city.region_name,
    r_district.region_name,
    h.id,
    h.hotel_name,
    h.address,
    h.star_level,
    h.telephone,
    rm.id,
    rm.type_name,
    rm.total_quantity,
    rm.available_quantity,
    rm.price,
    rm.weekend_price,
    rm.description,
    rev_stats.avg_rating,
    rev_stats.review_count;

-- 4.2 从身份证号推导性别与年龄
--     身份证第 7-14 位为出生日期（YYYYMMDD），第 17 位奇数为男、偶数为女
CREATE VIEW view_person_info AS
SELECT id_card,
       name,
       phone,
       TO_DATE(SUBSTRING(id_card, 7, 8), 'YYYYMMDD') AS birth_date,
       CASE
           WHEN SUBSTRING(id_card, 17, 1) ~ '\d'
                AND SUBSTRING(id_card, 17, 1)::INT % 2 = 1 THEN '男'
           WHEN SUBSTRING(id_card, 17, 1) ~ '\d'
                AND SUBSTRING(id_card, 17, 1)::INT % 2 = 0 THEN '女'
           ELSE NULL
       END AS gender,
       DATE_PART('year', AGE(TO_DATE(SUBSTRING(id_card, 7, 8), 'YYYYMMDD')))::INT AS age
FROM person_1718;
-- ============================================================
-- 5. 触发器
-- ============================================================

-- 评价 user_id 必须与订单的预订人一致
CREATE OR REPLACE FUNCTION fn_validate_review_user ()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.user_id != (SELECT user_id FROM order_1718 WHERE id = NEW.order_id) THEN
        RAISE EXCEPTION '评价用户与订单预订人不一致';
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_review_user_match
    BEFORE INSERT OR UPDATE OF user_id, order_id
    ON review_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_validate_review_user ();

-- 通用：更新行时自动刷新 update_at
CREATE OR REPLACE FUNCTION fn_set_update_at ()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    NEW.update_at = now ();
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_user_update_at
    BEFORE UPDATE ON user_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_set_update_at ();

CREATE TRIGGER trg_hotel_update_at
    BEFORE UPDATE ON hotel_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_set_update_at ();

CREATE TRIGGER trg_room_update_at
    BEFORE UPDATE ON room_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_set_update_at ();

CREATE TRIGGER trg_order_update_at
    BEFORE UPDATE ON order_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_set_update_at ();

CREATE TRIGGER trg_review_update_at
    BEFORE UPDATE ON review_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_set_update_at ();

-- 客房可预订数量自动更新
CREATE OR REPLACE FUNCTION fn_update_room_quantity ()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
DECLARE
    v_available INT;
BEGIN
    -- 预订或入住 → 扣减可预订数量
    IF NEW.status IN ('booked', 'checked_in')
       AND (TG_OP = 'INSERT' OR OLD.status = 'pending')
    THEN
        SELECT available_quantity INTO v_available
        FROM room_1718 WHERE id = NEW.room_id;

        IF v_available < NEW.quantity THEN
            RAISE EXCEPTION '客房库存不足：可用 % 间，请求 % 间', v_available, NEW.quantity;
        END IF;

        UPDATE room_1718
           SET available_quantity = available_quantity - NEW.quantity,
               update_at = now ()
         WHERE id = NEW.room_id;

    -- 取消 → 恢复可预订数量（仅限已确认但未完成的订单）
    ELSIF NEW.status = 'cancelled'
          AND TG_OP = 'UPDATE'
          AND OLD.status IN ('booked', 'checked_in')
    THEN
        UPDATE room_1718
           SET available_quantity = available_quantity + OLD.quantity,
               update_at = now ()
         WHERE id = OLD.room_id;
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_order_room_quantity
    AFTER INSERT OR UPDATE OF status
    ON order_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_room_quantity ();

-- 5.2 用户积分自动更新（订单一完成增减积分）
CREATE OR REPLACE FUNCTION fn_update_user_points ()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    -- 完成 → 增加积分
    IF NEW.status = 'completed'
       AND (OLD.status IS NULL OR OLD.status != 'completed')
    THEN
        UPDATE user_1718
           SET points = points + FLOOR(NEW.actual_price)::INT,
               update_at = now ()
         WHERE id = NEW.user_id;
    -- 离开 completed 状态 → 扣回积分
    ELSIF TG_OP = 'UPDATE'
          AND OLD.status = 'completed'
          AND NEW.status != 'completed'
    THEN
        UPDATE user_1718
           SET points = GREATEST(0, points - FLOOR(OLD.actual_price)::INT),
               update_at = now ()
         WHERE id = OLD.user_id;
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_order_user_points
    AFTER INSERT OR UPDATE OF status
    ON order_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_user_points ();
-- ============================================================
-- 6. 函数（用于查询和统计）
-- ============================================================

-- 6.1 按用户信息查询订单
--     支持按用户名、身份证号模糊匹配，返回该用户所有订单
CREATE OR REPLACE FUNCTION fn_query_orders_by_user (
    p_username TEXT DEFAULT NULL,
    p_id_card  VARCHAR(18) DEFAULT NULL
)
RETURNS TABLE (
    order_id       UUID,
    check_in_date  DATE,
    check_out_date DATE,
    status         VARCHAR(20),
    total_price    DECIMAL(10, 2),
    discount       DECIMAL(10, 2),
    actual_price   DECIMAL(10, 2),
    create_at      TIMESTAMPTZ,
    hotel_name     TEXT,
    room_type      TEXT,
    quantity       INT
)
LANGUAGE plpgsql AS $$
DECLARE
    v_user_ids UUID[];
BEGIN
    -- 通过身份证号查找关联的所有用户
    IF p_id_card IS NOT NULL THEN
        SELECT ARRAY_AGG(DISTINCT o.user_id) INTO v_user_ids
        FROM order_1718 o
        JOIN order_guest_1718 og ON og.order_id = o.id
        WHERE og.id_card = p_id_card;

        IF v_user_ids IS NULL THEN
            RAISE NOTICE '未找到身份证号 % 关联的订单', p_id_card;
            RETURN;
        END IF;
    END IF;

    -- 通过用户名模糊查找所有匹配用户
    IF p_username IS NOT NULL THEN
        SELECT ARRAY_AGG(id) INTO v_user_ids
        FROM (
            SELECT id FROM user_1718 WHERE username ILIKE '%' || p_username || '%'
            UNION
            SELECT UNNEST(v_user_ids)
        ) t;
    END IF;

    IF v_user_ids IS NULL THEN
        RAISE NOTICE '未找到匹配的用户';
        RETURN;
    END IF;

    RETURN QUERY
    SELECT o.id,
           o.check_in_date,
           o.check_out_date,
           o.status,
           o.total_price,
           o.discount,
           o.actual_price,
           o.create_at,
           h.hotel_name,
           rm.type_name,
           o.quantity::INT
    FROM order_1718 o
    JOIN room_1718 rm ON rm.id = o.room_id
    JOIN hotel_1718 h  ON h.id = rm.hotel_id
    WHERE o.user_id = ANY (v_user_ids)
    ORDER BY o.create_at DESC;
END;
$$;

-- 6.2 统计指定酒店在指定时间段的客房入住率
--     入住率 = 已预订房晚数 / (总客房数 × 统计天数) × 100
--     房晚数按订单在统计期内的实际重叠天数计算
CREATE OR REPLACE FUNCTION fn_room_occupancy_rate (
    p_hotel_id   UUID,
    p_start_date DATE,
    p_end_date   DATE
)
RETURNS TABLE (
    room_id                UUID,
    room_type              TEXT,
    total_quantity         INT,
    available_quantity     INT,
    booked_nights          BIGINT,
    occupancy_rate_percent DECIMAL(6, 2),
    price                  DECIMAL(10, 2),
    total_revenue          DECIMAL(10, 2)
)
LANGUAGE plpgsql AS $$
DECLARE
    v_total_days INT;
BEGIN
    IF p_start_date > p_end_date THEN
        RAISE EXCEPTION '开始日期不能晚于结束日期';
    END IF;

    v_total_days := p_end_date - p_start_date + 1;

    RETURN QUERY
    SELECT rm.id,
           rm.type_name,
           rm.total_quantity,
           rm.available_quantity,
           COALESCE(
               SUM(
                   o.quantity * GREATEST(
                       0,
                       LEAST(o.check_out_date, p_end_date)
                       - GREATEST(o.check_in_date, p_start_date)
                   )
               ),
               0
           )::BIGINT,
           CASE
               WHEN rm.total_quantity > 0 AND v_total_days > 0
               THEN ROUND(
                        (
                            COALESCE(
                                SUM(
                                    o.quantity * GREATEST(
                                        0,
                                        LEAST(o.check_out_date, p_end_date)
                                        - GREATEST(o.check_in_date, p_start_date)
                                    )
                                ),
                                0
                            )::DECIMAL
                            / (rm.total_quantity * v_total_days)
                        ) * 100,
                        2
                    )
               ELSE 0
           END,
           rm.price,
            COALESCE(
                SUM(
                    o.actual_price
                    * GREATEST(
                          0,
                          LEAST(o.check_out_date, p_end_date)
                          - GREATEST(o.check_in_date, p_start_date)
                      )
                    / NULLIF(o.check_out_date - o.check_in_date, 0)
                ),
                0
            )
    FROM room_1718 rm
    LEFT JOIN order_1718 o ON o.room_id = rm.id
        AND o.status IN ('booked', 'checked_in', 'completed')
        AND o.check_in_date < p_end_date
        AND o.check_out_date > p_start_date
    WHERE rm.hotel_id = p_hotel_id AND rm.status = 1
    GROUP BY rm.id, rm.type_name, rm.total_quantity, rm.available_quantity, rm.price
    ORDER BY rm.type_name;
END;
$$;