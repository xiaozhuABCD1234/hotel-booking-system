-- 酒店预订管理系统 — 视图定义
-- PostgreSQL 18
-- 依赖 view_person_info_1718 的视图需按顺序执行
-- 按城市 → 区域 → 酒店查看所有客房详细信息
CREATE OR REPLACE VIEW view_room_details_1718 AS
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
    STRING_AGG(
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
-- 从身份证号推导性别、年龄、出生日期
-- 被视图 view_guest_booking_stats_1718 依赖
CREATE OR REPLACE VIEW view_person_info_1718 AS
SELECT id_card,
    name,
    phone,
    occupation,
    education,
    income,
    TO_DATE(SUBSTRING(id_card, 7, 8), 'YYYYMMDD') AS birth_date,
    CASE
        WHEN SUBSTRING(id_card, 17, 1) ~ '\d'
        AND SUBSTRING(id_card, 17, 1)::INT % 2 = 1 THEN '男'
        WHEN SUBSTRING(id_card, 17, 1) ~ '\d'
        AND SUBSTRING(id_card, 17, 1)::INT % 2 = 0 THEN '女'
        ELSE NULL
    END AS gender,
    DATE_PART(
        'year',
        AGE(TO_DATE(SUBSTRING(id_card, 7, 8), 'YYYYMMDD'))
    )::INT AS age
FROM person_1718;
-- 用户 VIP 信息 — 个人中心、下单折扣计算
CREATE OR REPLACE VIEW view_user_vip_1718 AS
SELECT u.id AS user_id,
    u.username,
    u.phone,
    u.email,
    u.real_name,
    u.id_card,
    u.role,
    u.points,
    u.vip_level,
    vl.level_name AS vip_level_name,
    vl.discount_rate,
    (
        SELECT vl2.min_points - u.points
        FROM vip_level_1718 vl2
        WHERE vl2.min_points > u.points
        ORDER BY vl2.min_points ASC
        LIMIT 1
    ) AS points_to_next_level,
    u.create_at
FROM user_1718 u
    JOIN vip_level_1718 vl ON vl.level = u.vip_level
WHERE u.status = 1;
-- 酒店摘要列表 — 搜索页、首页推荐
CREATE OR REPLACE VIEW view_hotel_summary_1718 AS
SELECT h.id AS hotel_id,
    h.hotel_name,
    r_province.region_name AS province,
    r_city.region_name AS city,
    r_district.region_name AS district,
    h.address,
    h.telephone,
    h.star_level,
    h.description,
    (
        SELECT hi.image_url
        FROM hotel_image_1718 hi
        WHERE hi.hotel_id = h.id
        ORDER BY hi.image_url
        LIMIT 1
    ) AS main_image,
    COALESCE(MIN(rm.price), 0) AS min_price,
    COUNT(DISTINCT rm.id) AS room_count,
    COALESCE(SUM(rm.total_quantity), 0) AS total_rooms,
    COALESCE(rev_stats.avg_rating, 0) AS avg_rating,
    COALESCE(rev_stats.review_count, 0) AS review_count,
    h.status
FROM hotel_1718 h
    JOIN region_1718 r_district ON h.region_id = r_district.id
    LEFT JOIN region_1718 r_city ON r_district.parents_id = r_city.id
    LEFT JOIN region_1718 r_province ON r_city.parents_id = r_province.id
    LEFT JOIN room_1718 rm ON rm.hotel_id = h.id
    AND rm.status = 1
    LEFT JOIN (
        SELECT hotel_id,
            AVG(rating)::DECIMAL(3, 2) AS avg_rating,
            COUNT(*) AS review_count
        FROM review_1718
        GROUP BY hotel_id
    ) rev_stats ON rev_stats.hotel_id = h.id
WHERE h.status = 1
GROUP BY h.id,
    h.hotel_name,
    r_province.region_name,
    r_city.region_name,
    r_district.region_name,
    h.address,
    h.telephone,
    h.star_level,
    h.description,
    h.status,
    rev_stats.avg_rating,
    rev_stats.review_count;
-- 我的订单列表 — 每个订单一行，不展开入住人
CREATE OR REPLACE VIEW view_my_orders_1718 AS
SELECT o.id AS order_id,
    o.user_id,
    h.hotel_name,
    r_city.region_name AS city,
    rm.type_name AS room_type,
    o.quantity,
    o.check_in_date,
    o.check_out_date,
    (o.check_out_date - o.check_in_date) AS nights,
    o.actual_price,
    o.status AS order_status,
    COUNT(og.id_card) AS guest_count,
    o.create_at
FROM order_1718 o
    JOIN room_1718 rm ON rm.id = o.room_id
    JOIN hotel_1718 h ON h.id = rm.hotel_id
    JOIN region_1718 r_district ON h.region_id = r_district.id
    LEFT JOIN region_1718 r_city ON r_district.parents_id = r_city.id
    LEFT JOIN order_guest_1718 og ON og.order_id = o.id
GROUP BY o.id,
    o.user_id,
    h.hotel_name,
    r_city.region_name,
    rm.type_name,
    o.quantity,
    o.check_in_date,
    o.check_out_date,
    o.actual_price,
    o.status,
    o.create_at;
-- 评价详情 — 酒店评价列表页、用户评价记录
CREATE OR REPLACE VIEW view_review_full_1718 AS
SELECT rv.id AS review_id,
    rv.user_id,
    u.username,
    rv.hotel_id,
    h.hotel_name,
    rv.order_id,
    rm.type_name AS room_type,
    o.check_in_date,
    o.check_out_date,
    rv.rating,
    rv.content,
    rv.create_at
FROM review_1718 rv
    JOIN user_1718 u ON u.id = rv.user_id
    JOIN hotel_1718 h ON h.id = rv.hotel_id
    JOIN order_1718 o ON o.id = rv.order_id
    JOIN room_1718 rm ON rm.id = o.room_id;
-- 入住人预订统计分析 — 按年龄/性别/职业/学历/收入等维度分析客户偏好
-- 依赖 view_person_info_1718
CREATE OR REPLACE VIEW view_guest_booking_stats_1718 AS
SELECT vpi.id_card AS person_id_card,
    vpi.name AS person_name,
    vpi.gender,
    vpi.age,
    vpi.occupation,
    vpi.education,
    vpi.income,
    CASE
        WHEN vpi.age < 18 THEN '18岁以下'
        WHEN vpi.age BETWEEN 18 AND 25 THEN '18-25岁'
        WHEN vpi.age BETWEEN 26 AND 35 THEN '26-35岁'
        WHEN vpi.age BETWEEN 36 AND 45 THEN '36-45岁'
        WHEN vpi.age BETWEEN 46 AND 55 THEN '46-55岁'
        WHEN vpi.age > 55 THEN '55岁以上'
        ELSE '未知'
    END AS age_group,
    COUNT(DISTINCT o.id) AS total_orders,
    COALESCE(SUM(o.check_out_date - o.check_in_date), 0) AS total_nights,
    COALESCE(SUM(o.actual_price), 0) AS total_amount,
    CASE
        WHEN COUNT(DISTINCT o.id) > 0 THEN ROUND(
            SUM(o.actual_price) / COUNT(DISTINCT o.id),
            2
        )
        ELSE 0
    END AS avg_order_amount,
    (
        SELECT r_city.region_name
        FROM order_guest_1718 og2
            JOIN order_1718 o2 ON o2.id = og2.order_id
            JOIN room_1718 rm2 ON rm2.id = o2.room_id
            JOIN hotel_1718 h2 ON h2.id = rm2.hotel_id
            JOIN region_1718 r_district2 ON r_district2.id = h2.region_id
            JOIN region_1718 r_city ON r_city.id = r_district2.parents_id
        WHERE og2.id_card = vpi.id_card
            AND o2.status IN ('booked', 'checked_in', 'completed')
        GROUP BY r_city.region_name
        ORDER BY COUNT(*) DESC
        LIMIT 1
    ) AS fav_city,
    (
        SELECT h2.hotel_name
        FROM order_guest_1718 og2
            JOIN order_1718 o2 ON o2.id = og2.order_id
            JOIN room_1718 rm2 ON rm2.id = o2.room_id
            JOIN hotel_1718 h2 ON h2.id = rm2.hotel_id
        WHERE og2.id_card = vpi.id_card
            AND o2.status IN ('booked', 'checked_in', 'completed')
        GROUP BY h2.hotel_name
        ORDER BY COUNT(*) DESC
        LIMIT 1
    ) AS fav_hotel,
    (
        SELECT rm2.type_name
        FROM order_guest_1718 og2
            JOIN order_1718 o2 ON o2.id = og2.order_id
            JOIN room_1718 rm2 ON rm2.id = o2.room_id
        WHERE og2.id_card = vpi.id_card
            AND o2.status IN ('booked', 'checked_in', 'completed')
        GROUP BY rm2.type_name
        ORDER BY COUNT(*) DESC
        LIMIT 1
    ) AS fav_room_type,
    MAX(o.create_at) AS last_order_date
FROM view_person_info_1718 vpi
    JOIN person_1718 p ON p.id_card = vpi.id_card
    LEFT JOIN order_guest_1718 og ON og.id_card = vpi.id_card
    LEFT JOIN order_1718 o ON o.id = og.order_id
    AND o.status IN ('booked', 'checked_in', 'completed')
GROUP BY vpi.id_card,
    vpi.name,
    vpi.gender,
    vpi.age,
    vpi.occupation,
    vpi.education,
    vpi.income;
-- 订单详情（下单人与入住人明确区分，入住人聚合）
CREATE OR REPLACE VIEW view_order_detail_1718 AS
SELECT o.id AS order_id,
    o.status,
    o.quantity,
    o.check_in_date,
    o.check_out_date,
    (o.check_out_date - o.check_in_date) AS nights,
    o.total_price,
    o.actual_price,
    o.create_at,
    u.username AS order_user,
    u.real_name AS order_user_name,
    u.phone AS order_user_phone,
    h.hotel_name,
    r_city.region_name AS city,
    rm.type_name AS room_type,
    rm.price AS room_price,
    COALESCE(g.guest_count, 0) AS guest_count,
    COALESCE(g.guest_names, '') AS guest_names
FROM order_1718 o
    JOIN user_1718 u ON u.id = o.user_id
    JOIN room_1718 rm ON rm.id = o.room_id
    JOIN hotel_1718 h ON h.id = rm.hotel_id
    JOIN region_1718 r_district ON h.region_id = r_district.id
    LEFT JOIN region_1718 r_city ON r_district.parents_id = r_city.id
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
        GROUP BY og.order_id
    ) g ON g.order_id = o.id
