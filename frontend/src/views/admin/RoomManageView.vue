<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Plus, Pencil, Trash2 } from '@lucide/vue'
import { roomApi, hotelApi } from '@/api'
import { toast } from 'vue-sonner'
import { getApiErrorMessage } from '@/lib/utils'
import type { Room, Hotel } from '@/types'

const loading = ref(true)
const rooms = ref<Room[]>([])
const hotels = ref<Hotel[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(1)

const dialogOpen = ref(false)
const editingRoom = ref<Room | null>(null)
const submitting = ref(false)

const form = ref({
  hotelId: null as string | null,
  roomType: '',
  price: null as number | null,
  totalCount: null as number | null,
  availableCount: null as number | null,
  description: '',
})

const isEditing = computed(() => editingRoom.value !== null)

function resetForm() {
  form.value = {
    hotelId: null,
    roomType: '',
    price: null,
    totalCount: null,
    availableCount: null,
    description: '',
  }
}

function openCreateDialog() {
  editingRoom.value = null
  resetForm()
  dialogOpen.value = true
}

function openEditDialog(room: Room) {
  editingRoom.value = room
  form.value = {
    hotelId: room.hotelId,
    roomType: room.roomType,
    price: room.price,
    totalCount: room.totalCount,
    availableCount: room.availableCount,
    description: room.description ?? '',
  }
  dialogOpen.value = true
}

async function loadRooms() {
  loading.value = true
  try {
    const res = await roomApi.list({ page: currentPage.value, pageSize: pageSize.value })
    rooms.value = res.data.data ?? []
    totalPages.value = res.data.pagination?.totalPages ?? 1
  } catch (e: unknown) {
    console.error('Failed to load rooms:', e)
    toast.error(getApiErrorMessage(e, '加载客房列表失败'))
  } finally {
    loading.value = false
  }
}

async function loadHotels() {
  try {
    const res = await hotelApi.list({ page: 1, pageSize: 1000 })
    hotels.value = res.data.data ?? []
  } catch (e: unknown) {
    console.error('Failed to load hotels:', e)
  }
}

async function handleSubmit() {
  if (!form.value.hotelId || !form.value.roomType || form.value.price === null || 
      form.value.totalCount === null || form.value.availableCount === null) {
    toast.error('请填写必填字段')
    return
  }

  submitting.value = true
  try {
    const data = {
      hotelId: form.value.hotelId,
      roomType: form.value.roomType,
      price: form.value.price,
      totalCount: form.value.totalCount,
      availableCount: form.value.availableCount,
      description: form.value.description,
    }

    if (isEditing.value && editingRoom.value) {
      await roomApi.update(editingRoom.value.id, data)
      toast.success('客房更新成功')
    } else {
      await roomApi.create(data)
      toast.success('客房创建成功')
    }

    dialogOpen.value = false
    await loadRooms()
  } catch (e: unknown) {
    console.error('Failed to save room:', e)
    toast.error(getApiErrorMessage(e, '保存客房失败'))
  } finally {
    submitting.value = false
  }
}

async function handleDelete(room: Room) {
  if (!confirm(`确定要删除客房"${room.roomType}"吗？`)) {
    return
  }

  try {
    await roomApi.delete(room.id)
    toast.success('客房删除成功')
    await loadRooms()
  } catch (e: unknown) {
    console.error('Failed to delete room:', e)
    toast.error(getApiErrorMessage(e, '删除客房失败'))
  }
}

function getHotelName(hotelId: string): string {
  const hotel = hotels.value.find((h) => h.id === hotelId)
  return hotel?.hotelName ?? hotelId
}

function formatPrice(price: number): string {
  return `¥${price.toFixed(2)}`
}

onMounted(() => {
  loadRooms()
  loadHotels()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold">客房管理</h1>
      <Button @click="openCreateDialog">
        <Plus class="mr-2 h-4 w-4" />
        添加客房
      </Button>
    </div>

    <Card class="shadow-sm">
      <CardContent class="pt-6">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>酒店名称</TableHead>
              <TableHead>客房类型</TableHead>
              <TableHead>价格</TableHead>
              <TableHead>可用/总数</TableHead>
              <TableHead class="text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <template v-if="loading">
              <TableRow v-for="i in 5" :key="i">
                <TableCell v-for="j in 5" :key="j">
                  <Skeleton class="h-4 w-full" />
                </TableCell>
              </TableRow>
            </template>
            <template v-else-if="rooms.length === 0">
              <TableRow>
                <TableCell :colspan="5" class="text-center text-muted-foreground">
                  暂无数据
                </TableCell>
              </TableRow>
            </template>
            <template v-else>
              <TableRow v-for="room in rooms" :key="room.id">
                <TableCell class="font-medium">{{ getHotelName(room.hotelId) }}</TableCell>
                <TableCell>{{ room.roomType }}</TableCell>
                <TableCell>{{ formatPrice(room.price) }}</TableCell>
                <TableCell>
                  <span :class="room.availableCount > 0 ? 'text-green-600' : 'text-red-600'">
                    {{ room.availableCount }}
                  </span>
                  / {{ room.totalCount }}
                </TableCell>
                <TableCell class="text-right">
                  <div class="flex justify-end gap-2">
                    <Button variant="outline" size="sm" @click="openEditDialog(room)">
                      <Pencil class="mr-1 h-3 w-3" />
                      编辑
                    </Button>
                    <Button variant="destructive" size="sm" @click="handleDelete(room)">
                      <Trash2 class="mr-1 h-3 w-3" />
                      删除
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            </template>
          </TableBody>
        </Table>

        <!-- Pagination -->
        <div v-if="!loading && rooms.length > 0" class="mt-4 flex items-center justify-between">
          <div class="text-sm text-muted-foreground">
            第 {{ currentPage }} / {{ totalPages }} 页
          </div>
          <div class="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage <= 1"
              @click="currentPage--; loadRooms()"
            >
              上一页
            </Button>
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage >= totalPages"
              @click="currentPage++; loadRooms()"
            >
              下一页
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:open="dialogOpen">
      <DialogContent class="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{{ isEditing ? '编辑客房' : '添加客房' }}</DialogTitle>
        </DialogHeader>

        <div class="grid gap-4 py-4">
          <div class="grid gap-2">
            <Label>所属酒店 *</Label>
            <Select
              :model-value="form.hotelId"
              @update:model-value="(v) => form.hotelId = v as string"
            >
              <SelectTrigger>
                <SelectValue placeholder="请选择酒店" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="hotel in hotels" :key="hotel.id" :value="hotel.id">
                  {{ hotel.hotelName }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="grid gap-2">
            <Label for="roomType">客房类型 *</Label>
            <Input
              id="roomType"
              v-model="form.roomType"
              placeholder="请输入客房类型"
            />
          </div>

          <div class="grid gap-2">
            <Label for="price">价格 *</Label>
            <Input
              id="price"
              type="number"
              :model-value="form.price?.toString() ?? ''"
              @update:model-value="(v) => form.price = v ? Number(v) : null"
              placeholder="请输入价格"
              :min="0"
              :step="0.01"
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="grid gap-2">
              <Label for="totalCount">总数 *</Label>
              <Input
                id="totalCount"
                type="number"
                :model-value="form.totalCount?.toString() ?? ''"
                @update:model-value="(v) => form.totalCount = v ? Number(v) : null"
                placeholder="总数"
                :min="0"
              />
            </div>

            <div class="grid gap-2">
              <Label for="availableCount">可用数 *</Label>
              <Input
                id="availableCount"
                type="number"
                :model-value="form.availableCount?.toString() ?? ''"
                @update:model-value="(v) => form.availableCount = v ? Number(v) : null"
                placeholder="可用数"
                :min="0"
              />
            </div>
          </div>

          <div class="grid gap-2">
            <Label for="description">描述</Label>
            <Textarea
              id="description"
              v-model="form.description"
              placeholder="请输入客房描述"
              :rows="3"
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="dialogOpen = false">取消</Button>
          <Button @click="handleSubmit" :disabled="submitting">
            {{ submitting ? '保存中...' : '保存' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
