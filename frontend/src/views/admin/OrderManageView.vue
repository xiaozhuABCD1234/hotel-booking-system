<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { orderApi } from "@/api/order";
import type { OrderDetail, OrderStatus, OrderSummary } from "@/types";
import { toast } from "vue-sonner";
import { getApiErrorMessage } from "@/lib/utils";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import { ShoppingBag, Filter, Search, Check, Trash2 } from "@lucide/vue";

const orders = ref<OrderSummary[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const totalItems = ref(0);
const statusFilter = ref<string>("all");

const detailDialogOpen = ref(false);
const detailLoading = ref(false);
const selectedOrder = ref<OrderDetail | null>(null);

const statusDialogOpen = ref(false);
const statusUpdateOrder = ref<OrderSummary | null>(null);
const newStatus = ref<OrderStatus>("pending");

const deleteDialogOpen = ref(false);
const deleteTargetOrder = ref<OrderSummary | null>(null);

const statusOptions: { value: string; label: string }[] = [
  { value: "all", label: "全部状态" },
  { value: "pending", label: "待确认" },
  { value: "booked", label: "已预订" },
  { value: "checked_in", label: "已入住" },
  { value: "cancelled", label: "已取消" },
  { value: "completed", label: "已完成" },
];

const statusBadgeClass = (status: OrderStatus): string => {
  const map: Record<OrderStatus, string> = {
    pending: "bg-yellow-100 text-yellow-800 border-yellow-200",
    booked: "bg-blue-100 text-blue-800 border-blue-200",
    checked_in: "bg-green-100 text-green-800 border-green-200",
    cancelled: "bg-red-100 text-red-800 border-red-200",
    completed: "bg-gray-100 text-gray-800 border-gray-200",
  };
  return map[status] ?? "bg-gray-100 text-gray-800 border-gray-200";
};

const statusLabel = (status: OrderStatus): string => {
  const map: Record<OrderStatus, string> = {
    pending: "待确认",
    booked: "已预订",
    checked_in: "已入住",
    cancelled: "已取消",
    completed: "已完成",
  };
  return map[status] ?? status;
};

const filteredOrders = computed(() => {
  if (statusFilter.value === "all") return orders.value;
  return orders.value.filter((o) => o.status === statusFilter.value);
});

async function fetchOrders() {
  loading.value = true;
  try {
    const res = await orderApi.list(currentPage.value, pageSize.value);
    if (res.data.data) {
      orders.value = res.data.data ?? [];
      const p = res.data.pagination;
      if (p) {
        totalPages.value = p.totalPages;
        totalItems.value = p.totalItems;
      }
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "获取订单列表失败"));
  } finally {
    loading.value = false;
  }
}

async function openDetail(orderId: string) {
  detailDialogOpen.value = true;
  detailLoading.value = true;
  try {
    const res = await orderApi.getDetail(orderId);
    selectedOrder.value = res.data.data ?? null;
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "获取订单详情失败"));
    detailDialogOpen.value = false;
  } finally {
    detailLoading.value = false;
  }
}

function openStatusUpdate(order: OrderSummary) {
  statusUpdateOrder.value = order;
  newStatus.value = order.status as OrderStatus;
  statusDialogOpen.value = true;
}

async function confirmStatusUpdate() {
  if (!statusUpdateOrder.value) return;
  try {
    await orderApi.updateStatus(statusUpdateOrder.value.orderId, {
      status: newStatus.value,
    });
    toast.success("订单状态已更新");
    statusDialogOpen.value = false;
    await fetchOrders();
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "更新状态失败"));
  }
}

function openDelete(order: OrderSummary) {
  deleteTargetOrder.value = order;
  deleteDialogOpen.value = true;
}

async function confirmDelete() {
  if (!deleteTargetOrder.value) return;
  try {
    await orderApi.delete(deleteTargetOrder.value.orderId);
    toast.success("订单已删除");
    deleteDialogOpen.value = false;
    await fetchOrders();
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "删除订单失败"));
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return "-";
  const d = new Date(dateStr);
  return d.toLocaleDateString("zh-CN");
}

function goToPage(page: number) {
  currentPage.value = page;
  fetchOrders();
}

