-- 酒店预订管理系统 — 触发器定义
-- PostgreSQL 18

-- 评价用户必须与订单预订人一致
CREATE OR REPLACE FUNCTION fn_validate_review_user_1718()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.user_id != (SELECT user_id FROM order_1718 WHERE id = NEW.order_id) THEN
        RAISE EXCEPTION '评价用户与订单预订人不一致';
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_review_user_match_1718
    BEFORE INSERT OR UPDATE OF user_id, order_id
    ON review_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_validate_review_user_1718();

-- update_at 自动刷新
CREATE OR REPLACE FUNCTION fn_set_update_at_1718()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    NEW.update_at = now();
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_user_update_at_1718
    BEFORE UPDATE ON user_1718
    FOR EACH ROW EXECUTE FUNCTION fn_set_update_at_1718();

CREATE TRIGGER trg_hotel_update_at_1718
    BEFORE UPDATE ON hotel_1718
    FOR EACH ROW EXECUTE FUNCTION fn_set_update_at_1718();

CREATE TRIGGER trg_room_update_at_1718
    BEFORE UPDATE ON room_1718
    FOR EACH ROW EXECUTE FUNCTION fn_set_update_at_1718();

CREATE TRIGGER trg_order_update_at_1718
    BEFORE UPDATE ON order_1718
    FOR EACH ROW EXECUTE FUNCTION fn_set_update_at_1718();

CREATE TRIGGER trg_review_update_at_1718
    BEFORE UPDATE ON review_1718
    FOR EACH ROW EXECUTE FUNCTION fn_set_update_at_1718();

-- 客房可预订数量自动更新
CREATE OR REPLACE FUNCTION fn_update_room_quantity_1718()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
DECLARE
    v_available INT;
BEGIN
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
                update_at = now()
        WHERE id = NEW.room_id;

    ELSIF NEW.status = 'cancelled'
        AND TG_OP = 'UPDATE'
        AND OLD.status IN ('booked', 'checked_in')
    THEN
        UPDATE room_1718
            SET available_quantity = available_quantity + OLD.quantity,
                update_at = now()
        WHERE id = OLD.room_id;
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_order_room_quantity_1718
    AFTER INSERT OR UPDATE OF status
    ON order_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_room_quantity_1718();

-- 用户积分自动更新
CREATE OR REPLACE FUNCTION fn_update_user_points_1718()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.status = 'completed'
        AND (OLD.status IS NULL OR OLD.status != 'completed')
    THEN
        UPDATE user_1718
            SET points = points + FLOOR(NEW.actual_price)::INT,
                update_at = now()
        WHERE id = NEW.user_id;

    ELSIF TG_OP = 'UPDATE'
        AND OLD.status = 'completed'
        AND NEW.status != 'completed'
    THEN
        UPDATE user_1718
            SET points = GREATEST(0, points - FLOOR(OLD.actual_price)::INT),
                update_at = now()
        WHERE id = OLD.user_id;
    END IF;

    -- 通过积分自动更新vip等级
    UPDATE user_1718
    SET vip_level = (
        SELECT COALESCE(MAX(level), 0)
        FROM vip_level_1718
        WHERE min_points <= user_1718.points
    ),
    update_at = now()
    WHERE id = COALESCE(NEW.user_id, OLD.user_id);

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_order_user_points_1718
    AFTER INSERT OR UPDATE OF status
    ON order_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_user_points_1718();

-- 下单时自动应用 VIP 折扣（BEFORE INSERT，原子性，无 TOCTOU）
CREATE OR REPLACE FUNCTION fn_apply_vip_discount_1718()
RETURNS TRIGGER AS $$
DECLARE
    v_discount_rate DECIMAL(3,2);
BEGIN
    SELECT COALESCE(vl.discount_rate, 1.0)
    INTO v_discount_rate
    FROM user_1718 u
    JOIN vip_level_1718 vl ON vl.level = u.vip_level
    WHERE u.id = NEW.user_id;

    NEW.actual_price := NEW.total_price * v_discount_rate;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_order_vip_discount_1718
    BEFORE INSERT ON order_1718
    FOR EACH ROW
    EXECUTE FUNCTION fn_apply_vip_discount_1718();