ORDER BY o.create_at DESC;
-- 订单概览（简略版，适合列表页）
CREATE OR REPLACE VIEW view_order_summary_1718 AS
SELECT o.id AS order_id,
    o.status,
    o.quantity,
    o.check_in_date,
    o.check_out_date,
    (o.check_out_date - o.check_in_date) AS nights,
    o.actual_price,
    o.create_at,
    u.real_name AS order_user_name,
    h.hotel_name,
    rm.type_name AS room_type,
    COALESCE(g.guest_count, 0) AS guest_count
FROM order_1718 o
    JOIN user_1718 u ON u.id = o.user_id
    JOIN room_1718 rm ON rm.id = o.room_id
    JOIN hotel_1718 h ON h.id = rm.hotel_id
    LEFT JOIN (
        SELECT og.order_id,
            COUNT(*) AS guest_count
        FROM order_guest_1718 og
        GROUP BY og.order_id
    ) g ON g.order_id = o.id
ORDER BY o.create_at DESC;
-- 订单全量信息（一行一个入住人 兼容）
CREATE OR REPLACE VIEW view_order_full_1718 AS
SELECT o.id AS order_id,
    o.user_id,
    u.username,
    u.phone AS user_phone,
    u.real_name AS user_real_name,
    h.id AS hotel_id,
    h.hotel_name,
    r_city.region_name AS city,
    r_district.region_name AS district,
    h.telephone AS hotel_telephone,
    rm.id AS room_id,
    rm.type_name AS room_type,
    o.quantity,
    o.check_in_date,
    o.check_out_date,
    (o.check_out_date - o.check_in_date) AS nights,
    o.total_price,
    0::DECIMAL(10, 2) AS discount,
    o.actual_price,
    vl.discount_rate AS vip_discount_rate,
    o.status AS order_status,
    og.id_card AS guest_id_card,
    p.name AS guest_name,
    vpi.gender AS guest_gender,
    vpi.age AS guest_age,
    vpi.occupation AS guest_occupation,
    vpi.education AS guest_education,
    vpi.income AS guest_income,
    o.create_at
FROM order_1718 o
    JOIN user_1718 u ON u.id = o.user_id
    JOIN vip_level_1718 vl ON vl.level = u.vip_level
    JOIN room_1718 rm ON rm.id = o.room_id
    JOIN hotel_1718 h ON h.id = rm.hotel_id
    JOIN region_1718 r_district ON h.region_id = r_district.id
    LEFT JOIN region_1718 r_city ON r_district.parents_id = r_city.id
    LEFT JOIN order_guest_1718 og ON og.order_id = o.id
    LEFT JOIN person_1718 p ON p.id_card = og.id_card
    LEFT JOIN view_person_info_1718 vpi ON vpi.id_card = og.id_card;