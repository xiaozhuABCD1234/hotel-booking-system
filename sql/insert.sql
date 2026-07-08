-- VIP 等级
INSERT INTO vip_level_1718 (level, level_name, min_points, discount_rate) VALUES
    (0, '普通会员', 0, 1.00),
    (1, '白银会员', 100, 0.95),
    (2, '黄金会员', 500, 0.90),
    (3, '钻石会员', 2000, 0.85);

-- 管理员
INSERT INTO user_1718 (username, password, role) VALUES
    ('admin', '$2a$10$.vFlQ9RjUGDgEzNSnlojUuASlMgHaHHCKQNd1h7M7YYbhuv66yDn2', 'admin');
