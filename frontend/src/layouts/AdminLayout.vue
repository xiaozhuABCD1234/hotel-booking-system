<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores'
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  LayoutDashboard,
  Hotel,
  DoorOpen,
  ShoppingBag,
  Users,
  BarChart3,
  LogOut,
  Building2,
  Menu,
  X,
} from '@lucide/vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const sidebarOpen = ref(false)

const navItems = [
  { to: '/admin', label: '仪表盘', icon: LayoutDashboard },
  { to: '/admin/hotels', label: '酒店管理', icon: Hotel },
  { to: '/admin/rooms', label: '客房管理', icon: DoorOpen },
  { to: '/admin/orders', label: '订单管理', icon: ShoppingBag },
  { to: '/admin/users', label: '用户管理', icon: Users },
  { to: '/admin/reports', label: '数据报表', icon: BarChart3 },
]

function isActive(path: string) {
  if (path === '/admin') return route.path === '/admin'
  return route.path.startsWith(path)
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Sidebar overlay (mobile) -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 z-40 bg-black/50 lg:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Sidebar -->
    <aside
      class="fixed inset-y-0 left-0 z-50 w-64 border-r bg-card flex flex-col transition-transform lg:translate-x-0"
      :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Logo -->
      <div class="flex h-16 items-center gap-2 px-6 border-b">
        <Building2 class="h-6 w-6 text-primary" />
        <span class="font-semibold">后台管理</span>
        <Button variant="ghost" size="icon" class="ml-auto lg:hidden" @click="sidebarOpen = false">
          <X class="h-4 w-4" />
        </Button>
      </div>

      <!-- Nav -->
      <nav class="flex-1 overflow-auto p-4 space-y-1">
        <router-link
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors"
          :class="isActive(item.to) ? 'bg-primary/10 text-primary font-medium' : 'text-muted-foreground hover:text-foreground hover:bg-accent'"
          @click="sidebarOpen = false"
        >
          <component :is="item.icon" class="h-4 w-4" />
          {{ item.label }}
        </router-link>
      </nav>

      <Separator />
      <!-- User footer -->
      <div class="p-4">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" class="w-full justify-start gap-3 px-2">
              <Avatar class="h-8 w-8">
                <AvatarFallback class="bg-primary/10 text-primary text-xs">
                  {{ auth.user?.username?.charAt(0)?.toUpperCase() || 'A' }}
                </AvatarFallback>
              </Avatar>
              <div class="flex flex-col items-start text-sm">
                <span class="font-medium">{{ auth.user?.username }}</span>
                <span class="text-xs text-muted-foreground">管理员</span>
              </div>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" class="w-48">
            <DropdownMenuItem @click="router.push('/')">
              <Building2 class="mr-2 h-4 w-4" />
              返回前台
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="handleLogout" class="text-destructive">
              <LogOut class="mr-2 h-4 w-4" />
              退出登录
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>

    <!-- Main content -->
    <div class="lg:pl-64">
      <!-- Top bar -->
      <header class="sticky top-0 z-30 flex h-16 items-center gap-4 border-b bg-background/95 backdrop-blur px-4 lg:px-6">
        <Button variant="ghost" size="icon" class="lg:hidden" @click="sidebarOpen = true">
          <Menu class="h-5 w-5" />
        </Button>
        <div class="flex-1" />
        <Button variant="ghost" size="sm" @click="router.push('/')">
          <Building2 class="mr-2 h-4 w-4" />
          返回前台
        </Button>
      </header>

      <!-- Page content -->
      <main class="p-4 lg:p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>
