<template>
  <div class="brand-orders">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">订单管理</h1>
      <div class="nav-stats">
        <span class="total-orders">总计: {{ totalOrders }}</span>
      </div>
    </div>

    <!-- 筛选器 -->
    <div class="filters">
      <div class="filter-tabs">
        <button
          v-for="status in statusTabs"
          :key="status.value"
          @click="currentStatus = status.value"
          :class="['filter-tab', { active: currentStatus === status.value }]"
        >
          {{ status.label }}
        </button>
      </div>

      <div class="filter-actions">
        <button class="export-all-btn" @click="exportFilteredOrders">导出当前筛选</button>
      </div>
      
      <div class="date-filter">
        <input
          v-model="dateRange.start"
          type="date"
          class="date-input"
          placeholder="开始日期"
        >
        <span class="date-separator">至</span>
        <input
          v-model="dateRange.end"
          type="date"
          class="date-input"
          placeholder="结束日期"
        >
      </div>
    </div>

    <!-- 订单列表 -->
    <div class="orders-list">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="filteredOrders.length === 0" class="empty-state">
        <div class="empty-icon">📋</div>
        <p class="empty-text">暂无订单数据</p>
      </div>

      <div v-else class="order-cards">
        <div
          v-for="order in filteredOrders"
          :key="order.id"
          class="order-card"
          @click="viewOrderDetail(order)"
        >
          <div class="card-header">
            <div class="order-info">
              <h3 class="order-id">订单 #{{ order.id }}</h3>
              <span class="order-time">{{ formatDateTime(order.createdAt) }}</span>
            </div>
            <span :class="['status-badge', order.status]">
              {{ getStatusText(order.status) }}
            </span>
          </div>

          <div class="campaign-info">
            <h4 class="campaign-name">{{ order.campaignName }}</h4>
            <p class="user-info">用户: {{ order.phone }}</p>
          </div>

          <div class="order-details">
            <div class="detail-row">
              <span class="detail-label">订单金额:</span>
              <span class="detail-value amount">¥{{ order.amount }}</span>
            </div>
            <div v-if="order.referrerId" class="detail-row">
              <span class="detail-label">推荐人:</span>
              <span class="detail-value">{{ order.referrerName || `用户${order.referrerId}` }}</span>
            </div>
            <div v-if="order.rewardAmount" class="detail-row">
              <span class="detail-label">奖励金额:</span>
              <span class="detail-value reward">¥{{ order.rewardAmount }}</span>
            </div>
          </div>

          <div class="form-data" v-if="order.formData && Object.keys(order.formData).length > 0">
            <div class="form-data-title">用户信息:</div>
            <div class="form-data-content">
              <span
                v-for="(value, key) in order.formData"
                :key="key"
                class="form-data-item"
              >
                {{ key }}: {{ value }}
              </span>
            </div>
          </div>

          <div class="card-actions">
            <button
              v-if="order.status === 'pending'"
              @click.stop="processOrder(order, 'paid')"
              class="action-btn confirm"
            >
              确认支付
            </button>
            <button
              @click.stop="exportOrder(order)"
              class="action-btn export"
            >
              导出
            </button>
            <button
              @click.stop="viewOrderDetail(order)"
              class="action-btn detail"
            >
              详情
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="stats-section">
      <h2 class="stats-title">订单统计</h2>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-number">{{ orderStats.total }}</div>
          <div class="stat-label">总订单数</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">¥{{ orderStats.totalAmount }}</div>
          <div class="stat-label">总金额</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">¥{{ orderStats.totalRewards }}</div>
          <div class="stat-label">总奖励</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ orderStats.todayOrders }}</div>
          <div class="stat-label">今日订单</div>
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <router-link to="/brand/dashboard" class="nav-item">
        <div class="nav-icon">🏠</div>
        <div class="nav-text">工作台</div>
      </router-link>
      <router-link to="/brand/campaigns" class="nav-item">
        <div class="nav-icon">🎯</div>
        <div class="nav-text">活动</div>
      </router-link>
      <router-link to="/brand/orders" class="nav-item active">
        <div class="nav-icon">📋</div>
        <div class="nav-text">订单</div>
      </router-link>
      <router-link to="/brand/distributors" class="nav-item">
        <div class="nav-icon">🧭</div>
        <div class="nav-text">分销</div>
      </router-link>
      <router-link to="/brand/members" class="nav-item">
        <div class="nav-icon">👤</div>
        <div class="nav-text">会员</div>
      </router-link>
      <router-link to="/brand/promoters" class="nav-item">
        <div class="nav-icon">👥</div>
        <div class="nav-text">推广员</div>
      </router-link>
      <router-link to="/brand/settings" class="nav-item">
        <div class="nav-icon">⚙️</div>
        <div class="nav-text">设置</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  applyOrderStatus,
  buildOrdersCsv,
  buildExportOrderData,
  calculateOrderStats,
  filterAndSortOrders,
  formatOrderDateTime,
  getOrderStatusText,
} from './orders.logic.js'
import { orderApi } from '@/services/brandApi.js'

const router = useRouter()

const orders = ref([])
const loading = ref(false)
const currentStatus = ref('all')

const dateRange = reactive({
  start: '',
  end: ''
})

const statusTabs = [
  { value: 'all', label: '全部' },
  { value: 'pending', label: '待支付' },
  { value: 'paid', label: '已支付' },
  { value: 'cancelled', label: '已取消' }
]

const orderStats = reactive({
  total: 0,
  totalAmount: 0,
  totalRewards: 0,
  todayOrders: 0
})

