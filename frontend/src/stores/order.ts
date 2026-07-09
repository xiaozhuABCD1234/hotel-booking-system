import { defineStore } from "pinia";
import { ref } from "vue";
import { orderApi } from "@/api";
import type { Order, CreateOrderRequest, Pagination } from "@/types";

export const useOrderStore = defineStore("order", () => {
  const orders = ref<Order[]>([]);
  const pagination = ref<Pagination | null>(null);
  const loading = ref(false);

  async function fetchMyOrders(userId: string, page = 1, pageSize = 10) {
    loading.value = true;
    try {
      const res = await orderApi.listByUser(userId, page, pageSize);
      if (res.data.success && res.data.data) {
        orders.value = res.data.data;
        pagination.value = res.data.pagination ?? null;
      } else {
        orders.value = [];
        pagination.value = null;
      }
    } catch {
      orders.value = [];
      pagination.value = null;
    } finally {
      loading.value = false;
    }
  }

  async function fetchAllOrders(page = 1, pageSize = 10) {
    loading.value = true;
    try {
      const res = await orderApi.list(page, pageSize);
      if (res.data.success && res.data.data) {
        orders.value = res.data.data;
        pagination.value = res.data.pagination ?? null;
      } else {
        orders.value = [];
        pagination.value = null;
      }
    } catch {
      orders.value = [];
      pagination.value = null;
    } finally {
      loading.value = false;
    }
  }

  async function createOrder(data: CreateOrderRequest) {
    const res = await orderApi.create(data);
    return res.data;
  }

  async function cancelOrder(id: string) {
    const res = await orderApi.updateStatus(id, { status: "cancelled" });
    return res.data;
  }

  async function updateOrderStatus(id: string, status: Order["status"]) {
    const res = await orderApi.updateStatus(id, { status });
    return res.data;
  }

  return {
    orders,
    pagination,
    loading,
    fetchMyOrders,
    fetchAllOrders,
    createOrder,
    cancelOrder,
    updateOrderStatus,
  };
});
