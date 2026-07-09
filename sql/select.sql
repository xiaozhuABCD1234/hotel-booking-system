-- 订单详细数据（按订单分组，入住人聚合）
SELECT o.id              AS order_id,
       o.status,
       o.quantity,
       o.check_in_date,
       o.check_out_date,
       (o.check_out_date - o.check_in_date) AS nights,
       o.total_price,
       o.actual_price,
       o.create_at,
       u.username        AS order_user,
       u.real_name       AS order_user_name,
       u.phone           AS order_user_phone,
       h.hotel_name,
       rm.type_name      AS room_type,
       rm.price          AS room_price,
       COALESCE(g.guest_count, 0) AS guest_count,
       COALESCE(g.guest_names, '') AS guest_names
FROM order_1718 o
JOIN user_1718 u ON u.id = o.user_id
JOIN room_1718 rm ON rm.id = o.room_id
JOIN hotel_1718 h ON h.id = rm.hotel_id
LEFT JOIN (
    SELECT og.order_id,
           COUNT(*) AS guest_count,
           STRING_AGG(p.name, ', ' ORDER BY p.id_card) AS guest_names
    FROM order_guest_1718 og
    JOIN person_1718 p ON p.id_card = og.id_card
    GROUP BY og.order_id
) g ON g.order_id = o.id
ORDER BY o.create_at DESC;
