<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
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
import { Plus, Pencil, Trash2, Star } from '@lucide/vue'
import { hotelApi, regionApi } from '@/api'
import { toast } from 'vue-sonner'
import { getApiErrorMessage } from '@/lib/utils'
import type { Hotel, Region } from '@/types'

const loading = ref(true)
const hotels = ref<Hotel[]>([])
const regions = ref<Region[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(1)

const dialogOpen = ref(false)
const editingHotel = ref<Hotel | null>(null)
const submitting = ref(false)

const form = ref({
  hotelName: '',
  regionId: null as number | null,
  starLevel: null as number | null,
  address: '',
  telephone: '',
  description: '',
})

const isEditing = computed(() => editingHotel.value !== null)

const starLevels = [1, 2, 3, 4, 5]

function resetForm() {
  form.value = {
    hotelName: '',
    regionId: null,
    starLevel: null,
    address: '',
    telephone: '',
    description: '',
  }
}

function openCreateDialog() {
  editingHotel.value = null
  resetForm()
  dialogOpen.value = true
}

function openEditDialog(hotel: Hotel) {
  editingHotel.value = hotel
  form.value = {
    hotelName: hotel.hotelName,
    regionId: hotel.regionId,
    starLevel: hotel.starLevel ?? null,
    address: hotel.address,
    telephone: hotel.telephone,
    description: hotel.description ?? '',
  }
  dialogOpen.value = true
}

async function loadHotels() {
  loading.value = true
  try {
    const res = await hotelApi.list({ page: currentPage.value, pageSize: pageSize.value })
    hotels.value = res.data.data?.items ?? []
    totalPages.value = res.data.pagination?.totalPages ?? 1
  } catch (error) {
    console.error('Failed to load hotels:', error)
    toast.error(getApiErrorMessage(error, '加载酒店列表失败'))
  } finally {
    loading.value = false
  }
}

async function loadRegions() {
  try {
    const res = await regionApi.list()
    regions.value = (res.data.data as Region[] | undefined) ?? []
  } catch (error) {
    console.error('Failed to load regions:', error)
  }
}

async function handleSubmit() {
  if (!form.value.hotelName || !form.value.regionId || !form.value.starLevel) {
    toast.error('请填写必填字段')
    return
  }

  submitting.value = true
  try {
    const data = {
      hotelName: form.value.hotelName,
      regionId: form.value.regionId,
      starLevel: form.value.starLevel,
      address: form.value.address,
      telephone: form.value.telephone,
      description: form.value.description,
    }

    if (isEditing.value && editingHotel.value) {
      await hotelApi.update(editingHotel.value.id, data)
      toast.success('酒店更新成功')
    } else {
      await hotelApi.create(data)
      toast.success('酒店创建成功')
    }

    dialogOpen.value = false
    await loadHotels()
  } catch (error) {
    console.error('Failed to save hotel:', error)
    toast.error(getApiErrorMessage(error, '保存酒店失败'))
  } finally {
    submitting.value = false
  }
}

async function handleDelete(hotel: Hotel) {
  if (!confirm(`确定要删除酒店"${hotel.hotelName}"吗？`)) {
    return
  }

  try {
    await hotelApi.delete(hotel.id)
    toast.success('酒店删除成功')
    await loadHotels()
  } catch (error) {
    console.error('Failed to delete hotel:', error)
    toast.error(getApiErrorMessage(error, '删除酒店失败'))
  }
}

function getRegionName(regionId: number): string {
  const region = regions.value.find((r) => r.id === regionId)
  return region?.name ?? String(regionId)
}

onMounted(() => {
  loadHotels()
  loadRegions()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold">酒店管理</h1>
      <Button @click="openCreateDialog">
        <Plus class="mr-2 h-4 w-4" />
        添加酒店
      </Button>
    </div>

    <Card class="shadow-sm">
      <CardContent class="pt-6">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>酒店名称</TableHead>
              <TableHead>地区</TableHead>
              <TableHead>星级</TableHead>
              <TableHead>地址</TableHead>
              <TableHead>电话</TableHead>
              <TableHead class="text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <template v-if="loading">
              <TableRow v-for="i in 5" :key="i">
                <TableCell v-for="j in 6" :key="j">
                  <Skeleton class="h-4 w-full" />
                </TableCell>
              </TableRow>
            </template>
            <template v-else-if="hotels.length === 0">
              <TableRow>
                <TableCell :colspan="6" class="text-center text-muted-foreground">
                  暂无数据
                </TableCell>
              </TableRow>
            </template>
            <template v-else>
              <TableRow v-for="hotel in hotels" :key="hotel.id">
                <TableCell class="font-medium">{{ hotel.hotelName }}</TableCell>
                <TableCell>{{ getRegionName(hotel.regionId) }}</TableCell>
                <TableCell>
                  <Badge v-if="hotel.starLevel" variant="secondary">
                    <Star class="mr-1 h-3 w-3" />
                    {{ hotel.starLevel }}星
                  </Badge>
                </TableCell>
                <TableCell>{{ hotel.address }}</TableCell>
                <TableCell>{{ hotel.telephone }}</TableCell>
                <TableCell class="text-right">
                  <div class="flex justify-end gap-2">
                    <Button variant="outline" size="sm" @click="openEditDialog(hotel)">
                      <Pencil class="mr-1 h-3 w-3" />
                      编辑
                    </Button>
                    <Button variant="destructive" size="sm" @click="handleDelete(hotel)">
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
        <div v-if="!loading && hotels.length > 0" class="mt-4 flex items-center justify-between">
          <div class="text-sm text-muted-foreground">
            第 {{ currentPage }} / {{ totalPages }} 页
          </div>
          <div class="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage <= 1"
              @click="currentPage--; loadHotels()"
            >
              上一页
            </Button>
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage >= totalPages"
              @click="currentPage++; loadHotels()"
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
          <DialogTitle>{{ isEditing ? '编辑酒店' : '添加酒店' }}</DialogTitle>
        </DialogHeader>

        <div class="grid gap-4 py-4">
          <div class="grid gap-2">
            <Label for="hotelName">酒店名称 *</Label>
            <Input
              id="hotelName"
              v-model="form.hotelName"
              placeholder="请输入酒店名称"
            />
          </div>

          <div class="grid gap-2">
            <Label>地区 *</Label>
            <Select
              :model-value="form.regionId?.toString()"
              @update:model-value="form.regionId = Number($event)"
            >
              <SelectTrigger>
                <SelectValue placeholder="请选择地区" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="region in regions" :key="region.id" :value="region.id.toString()">
                  {{ region.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="grid gap-2">
            <Label>星级 *</Label>
            <Select
              :model-value="form.starLevel?.toString()"
              @update:model-value="form.starLevel = Number($event)"
            >
              <SelectTrigger>
                <SelectValue placeholder="请选择星级" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="level in starLevels" :key="level" :value="level.toString()">
                  {{ level }}星
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="grid gap-2">
            <Label for="address">地址</Label>
            <Input
              id="address"
              v-model="form.address"
              placeholder="请输入地址"
            />
          </div>

          <div class="grid gap-2">
            <Label for="telephone">电话</Label>
            <Input
              id="telephone"
              v-model="form.telephone"
              placeholder="请输入电话"
            />
          </div>

          <div class="grid gap-2">
            <Label for="description">描述</Label>
            <Textarea
              id="description"
              v-model="form.description"
              placeholder="请输入酒店描述"
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
