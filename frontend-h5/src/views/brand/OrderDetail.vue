<template>
  <div class="order-detail-page">
    <div class="top-nav">
      <button class="back-btn" @click="goBack">返回</button>
      <h1 class="nav-title">订单详情</h1>
      <span class="nav-placeholder"></span>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="errorMessage" class="error-box">
      <p>{{ errorMessage }}</p>
      <button class="retry-btn" @click="loadOrderDetail">重试</button>
    </div>

    <div v-else class="content">
      <div class="detail-card">
        <div v-for="item in detailItems" :key="item.label" class="detail-item">
          <span class="label">{{ item.label }}</span>
          <span class="value">{{ item.value }}</span>
        </div>
      </div>

      <div class="detail-card" v-if="formEntries.length">
        <h2 class="section-title">表单信息</h2>
        <div v-for="entry in formEntries" :key="entry.key" class="detail-item">
          <span class="label">{{ entry.key }}</span>
          <span class="value">{{ entry.value }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { orderApi } from '../../services/brandApi.js'
import { getOrderDetailItems } from './orders.logic.js'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const errorMessage = ref('')
const orderDetail = ref(null)

const detailItems = computed(() => getOrderDetailItems(orderDetail.value || {}))

const formEntries = computed(() => {
  const formData = orderDetail.value?.formData || {}
  return Object.entries(formData).map(([key, value]) => ({ key, value }))
})

const goBack = () => {
  router.back()
}

const loadOrderDetail = async () => {
  const orderId = route.params.id
  if (!orderId) {
    errorMessage.value = '订单 ID 无效'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await orderApi.getOrder(orderId)
    const payload = response?.data || response
    orderDetail.value = payload || {}
  } catch (error) {
    console.error('加载订单详情失败:', error)
    errorMessage.value = '加载订单详情失败，请重试'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOrderDetail()
})
</script>

<style scoped>
.order-detail-page { min-height: 100vh; background: #f5f7fa; }
.top-nav { background: #fff; padding: 16px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #eee; }
.back-btn { border: none; background: none; color: #667eea; }
.nav-title { margin: 0; font-size: 18px; }
.nav-placeholder { width: 32px; }
.loading, .error-box { padding: 24px; text-align: center; color: #666; }
.retry-btn { margin-top: 12px; border: 1px solid #667eea; color: #667eea; background: #fff; border-radius: 6px; padding: 8px 14px; }
.content { padding: 16px; display: grid; gap: 12px; }
.detail-card { background: #fff; border-radius: 10px; padding: 14px; }
.section-title { margin: 0 0 10px; font-size: 15px; }
.detail-item { display: flex; justify-content: space-between; gap: 12px; padding: 8px 0; border-bottom: 1px solid #f2f2f2; }
.detail-item:last-child { border-bottom: none; }
.label { color: #666; font-size: 13px; }
.value { color: #333; font-size: 13px; text-align: right; }
</style>
