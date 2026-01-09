<template>
  <div class="container">
    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="campaign" class="content">
      <div class="banner">
        <h1 class="title">{{ campaign.name }}</h1>
        <p class="desc">{{ campaign.description }}</p>
      </div>

      <div class="section">
        <div class="section-title">活动时间</div>
        <div class="section-content">
          {{ formatTime(campaign.startTime) }} ~ {{ formatTime(campaign.endTime) }}
        </div>
      </div>

      <div class="section">
        <div class="section-title">报名奖励</div>
        <div class="reward">每成功报名奖励 {{ campaign.rewardRule?.toFixed(2) || 0 }} 元</div>
      </div>

      <div class="footer">
        <button class="btn-primary" @click="goForm">立即报名</button>
        <button class="btn-secondary" @click="share">分享给好友</button>
      </div>
    </div>

    <div v-else class="error">活动不存在或已下线</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const campaign = ref(null)
const source = ref({
  c_id: route.query.c_id || '',
  u_id: route.query.u_id || ''
})

// 存储来源信息到 localStorage
const saveSource = () => {
  try {
    localStorage.setItem('dmh_source', JSON.stringify(source.value))
  } catch (e) {
    console.error('保存来源信息失败', e)
  }
}

// 获取活动详情
const fetchCampaign = async () => {
  const campaignId = route.params.id
  if (!campaignId) {
    loading.value = false
    return
  }

  try {
    const response = await fetch(`/api/v1/h5/campaigns/${campaignId}`)
    if (response.ok) {
      campaign.value = await response.json()
    }
  } catch (error) {
    console.error('加载活动失败', error)
  } finally {
    loading.value = false
  }
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return time.replace('T', ' ').replace('Z', '').substring(0, 16)
}

// 跳转到报名表单
const goForm = () => {
  router.push(`/campaign/${route.params.id}/form`)
}

// 分享功能
const share = () => {
  alert('请使用浏览器的分享功能')
}

onMounted(() => {
  saveSource()
  fetchCampaign()
})
</script>

<style scoped>
.container {
  padding: 16px;
  padding-bottom: 80px;
}

.loading, .error {
  text-align: center;
  padding: 60px 20px;
  color: #999;
  font-size: 14px;
}

.banner {
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  border-radius: 16px;
  padding: 24px;
  color: #fff;
  margin-bottom: 24px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
}

.desc {
  font-size: 14px;
  opacity: 0.9;
  line-height: 1.6;
}

.section {
  background-color: #fff;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 8px;
  color: #666;
}

.section-content {
  font-size: 14px;
  color: #333;
}

.reward {
  font-size: 20px;
  color: #16a34a;
  font-weight: 600;
}

.footer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 12px 16px;
  display: flex;
  gap: 12px;
  background-color: #fff;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
}

.btn-primary, .btn-secondary {
  flex: 1;
  border-radius: 999px;
  padding: 14px 0;
  text-align: center;
  font-size: 16px;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background-color: #4f46e5;
  color: #fff;
}

.btn-primary:active {
  background-color: #4338ca;
}

.btn-secondary {
  background-color: #fff;
  color: #666;
  border: 1px solid #e5e7eb;
}

.btn-secondary:active {
  background-color: #f9fafb;
}
</style>
