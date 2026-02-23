<template>
  <div class="brand-promoters">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">推广员管理</h1>
      <div class="nav-stats">
        <span class="total-promoters">总计: {{ totalPromoters }}</span>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-number">{{ promoterStats.active }}</div>
          <div class="stat-label">活跃推广员</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">¥{{ promoterStats.totalRewards }}</div>
          <div class="stat-label">总奖励发放</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ promoterStats.todayOrders }}</div>
          <div class="stat-label">今日订单</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ promoterStats.conversionRate }}%</div>
          <div class="stat-label">转化率</div>
        </div>
      </div>
    </div>

    <!-- 筛选器 -->
    <div class="filters">
      <div class="filter-row">
        <select v-model="currentFilter" class="filter-select">
          <option value="all">全部推广员</option>
          <option value="active">活跃推广员</option>
          <option value="top">优秀推广员</option>
          <option value="new">新注册</option>
        </select>
        
        <div class="search-box">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索推广员"
            class="search-input"
          >
          <button class="search-btn">🔍</button>
        </div>
      </div>
    </div>

    <!-- 推广员列表 -->
    <div class="promoters-list">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="filteredPromoters.length === 0" class="empty-state">
        <div class="empty-icon">👥</div>
        <p class="empty-text">暂无推广员数据</p>
      </div>

      <div v-else class="promoter-cards">
        <div
          v-for="promoter in filteredPromoters"
          :key="promoter.id"
          class="promoter-card"
          @click="viewPromoterDetail(promoter)"
        >
          <div class="card-header">
            <div class="promoter-avatar">
              <img :src="promoter.avatar" :alt="promoter.name" class="avatar-img">
            </div>
            <div class="promoter-info">
              <h3 class="promoter-name">{{ promoter.name }}</h3>
              <p class="promoter-phone">{{ promoter.phone }}</p>
              <span :class="['status-badge', promoter.status]">
                {{ getStatusText(promoter.status) }}
              </span>
            </div>
            <div class="promoter-level">
              <div class="level-badge">{{ promoter.level }}</div>
            </div>
          </div>

          <div class="promoter-stats">
            <div class="stat-item">
              <span class="stat-value">{{ promoter.totalOrders }}</span>
              <span class="stat-label">推广订单</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">¥{{ promoter.totalRewards }}</span>
              <span class="stat-label">累计奖励</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ promoter.conversionRate }}%</span>
              <span class="stat-label">转化率</span>
            </div>
          </div>

          <div class="recent-activity">
            <div class="activity-title">最近活动</div>
            <div class="activity-list">
              <div
                v-for="activity in promoter.recentActivities"
                :key="activity.id"
                class="activity-item"
              >
                <span class="activity-desc">{{ activity.description }}</span>
                <span class="activity-time">{{ formatTime(activity.time) }}</span>
              </div>
            </div>
          </div>

          <div class="card-actions">
            <button
              @click.stop="generateLink(promoter)"
              class="action-btn link"
            >
              生成链接
            </button>
            <button
              @click.stop="viewRewards(promoter)"
              class="action-btn rewards"
            >
              奖励记录
            </button>
            <button
              @click.stop="contactPromoter(promoter)"
              class="action-btn contact"
            >
              联系
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 推广链接生成模态框 -->
    <div v-if="showLinkModal" class="modal-overlay" @click="showLinkModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>生成推广链接</h3>
          <button @click="showLinkModal = false" class="close-btn">✕</button>
        </div>
        
        <div class="link-form">
          <div class="form-group">
            <label>选择活动</label>
            <select v-model="linkForm.campaignId" class="form-select">
              <option value="">请选择活动</option>
              <option v-for="campaign in campaigns" :key="campaign.id" :value="campaign.id">
                {{ campaign.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label>推广员</label>
            <input :value="linkForm.promoterName" readonly class="form-input">
          </div>

          <div v-if="generatedLink" class="generated-link">
            <label>推广链接</label>
            <div class="link-container">
              <input :value="generatedLink" readonly class="link-input">
              <button @click="copyLink" class="copy-btn">复制</button>
            </div>
            <div class="qr-code">
              <div class="qr-placeholder">二维码</div>
            </div>
          </div>
        </div>

        <div class="modal-actions">
          <button @click="showLinkModal = false" class="cancel-btn">取消</button>
          <button @click="generatePromoLink" :disabled="!linkForm.campaignId" class="confirm-btn">
            生成链接
          </button>
        </div>
      </div>
    </div>

    <!-- 联系推广员模态框 -->
    <div v-if="showContactModal" class="modal-overlay" @click="showContactModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>联系推广员</h3>
          <button @click="showContactModal = false" class="close-btn">✕</button>
        </div>
        
        <div class="contact-info">
          <p class="contact-name">{{ currentContactPromoter?.name }}</p>
          <p class="contact-phone">{{ currentContactPromoter?.phone }}</p>
        </div>

        <div class="contact-form">
          <div class="form-group">
            <label>联系方式</label>
            <div class="radio-group">
              <label class="radio-item">
                <input type="radio" v-model="contactForm.method" value="phone">
                <span>电话</span>
              </label>
              <label class="radio-item">
                <input type="radio" v-model="contactForm.method" value="sms">
                <span>短信</span>
              </label>
            </div>
          </div>

          <div v-if="contactForm.method === 'sms'" class="form-group">
            <label>短信内容</label>
            <textarea
              v-model="contactForm.message"
              class="form-textarea"
              placeholder="请输入短信内容"
              rows="3"
            ></textarea>
          </div>
        </div>

        <div class="modal-actions">
          <button @click="showContactModal = false" class="cancel-btn">取消</button>
          <button @click="executeContact" class="confirm-btn">
            {{ contactForm.method === 'phone' ? '拨打电话' : '发送短信' }}
          </button>
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
      <router-link to="/brand/promoters" class="nav-item active">
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
import { ref, computed, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  buildPromoterLink,
  buildPromoterLinkForm,
  calculatePromoterStats,
  filterAndSortPromoters,
  formatPromoterTime,
  getPromoterStatusText,
  getDefaultContactForm,
  validateContactForm,
  buildContactAction,
} from './promoters.logic.js'

const router = useRouter()
const promoters = ref([])
const campaigns = ref([])
const loading = ref(false)
const currentFilter = ref('all')
const searchKeyword = ref('')
const showLinkModal = ref(false)
const generatedLink = ref('')
const showContactModal = ref(false)
const currentContactPromoter = ref(null)
const contactForm = reactive(getDefaultContactForm())

const promoterStats = reactive({
  active: 0,
  totalRewards: 0,
  todayOrders: 0,
  conversionRate: 0
})

const linkForm = reactive({
  promoterId: null,
  promoterName: '',
  campaignId: ''
})

const totalPromoters = computed(() => promoters.value.length)

const filteredPromoters = computed(() => {
  return filterAndSortPromoters(promoters.value, currentFilter.value, searchKeyword.value)
})

const getStatusText = (status) => {
  return getPromoterStatusText(status)
}

const formatTime = (timeString) => {
  return formatPromoterTime(timeString)
}

const loadPromoters = async () => {
  loading.value = true
  try {
    const [promotersRes, campaignsRes] = await Promise.all([
      promoterApi.getPromoters(),
      campaignApi.getCampaigns()
    ])
    
    const promotersData = promotersRes.data || promotersRes
    promoters.value = Array.isArray(promotersData) ? promotersData : (promotersData.list || [])
    
    const campaignsData = campaignsRes.data || campaignsRes
    campaigns.value = Array.isArray(campaignsData) ? campaignsData : (campaignsData.list || [])

    calculateStats()
  } catch (error) {
    console.error('加载推广员失败:', error)
  } finally {
    loading.value = false
  }
}

const calculateStats = () => {
  const stats = calculatePromoterStats(promoters.value)
  promoterStats.active = stats.active
  promoterStats.totalRewards = stats.totalRewards
  promoterStats.todayOrders = stats.todayOrders
  promoterStats.conversionRate = stats.conversionRate
}

const viewPromoterDetail = (promoter) => {
  router.push(`/brand/promoter-detail/${promoter.id}`)
}

const generateLink = (promoter) => {
  const form = buildPromoterLinkForm(promoter)
  linkForm.promoterId = form.promoterId
  linkForm.promoterName = form.promoterName
  linkForm.campaignId = form.campaignId
  generatedLink.value = ''
  showLinkModal.value = true
}

const generatePromoLink = () => {
  generatedLink.value = buildPromoterLink(window.location.origin, linkForm.campaignId, linkForm.promoterId)
}

const copyLink = async () => {
  try {
    await navigator.clipboard.writeText(generatedLink.value)
    alert('链接已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    alert('复制失败，请手动复制')
  }
}

const viewRewards = (promoter) => {
  router.push(`/brand/reward-records/${promoter.id}`)
}

const contactPromoter = (promoter) => {
  currentContactPromoter.value = promoter
  Object.assign(contactForm, getDefaultContactForm())
  showContactModal.value = true
}

const executeContact = () => {
  const errorMsg = validateContactForm(contactForm, currentContactPromoter.value)
  if (errorMsg) {
    alert(errorMsg)
    return
  }

  const action = buildContactAction(contactForm, currentContactPromoter.value)
  if (action?.link) {
    window.location.href = action.link
    showContactModal.value = false
  }
}

onMounted(() => {
  loadPromoters()
})
</script>

<style scoped>
.brand-promoters {
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

.stats-section {
  padding: 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-card {
  background: white;
  padding: 16px;
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

.filters {
  background: white;
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.filter-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.filter-select {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
}

.search-box {
  flex: 1;
  display: flex;
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
}

.search-input {
  flex: 1;
  padding: 8px 12px;
  border: none;
  font-size: 14px;
}

.search-input:focus {
  outline: none;
}

.search-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 8px 12px;
  cursor: pointer;
}

.promoters-list {
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

.promoter-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.promoter-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s;
}

.promoter-card:hover {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.promoter-avatar {
  margin-right: 12px;
}

.avatar-img {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.promoter-info {
  flex: 1;
}

.promoter-name {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 4px 0;
  color: #333;
}

.promoter-phone {
  font-size: 12px;
  color: #666;
  margin: 0 0 8px 0;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.status-badge.active {
  background: #e8f5e8;
  color: #4caf50;
}

.status-badge.inactive {
  background: #fff3e0;
  color: #ff9800;
}

.status-badge.blocked {
  background: #fce4ec;
  color: #e91e63;
}

.promoter-level {
  margin-left: 12px;
}

.level-badge {
  background: linear-gradient(135deg, #ffd700 0%, #ffb347 100%);
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: bold;
}

.promoter-stats {
  display: flex;
  justify-content: space-around;
  margin-bottom: 16px;
  padding: 12px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 11px;
  color: #666;
}

.recent-activity {
  margin-bottom: 16px;
}

.activity-title {
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
  font-weight: 500;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.activity-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
}

.activity-desc {
  color: #333;
  flex: 1;
}

.activity-time {
  color: #999;
  margin-left: 8px;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  flex: 1;
  padding: 6px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn.link {
  border-color: #4caf50;
  color: #4caf50;
}

.action-btn.rewards {
  border-color: #ff9800;
  color: #ff9800;
}

.action-btn.contact {
  border-color: #2196f3;
  color: #2196f3;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 20px;
  width: 100%;
  max-width: 400px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  padding: 20px 20px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #999;
}

.link-form {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-input,
.form-select {
  width: 100%;
  padding: 12px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #667eea;
}

.generated-link {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.link-container {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.link-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 12px;
  background: #f8f9fa;
}

.copy-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 12px;
  cursor: pointer;
}

.qr-code {
  text-align: center;
}

.qr-placeholder {
  width: 120px;
  height: 120px;
  border: 2px dashed #ddd;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: #666;
  font-size: 14px;
}

.modal-actions {
  padding: 20px;
  display: flex;
  gap: 12px;
}

.cancel-btn,
.confirm-btn {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.confirm-btn {
  background: #667eea;
  color: white;
}

.confirm-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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

.contact-info {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #eee;
}

.contact-name {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin: 0 0 4px;
}

.contact-phone {
  font-size: 14px;
  color: #667eea;
  margin: 0;
}

.contact-form {
  padding: 20px;
}

.radio-group {
  display: flex;
  gap: 16px;
}

.radio-item {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.radio-item input {
  width: 18px;
  height: 18px;
}

.form-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  resize: none;
  font-family: inherit;
  box-sizing: border-box;
}

.form-textarea:focus {
  outline: none;
  border-color: #667eea;
}
</style>
