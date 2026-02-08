<template>
  <div class="container">
    <div class="header">
      <div>
        <h1 class="title">营销活动</h1>
        <p class="subtitle">选择感兴趣的活动立即报名</p>
      </div>
      <button class="my-orders-btn" @click="goMyOrders">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/></svg>
        <span>我的报名</span>
      </button>
      <button class="distributor-btn" @click="goDistributor">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-3 3-3-3a2 2 0 0 1 2.83-2.83L12 7Z"/><path d="M2 12h20"/></svg>
        <span>分销中心</span>
      </button>
      <button class="feedback-btn" @click="goFeedback">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
        <span>帮助反馈</span>
      </button>
    </div>

    <!-- 筛选标签 -->
    <div class="filter-tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.key"
        :class="['tab-item', { active: activeTab === tab.key }]"
        @click="switchTab(tab.key)"
      >
        <span class="tab-label">{{ tab.label }}</span>
        <span v-if="tab.count !== undefined" class="tab-count">{{ tab.count }}</span>
      </button>

      <!-- 只看未报名开关，仅在“进行中”时展示 -->
      <div v-if="activeTab === 'ongoing'" class="only-unreg-toggle">
        <label class="toggle-label">
          <input type="checkbox" v-model="onlyUnregistered" />
          <span>只看未报名</span>
        </label>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="filteredCampaigns.length > 0" class="campaign-list">
      <div 
        v-for="campaign in filteredCampaigns" 
        :key="campaign.id"
        class="campaign-card"
        @click="goDetail(campaign.id)"
      >
        <div class="card-header">
          <h3 class="campaign-name">{{ campaign.name }}</h3>
          <div class="badges">
            <span v-if="campaign.isRegistered" class="badge registered">已报名</span>
            <span class="status" :class="campaign.status">
              {{ statusText(campaign.status) }}
            </span>
          </div>
        </div>
        <p class="campaign-desc">{{ campaign.description }}</p>
        <div class="campaign-info">
          <div class="info-item">
            <span class="label">活动时间</span>
            <span class="value">{{ formatDate(campaign.startTime) }} ~ {{ formatDate(campaign.endTime) }}</span>
          </div>
          <div class="info-item reward">
            <span class="label">报名奖励</span>
            <span class="value">{{ campaign.rewardRule?.toFixed(2) || 0 }} 元</span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty">
      <div class="empty-icon">📝</div>
      <p>暂无{{ emptyText }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const campaigns = ref([])
const myOrders = ref([])
const loading = ref(true)
const activeTab = ref('all')
const myPhone = ref('')
const onlyUnregistered = ref(false)

// 获取本地存储的手机号
const loadMyPhone = () => {
  try {
    const saved = localStorage.getItem('dmh_my_phone')
    if (saved) {
      myPhone.value = saved
    }
  } catch (e) {
    console.error('读取手机号失败', e)
  }
}

// 筛选标签
const tabs = computed(() => {
  const ongoing = campaigns.value.filter(c => c.status === 'active')
  const ended = campaigns.value.filter(c => c.status === 'ended' || c.status === 'paused')
  
  return [
    { key: 'all', label: '全部', count: campaigns.value.length },
    { key: 'ongoing', label: '进行中', count: ongoing.length },
    { key: 'ended', label: '已结束', count: ended.length }
  ]
})

// 筛选后的活动列表
const filteredCampaigns = computed(() => {
  if (activeTab.value === 'all') {
    return campaigns.value
  } else if (activeTab.value === 'ongoing') {
    return campaigns.value
      .filter(c => c.status === 'active')
      .filter(c => !onlyUnregistered.value || !c.isRegistered)
  } else if (activeTab.value === 'ended') {
    return campaigns.value.filter(c => c.status === 'ended' || c.status === 'paused')
  }
  return campaigns.value
})

// 空状态提示文字
const emptyText = computed(() => {
  const map = {
    all: '活动',
    ongoing: '进行中的活动',
    ended: '已结束的活动'
  }
  return map[activeTab.value] || '活动'
})

// 切换标签
const switchTab = (key) => {
  activeTab.value = key
}

// 存储来源信息
const saveSource = () => {
  const source = {
    c_id: route.query.c_id || '',
    u_id: route.query.u_id || ''
  }
  try {
    localStorage.setItem('dmh_source', JSON.stringify(source))
  } catch (e) {
    console.error('保存来源信息失败', e)
  }
}

// 获取我的订单
const fetchMyOrders = async () => {
  if (!myPhone.value) return
  
  try {
    const response = await fetch(`/api/v1/orders?pageSize=100`)
    if (response.ok) {
      const data = await response.json()
      // 过滤出当前手机号的订单
      myOrders.value = (data.orders || []).filter(order => order.phone === myPhone.value)
    }
  } catch (error) {
    console.error('加载订单失败', error)
  }
}

// 获取活动列表
const fetchCampaigns = async () => {
  try {
    const response = await fetch('/api/v1/campaigns')
    if (response.ok) {
      const data = await response.json()
      campaigns.value = (data.campaigns || []).map(campaign => {
        // 标记是否已报名
        const isRegistered = myOrders.value.some(order => order.campaignId === campaign.id)
        return {
          ...campaign,
          isRegistered
        }
      })
    } else {
      throw new Error(`HTTP ${response.status}`)
    }
  } catch (error) {
    console.error('加载活动列表失败', error)
    // API失败时使用示例数据
    campaigns.value = [
      {
        id: 1,
        name: '春节特惠活动',
        description: '新春佳节，推荐好友享双重奖励，完成任务即可获得现金奖励',
        status: 'active',
        rewardRule: 88,
        startTime: '2026-02-01 00:00:00',
        endTime: '2026-02-15 23:59:59',
        isRegistered: false
      },
      {
        id: 2,
        name: '会员招募计划',
        description: '招募品牌会员，享受专属优惠和会员特权',
        status: 'active',
        rewardRule: 66,
        startTime: '2026-01-01 00:00:00',
        endTime: '2026-12-31 23:59:59',
        isRegistered: false
      },
      {
        id: 3,
        name: '元宵节活动',
        description: '元宵佳节，猜灯谜赢大奖，传统节日新玩法',
        status: 'active',
        rewardRule: 50,
        startTime: '2026-02-28 00:00:00',
        endTime: '2026-03-01 23:59:59',
        isRegistered: false
      },
      {
        id: 4,
        name: '双11狂欢预热',
        description: '双11预热活动，提前报名享受专属优惠',
        status: 'ended',
        rewardRule: 20,
        startTime: '2024-11-01 00:00:00',
        endTime: '2024-11-11 23:59:59',
        isRegistered: true
      }
    ]
  } finally {
    loading.value = false
  }
}

// 初始化加载
const init = async () => {
  loadMyPhone()
  await fetchMyOrders()
  await fetchCampaigns()
}

// 格式化日期
const formatDate = (time) => {
  if (!time) return ''
  return time.substring(0, 10)
}

// 状态文本
const statusText = (status) => {
  const map = {
    active: '进行中',
    paused: '已暂停',
    ended: '已结束'
  }
  return map[status] || status
}

// 跳转到活动详情
const goDetail = (id) => {
  router.push(`/campaign/${id}`)
}

// 跳转到我的报名
const goMyOrders = () => {
  router.push('/orders')
}

const goDistributor = () => {
  router.push('/distributor')
}

const goFeedback = () => {
  router.push('/feedback')
}

onMounted(() => {
  saveSource()
  init()
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.header {
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  padding: 32px 16px 24px;
  color: #fff;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.my-orders-btn {
  background-color: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
}

.my-orders-btn:active {
  background-color: rgba(255, 255, 255, 0.3);
  transform: scale(0.95);
}

.distributor-btn {
  background-color: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
}

.distributor-btn:active {
  background-color: rgba(255, 255, 255, 0.3);
  transform: scale(0.95);
}

.feedback-btn {
  background-color: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
}

.feedback-btn:active {
  background-color: rgba(255, 255, 255, 0.3);
  transform: scale(0.95);
}

.title {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  opacity: 0.9;
}

.filter-tabs {
  display: flex;
  align-items: center;
  background-color: #fff;
  padding: 8px;
  gap: 8px;
  overflow-x: auto;
  border-bottom: 1px solid #f3f4f6;
}

.tab-item {
  flex-shrink: 0;
  padding: 8px 16px;
  border: none;
  background-color: #f9fafb;
  border-radius: 8px;
  font-size: 14px;
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
}

.only-unreg-toggle {
  margin-left: auto;
  padding-left: 8px;
  border-left: 1px solid #e5e7eb;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #6b7280;
}

.toggle-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
}

.tab-item.active {
  background-color: #4f46e5;
  color: #fff;
  font-weight: 500;
}

.tab-label {
  white-space: nowrap;
}

.tab-count {
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 10px;
  background-color: rgba(0, 0, 0, 0.1);
  min-width: 20px;
  text-align: center;
}

.tab-item.active .tab-count {
  background-color: rgba(255, 255, 255, 0.2);
}

.loading, .empty {
  text-align: center;
  padding: 60px 20px;
  color: #999;
  font-size: 14px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.campaign-list {
  padding: 16px;
}

.campaign-card {
  background-color: #fff;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.campaign-card:active {
  transform: scale(0.98);
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.badges {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-end;
}

.badge {
  font-size: 11px;
  padding: 3px 8px;
  border-radius: 4px;
  white-space: nowrap;
  font-weight: 500;
}

.badge.registered {
  background-color: #dbeafe;
  color: #1e40af;
}

.campaign-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  flex: 1;
  margin-right: 12px;
}

.status {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  white-space: nowrap;
}

.status.active {
  background-color: #dcfce7;
  color: #16a34a;
}

.status.paused {
  background-color: #fef3c7;
  color: #d97706;
}

.status.ended {
  background-color: #f3f4f6;
  color: #6b7280;
}

.campaign-desc {
  font-size: 14px;
  color: #666;
  margin-bottom: 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.campaign-info {
  padding-top: 12px;
  border-top: 1px solid #f3f4f6;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  margin-bottom: 8px;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-item .label {
  color: #999;
}

.info-item .value {
  color: #333;
  font-weight: 500;
}

.info-item.reward .value {
  color: #16a34a;
  font-size: 16px;
  font-weight: 600;
}
</style>