onMounted(() => {
  fetchOrders();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold flex items-center gap-2">
        <ShoppingBag class="h-6 w-6" />
        订单管理
      </h1>
    </div>

    <Card>
      <CardHeader class="pb-4">
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <Filter class="h-4 w-4 text-muted-foreground" />
            <Select v-model="statusFilter">
              <SelectTrigger class="w-[160px]">
                <SelectValue placeholder="筛选状态" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="opt in statusOptions"
                  :key="opt.value"
                  :value="opt.value"
                >
                  {{ opt.label }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="text-sm text-muted-foreground">
            共 {{ totalItems }} 条订单
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="space-y-3">
          <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
        </div>
        <div v-else class="rounded-md border overflow-x-auto">
          <Table class="w-full">
            <TableHeader>
              <TableRow>
                <TableHead class="w-[90px] whitespace-nowrap">订单ID</TableHead>
                <TableHead class="whitespace-nowrap">酒店</TableHead>
                <TableHead class="whitespace-nowrap">房型</TableHead>
                <TableHead class="whitespace-nowrap">下单人</TableHead>
                <TableHead class="whitespace-nowrap">入住人数</TableHead>
                <TableHead class="whitespace-nowrap">入住日期</TableHead>
                <TableHead class="whitespace-nowrap">离店日期</TableHead>
                <TableHead class="text-right whitespace-nowrap">总价</TableHead>
                <TableHead class="whitespace-nowrap w-[80px]">状态</TableHead>
                <TableHead class="text-right whitespace-nowrap w-[120px]"
                  >操作</TableHead
                >
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="filteredOrders.length === 0">
                <TableCell
                  colspan="10"
                  class="text-center py-8 text-muted-foreground"
                >
                  暂无数据
                </TableCell>
              </TableRow>
              <TableRow
                v-for="order in filteredOrders"
                :key="order.orderId"
                class="hover:bg-muted/50"
              >
                <TableCell class="font-mono text-xs whitespace-nowrap">
                  {{ order.orderId.slice(0, 8) }}
                </TableCell>
                <TableCell class="whitespace-nowrap">{{
                  order.hotelName
                }}</TableCell>
                <TableCell class="whitespace-nowrap">{{
                  order.roomType
                }}</TableCell>
                <TableCell class="whitespace-nowrap">{{
                  order.orderUserName || "-"
                }}</TableCell>
                <TableCell class="whitespace-nowrap">{{
                  order.guestCount
                }}</TableCell>
                <TableCell class="whitespace-nowrap">{{
                  formatDate(order.checkInDate)
                }}</TableCell>
                <TableCell class="whitespace-nowrap">{{
                  formatDate(order.checkOutDate)
                }}</TableCell>
                <TableCell class="text-right font-medium whitespace-nowrap">
                  ¥{{ order.actualPrice.toFixed(2) }}
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  <Badge
                    :class="statusBadgeClass(order.status as OrderStatus)"
                    variant="outline"
                    class="whitespace-nowrap"
                  >
                    {{ statusLabel(order.status as OrderStatus) }}
                  </Badge>
                </TableCell>
                <TableCell class="text-right whitespace-nowrap">
                  <div class="flex justify-end gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      @click="openDetail(order.orderId)"
                    >
                      <Search class="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      @click="openStatusUpdate(order)"
                    >
                      <Check class="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      @click="openDelete(order)"
                    >
                      <Trash2 class="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <div
          v-if="!loading && totalPages > 1"
          class="flex items-center justify-between mt-4"
        >
          <div class="text-sm text-muted-foreground">
            第 {{ currentPage }} / {{ totalPages }} 页
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage <= 1"
              @click="goToPage(currentPage - 1)"
            >
              上一页
            </Button>
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage >= totalPages"
              @click="goToPage(currentPage + 1)"
            >
              下一页
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Detail Dialog -->
    <Dialog v-model:open="detailDialogOpen">
      <DialogContent class="max-w-lg">
        <DialogHeader>
          <DialogTitle>订单详情</DialogTitle>
          <DialogDescription>
            订单编号: {{ selectedOrder?.orderId }}
          </DialogDescription>
        </DialogHeader>

        <!-- Loading -->
        <div v-if="detailLoading" class="space-y-3 py-4">
          <Skeleton v-for="i in 6" :key="i" class="h-5 w-full" />
        </div>

        <!-- Content -->
        <div v-else-if="selectedOrder" class="space-y-3">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <Label class="text-muted-foreground">酒店</Label>
              <p class="font-medium">{{ selectedOrder.hotelName }}</p>
              <p
                v-if="selectedOrder.city"
                class="text-sm text-muted-foreground"
              >
                {{ selectedOrder.city }}
              </p>
            </div>
            <div>
              <Label class="text-muted-foreground">房型</Label>
              <p class="font-medium">{{ selectedOrder.roomType }}</p>
              <p class="text-sm text-muted-foreground">
                ¥{{ selectedOrder.roomPrice.toFixed(2) }}/晚
              </p>
            </div>
            <div>
              <Label class="text-muted-foreground">下单人</Label>
              <p class="font-medium">
                {{ selectedOrder.orderUserName || selectedOrder.orderUser }}
              </p>
            </div>
            <div>
              <Label class="text-muted-foreground">联系电话</Label>
              <p class="font-medium">
                {{ selectedOrder.orderUserPhone || "-" }}
              </p>
            </div>
            <div class="col-span-2">
              <Label class="text-muted-foreground">入住人</Label>
              <p class="font-medium">{{ selectedOrder.guestNames || "-" }}</p>
            </div>
            <div>
              <Label class="text-muted-foreground">房间数 / 入住人数</Label>
              <p class="font-medium">
                {{ selectedOrder.quantity }} 间 /
                {{ selectedOrder.guestCount }} 人
              </p>
            </div>
            <div>
              <Label class="text-muted-foreground">晚数</Label>
              <p class="font-medium">{{ selectedOrder.nights }} 晚</p>
            </div>
            <div>
              <Label class="text-muted-foreground">入住日期</Label>
              <p class="font-medium">
                {{ formatDate(selectedOrder.checkInDate) }}
              </p>
            </div>
            <div>
              <Label class="text-muted-foreground">离店日期</Label>
              <p class="font-medium">
                {{ formatDate(selectedOrder.checkOutDate) }}
              </p>
            </div>
            <div>
              <Label class="text-muted-foreground">总价 / 实付</Label>
              <p class="font-medium">
                <span class="text-muted-foreground line-through mr-2"
                  >¥{{ selectedOrder.totalPrice.toFixed(2) }}</span
                >
                ¥{{ selectedOrder.actualPrice.toFixed(2) }}
              </p>
            </div>
            <div>
              <Label class="text-muted-foreground">状态</Label>
              <div class="mt-1">
                <Badge
                  :class="statusBadgeClass(selectedOrder.status as OrderStatus)"
                  variant="outline"
                >
                  {{ statusLabel(selectedOrder.status as OrderStatus) }}
                </Badge>
              </div>
            </div>
          </div>
          <div>
            <Label class="text-muted-foreground">创建时间</Label>
            <p class="font-medium">{{ formatDate(selectedOrder.createAt) }}</p>
          </div>
        </div>
      </DialogContent>
    </Dialog>

    <!-- Status Update Dialog -->
    <Dialog v-model:open="statusDialogOpen">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>更新订单状态</DialogTitle>
          <DialogDescription>
            修改订单 {{ statusUpdateOrder?.orderId?.slice(0, 8) }} 的状态
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div>
            <Label>当前状态</Label>
            <Badge
              v-if="statusUpdateOrder"
              :class="
                statusBadgeClass(statusUpdateOrder!.status as OrderStatus)
              "
              variant="outline"
              class="ml-2"
            >
              {{ statusLabel(statusUpdateOrder!.status as OrderStatus) }}
            </Badge>
          </div>
          <div>
            <Label>新状态</Label>
            <Select v-model="newStatus">
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="pending">待确认</SelectItem>
                <SelectItem value="booked">已预订</SelectItem>
                <SelectItem value="checked_in">已入住</SelectItem>
                <SelectItem value="cancelled">已取消</SelectItem>
                <SelectItem value="completed">已完成</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="statusDialogOpen = false"
            >取消</Button
          >
          <Button @click="confirmStatusUpdate">确认更新</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirm Dialog -->
    <Dialog v-model:open="deleteDialogOpen">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>确认删除</DialogTitle>
          <DialogDescription>
            确定要删除订单
            {{ deleteTargetOrder?.orderId?.slice(0, 8) }} 吗？此操作不可撤销。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="deleteDialogOpen = false"
            >取消</Button
          >
          <Button variant="destructive" @click="confirmDelete">删除</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
