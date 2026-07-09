import type { Room } from "./room";
import type { User } from "./user";

export type OrderStatus =
  "pending" | "booked" | "checked_in" | "cancelled" | "completed";

/** 订单（GORM Preload 完整模型） */
export interface Order {
  id: string;
  userId: string;
  roomId: string;
  quantity: number;
  checkInDate: string;
  checkOutDate: string;
  totalPrice: number;
  actualPrice: number;
  status: OrderStatus;
  createAt: string;
  updateAt: string;
  room?: Room;
  user?: User; // 下单人
  guests?: OrderGuest[]; // 入住人
}

/** 入住人 */
export interface OrderGuest {
  orderId: string;
  idCard: string;
  person?: {
    idCard: string;
    name: string;
    phone?: string;
  };
}

/** 订单完整视图（一行一个入住人，兼容旧接口） */
export interface OrderFull {
  orderId: string;
  userId: string;
  username: string; // 下单人用户名
  hotelName: string;
  roomType: string;
  checkInDate: string;
  checkOutDate: string;
  guestName: string; // 入住人姓名
  guestIdCard: string; // 入住人身份证号
  quantity: number;
  totalPrice: number;
  orderStatus: string;
  createAt: string;
}

/** 我的订单列表（用户端） */
export interface MyOrder {
  orderId: string;
  hotelName: string;
  roomType: string;
  checkInDate: string;
  checkOutDate: string;
  quantity: number;
  actualPrice: number;
  orderStatus: string;
  createAt: string;
}

/** 下单请求 */
export interface CreateOrderRequest {
  roomId: string;
  checkInDate: string;
  checkOutDate: string;
  guestName: string; // 入住人姓名
  guestPhone: string; // 入住人电话
  guestIdCard: string; // 入住人身份证号
  roomCount: number;
  totalPrice: number;
  actualPrice: number;
}

export interface UpdateOrderStatusRequest {
  status: OrderStatus;
}

/** 订单详情（管理端详情，走 view_order_detail_1718） */
export interface OrderDetail {
  orderId: string;
  status: string;
  quantity: number;
  checkInDate: string;
  checkOutDate: string;
  nights: number;
  totalPrice: number;
  actualPrice: number;
  createAt: string;
  orderUser: string; // 下单人用户名
  orderUserName?: string; // 下单人真实姓名
  orderUserPhone?: string; // 下单人手机号
  hotelName: string;
  city?: string;
  roomType: string;
  roomPrice: number;
  guestCount: number; // 入住人数
  guestNames: string; // 入住人姓名（逗号分隔）
}

/** 订单概览（管理端列表，走 view_order_summary_1718） */
export interface OrderSummary {
  orderId: string;
  status: string;
  quantity: number;
  checkInDate: string;
  checkOutDate: string;
  nights: number;
  actualPrice: number;
  createAt: string;
  orderUserName?: string; // 下单人姓名
  hotelName: string;
  roomType: string;
  guestCount: number; // 入住人数
}
