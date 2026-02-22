<template>
  <div class="brand-analytics">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">数据分析</h1>
      <div class="date-selector">
        <select v-model="selectedPeriod" class="period-select">
          <option value="today">今日</option>
          <option value="week">本周</option>
          <option value="month">本月</option>
          <option value="quarter">本季度</option>
        </select>
      </div>
    </div>

    <!-- 核心指标 -->
    <div class="metrics-section">
      <h2 class="section-title">核心指标</h2>
      <div class="metrics-grid">
        <div class="metric-card">
          <div class="metric-icon">📊</div>
          <div class="metric-content">
            <div class="metric-value">{{ coreMetrics.totalRevenue }}</div>
            <div class="metric-label">总收入 (元)</div>
            <div class="metric-change positive">+12.5%</div>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon">🎯</div>
          <div class="metric-content">
            <div class="metric-value">{{ coreMetrics.totalOrders }}</div>
            <div class="metric-label">总订单数</div>
            <div class="metric-change positive">+8.3%</div>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon">👥</div>
          <div class="metric-content">
            <div class="metric-value">{{ coreMetrics.activePromoters }}</div>
            <div class="metric-label">活跃推广员</div>
            <div class="metric-change positive">+15.2%</div>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon">💰</div>
          <div class="metric-content">
            <div class="metric-value">{{ coreMetrics.avgOrderValue }}</div>
            <div class="metric-label">客单价 (元)</div>
            <div class="metric-change negative">-2.1%</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 趋势图表 -->
    <div class="charts-section">
      <h2 class="section-title">趋势分析</h2>
      
      <!-- 订单趋势 -->
      <div class="chart-card">
        <div class="chart-header">
          <h3 class="chart-title">订单趋势</h3>
          <div class="chart-legend">
            <span class="legend-item">
              <span class="legend-color orders"></span>
              订单数量
            </span>
            <span class="legend-item">
              <span class="legend-color revenue"></span>
              收入金额
            </span>
          </div>
        </div>
        <div class="chart-container">
          <div class="chart-placeholder">
            <div class="chart-bars">
              <div v-for="(data, index) in chartData.orders" :key="index" class="chart-bar">
                <div class="bar orders" :style="{ height: `${data.orders / 10}px` }"></div>
                <div class="bar revenue" :style="{ height: `${data.revenue / 100}px` }"></div>
                <div class="bar-label">{{ data.date }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 转化漏斗 -->
      <div class="chart-card">
        <div class="chart-header">
          <h3 class="chart-title">转化漏斗</h3>
        </div>
        <div class="funnel-container">
          <div
            v-for="(step, index) in funnelData"
            :key="index"
            class="funnel-step"
            :style="{ width: `${step.percentage}%` }"
          >
            <div class="funnel-content">
              <span class="funnel-label">{{ step.label }}</span>
              <span class="funnel-value">{{ step.value }}</span>
              <span class="funnel-rate">{{ step.percentage }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 活动排行 -->
    <div class="ranking-section">
      <h2 class="section-title">活动排行</h2>
      <div class="ranking-list">
        <div
          v-for="(campaign, index) in topCampaigns"
          :key="campaign.id"
          class="ranking-item"
        >
          <div class="ranking-number">{{ index + 1 }}</div>
          <div class="ranking-info">
            <h4 class="ranking-name">{{ campaign.name }}</h4>
            <p class="ranking-desc">{{ campaign.description }}</p>
          </div>
          <div class="ranking-stats">
            <div class="ranking-stat">
              <span class="stat-value">{{ campaign.orders }}</span>
              <span class="stat-label">订单</span>
            </div>
            <div class="ranking-stat">
              <span class="stat-value">¥{{ campaign.revenue }}</span>
              <span class="stat-label">收入</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 推广员排行 -->
    <div class="ranking-section">
      <h2 class="section-title">推广员排行</h2>
      <div class="ranking-list">
        <div
          v-for="(promoter, index) in topPromoters"
          :key="promoter.id"
          class="ranking-item"
        >
          <div class="ranking-number">{{ index + 1 }}</div>
          <div class="promoter-avatar">
            <img :src="promoter.avatar" :alt="promoter.name" class="avatar-img">
          </div>
          <div class="ranking-info">
            <h4 class="ranking-name">{{ promoter.name }}</h4>
            <p class="ranking-desc">{{ promoter.phone }}</p>
          </div>
          <div class="ranking-stats">
            <div class="ranking-stat">
              <span class="stat-value">{{ promoter.orders }}</span>
              <span class="stat-label">推广</span>
            </div>
            <div class="ranking-stat">
              <span class="stat-value">¥{{ promoter.rewards }}</span>
              <span class="stat-label">奖励</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 导出功能 -->
    <div class="export-section">
      <h2 class="section-title">数据导出</h2>
      <div class="export-options">
        <button @click="exportData('orders')" class="export-btn">
          📋 导出订单数据
        </button>
        <button @click="exportData('promoters')" class="export-btn">
          👥 导出推广员数据
        </button>
        <button @click="exportData('campaigns')" class="export-btn">
          🎯 导出活动数据
        </button>
        <button @click="exportData('all')" class="export-btn primary">
          📊 导出完整报表
        </button>
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
      <router-link to="/brand/orders" class="nav-item">
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
import { ref, reactive, onMounted, watch } from 'vue'
import {
  PERIOD_OPTIONS,
  getDefaultCoreMetrics,
  getDefaultChartData,
  getDefaultFunnelData,
  EXPORT_TYPES,
  getExportTypeLabel
} from './analytics.logic.js'
import { analyticsApi, campaignApi, promoterApi } from '@/services/brandApi.js'

const selectedPeriod = ref('month')

const coreMetrics = reactive(getDefaultCoreMetrics())

const chartData = reactive(getDefaultChartData())

const funnelData = ref(getDefaultFunnelData())

const topCampaigns = ref([])
const topPromoters = ref([])

const loadAnalyticsData = async () => {
  try {
    // 并行调用所有API
    const [metricsRes, trendsRes, campaignRes, promoterRes] = await Promise.all([
      analyticsApi.getMetrics(selectedPeriod.value),
      analyticsApi.getTrends(selectedPeriod.value),
      analyticsApi.getCampaignRanking(selectedPeriod.value),
      analyticsApi.getPromoterRanking(selectedPeriod.value)
    ])
    
    // 更新核心指标
    const metricsData = metricsRes.data || metricsRes
    Object.assign(coreMetrics, {
      totalRevenue: metricsData.totalRevenue || 0,
      totalOrders: metricsData.totalOrders || 0,
      activePromoters: metricsData.activePromoters || 0,
      avgOrderValue: metricsData.avgOrderValue || 0
    })

    // 更新图表数据
    const trendsData = trendsRes.data || trendsRes
    if (trendsData.orders && Array.isArray(trendsData.orders)) {
      chartData.orders = trendsData.orders
    }

    // 更新活动排行
    const campaignData = campaignRes.data || campaignRes
    topCampaigns.value = Array.isArray(campaignData) ? campaignData : (campaignData.list || [])

    // 更新推广员排行
    const promoterData = promoterRes.data || promoterRes
    topPromoters.value = Array.isArray(promoterData) ? promoterData : (promoterData.list || [])
  } catch (error) {
    console.error('加载分析数据失败:', error)
  }
}

const exportData = (type) => {
  alert(`导出${getExportTypeLabel(type)}功能开发中...`)
}

// 监听时间周期变化
watch(selectedPeriod, (newPeriod) => {
  console.log('切换到:', newPeriod)
  loadAnalyticsData()
})

onMounted(() => {
  loadAnalyticsData()
})
</script>
<style scoped>
.brand-analytics {
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

.date-selector {
  display: flex;
  align-items: center;
}

.period-select {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 16px;
  font-size: 14px;
  background: white;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 16px 0;
  color: #333;
}

.metrics-section {
  padding: 16px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.metric-card {
  background: white;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 12px;
}

.metric-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9ff;
  border-radius: 12px;
}

.metric-content {
  flex: 1;
}

.metric-value {
  font-size: 20px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.metric-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 4px;
}

.metric-change {
  font-size: 11px;
  font-weight: 500;
}

.metric-change.positive {
  color: #4caf50;
}

.metric-change.negative {
  color: #f44336;
}

.charts-section {
  padding: 16px;
}

.chart-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.chart-title {
  font-size: 16px;
  font-weight: bold;
  margin: 0;
  color: #333;
}

.chart-legend {
  display: flex;
  gap: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-color.orders {
  background: #667eea;
}

.legend-color.revenue {
  background: #f093fb;
}

.chart-container {
  height: 200px;
  display: flex;
  align-items: end;
  justify-content: center;
}

.chart-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: end;
  justify-content: center;
}

.chart-bars {
  display: flex;
  align-items: end;
  gap: 8px;
  height: 100%;
}

.chart-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.bar {
  width: 20px;
  border-radius: 4px 4px 0 0;
  min-height: 10px;
}

.bar.orders {
  background: #667eea;
}

.bar.revenue {
  background: #f093fb;
}

.bar-label {
  font-size: 10px;
  color: #666;
  margin-top: 8px;
}

.funnel-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.funnel-step {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px;
  border-radius: 8px;
  min-width: 200px;
  text-align: center;
  position: relative;
}

.funnel-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.funnel-label {
  font-size: 14px;
  font-weight: 500;
}

.funnel-value {
  font-size: 16px;
  font-weight: bold;
}

.funnel-rate {
  font-size: 12px;
  opacity: 0.8;
}

.ranking-section {
  padding: 16px;
}

.ranking-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ranking-item {
  background: white;
  border-radius: 16px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.ranking-number {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #ffd700 0%, #ffb347 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
}

.promoter-avatar {
  width: 40px;
  height: 40px;
}

.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.ranking-info {
  flex: 1;
}

.ranking-name {
  font-size: 14px;
  font-weight: bold;
  margin: 0 0 4px 0;
  color: #333;
}

.ranking-desc {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.ranking-stats {
  display: flex;
  gap: 16px;
}

.ranking-stat {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 14px;
  font-weight: bold;
  color: #333;
  margin-bottom: 2px;
}

.stat-label {
  font-size: 10px;
  color: #666;
}

.export-section {
  padding: 16px;
}

.export-options {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.export-btn {
  background: white;
  border: 2px solid #e1e5e9;
  border-radius: 12px;
  padding: 16px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
}

.export-btn:hover {
  border-color: #667eea;
  background: #f8f9ff;
}

.export-btn.primary {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.export-btn.primary:hover {
  background: #5a6fd8;
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
