import api from "./client";
import type {
  ApiResponse,
  Order,
  OrderSummary,
  OrderDetail,
  CreateOrderRequest,
  UpdateOrderStatusRequest,
  OrderFull,
  MyOrder,
} from "@/types";

export const orderApi = {
  /** Get all orders with pagination */
  list(page = 1, pageSize = 10) {
    return api.get<ApiResponse<OrderSummary[]>>("/orders", {
      params: { page, pageSize },
    });
  },

  /** Get one order by ID */
  getById(id: string) {
    return api.get<ApiResponse<Order>>(`/orders/${id}`);
  },

  /** Get order detail (下单人/入住人区分, view_order_detail_1718) */
  getDetail(id: string) {
    return api.get<ApiResponse<OrderDetail>>(`/orders/${id}/detail`);
  },

  /** Create order */
  create(data: CreateOrderRequest) {
    return api.post<ApiResponse<Order>>("/orders", data);
  },

  /** Update order status */
  updateStatus(id: string, data: UpdateOrderStatusRequest) {
    return api.put<ApiResponse<Order>>(`/orders/${id}/status`, data);
  },

  /** Delete order */
  delete(id: string) {
    return api.delete<ApiResponse>(`/orders/${id}`);
  },

  /** Get orders by current user */
  listByUser(userId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<Order[]>>(`/orders/by-user`, {
      params: { userID: userId, page, pageSize },
    });
  },

  /** Get orders by hotel */
  listByHotel(hotelId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<Order[]>>(`/orders/by-hotel`, {
      params: { hotelID: hotelId, page, pageSize },
    });
  },

  /** Report: order full by user */
  orderFullByUser(userId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<OrderFull[]>>(`/reports/order-full/by-user`, {
      params: { userID: userId, page, pageSize },
    });
  },

  /** Report: order full by hotel */
  orderFullByHotel(hotelId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<OrderFull[]>>(`/reports/order-full/by-hotel`, {
      params: { hotelID: hotelId, page, pageSize },
    });
  },

  /** Report: my orders */
  myOrders(userId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<MyOrder[]>>(`/reports/my-orders`, {
      params: { userID: userId, page, pageSize },
    });
  },
};
