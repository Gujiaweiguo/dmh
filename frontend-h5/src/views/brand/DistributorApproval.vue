<template>
  <div class="distributor-approval">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">分销商审批</h1>
      <div class="nav-tabs">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          :class="['tab-btn', { active: activeTab === tab.key }]"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
          <span v-if="tab.count > 0" class="count-badge">{{ tab.count }}</span>
        </button>
      </div>
    </div>

    <!-- 待审批列表 -->
    <div v-if="activeTab === 'pending'" class="approval-list">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="pendingList.length === 0" class="empty-state">
        <div class="empty-icon">📝</div>
        <p class="empty-text">暂无待审批申请</p>
      </div>

      <div v-else class="application-cards">
        <div
          v-for="app in pendingList"
          :key="app.id"
          class="application-card"
        >
          <div class="card-header">
            <div class="applicant-info">
              <div class="avatar">
                {{ app.username?.charAt(0) || '申' }}
              </div>
              <div class="info">
                <h3 class="name">{{ app.username }}</h3>
                <p class="time">{{ formatTime(app.createdAt) }}</p>
              </div>
            </div>
            <span class="badge pending">待审批</span>
          </div>

          <div class="card-body">
            <div class="info-row">
              <span class="label">申请品牌:</span>
              <span class="value">{{ app.brandName }}</span>
            </div>
            <div class="info-row">
              <span class="label">申请理由:</span>
              <span class="value reason">{{ app.reason || '无' }}</span>
            </div>
            <div v-if="app.phone" class="info-row">
              <span class="label">联系电话:</span>
              <span class="value">{{ app.phone }}</span>
            </div>
          </div>

          <div class="card-actions">
            <button
              @click="openApprovalModal(app)"
              class="action-btn approve"
            >
              审批
            </button>
            <button
              @click="viewDetail(app)"
              class="action-btn detail"
            >
              详情
            </button>
          </div>
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMore && !loading" class="load-more">
        <button @click="loadMore" class="load-more-btn">加载更多</button>
      </div>
    </div>

    <!-- 已处理列表 -->
    <div v-if="activeTab === 'processed'" class="processed-list">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="processedList.length === 0" class="empty-state">
        <div class="empty-icon">📋</div>
        <p class="empty-text">暂无已处理申请</p>
      </div>

      <div v-else class="application-cards">
        <div
          v-for="app in processedList"
          :key="app.id"
          class="application-card processed"
        >
          <div class="card-header">
            <div class="applicant-info">
              <div class="avatar" :class="app.status">
                {{ app.username?.charAt(0) || '申' }}
              </div>
              <div class="info">
                <h3 class="name">{{ app.username }}</h3>
                <p class="time">{{ formatTime(app.createdAt) }}</p>
              </div>
            </div>
            <span :class="['badge', app.status]">
              {{ app.status === 'approved' ? '已通过' : '已拒绝' }}
            </span>
          </div>

          <div class="card-body">
            <div class="info-row">
              <span class="label">品牌:</span>
              <span class="value">{{ app.brandName }}</span>
            </div>
            <div v-if="app.reviewer" class="info-row">
              <span class="label">审批人:</span>
              <span class="value">{{ app.reviewer }}</span>
            </div>
            <div v-if="app.reviewedAt" class="info-row">
              <span class="label">审批时间:</span>
              <span class="value">{{ formatTime(app.reviewedAt) }}</span>
            </div>
            <div v-if="app.reviewNotes" class="info-row">
              <span class="label">审批备注:</span>
              <span class="value">{{ app.reviewNotes }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 审批模态框 -->
    <div v-if="showApprovalModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>审批分销商申请</h3>
          <button @click="closeModal" class="close-btn">✕</button>
        </div>

        <div class="modal-body">
          <div class="applicant-detail">
            <h4>申请人信息</h4>
            <div class="detail-grid">
              <div class="detail-item">
                <span class="label">用户名:</span>
                <span class="value">{{ currentApp?.username }}</span>
              </div>
              <div class="detail-item">
                <span class="label">手机号:</span>
                <span class="value">{{ currentApp?.phone }}</span>
              </div>
              <div class="detail-item">
                <span class="label">申请品牌:</span>
                <span class="value">{{ currentApp?.brandName }}</span>
              </div>
              <div class="detail-item">
                <span class="label">申请时间:</span>
                <span class="value">{{ formatTime(currentApp?.createdAt) }}</span>
              </div>
            </div>
            <div class="reason-box">
              <span class="label">申请理由:</span>
              <p class="reason-text">{{ currentApp?.reason || '无' }}</p>
            </div>
          </div>

          <div class="approval-form">
            <h4>审批操作</h4>
            <div class="form-group">
              <label>设置级别</label>
              <select v-model="approvalForm.level" class="form-select">
                <option :value="1">一级分销商</option>
                <option :value="2">二级分销商</option>
                <option :value="3">三级分销商</option>
              </select>
              <p class="hint">级别决定了分销商的奖励比例</p>
            </div>

            <div class="form-group">
              <label>审批备注</label>
              <textarea
                v-model="approvalForm.notes"
                class="form-textarea"
                placeholder="请输入审批备注（可选）"
                rows="3"
              ></textarea>
            </div>
          </div>
        </div>

        <div class="modal-actions">
          <button
            @click="rejectApplication"
            class="action-btn reject"
          >
            拒绝
          </button>
          <button
            @click="approveApplication"
            class="action-btn approve"
          >
            通过
          </button>
        </div>
      </div>
    </div>

    <!-- 详情模态框 -->
    <div v-if="showDetailModal" class="modal-overlay" @click="showDetailModal = false">
      <div class="modal-content detail-modal" @click.stop>
        <div class="modal-header">
          <h3>申请详情</h3>
          <button @click="showDetailModal = false" class="close-btn">✕</button>
        </div>

        <div class="modal-body">
          <div class="detail-section">
            <h4>申请人信息</h4>
            <div class="detail-grid">
              <div class="detail-item full">
                <span class="label">用户名:</span>
                <span class="value">{{ currentApp?.username }}</span>
              </div>
              <div class="detail-item full">
                <span class="label">手机号:</span>
                <span class="value">{{ currentApp?.phone || '未填写' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">申请品牌:</span>
                <span class="value">{{ currentApp?.brandName }}</span>
              </div>
              <div class="detail-item">
                <span class="label">申请状态:</span>
                <span :class="['value', 'status', currentApp?.status]">
                  {{ getStatusText(currentApp?.status) }}
                </span>
              </div>
              <div class="detail-item">
                <span class="label">申请时间:</span>
                <span class="value">{{ formatTime(currentApp?.createdAt) }}</span>
              </div>
              <div v-if="currentApp?.reviewedAt" class="detail-item">
                <span class="label">审批时间:</span>
                <span class="value">{{ formatTime(currentApp?.reviewedAt) }}</span>
              </div>
            </div>
          </div>

          <div class="detail-section">
            <h4>申请理由</h4>
            <p class="reason-text">{{ currentApp?.reason || '无' }}</p>
          </div>

          <div v-if="currentApp?.reviewNotes" class="detail-section">
            <h4>审批备注</h4>
            <p class="reason-text">{{ currentApp.reviewNotes }}</p>
          </div>
        </div>

        <div class="modal-actions">
          <button
            v-if="currentApp?.status === 'pending'"
            @click="showDetailModal = false; openApprovalModal(currentApp)"
            class="action-btn approve"
          >
            去审批
          </button>
          <button
            @click="showDetailModal = false"
            class="action-btn close"
          >
            关闭
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
      <router-link to="/brand/distributor-approval" class="nav-item active">
        <div class="nav-icon">✅</div>
        <div class="nav-text">审批</div>
      </router-link>
      <router-link to="/brand/promoters" class="nav-item">
        <div class="nav-icon">👥</div>
        <div class="nav-text">推广员</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  formatTime,
  getStatusText,
  getDefaultApprovalForm,
  buildTabs,
  validateRejection,
  getMockPendingApplications,
  getAvatarText,
  buildPaginationParams,
  calculateHasMore,
  mergePaginatedResults
} from './distributorApproval.logic.js'

const router = useRouter()
const activeTab = ref('pending')
const loading = ref(false)
const showApprovalModal = ref(false)
const showDetailModal = ref(false)

const currentApp = ref(null)
const approvalForm = ref(getDefaultApprovalForm())

const pendingList = ref([])
const processedList = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const hasMore = ref(false)

const tabs = computed(() => buildTabs(pendingList.value.length, processedList.value.length))

const getCurrentBrandId = () => {
  const fromStorage = Number(localStorage.getItem('dmh_current_brand_id'))
  if (Number.isFinite(fromStorage) && fromStorage > 0) return fromStorage

  try {
    const info = JSON.parse(localStorage.getItem('dmh_user_info') || '{}')
    const firstBrandId = Array.isArray(info.brandIds) && info.brandIds.length > 0 ? Number(info.brandIds[0]) : 0
    if (Number.isFinite(firstBrandId) && firstBrandId > 0) {
      localStorage.setItem('dmh_current_brand_id', String(firstBrandId))
      return firstBrandId
    }
  } catch {
    // ignore
  }

  return 0
}



// 加载申请列表
const loadApplications = async (append = false) => {
  if (!append) {
    currentPage.value = 1
    pendingList.value = []
  }
  loading.value = true
  try {
    const token = localStorage.getItem('dmh_token')
    const currentBrandId = getCurrentBrandId()
    if (!currentBrandId) {
      alert('未选择品牌，请重新登录')
      router.push('/brand/login')
      return
    }

    const params = buildPaginationParams(currentPage.value, pageSize.value, 'pending')
    const response = await fetch(`/api/v1/brands/${currentBrandId}/distributor/applications?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      const newList = data.applications || []
      pendingList.value = mergePaginatedResults(pendingList.value, newList, append)
      hasMore.value = calculateHasMore(pendingList.value.length, pageSize.value, data.total)
    } else {
      if (!append) {
        pendingList.value = getMockPendingApplications()
        hasMore.value = false
      }
    }
  } catch (error) {
    console.error('加载申请失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载已处理列表
const loadProcessed = async (append = false) => {
  if (!append) {
    currentPage.value = 1
    processedList.value = []
  }
  loading.value = true
  try {
    const token = localStorage.getItem('dmh_token')
    const currentBrandId = getCurrentBrandId()
    if (!currentBrandId) {
      alert('未选择品牌，请重新登录')
      router.push('/brand/login')
      return
    }

    const params = buildPaginationParams(currentPage.value, pageSize.value, 'approved,rejected')
    const response = await fetch(`/api/v1/brands/${currentBrandId}/distributor/applications?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      const newList = data.applications || []
      processedList.value = mergePaginatedResults(processedList.value, newList, append)
      hasMore.value = calculateHasMore(processedList.value.length, pageSize.value, data.total)
    } else {
      if (!append) {
        processedList.value = []
        hasMore.value = false
      }
    }
  } catch (error) {
    console.error('加载已处理失败:', error)
  } finally {
    loading.value = false
  }
}

// 打开审批模态框
const openApprovalModal = (app) => {
  currentApp.value = app
  approvalForm.value = getDefaultApprovalForm()
  showApprovalModal.value = true
}

// 查看详情
const viewDetail = (app) => {
  currentApp.value = app
  showDetailModal.value = true
}

// 关闭模态框
const closeModal = () => {
  showApprovalModal.value = false
  showDetailModal.value = false
  currentApp.value = null
}

// 通过申请
const approveApplication = async () => {
  if (!currentApp.value) return

  try {
    const token = localStorage.getItem('dmh_token')
    const currentBrandId = getCurrentBrandId()
    if (!currentBrandId) {
      alert('未选择品牌，请重新登录')
      router.push('/brand/login')
      return
    }

    const response = await fetch(`/api/v1/brands/${currentBrandId}/distributor/approve/${currentApp.value.id}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        action: 'approve',
        level: approvalForm.value.level,
        reason: approvalForm.value.notes
      })
    })

    if (response.ok) {
      alert('审批通过！')
      closeModal()
      loadApplications()
    } else {
      const data = await response.json()
      alert(`审批失败: ${data.message || '未知错误'}`)
    }
  } catch (error) {
    console.error('审批失败:', error)
    alert('审批失败，请重试')
  }
}

// 拒绝申请
const rejectApplication = async () => {
  if (!currentApp.value) return

  if (!validateRejection(approvalForm.value)) {
    alert('请填写拒绝理由')
    return
  }

  try {
    const token = localStorage.getItem('dmh_token')
    const currentBrandId = getCurrentBrandId()
    if (!currentBrandId) {
      alert('未选择品牌，请重新登录')
      router.push('/brand/login')
      return
    }

    const response = await fetch(`/api/v1/brands/${currentBrandId}/distributor/approve/${currentApp.value.id}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        action: 'reject',
        reason: approvalForm.value.notes
      })
    })

    if (response.ok) {
      alert('已拒绝申请')
      closeModal()
      loadApplications()
    } else {
      const data = await response.json()
      alert(`操作失败: ${data.message || '未知错误'}`)
    }
  } catch (error) {
    console.error('操作失败:', error)
    alert('操作失败，请重试')
  }
}

