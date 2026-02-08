<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RefreshCw } from 'lucide-vue-next';
import { feedbackApi, type FeedbackItem, type FeedbackStatsResponse } from '../services/feedbackApi';

const loading = ref(false);
const savingId = ref<number | null>(null);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const list = ref<FeedbackItem[]>([]);
const selectedId = ref<number | null>(null);

const filters = ref({
  category: '',
  status: '',
  priority: '',
});

const stats = ref<FeedbackStatsResponse>({
  totalFeedbacks: 0,
  byCategory: {},
  byStatus: {},
  byPriority: {},
  averageRating: 0,
  resolutionRate: 0,
  avgResolutionTime: 0,
  byRating: {},
});

const selected = computed(() => list.value.find((item) => item.id === selectedId.value) || null);

const formatDateTime = (value: string | null) => {
  if (!value) return '-';
  return new Date(value).toLocaleString('zh-CN');
};

const statusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    reviewing: '处理中',
    resolved: '已解决',
    closed: '已关闭',
  };
  return map[status] || status;
};

const priorityText = (priority: string) => {
  const map: Record<string, string> = {
    low: '低',
    medium: '中',
    high: '高',
  };
  return map[priority] || priority;
};

const fetchData = async () => {
  loading.value = true;
  try {
    const [listResp, statsResp] = await Promise.all([
      feedbackApi.getFeedbackList({
        page: page.value,
        pageSize: pageSize.value,
        category: filters.value.category || undefined,
        status: filters.value.status || undefined,
        priority: filters.value.priority || undefined,
      }),
      feedbackApi.getFeedbackStatistics({
        category: filters.value.category || undefined,
      }),
    ]);

    list.value = listResp.feedbacks || [];
    total.value = listResp.total || 0;
    stats.value = statsResp;
    if (selectedId.value == null && list.value.length > 0) {
      selectedId.value = list.value[0].id;
    }
  } catch (error) {
    console.error('加载反馈数据失败:', error);
    alert('加载反馈数据失败');
  } finally {
    loading.value = false;
  }
};

