-- VIP 等级
INSERT INTO vip_level_1718 (level, level_name, min_points, discount_rate)
VALUES (0, '普通会员', 0, 1.00),
    (1, '白银会员', 100, 0.95),
    (2, '黄金会员', 500, 0.90),
    (3, '钻石会员', 2000, 0.85);
-- 管理员
INSERT INTO user_1718 (username, password, role)
VALUES (
        'admin',
        '$2a$10$.vFlQ9RjUGDgEzNSnlojUuASlMgHaHHCKQNd1h7M7YYbhuv66yDn2',
        'admin'
    );
-- 人员（入住人身份信息）
INSERT INTO person_1718 (id_card, name, phone)
VALUES ('320102200003150018', '张伟', '13800001111'),
    ('310115199807120026', '李娜', '13800002222'),
    ('110108198511090015', '王强', '13800003333'),
    ('440106197602280024', '赵敏', '13800004444'),
    ('330102196510220019', '陈建国', '13800005555'),
    ('510104199203080042', '刘芳', '13800006666');
-- 地区数据（省/市/区 三级层次）
-- 数据来源：https://cn.quhua.net
-- 省级：全部 34 个；区级：仅上海市 → 浦东新区；街道/镇级：浦东新区下辖全部
-- id 由 SERIAL 自动生成，parents_id 通过子查询关联
-- 省级
INSERT INTO region_1718 (region_name, parents_id)
VALUES ('北京市', NULL),
    ('天津市', NULL),
    ('河北省', NULL),
    ('山西省', NULL),
    ('内蒙古自治区', NULL),
    ('辽宁省', NULL),
    ('吉林省', NULL),
    ('黑龙江省', NULL),
    ('上海市', NULL),
    ('江苏省', NULL),
    ('浙江省', NULL),
    ('安徽省', NULL),
    ('福建省', NULL),
    ('江西省', NULL),
    ('山东省', NULL),
    ('河南省', NULL),
    ('湖北省', NULL),
    ('湖南省', NULL),
    ('广东省', NULL),
    ('广西壮族自治区', NULL),
    ('海南省', NULL),
    ('重庆市', NULL),
    ('四川省', NULL),
    ('贵州省', NULL),
    ('云南省', NULL),
    ('西藏自治区', NULL),
    ('陕西省', NULL),
    ('甘肃省', NULL),
    ('青海省', NULL),
    ('宁夏回族自治区', NULL),
    ('新疆维吾尔自治区', NULL),
    ('台湾省', NULL),
    ('香港特别行政区', NULL),
    ('澳门特别行政区', NULL);
-- 区级 — 上海市 → 浦东新区
INSERT INTO region_1718 (region_name, parents_id)
SELECT '浦东新区',
    id
FROM region_1718
WHERE region_name = '上海市';
-- 街道/镇级 — 浦东新区下辖
INSERT INTO region_1718 (region_name, parents_id)
SELECT m.name,
    pd.id
FROM region_1718 pd
    CROSS JOIN (
        VALUES ('潍坊新村街道'),
            ('陆家嘴街道'),
            ('周家渡街道'),
            ('塘桥街道'),
            ('上钢新村街道'),
            ('南码头路街道'),
            ('沪东新村街道'),
            ('金杨新村街道'),
            ('洋泾街道'),
            ('浦兴路街道'),
            ('东明路街道'),
            ('花木街道'),
            ('申港街道'),
            ('川沙新镇'),
            ('高桥镇'),
            ('北蔡镇'),
            ('合庆镇'),
            ('唐镇'),
            ('曹路镇'),
            ('金桥镇'),
            ('高行镇'),
            ('高东镇'),
            ('张江镇'),
            ('三林镇'),
            ('惠南镇'),
            ('周浦镇'),
            ('新场镇'),
            ('大团镇'),
            ('芦潮港镇'),
            ('康桥镇'),
            ('航头镇'),
            ('六灶镇'),
            ('祝桥镇'),
            ('泥城镇'),
            ('宣桥镇'),
            ('书院镇'),
            ('万祥镇'),
            ('老港镇'),
            ('芦潮港农场'),
            ('东海农场'),
            ('朝阳农场'),
            ('外高桥保税区'),
            ('金桥出口加工区')
    ) AS m(name)
