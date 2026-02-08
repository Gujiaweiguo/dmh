<template>
  <div class="feedback-center">
    <div class="header">
      <button class="back-btn" @click="goBack">返回</button>
      <h1>帮助与反馈</h1>
    </div>

    <div class="card">
      <h2>提交反馈</h2>
      <div class="form-grid">
        <label>
          分类
          <select v-model="form.category">
            <option value="poster">海报</option>
            <option value="payment">支付</option>
            <option value="verification">核销</option>
            <option value="other">其他</option>
          </select>
        </label>
        <label>
          优先级
          <select v-model="form.priority">
            <option value="medium">中</option>
            <option value="high">高</option>
            <option value="low">低</option>
          </select>
        </label>
      </div>

      <label class="block">
        标题
        <input v-model.trim="form.title" maxlength="100" placeholder="请简要描述问题" />
      </label>

      <label class="block">
        内容
        <textarea v-model.trim="form.content" rows="4" maxlength="1000" placeholder="请提供详细信息，便于快速定位问题"></textarea>
      </label>

      <label>
        评分（可选）
        <select v-model.number="form.rating">
          <option :value="null">不评分</option>
          <option :value="1">1 分</option>
          <option :value="2">2 分</option>
          <option :value="3">3 分</option>
          <option :value="4">4 分</option>
          <option :value="5">5 分</option>
        </select>
      </label>

      <button class="btn-primary" :disabled="submitting" @click="submitFeedback">
        {{ submitting ? '提交中...' : '提交反馈' }}
      </button>
    </div>

    <div class="card">
      <div class="card-header">
        <h2>我的反馈</h2>
        <button class="btn-link" @click="loadMyFeedback">刷新</button>
      </div>
      <div v-if="myFeedbackLoading" class="muted">加载中...</div>
      <div v-else-if="myFeedback.length === 0" class="muted">暂无反馈记录</div>
      <ul v-else class="list">
        <li v-for="item in myFeedback" :key="item.id">
          <div class="row">
            <strong>{{ item.title }}</strong>
            <span class="tag">{{ statusText(item.status) }}</span>
          </div>
          <div class="meta">{{ item.category }} · {{ formatDate(item.createdAt) }}</div>
          <p class="content">{{ item.content }}</p>
          <p v-if="item.response" class="response">回复：{{ item.response }}</p>
        </li>
      </ul>
    </div>

    <div class="card">
      <div class="card-header">
        <h2>常见问题</h2>
        <button class="btn-link" @click="loadFaq">刷新</button>
      </div>
      <div v-if="faqLoading" class="muted">加载中...</div>
      <div v-else-if="faqList.length === 0" class="muted">暂无 FAQ</div>
      <ul v-else class="list">
        <li v-for="faq in faqList" :key="faq.id">
          <strong>{{ faq.question }}</strong>
          <p class="content">{{ faq.answer }}</p>
          <div class="faq-actions">
            <button class="btn-link" @click="markHelpful(faq.id, 'helpful')">有帮助（{{ faq.helpfulCount || 0 }}）</button>
            <button class="btn-link" @click="markHelpful(faq.id, 'not_helpful')">没帮助（{{ faq.notHelpfulCount || 0 }}）</button>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { feedbackApi } from '@/services/api'

const router = useRouter()

const submitting = ref(false)
const myFeedbackLoading = ref(false)
const faqLoading = ref(false)

const myFeedback = ref([])
const faqList = ref([])

const form = ref({
  category: 'poster',
  priority: 'medium',
  title: '',
  content: '',
  rating: null,
})

const goBack = () => {
  router.back()
}

const statusText = (status) => {
  const map = {
    pending: '待处理',
    reviewing: '处理中',
    resolved: '已解决',
    closed: '已关闭',
  }
  return map[status] || status
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const submitFeedback = async () => {
  if (!form.value.title || !form.value.content) {
    showToast('请填写标题和内容')
    return
  }

  submitting.value = true
  try {
    await feedbackApi.createFeedback({
      category: form.value.category,
      priority: form.value.priority,
      title: form.value.title,
      content: form.value.content,
      rating: form.value.rating,
      featureUseCase: 'h5_feedback_center',
      deviceInfo: navigator.userAgent,
      browserInfo: navigator.userAgent,
    })
    showToast('反馈提交成功')
    form.value.title = ''
    form.value.content = ''
    form.value.rating = null
    await loadMyFeedback()
  } catch (error) {
    console.error('submit feedback failed:', error)
    showToast(error?.message || '提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

const loadMyFeedback = async () => {
  myFeedbackLoading.value = true
  try {
    const resp = await feedbackApi.listFeedback({ page: 1, pageSize: 20 })
    myFeedback.value = resp.feedbacks || []
  } catch (error) {
    console.error('load feedback failed:', error)
    myFeedback.value = []
  } finally {
    myFeedbackLoading.value = false
  }
}

const loadFaq = async () => {
  faqLoading.value = true
  try {
    const resp = await feedbackApi.listFaq({})
    faqList.value = resp.faqs || []
  } catch (error) {
    console.error('load faq failed:', error)
    faqList.value = []
  } finally {
    faqLoading.value = false
  }
}

const markHelpful = async (id, type) => {
  try {
    await feedbackApi.markFaqHelpful(id, type)
    await loadFaq()
  } catch (error) {
    console.error('mark faq helpful failed:', error)
    showToast('操作失败，请稍后重试')
  }
}

onMounted(async () => {
  await Promise.all([loadMyFeedback(), loadFaq()])
})
</script>

<style scoped>
.feedback-center { min-height: 100vh; background: #f7f8fa; padding: 12px; }
.header { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.header h1 { font-size: 20px; margin: 0; }
.back-btn { border: none; background: #fff; padding: 8px 12px; border-radius: 8px; }

.card { background: #fff; border-radius: 12px; padding: 14px; margin-bottom: 12px; }
.card h2 { margin: 0 0 10px; font-size: 16px; }
.card-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.block { display: block; margin-top: 10px; }
label { display: block; font-size: 13px; color: #374151; margin-bottom: 6px; }
input, select, textarea { width: 100%; box-sizing: border-box; margin-top: 6px; border: 1px solid #d1d5db; border-radius: 8px; padding: 8px 10px; font-size: 14px; }
textarea { resize: vertical; }
.btn-primary { width: 100%; border: 0; border-radius: 8px; background: #2563eb; color: #fff; padding: 10px; margin-top: 10px; }
.btn-primary:disabled { opacity: .6; }
.btn-link { border: 0; background: transparent; color: #2563eb; padding: 0; }

.list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 10px; }
.list li { border: 1px solid #eef2f7; border-radius: 10px; padding: 10px; }
.row { display: flex; justify-content: space-between; gap: 8px; }
.tag { font-size: 12px; color: #2563eb; background: #eff6ff; border-radius: 999px; padding: 2px 8px; }
.meta { font-size: 12px; color: #6b7280; margin-top: 4px; }
.content { margin: 8px 0 0; color: #111827; font-size: 14px; white-space: pre-wrap; }
.response { margin: 8px 0 0; color: #065f46; font-size: 13px; background: #ecfdf5; border-radius: 8px; padding: 6px 8px; }
.faq-actions { display: flex; gap: 12px; margin-top: 8px; }
.muted { color: #6b7280; font-size: 13px; }
</style>
