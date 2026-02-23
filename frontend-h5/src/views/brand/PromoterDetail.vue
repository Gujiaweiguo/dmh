<template>
  <div class="promoter-detail-page">
    <div class="top-nav">
      <button class="back-btn" @click="goBack">返回</button>
      <h1 class="nav-title">推广员详情</h1>
      <span class="nav-placeholder"></span>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="errorMessage" class="error-box">
      <p>{{ errorMessage }}</p>
      <button class="retry-btn" @click="loadPromoterDetail">重试</button>
    </div>

    <div v-else class="content">
      <!-- 基本信息 -->
      <div class="detail-card">
        <div class="profile-header">
          <img :src="promoter.avatar || defaultAvatar" :alt="promoter.name" class="avatar">
          <div class="profile-info">
            <h2 class="name">{{ promoter.name }}</h2>
            <p class="phone">{{ promoter.phone }}</p>
            <span :class="['status-badge', promoter.status]">
              {{ getStatusText(promoter.status) }}
            </span>
          </div>
          <div class="level-badge">{{ promoter.level || '普通' }}</div>
        </div>
      </div>

      <!-- 业绩统计 -->
      <div class="detail-card">
        <h3 class="section-title">业绩统计</h3>
        <div class="stats-grid">
          <div class="stat-item">
            <span class="stat-value">{{ promoter.totalOrders || 0 }}</span>
            <span class="stat-label">推广订单</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">¥{{ promoter.totalRewards || 0 }}</span>
            <span class="stat-label">累计奖励</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ promoter.conversionRate || 0 }}%</span>
            <span class="stat-label">转化率</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ promoter.campaignCount || 0 }}</span>
            <span class="stat-label">参与活动</span>
          </div>
        </div>
      </div>

      <!-- 详细信息 -->
      <div class="detail-card">
        <h3 class="section-title">详细信息</h3>
        <div class="info-list">
          <div class="info-item">
            <span class="info-label">注册时间</span>
            <span class="info-value">{{ promoter.createdAt || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">最近活跃</span>
            <span class="info-value">{{ promoter.lastActiveAt || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">邀请人</span>
            <span class="info-value">{{ promoter.inviterName || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- 推广链接 -->
      <div class="detail-card" v-if="promoter.links && promoter.links.length">
        <h3 class="section-title">推广链接</h3>
        <div class="link-list">
          <div v-for="link in promoter.links" :key="link.id" class="link-item">
            <div class="link-info">
              <span class="link-campaign">{{ link.campaignName }}</span>
              <span class="link-clicks">{{ link.clicks }} 次点击</span>
            </div>
            <button class="copy-link-btn" @click="copyLink(link.url)">复制</button>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <button class="action-btn primary" @click="viewRewardRecords">
          查看奖励记录
        </button>
        <button class="action-btn secondary" @click="generateLink">
          生成推广链接
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { promoterApi } from '../../services/brandApi.js'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const errorMessage = ref('')
const promoter = ref({})
const defaultAvatar = 'https://api.dicebear.com/7.x/initials/svg?seed=User'

const getStatusText = (status) => {
  const map = {
    active: '活跃',
    inactive: '不活跃',
    blocked: '已禁用'
  }
  return map[status] || status || '-'
}

const goBack = () => {
  router.back()
}

const loadPromoterDetail = async () => {
  const promoterId = route.params.id
  if (!promoterId) {
    errorMessage.value = '推广员 ID 无效'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await promoterApi.getPromoter(promoterId)
    const payload = response?.data || response
    promoter.value = payload || {}
  } catch (error) {
    console.error('加载推广员详情失败:', error)
    errorMessage.value = '加载推广员详情失败，请重试'
  } finally {
    loading.value = false
  }
}

const viewRewardRecords = () => {
  router.push(`/brand/reward-records/${promoter.value.id}`)
}

const generateLink = () => {
  router.push(`/brand/promoters?generateLink=${promoter.value.id}`)
}

const copyLink = async (url) => {
  try {
    await navigator.clipboard.writeText(url)
    alert('链接已复制')
  } catch (error) {
    console.error('复制失败:', error)
    alert('复制失败')
  }
}

onMounted(() => {
  loadPromoterDetail()
})
</script>

<style scoped>
.promoter-detail-page { min-height: 100vh; background: #f5f7fa; }
.top-nav { background: #fff; padding: 16px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #eee; }
.back-btn { border: none; background: none; color: #667eea; font-size: 14px; }
.nav-title { margin: 0; font-size: 18px; }
.nav-placeholder { width: 40px; }
.loading, .error-box { padding: 40px 20px; text-align: center; color: #666; }
.retry-btn { margin-top: 12px; border: 1px solid #667eea; color: #667eea; background: #fff; border-radius: 8px; padding: 8px 16px; }
.content { padding: 16px; }
.detail-card { background: #fff; border-radius: 12px; padding: 16px; margin-bottom: 12px; }
.section-title { margin: 0 0 12px; font-size: 15px; color: #333; }
.profile-header { display: flex; align-items: center; }
.avatar { width: 56px; height: 56px; border-radius: 50%; margin-right: 12px; }
.profile-info { flex: 1; }
.name { margin: 0 0 4px; font-size: 18px; }
.phone { margin: 0 0 8px; font-size: 13px; color: #666; }
.status-badge { padding: 2px 8px; border-radius: 12px; font-size: 11px; }
.status-badge.active { background: #e8f5e8; color: #4caf50; }
.status-badge.inactive { background: #fff3e0; color: #ff9800; }
.status-badge.blocked { background: #fce4ec; color: #e91e63; }
.level-badge { background: linear-gradient(135deg, #ffd700 0%, #ffb347 100%); color: white; padding: 4px 12px; border-radius: 12px; font-size: 12px; }
.stats-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; }
.stat-item { text-align: center; padding: 8px; background: #f8f9fa; border-radius: 8px; }
.stat-value { display: block; font-size: 18px; font-weight: bold; color: #333; }
.stat-label { font-size: 12px; color: #666; }
.info-list { display: flex; flex-direction: column; gap: 8px; }
.info-item { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #f0f0f0; }
.info-item:last-child { border-bottom: none; }
.info-label { color: #666; font-size: 13px; }
.info-value { color: #333; font-size: 13px; }
.link-list { display: flex; flex-direction: column; gap: 8px; }
.link-item { display: flex; justify-content: space-between; align-items: center; padding: 12px; background: #f8f9fa; border-radius: 8px; }
.link-campaign { font-size: 13px; color: #333; }
.link-clicks { font-size: 11px; color: #999; margin-left: 8px; }
.copy-link-btn { border: 1px solid #667eea; color: #667eea; background: #fff; padding: 4px 12px; border-radius: 6px; font-size: 12px; }
.action-buttons { display: flex; gap: 12px; margin-top: 16px; }
.action-btn { flex: 1; padding: 12px; border-radius: 8px; font-size: 14px; border: none; }
.action-btn.primary { background: #667eea; color: white; }
.action-btn.secondary { background: #fff; color: #667eea; border: 1px solid #667eea; }
</style>
