<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Search, Plus, Edit2, Trash2, Eye, Image } from 'lucide-vue-next';
import { campaignApi } from '../services/campaignApi';
import { posterApi } from '../services/posterApi';
import type { Campaign } from '../types';

const campaigns = ref<Campaign[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);
const keyword = ref('');
const statusFilter = ref('');
const loading = ref(false);

const fetchCampaigns = async () => {
  loading.value = true;
  try {
    const response = await campaignApi.getCampaigns(
      page.value,
      pageSize.value,
      statusFilter.value,
      keyword.value
    );
    campaigns.value = response.campaigns;
    total.value = response.total;
  } catch (error) {
    console.error('Failed to fetch campaigns:', error);
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  page.value = 1;
  fetchCampaigns();
};

const handleDelete = async (id: number) => {
  if (!confirm('确定要删除这个活动吗？')) return;
  try {
    await campaignApi.deleteCampaign(id);
    fetchCampaigns();
  } catch (error) {
    console.error('Failed to delete campaign:', error);
    alert('删除失败');
  }
};

const handleGeneratePoster = async (campaign: Campaign) => {
  try {
    const response = await posterApi.generateCampaignPoster(
      campaign.id,
      campaign.posterTemplateId
    );
    if (response?.posterUrl) {
      const shouldOpen = confirm('海报生成成功，是否打开预览？');
      if (shouldOpen) {
        window.open(response.posterUrl, '_blank');
      }
    } else {
      alert('海报生成成功，但未返回链接');
    }
  } catch (error) {
    console.error('Failed to generate poster:', error);
    alert('海报生成失败');
  }
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN');
};

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '进行中',
    paused: '已暂停',
    ended: '已结束'
  };
  return statusMap[status] || status;
};

const getStatusClass = (status: string) => {
  const classMap: Record<string, string> = {
    active: 'bg-green-100 text-green-800',
    paused: 'bg-yellow-100 text-yellow-800',
    ended: 'bg-gray-100 text-gray-800'
  };
  return classMap[status] || '';
};

onMounted(() => {
  fetchCampaigns();
});

const emit = defineEmits(['create', 'edit', 'view']);
</script>

<template>
  <div class="campaign-list-view">
    <div class="header">
      <h1 class="title">营销活动管理</h1>
      <button class="btn-primary" @click="emit('create')">
        <Plus :size="20" />
        <span>创建活动</span>
      </button>
    </div>

    <div class="filters">
      <div class="search-box">
        <Search :size="20" class="search-icon" />
        <input
          v-model="keyword"
          type="text"
          placeholder="搜索活动名称或描述"
          @keyup.enter="handleSearch"
        />
      </div>
      <select v-model="statusFilter" @change="handleSearch" class="status-filter">
        <option value="">全部状态</option>
        <option value="active">进行中</option>
        <option value="paused">已暂停</option>
        <option value="ended">已结束</option>
      </select>
      <button class="btn-secondary" @click="handleSearch">搜索</button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="table-container">
      <table class="campaign-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>活动名称</th>
            <th>描述</th>
            <th>奖励规则</th>
            <th>开始时间</th>
            <th>结束时间</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="campaign in campaigns" :key="campaign.id">
            <td>{{ campaign.id }}</td>
            <td class="campaign-name">{{ campaign.name }}</td>
            <td class="description">{{ campaign.description }}</td>
            <td class="reward">¥{{ campaign.rewardRule.toFixed(2) }}</td>
            <td>{{ formatDate(campaign.startTime) }}</td>
            <td>{{ formatDate(campaign.endTime) }}</td>
            <td>
              <span class="status-badge" :class="getStatusClass(campaign.status)">
                {{ getStatusText(campaign.status) }}
              </span>
            </td>
            <td class="actions">
              <button class="btn-icon" @click="emit('view', campaign.id)" title="查看">
                <Eye :size="18" />
              </button>
              <button class="btn-icon" @click="handleGeneratePoster(campaign)" title="生成海报">
                <Image :size="18" />
              </button>
              <button class="btn-icon" @click="emit('edit', campaign.id)" title="编辑">
                <Edit2 :size="18" />
              </button>
              <button class="btn-icon danger" @click="handleDelete(campaign.id)" title="删除">
                <Trash2 :size="18" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination">
      <div class="pagination-info">
        共 {{ total }} 条记录，第 {{ page }} / {{ Math.ceil(total / pageSize) }} 页
      </div>
      <div class="pagination-controls">
        <button :disabled="page <= 1" @click="page--; fetchCampaigns()">上一页</button>
        <button :disabled="page >= Math.ceil(total / pageSize)" @click="page++; fetchCampaigns()">下一页</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.campaign-list-view {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-primary:hover {
  background: #2563eb;
}

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.search-box {
  flex: 1;
  position: relative;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
}

.search-box input {
  width: 100%;
  padding: 10px 12px 10px 42px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
}

.status-filter {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  background: white;
}

.btn-secondary {
  padding: 10px 20px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-secondary:hover {
  background: #f9fafb;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}

.table-container {
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  overflow-x: auto;
}

.campaign-table {
  width: 100%;
  border-collapse: collapse;
}

.campaign-table th {
  background: #f9fafb;
  padding: 12px;
  text-align: left;
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  border-bottom: 1px solid #e5e7eb;
}

.campaign-table td {
  padding: 12px;
  border-bottom: 1px solid #f3f4f6;
  font-size: 14px;
  color: #374151;
}

.campaign-name {
  font-weight: 500;
  color: #1a1a1a;
}

.description {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.reward {
  font-weight: 600;
  color: #f59e0b;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
  padding: 6px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: #6b7280;
  border-radius: 4px;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: #f3f4f6;
  color: #3b82f6;
}

.btn-icon.danger:hover {
  color: #ef4444;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
}

.pagination-info {
  color: #6b7280;
  font-size: 14px;
}

.pagination-controls {
  display: flex;
  gap: 8px;
}

.pagination-controls button {
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.pagination-controls button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-controls button:not(:disabled):hover {
  background: #f9fafb;
}
</style>
