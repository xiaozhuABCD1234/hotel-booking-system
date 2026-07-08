<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores'
import type { LoginRequest } from '@/types'
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

const form = reactive<LoginRequest>({
  username: '',
  password: '',
})

const errors = reactive({
  username: '',
  password: '',
})

const loading = ref(false)

function validate(): boolean {
  errors.username = ''
  errors.password = ''
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

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  loading.value = true
  try {
    const res = await auth.login(form)
    if (res.success) {
      const redirect = (route.query.redirect as string) || '/'
      await router.push(redirect)
    } else {
      toast.error(res.error?.message || '登录失败，请重试')
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '登录失败，请重试'))
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
          <CardTitle class="text-2xl">登录</CardTitle>
          <CardDescription>欢迎回来，请登录您的账号</CardDescription>
        </CardHeader>

        <form @submit.prevent="handleSubmit">
          <CardContent class="space-y-4">
            <!-- Username -->
            <div class="space-y-2">
              <Label for="login-username">
                用户名 <span class="text-destructive">*</span>
              </Label>
              <Input
                id="login-username"
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
              <Label for="login-password">
                密码 <span class="text-destructive">*</span>
              </Label>
              <Input
                id="login-password"
                v-model="form.password"
                type="password"
                placeholder="请输入密码"
                autocomplete="current-password"
                :disabled="loading"
              />
              <p v-if="errors.password" class="text-sm text-destructive">
                {{ errors.password }}
              </p>
            </div>
          </CardContent>

          <CardFooter class="flex flex-col gap-4">
            <Button type="submit" class="w-full" :disabled="loading">
              <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
              {{ loading ? '登录中...' : '登录' }}
            </Button>
            <p class="text-sm text-muted-foreground">
              没有账号？
              <router-link
                to="/register"
                class="font-medium text-blue-900 hover:underline"
              >
                立即注册
              </router-link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  </div>
</template>
