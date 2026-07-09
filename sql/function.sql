-- 酒店预订管理系统 — 存储函数/过程
-- PostgreSQL 18
-- 按用户信息查询订单
-- 支持按用户名、身份证号模糊匹配
CREATE OR REPLACE FUNCTION fn_query_orders_by_user_1718(
        p_username TEXT DEFAULT NULL,
        p_id_card CHAR(18) DEFAULT NULL
    ) RETURNS TABLE (
        order_id UUID,
        check_in_date DATE,
        check_out_date DATE,
        status VARCHAR(20),
    total_price DECIMAL(10, 2),
    actual_price DECIMAL(10, 2),
        create_at TIMESTAMPTZ,
        hotel_name TEXT,
        room_type TEXT,
        quantity INT
    ) LANGUAGE plpgsql AS $$
DECLARE v_user_ids UUID [];
BEGIN IF p_id_card IS NOT NULL THEN
SELECT ARRAY_AGG(DISTINCT o.user_id) INTO v_user_ids
FROM order_1718 o
    JOIN order_guest_1718 og ON og.order_id = o.id
WHERE og.id_card = p_id_card;
IF v_user_ids IS NULL THEN RAISE NOTICE '未找到身份证号 % 关联的订单',
p_id_card;
RETURN;
END IF;
END IF;
IF p_username IS NOT NULL THEN
SELECT ARRAY_AGG(id) INTO v_user_ids
FROM (
        SELECT id
        FROM user_1718
        WHERE username ILIKE '%' || p_username || '%'
        UNION
        SELECT UNNEST(v_user_ids)
    ) t;
END IF;
IF v_user_ids IS NULL THEN RAISE NOTICE '未找到匹配的用户';
RETURN;
END IF;
RETURN QUERY
SELECT o.id,
    o.check_in_date,
    o.check_out_date,
    o.status,
    o.total_price,
    o.actual_price,
    o.create_at,
    h.hotel_name,
    rm.type_name,
    o.quantity::INT
FROM order_1718 o
    JOIN room_1718 rm ON rm.id = o.room_id
    JOIN hotel_1718 h ON h.id = rm.hotel_id
WHERE o.user_id = ANY (v_user_ids)
ORDER BY o.create_at DESC;
END;
$$;
-- 统计指定酒店在指定时间段的客房入住率
-- 入住率 = 已预订房晚数 / (总客房数 × 统计天数) × 100
CREATE OR REPLACE FUNCTION fn_room_occupancy_rate_1718(
        p_hotel_id UUID,
        p_start_date DATE,
        p_end_date DATE
    ) RETURNS TABLE (
        room_id UUID,
        room_type TEXT,
        total_quantity INT,
        available_quantity INT,
        booked_nights BIGINT,
        occupancy_rate_percent DECIMAL(6, 2),
        price DECIMAL(10, 2),
        total_revenue DECIMAL(10, 2)
    ) LANGUAGE plpgsql AS $$
DECLARE v_total_days INT;
BEGIN IF p_start_date > p_end_date THEN RAISE EXCEPTION '开始日期不能晚于结束日期';
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
                LEAST(o.check_out_date, p_end_date) - GREATEST(o.check_in_date, p_start_date)
            )
        ),
        0
    )::BIGINT,
    CASE
        WHEN rm.total_quantity > 0
        AND v_total_days > 0 THEN ROUND(
            (
                COALESCE(
                    SUM(
                        o.quantity * GREATEST(
                            0,
                            LEAST(o.check_out_date, p_end_date) - GREATEST(o.check_in_date, p_start_date)
                        )
                    ),
                    0
                )::DECIMAL / (rm.total_quantity * v_total_days)
            ) * 100,
            2
        )
        ELSE 0
    END,
    rm.price,
    COALESCE(
        SUM(
            o.actual_price * GREATEST(
                0,
                LEAST(o.check_out_date, p_end_date) - GREATEST(o.check_in_date, p_start_date)
            ) / NULLIF(o.check_out_date - o.check_in_date, 0)
        ),
        0
    )
FROM room_1718 rm
    LEFT JOIN order_1718 o ON o.room_id = rm.id
    AND o.status IN ('booked', 'checked_in', 'completed')
    AND o.check_in_date < p_end_date
    AND o.check_out_date > p_start_date
WHERE rm.hotel_id = p_hotel_id
    AND rm.status = 1
GROUP BY rm.id,
    rm.type_name,
    rm.total_quantity,
    rm.available_quantity,
    rm.price
ORDER BY rm.type_name;
END;
$$;
-- 确保入住人存在于人员表，不存在则插入，已存在则更新 phone/occupation/education/income
-- 下订单前调用，避免 order_guest_1718 外键约束失败
CREATE OR REPLACE PROCEDURE sp_ensure_person_1718(
        p_id_card CHAR(18),
        p_name TEXT,
        p_phone VARCHAR(20) DEFAULT NULL
    ) LANGUAGE plpgsql AS $$ BEGIN
INSERT INTO person_1718 (id_card, name, phone)
VALUES (p_id_card, p_name, p_phone) ON CONFLICT (id_card) DO
UPDATE
SET phone = COALESCE(EXCLUDED.phone, person_1718.phone);
END;
$$;
-- 查询指定订单的详细信息
CREATE OR REPLACE FUNCTION fn_order_detail_1718(p_order_id UUID) RETURNS TABLE (
        order_id UUID,
        status VARCHAR(20),
        quantity INT,
        check_in_date DATE,
        check_out_date DATE,
        nights INT,
        total_price DECIMAL(10, 2),
        actual_price DECIMAL(10, 2),
        create_at TIMESTAMPTZ,
        order_user TEXT,
        order_user_name TEXT,
        order_user_phone VARCHAR(20),
        hotel_name TEXT,
        room_type TEXT,
        room_price DECIMAL(10, 2),
        guest_count BIGINT,
        guest_names TEXT
    ) LANGUAGE plpgsql STABLE AS $$ BEGIN RETURN QUERY
SELECT o.id,
    o.status::VARCHAR(20),
    o.quantity,
    o.check_in_date,
    o.check_out_date,
    (o.check_out_date - o.check_in_date)::INT,
    o.total_price,
    o.actual_price,
    o.create_at,
    u.username,
    u.real_name,
    u.phone,
    h.hotel_name,
    rm.type_name,
    rm.price,
    COALESCE(g.guest_count, 0),
    g.guest_names
FROM order_1718 o
    JOIN user_1718 u ON u.id = o.user_id
    JOIN room_1718 rm ON rm.id = o.room_id
    JOIN hotel_1718 h ON h.id = rm.hotel_id
    LEFT JOIN (
        SELECT og.order_id,
            COUNT(*) AS guest_count,
            STRING_AGG(
                p.name,
                ', '
                ORDER BY p.id_card
            ) AS guest_names
        FROM order_guest_1718 og
            JOIN person_1718 p ON p.id_card = og.id_card
        WHERE og.order_id = p_order_id
        GROUP BY og.order_id
    ) g ON g.order_id = o.id
WHERE o.id = p_order_id;
END;
$$;