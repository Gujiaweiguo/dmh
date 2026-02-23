<template>
  <div class="reward-records-page">
    <div class="top-nav">
      <button class="back-btn" @click="goBack">返回</button>
      <h1 class="nav-title">奖励记录</h1>
      <span class="nav-placeholder"></span>
    </div>

    <!-- 筛选器 -->
    <div class="filters">
      <div class="filter-row">
        <select v-model="filterType" class="filter-select">
          <option value="all">全部类型</option>
          <option value="commission">佣金</option>
          <option value="bonus">奖金</option>
          <option value="penalty">扣款</option>
        </select>
        <select v-model="filterStatus" class="filter-select">
          <option value="all">全部状态</option>
          <option value="pending">待发放</option>
          <option value="paid">已发放</option>
          <option value="cancelled">已取消</option>
        </select>
      </div>
    </div>

    <!-- 统计摘要 -->
    <div class="stats-summary">
      <div class="stat-item">
        <span class="stat-value">¥{{ summary.totalAmount }}</span>
        <span class="stat-label">累计奖励</span>
      </div>
      <div class="stat-item">
        <span class="stat-value">¥{{ summary.pendingAmount }}</span>
        <span class="stat-label">待发放</span>
      </div>
      <div class="stat-item">
        <span class="stat-value">{{ summary.totalCount }}</span>
        <span class="stat-label">记录数</span>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="filteredRecords.length === 0" class="empty-state">
      <p>暂无奖励记录</p>
    </div>

    <div v-else class="records-list">
      <div
        v-for="record in filteredRecords"
        :key="record.id"
        class="record-card"
      >
        <div class="record-header">
          <span :class="['type-badge', record.type]">
            {{ getTypeText(record.type) }}
          </span>
          <span :class="['status-badge', record.status]">
            {{ getStatusText(record.status) }}
          </span>
        </div>
        <div class="record-body">
          <div class="record-amount">
            <span :class="['amount', record.amount >= 0 ? 'positive' : 'negative']">
              {{ record.amount >= 0 ? '+' : '' }}¥{{ record.amount }}
            </span>
          </div>
          <div class="record-info">
            <p class="record-desc">{{ record.description }}</p>
            <p class="record-meta">
              <span v-if="record.promoterName">推广员: {{ record.promoterName }}</span>
              <span v-if="record.campaignName">活动: {{ record.campaignName }}</span>
            </p>
            <p class="record-time">{{ record.createdAt }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { promoterApi } from '../../services/brandApi.js'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const records = ref([])
const filterType = ref('all')
const filterStatus = ref('all')

const summary = computed(() => {
  const totalAmount = records.value.reduce((sum, r) => sum + (r.amount || 0), 0)
  const pendingAmount = records.value
    .filter(r => r.status === 'pending')
    .reduce((sum, r) => sum + (r.amount || 0), 0)
  return {
    totalAmount: totalAmount.toFixed(2),
    pendingAmount: pendingAmount.toFixed(2),
    totalCount: records.value.length
  }
})

const filteredRecords = computed(() => {
  let result = records.value
  if (filterType.value !== 'all') {
    result = result.filter(r => r.type === filterType.value)
  }
  if (filterStatus.value !== 'all') {
    result = result.filter(r => r.status === filterStatus.value)
  }
  return result
})

const getTypeText = (type) => {
  const map = {
    commission: '佣金',
    bonus: '奖金',
    penalty: '扣款'
  }
  return map[type] || type
}

const getStatusText = (status) => {
  const map = {
    pending: '待发放',
    paid: '已发放',
    cancelled: '已取消'
  }
  return map[status] || status
}

const goBack = () => {
  router.back()
}

const loadRecords = async () => {
  const promoterId = route.params.promoterId
  loading.value = true

  try {
    let response
    if (promoterId) {
      response = await promoterApi.getRewards(promoterId)
    } else {
      // 获取所有推广员的奖励记录
      response = await promoterApi.getRewards('all')
    }

    const payload = response?.data || response
    records.value = Array.isArray(payload) ? payload : (payload?.list || [])
  } catch (error) {
    console.error('加载奖励记录失败:', error)
    records.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRecords()
})
</script>

<style scoped>
.reward-records-page { min-height: 100vh; background: #f5f7fa; padding-bottom: 20px; }
.top-nav { background: #fff; padding: 16px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #eee; }
.back-btn { border: none; background: none; color: #667eea; font-size: 14px; }
.nav-title { margin: 0; font-size: 18px; }
.nav-placeholder { width: 40px; }
.filters { background: #fff; padding: 12px 16px; border-bottom: 1px solid #eee; }
.filter-row { display: flex; gap: 12px; }
.filter-select { flex: 1; padding: 8px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; }
.stats-summary { display: flex; padding: 16px; gap: 12px; }
.stat-item { flex: 1; background: #fff; padding: 12px; border-radius: 12px; text-align: center; }
.stat-value { display: block; font-size: 18px; font-weight: bold; color: #667eea; }
.stat-label { font-size: 12px; color: #666; }
.loading, .empty-state { padding: 40px 20px; text-align: center; color: #666; }
.records-list { padding: 0 16px; }
.record-card { background: #fff; border-radius: 12px; padding: 16px; margin-bottom: 12px; }
.record-header { display: flex; justify-content: space-between; margin-bottom: 12px; }
.type-badge { padding: 2px 8px; border-radius: 12px; font-size: 11px; }
.type-badge.commission { background: #e3f2fd; color: #1976d2; }
.type-badge.bonus { background: #fff3e0; color: #f57c00; }
.type-badge.penalty { background: #fce4ec; color: #e91e63; }
.status-badge { padding: 2px 8px; border-radius: 12px; font-size: 11px; }
.status-badge.pending { background: #fff8e1; color: #ff8f00; }
.status-badge.paid { background: #e8f5e8; color: #4caf50; }
.status-badge.cancelled { background: #f5f5f5; color: #999; }
.record-body { display: flex; align-items: flex-start; }
.record-amount { min-width: 80px; }
.amount { font-size: 18px; font-weight: bold; }
.amount.positive { color: #4caf50; }
.amount.negative { color: #e91e63; }
.record-info { flex: 1; margin-left: 12px; }
.record-desc { margin: 0 0 4px; font-size: 14px; color: #333; }
.record-meta { margin: 0 0 4px; font-size: 12px; color: #666; }
.record-meta span { margin-right: 12px; }
.record-time { margin: 0; font-size: 11px; color: #999; }
</style>
