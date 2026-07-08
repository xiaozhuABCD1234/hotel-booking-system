<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores'
import type { RegisterRequest } from '@/types'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Building2, Loader2 } from '@lucide/vue'
import { toast } from 'vue-sonner'
import { getApiErrorMessage } from '@/lib/utils'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  phone: '',
  email: '',
})

const errors = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  phone: '',
  email: '',
})

const loading = ref(false)

function validate(): boolean {
  errors.username = ''
  errors.password = ''
  errors.confirmPassword = ''
  errors.phone = ''
  errors.email = ''
  let valid = true

  if (!form.username.trim()) {
    errors.username = '请输入用户名'
    valid = false
  }

  if (!form.password) {
    errors.password = '请输入密码'
    valid = false
  } else if (form.password.length < 6) {
    errors.password = '密码至少6个字符'
    valid = false
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = '请确认密码'
    valid = false
  } else if (form.confirmPassword !== form.password) {
    errors.confirmPassword = '两次输入的密码不一致'
    valid = false
  }

  if (form.phone && !/^1\d{10}$/.test(form.phone)) {
    errors.phone = '请输入有效的手机号'
    valid = false
  }

  if (form.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = '请输入有效的邮箱地址'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  loading.value = true
  try {
    const data: RegisterRequest = {
      username: form.username,
      password: form.password,
    }
    if (form.phone) data.phone = form.phone
    if (form.email) data.email = form.email

    const res = await auth.register(data)
    if (res.success) {
      const redirect = (route.query.redirect as string) || '/'
      router.push(redirect)
    } else {
      toast.error(res.error?.message || '注册失败，请重试')
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '注册失败，请重试'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-dvh items-center justify-center bg-gray-50 px-4 py-12">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="mb-8 flex flex-col items-center gap-2">
        <Building2 class="h-10 w-10 text-blue-900" />
        <h1 class="text-xl font-bold text-blue-900">酒店预订管理系统</h1>
      </div>

      <Card class="shadow-lg">
        <CardHeader class="text-center">
          <CardTitle class="text-2xl">注册</CardTitle>
          <CardDescription>创建您的账号，开始预订酒店</CardDescription>
        </CardHeader>

        <form @submit.prevent="handleSubmit">
          <CardContent class="space-y-4">
            <!-- Username -->
            <div class="space-y-2">
              <Label for="reg-username">
                用户名 <span class="text-destructive">*</span>
              </Label>
              <Input
                id="reg-username"
                v-model="form.username"
                type="text"
                placeholder="请输入用户名"
                autocomplete="username"
                :disabled="loading"
              />
              <p v-if="errors.username" class="text-sm text-destructive">
                {{ errors.username }}
              </p>
            </div>

            <!-- Password -->
            <div class="space-y-2">
              <Label for="reg-password">
                密码 <span class="text-destructive">*</span>
              </Label>
              <Input
                id="reg-password"
                v-model="form.password"
                type="password"
                placeholder="请输入密码（至少6个字符）"
                autocomplete="new-password"
                :disabled="loading"
              />
              <p v-if="errors.password" class="text-sm text-destructive">
                {{ errors.password }}
              </p>
            </div>

            <!-- Confirm Password -->
            <div class="space-y-2">
              <Label for="reg-confirm-password">
                确认密码 <span class="text-destructive">*</span>
              </Label>
              <Input
                id="reg-confirm-password"
                v-model="form.confirmPassword"
                type="password"
                placeholder="请再次输入密码"
                autocomplete="new-password"
                :disabled="loading"
              />
              <p v-if="errors.confirmPassword" class="text-sm text-destructive">
                {{ errors.confirmPassword }}
              </p>
            </div>

            <!-- Phone (optional) -->
            <div class="space-y-2">
              <Label for="reg-phone">手机号</Label>
              <Input
                id="reg-phone"
                v-model="form.phone"
                type="tel"
                placeholder="选填"
                autocomplete="tel"
                :disabled="loading"
              />
              <p v-if="errors.phone" class="text-sm text-destructive">
                {{ errors.phone }}
              </p>
            </div>

            <!-- Email (optional) -->
            <div class="space-y-2">
              <Label for="reg-email">邮箱</Label>
              <Input
                id="reg-email"
                v-model="form.email"
                type="email"
                placeholder="选填"
                autocomplete="email"
                :disabled="loading"
              />
              <p v-if="errors.email" class="text-sm text-destructive">
                {{ errors.email }}
              </p>
            </div>
          </CardContent>

          <CardFooter class="flex flex-col gap-4">
            <Button type="submit" class="w-full" :disabled="loading">
              <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
              {{ loading ? '注册中...' : '注册' }}
            </Button>
            <p class="text-sm text-muted-foreground">
              已有账号？
              <router-link
                to="/login"
                class="font-medium text-blue-900 hover:underline"
              >
                去登录
              </router-link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  </div>
</template>