WHERE pd.region_name = '浦东新区';
-- 测试数据：酒店、房型、订单
DO $$
DECLARE v_user_id UUID;
v_hotel_id UUID;
v_room_std UUID;
v_room_deluxe UUID;
v_room_suite UUID;
v_order1 UUID;
v_order2 UUID;
BEGIN -- 下单用户
INSERT INTO user_1718 (username, password, role, real_name, phone)
VALUES (
        'zhangsan',
        '$2a$10$dummyhashedpassword1234567890abcdef',
        'customer',
        '张三',
        '13900001111'
    )
RETURNING id INTO v_user_id;
-- 酒店
INSERT INTO hotel_1718 (
        hotel_name,
        region_id,
        address,
        telephone,
        star_level,
        description
    )
SELECT '陆家嘴国际酒店',
    id,
    '上海市浦东新区陆家嘴环路1000号',
    '021-58888888',
    5,
    '五星级商务酒店，坐落于陆家嘴金融贸易区核心地段，毗邻东方明珠、上海中心大厦'
FROM region_1718
WHERE region_name = '陆家嘴街道'
RETURNING id INTO v_hotel_id;
-- 房型
INSERT INTO room_1718 (
        hotel_id,
        type_name,
        total_quantity,
        available_quantity,
        price,
        description
    )
VALUES (
        v_hotel_id,
        '标准间',
        50,
        50,
        388.00,
        '舒适双床房，40㎡，城景'
    )
RETURNING id INTO v_room_std;
INSERT INTO room_1718 (
        hotel_id,
        type_name,
        total_quantity,
        available_quantity,
        price,
        weekend_price,
        description
    )
VALUES (
        v_hotel_id,
        '豪华大床房',
        30,
        30,
        688.00,
        788.00,
        '宽敞大床房，55㎡，江景'
    )
RETURNING id INTO v_room_deluxe;
INSERT INTO room_1718 (
        hotel_id,
        type_name,
        total_quantity,
        available_quantity,
        price,
        weekend_price,
        description
    )
VALUES (
        v_hotel_id,
        '行政套房',
        10,
        10,
        1288.00,
        1588.00,
        '一室一厅套房，80㎡，含行政酒廊'
    )
RETURNING id INTO v_room_suite;
-- 订单 1：单人入住，标准间 1 间，2 晚
INSERT INTO order_1718 (
        user_id,
        room_id,
        quantity,
        check_in_date,
        check_out_date,
        total_price,
        actual_price,
        status
    )
VALUES (
        v_user_id,
        v_room_std,
        1,
        '2026-07-15',
        '2026-07-17',
        388.00 * 2,
        388.00 * 2,
        'booked'
    )
RETURNING id INTO v_order1;
INSERT INTO order_guest_1718 (order_id, id_card)
VALUES (v_order1, '320102200003150018');
-- 张伟
-- 订单 2：多人入住，豪华大床房 2 间，3 晚
INSERT INTO order_1718 (
        user_id,
        room_id,
        quantity,
        check_in_date,
        check_out_date,
        total_price,
        actual_price,
        status
    )
VALUES (
        v_user_id,
        v_room_deluxe,
        2,
        '2026-08-01',
        '2026-08-04',
        688.00 * 3 * 2,
        688.00 * 3 * 2,
        'booked'
    )
RETURNING id INTO v_order2;
INSERT INTO order_guest_1718 (order_id, id_card)
VALUES (v_order2, '310115199807120026'),
    -- 李娜
    (v_order2, '110108198511090015'),
    -- 王强
    (v_order2, '440106197602280024');
-- 赵敏
END;
$$;