const updateStatus = async (item: FeedbackItem, nextStatus: string) => {
  savingId.value = item.id;
  try {
    await feedbackApi.updateFeedbackStatus({
      id: item.id,
      status: nextStatus,
    });
    await fetchData();
  } catch (error) {
    console.error('更新状态失败:', error);
    alert('更新状态失败');
  } finally {
    savingId.value = null;
  }
};

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="feedback-view">
    <div class="header">
      <h2 class="title">反馈管理</h2>
      <button class="btn-refresh" @click="fetchData" :disabled="loading">
        <RefreshCw :size="16" :class="{ spinning: loading }" />
        刷新
      </button>
    </div>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-label">反馈总数</div>
        <div class="stat-value">{{ stats.totalFeedbacks }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">平均评分</div>
        <div class="stat-value">{{ stats.averageRating.toFixed(2) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">解决率</div>
        <div class="stat-value">{{ (stats.resolutionRate * 100).toFixed(1) }}%</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">平均解决时长</div>
        <div class="stat-value">{{ stats.avgResolutionTime.toFixed(2) }}h</div>
      </div>
    </div>

    <div class="filters">
      <select v-model="filters.category" @change="fetchData">
        <option value="">全部分类</option>
        <option value="poster">海报</option>
        <option value="payment">支付</option>
        <option value="verification">核销</option>
        <option value="other">其他</option>
      </select>
      <select v-model="filters.status" @change="fetchData">
        <option value="">全部状态</option>
        <option value="pending">待处理</option>
        <option value="reviewing">处理中</option>
        <option value="resolved">已解决</option>
        <option value="closed">已关闭</option>
      </select>
      <select v-model="filters.priority" @change="fetchData">
        <option value="">全部优先级</option>
        <option value="high">高</option>
        <option value="medium">中</option>
        <option value="low">低</option>
      </select>
    </div>

    <div class="content">
      <div class="list-panel">
        <div class="panel-title">反馈列表（{{ total }}）</div>
        <div v-if="loading" class="loading">加载中...</div>
        <table v-else class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>标题</th>
              <th>分类</th>
              <th>优先级</th>
              <th>状态</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="item in list"
              :key="item.id"
              :class="{ active: selectedId === item.id }"
              @click="selectedId = item.id"
            >
              <td>{{ item.id }}</td>
              <td class="title-cell">{{ item.title }}</td>
              <td>{{ item.category }}</td>
              <td>{{ priorityText(item.priority) }}</td>
              <td>{{ statusText(item.status) }}</td>
              <td>{{ formatDateTime(item.createdAt) }}</td>
              <td>
                <div class="actions">
                  <button
                    class="btn-mini"
                    :disabled="savingId === item.id || item.status === 'reviewing'"
                    @click.stop="updateStatus(item, 'reviewing')"
                  >
                    处理
                  </button>
                  <button
                    class="btn-mini success"
                    :disabled="savingId === item.id || item.status === 'resolved'"
                    @click.stop="updateStatus(item, 'resolved')"
                  >
                    解决
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="detail-panel">
        <div class="panel-title">反馈详情</div>
        <div v-if="!selected" class="empty">请选择一条反馈</div>
        <div v-else class="detail">
          <div class="detail-row"><span>标题</span><strong>{{ selected.title }}</strong></div>
          <div class="detail-row"><span>内容</span><p>{{ selected.content }}</p></div>
          <div class="detail-row"><span>用户</span><strong>{{ selected.userName || selected.userId }}</strong></div>
          <div class="detail-row"><span>分类</span><strong>{{ selected.category }}</strong></div>
          <div class="detail-row"><span>状态</span><strong>{{ statusText(selected.status) }}</strong></div>
          <div class="detail-row"><span>优先级</span><strong>{{ priorityText(selected.priority) }}</strong></div>
          <div class="detail-row"><span>评分</span><strong>{{ selected.rating ?? '-' }}</strong></div>
          <div class="detail-row"><span>设备</span><strong>{{ selected.deviceInfo || '-' }}</strong></div>
          <div class="detail-row"><span>浏览器</span><strong>{{ selected.browserInfo || '-' }}</strong></div>
          <div class="detail-row"><span>处理回复</span><p>{{ selected.response || '-' }}</p></div>
          <div class="detail-row"><span>创建时间</span><strong>{{ formatDateTime(selected.createdAt) }}</strong></div>
          <div class="detail-row"><span>解决时间</span><strong>{{ formatDateTime(selected.resolvedAt) }}</strong></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.feedback-view { padding: 20px; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.title { font-size: 24px; font-weight: 700; color: #111827; }
.btn-refresh { display: inline-flex; align-items: center; gap: 6px; border: 0; background: #2563eb; color: #fff; border-radius: 8px; padding: 8px 12px; cursor: pointer; }
.btn-refresh:disabled { opacity: .6; cursor: not-allowed; }
.spinning { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.stats-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 16px; }
.stat-card { background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; padding: 12px; }
.stat-label { color: #6b7280; font-size: 12px; }
.stat-value { color: #111827; font-size: 22px; font-weight: 700; margin-top: 4px; }

.filters { display: flex; gap: 8px; margin-bottom: 16px; }
.filters select { height: 34px; border: 1px solid #d1d5db; border-radius: 8px; padding: 0 10px; background: #fff; }

.content { display: grid; grid-template-columns: 1.4fr 1fr; gap: 12px; }
.list-panel, .detail-panel { background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; overflow: hidden; }
.panel-title { padding: 12px 14px; font-weight: 700; border-bottom: 1px solid #f3f4f6; }
.loading, .empty { padding: 18px; color: #6b7280; }

.table { width: 100%; border-collapse: collapse; }
.table th, .table td { padding: 10px; border-bottom: 1px solid #f3f4f6; font-size: 13px; }
.table th { text-align: left; color: #6b7280; background: #f9fafb; }
.table tr.active { background: #eef2ff; }
.title-cell { max-width: 220px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.actions { display: flex; gap: 6px; }
.btn-mini { border: 0; background: #e5e7eb; color: #111827; border-radius: 6px; padding: 4px 8px; cursor: pointer; }
.btn-mini.success { background: #dcfce7; color: #166534; }
.btn-mini:disabled { opacity: .6; cursor: not-allowed; }

.detail { padding: 12px 14px; display: flex; flex-direction: column; gap: 10px; }
.detail-row { display: grid; grid-template-columns: 90px 1fr; align-items: start; gap: 8px; font-size: 13px; }
.detail-row span { color: #6b7280; }
.detail-row p { margin: 0; color: #111827; white-space: pre-wrap; }

@media (max-width: 1200px) {
  .stats-row { grid-template-columns: repeat(2, 1fr); }
  .content { grid-template-columns: 1fr; }
}
</style>