const totalOrders = computed(() => orders.value.length)

const filteredOrders = computed(() => {
  return filterAndSortOrders(orders.value, currentStatus.value, dateRange)
})

const getStatusText = (status) => {
  return getOrderStatusText(status)
}

const formatDateTime = (dateString) => {
  return formatOrderDateTime(dateString)
}

const loadOrders = async () => {
  loading.value = true
  try {
    const response = await orderApi.getOrders()
    // Map API response to component format
    const orderList = response.data?.list || response.list || response.data || []
    orders.value = orderList.map(order => ({
      id: order.id,
      campaignId: order.campaignId,
      campaignName: order.campaignName || order.campaign?.name || '未知活动',
      phone: order.phone,
      amount: order.amount || 0,
      status: order.status || 'pending',
      referrerId: order.referrerId,
      referrerName: order.referrerName || '',
      rewardAmount: order.rewardAmount || 0,
      createdAt: order.createdAt,
      formData: order.formData || {}
    }))
    calculateStats()
  } catch (error) {
    console.error('加载订单失败:', error)
    // 如果API调用失败，保留空列表
    orders.value = []
  } finally {
    loading.value = false
  }
}

const calculateStats = () => {
  const stats = calculateOrderStats(orders.value)
  orderStats.total = stats.total
  orderStats.totalAmount = stats.totalAmount
  orderStats.totalRewards = stats.totalRewards
  orderStats.todayOrders = stats.todayOrders
}

const processOrder = async (order, newStatus) => {
  try {
    await orderApi.updateOrderStatus(order.id, newStatus)
    const next = applyOrderStatus(order, newStatus)
    Object.assign(order, next)
    calculateStats()
  } catch (error) {
    console.error('处理订单失败:', error)
    alert('处理订单失败')
  }
}











const downloadCsv = (fileName, csvContent) => {
  const blob = new Blob([`\ufeff${csvContent}`], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

const exportOrder = (order) => {
  const row = buildExportOrderData(order)
  const headers = Object.keys(row)
  const values = headers.map((key) => JSON.stringify(row[key] ?? ''))
  const csv = `${headers.join(',')}\n${values.join(',')}`
  downloadCsv(`order_${order.id}.csv`, csv)
}

const viewOrderDetail = (order) => {
  router.push(`/brand/order-detail/${order.id}`)
}

const exportFilteredOrders = () => {
  if (filteredOrders.value.length === 0) {
    alert('当前筛选条件下没有可导出订单')
    return
  }

  const csv = buildOrdersCsv(filteredOrders.value)
  if (!csv) {
    alert('导出失败，请重试')
    return
  }

  const suffix = `${dateRange.start || 'all'}_${dateRange.end || 'all'}`
  downloadCsv(`orders_${currentStatus.value}_${suffix}.csv`, csv)
}

// 监听筛选条件变化
watch([currentStatus, dateRange], () => {
  // 可以在这里添加防抖逻辑
}, { deep: true })

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.brand-orders {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 80px;
}

.top-nav {
  background: white;
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 10;
}

.nav-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
  color: #333;
}

.nav-stats {
  font-size: 14px;
  color: #666;
}

.total-orders {
  font-weight: 500;
}

.filters {
  background: white;
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.filter-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.filter-tab {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.filter-tab.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.date-filter {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-actions {
  margin-bottom: 12px;
}

.export-all-btn {
  width: 100%;
  border: 1px solid #2196f3;
  color: #2196f3;
  background: #f4faff;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: 500;
}

.date-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
}

.date-separator {
  color: #666;
  font-size: 14px;
}

.orders-list {
  padding: 16px;
}

.loading {
  text-align: center;
  padding: 40px 20px;
  color: #666;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  color: #666;
  margin: 0;
}

.order-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.order-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s;
}

.order-card:hover {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.order-info h3 {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 4px 0;
  color: #333;
}

.order-time {
  font-size: 12px;
  color: #999;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.pending {
  background: #fff3e0;
  color: #ff9800;
}

.status-badge.paid {
  background: #e8f5e8;
  color: #4caf50;
}

.status-badge.cancelled {
  background: #fce4ec;
  color: #e91e63;
}

.campaign-info {
  margin-bottom: 12px;
}

.campaign-name {
  font-size: 14px;
  font-weight: 600;
  margin: 0 0 4px 0;
  color: #333;
}

.user-info {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.order-details {
  margin-bottom: 12px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.detail-label {
  font-size: 12px;
  color: #666;
}

.detail-value {
  font-size: 12px;
  color: #333;
  font-weight: 500;
}

.detail-value.amount {
  color: #f39c12;
  font-weight: bold;
}

.detail-value.reward {
  color: #27ae60;
  font-weight: bold;
}

.form-data {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.form-data-title {
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-data-content {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.form-data-item {
  font-size: 11px;
  background: white;
  padding: 4px 8px;
  border-radius: 12px;
  color: #333;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn.confirm {
  border-color: #4caf50;
  color: #4caf50;
}

.action-btn.export {
  border-color: #2196f3;
  color: #2196f3;
}

.action-btn.detail {
  border-color: #9c27b0;
  color: #9c27b0;
}

.stats-section {
  padding: 16px;
  margin-top: 20px;
}

.stats-title {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 16px 0;
  color: #333;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 20px;
  font-weight: bold;
  color: #667eea;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #666;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  border-top: 1px solid #eee;
  padding: 8px 0;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #999;
  padding: 8px;
}

.nav-item.active {
  color: #667eea;
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 4px;
}

.nav-text {
  font-size: 12px;
}
</style>
