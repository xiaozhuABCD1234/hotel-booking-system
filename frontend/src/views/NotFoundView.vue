<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Home } from '@lucide/vue'

const router = useRouter()
</script>

<template>
  <div class="flex min-h-[70vh] flex-col items-center justify-center px-4 text-center">
    <h1 class="mb-2 text-7xl font-bold text-blue-900">404</h1>
    <p class="mb-1 text-xl font-medium text-gray-800">页面未找到</p>
    <p class="mb-8 text-sm text-gray-500">
      您访问的页面不存在或已被移除
    </p>
    <Button
      class="bg-blue-900 hover:bg-blue-800"
      @click="router.push('/')"
    >
      <Home class="mr-2 h-4 w-4" />
      返回首页
    </Button>
  </div>
</template>
