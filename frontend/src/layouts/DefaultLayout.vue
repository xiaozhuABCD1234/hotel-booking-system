<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Building2, LogOut, User, ShoppingBag, Hotel } from '@lucide/vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const navLinks = [
  { to: '/', label: '首页', icon: Hotel },
  { to: '/orders', label: '我的订单', icon: ShoppingBag },
]

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Header -->
    <header class="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div class="container mx-auto flex h-16 items-center justify-between px-4">
        <!-- Logo -->
        <router-link to="/" class="flex items-center gap-2 font-semibold text-lg">
          <Building2 class="h-6 w-6 text-primary" />
          <span>酒店预订</span>
        </router-link>

        <!-- Desktop Nav -->
        <nav class="hidden md:flex items-center gap-6">
          <router-link
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
            :class="{ 'text-foreground': route.path === link.to }"
          >
            {{ link.label }}
          </router-link>
        </nav>

        <!-- User Menu -->
        <div class="flex items-center gap-3">
          <template v-if="auth.isLoggedIn">
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button variant="ghost" size="icon" class="rounded-full">
                  <Avatar class="h-8 w-8">
                    <AvatarFallback class="bg-primary/10 text-primary text-xs">
                      {{ auth.user?.username?.charAt(0)?.toUpperCase() || 'U' }}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" class="w-48">
                <div class="px-2 py-1.5">
                  <p class="text-sm font-medium">{{ auth.user?.username }}</p>
                  <p class="text-xs text-muted-foreground">{{ auth.user?.role === 'admin' ? '管理员' : '用户' }}</p>
                </div>
                <DropdownMenuSeparator />
                <DropdownMenuItem v-if="auth.isAdmin" @click="router.push('/admin')">
                  <User class="mr-2 h-4 w-4" />
                  后台管理
                </DropdownMenuItem>
                <DropdownMenuItem @click="router.push('/orders')">
                  <ShoppingBag class="mr-2 h-4 w-4" />
                  我的订单
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem @click="handleLogout" class="text-destructive">
                  <LogOut class="mr-2 h-4 w-4" />
                  退出登录
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </template>
          <template v-else>
            <Button variant="ghost" size="sm" @click="router.push('/login')">登录</Button>
            <Button size="sm" @click="router.push('/register')">注册</Button>
          </template>
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <main>
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="border-t py-8 mt-16">
      <div class="container mx-auto px-4 text-center text-sm text-muted-foreground">
        酒店预订管理系统 &copy; {{ new Date().getFullYear() }}
      </div>
    </footer>
  </div>
</template>