// 加载更多
const loadMore = async () => {
  currentPage.value++
  if (activeTab.value === 'pending') {
    await loadApplications(true)
  } else {
    await loadProcessed(true)
  }
}

// 切换标签时重新加载
const handleTabChange = () => {
  if (activeTab.value === 'pending') {
    loadApplications()
  } else {
    loadProcessed()
  }
}

onMounted(() => {
  handleTabChange()
})
</script>

<style scoped>
.distributor-approval {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 80px;
}

.top-nav {
  background: white;
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.nav-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0 0 12px 0;
  color: #333;
}

.nav-tabs {
  display: flex;
  gap: 8px;
}

.tab-btn {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.tab-btn.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.count-badge {
  background: rgba(0, 0, 0, 0.2);
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 12px;
  min-width: 18px;
  text-align: center;
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

.application-cards {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.application-card {
  background: white;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.application-card.processed {
  opacity: 0.8;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.applicant-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #667eea;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
}

.avatar.approved {
  background: #4caf50;
}

.avatar.rejected {
  background: #f44336;
}

.info h3 {
  margin: 0;
  font-size: 15px;
  font-weight: bold;
  color: #333;
}

.info .time {
  margin: 2px 0 0 0;
  font-size: 12px;
  color: #999;
}

.badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.badge.pending {
  background: #fff3e0;
  color: #ff9800;
}

.badge.approved {
  background: #e8f5e8;
  color: #4caf50;
}

.badge.rejected {
  background: #ffebee;
  color: #f44336;
}

.card-body {
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  margin-bottom: 8px;
  font-size: 14px;
}

.info-row .label {
  color: #666;
  min-width: 70px;
}

.info-row .value {
  color: #333;
  flex: 1;
}

.info-row .value.reason {
  color: #667eea;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.action-btn.approve {
  background: #667eea;
  color: white;
}

.action-btn.detail {
  background: #f5f5f5;
  color: #666;
}

.load-more {
  padding: 16px;
  text-align: center;
}

.load-more-btn {
  padding: 10px 24px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 20px;
  color: #667eea;
  font-size: 14px;
  cursor: pointer;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 20px 20px 0 0;
  width: 100%;
  max-height: 80vh;
  overflow-y: auto;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.detail-modal {
  border-radius: 20px;
  max-height: 90vh;
  margin: 20px;
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #eee;
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
  font-size: 24px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-body {
  padding: 20px;
}

.applicant-detail {
  margin-bottom: 24px;
}

.applicant-detail h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #666;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item.full {
  grid-column: 1 / -1;
}

.detail-item .label {
  font-size: 12px;
  color: #999;
}

.detail-item .value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.detail-item .value.status {
  color: #667eea;
}

.reason-box {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 8px;
}

.reason-box .label {
  display: block;
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
}

.reason-text {
  margin: 0;
  font-size: 14px;
  color: #333;
  line-height: 1.5;
}

.approval-form h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #666;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-select,
.form-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: #667eea;
}

.form-textarea {
  resize: none;
  font-family: inherit;
}

.hint {
  margin: 4px 0 0 0;
  font-size: 12px;
  color: #999;
}

.modal-actions {
  padding: 16px 20px;
  border-top: 1px solid #eee;
  display: flex;
  gap: 12px;
}

.modal-actions .action-btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.modal-actions .action-btn.reject {
  background: #f5f5f5;
  color: #666;
}

.modal-actions .action-btn.approve {
  background: #667eea;
  color: white;
}

.modal-actions .action-btn.close {
  background: #ddd;
  color: #333;
}

.detail-section {
  margin-bottom: 20px;
}

.detail-section h4 {
  margin: 0 0 12px 0;
  font-size: 15px;
  color: #333;
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
  z-index: 100;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #999;
  padding: 4px;
}

.nav-item.active {
  color: #667eea;
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 2px;
}

.nav-text {
  font-size: 11px;
}
</style>
