import { ref, onMounted } from 'vue';
import { memberApi } from '../services/memberApi';

// 公开接口类型定义 (Vue component instance 自动解包 Ref)
export interface MemberDetailViewInstance {
  member: any;
  loading: boolean;
  showTagDialog: boolean;
  availableTags: any[];
  selectedTags: number[];
  handleAddTags: () => Promise<void>;
  formatAmount: (amount: number) => string;
  formatDate: (date: string) => string;
  genderText: (gender: number) => string;
  formatFormData: (formData: Record<string, string>) => string;
  orderStatusText: (status: string) => string;
  goBack: () => void;
}

export default {
  name: 'MemberDetailView',
  setup() {
    const member = ref<any>(null);
    const loading = ref(false);
    const memberId = ref<number>(0);

    // 标签管理
    const showTagDialog = ref(false);
    const availableTags = ref<any[]>([]);
    const selectedTags = ref<number[]>([]);

    // 从 URL 获取会员 ID
    const getMemberIdFromUrl = () => {
      const hash = window.location.hash;
      const match = hash.match(/\/members\/(\d+)/);
      return match ? parseInt(match[1]) : 0;
    };

    // 加载会员详情
    const loadMemberDetail = async () => {
      loading.value = true;
      try {
        const response = await memberApi.getMember(memberId.value);
        member.value = response;
      } catch (error: any) {
        console.error('加载会员详情失败:', error);
        alert(error.message || '加载失败');
      } finally {
        loading.value = false;
      }
    };

    // 添加标签
    const handleAddTags = async () => {
      if (selectedTags.value.length === 0) {
        alert('请选择标签');
        return;
      }

      try {
        await memberApi.addMemberTags(memberId.value, selectedTags.value);
        alert('添加标签成功');
        showTagDialog.value = false;
        loadMemberDetail();
      } catch (error: any) {
        console.error('添加标签失败:', error);
        alert(error.message || '添加失败');
      }
    };

    // 格式化金额
    const formatAmount = (amount: number) => {
      return `¥${amount.toFixed(2)}`;
    };

    // 格式化日期
    const formatDate = (date: string) => {
      if (!date) return '-';
      return new Date(date).toLocaleString('zh-CN');
    };

    // 性别显示
    const genderText = (gender: number) => {
      const map: Record<number, string> = { 0: '未知', 1: '男', 2: '女' };
      return map[gender] || '未知';
    };

    const formatFormData = (formData: Record<string, string>) => {
      if (!formData) return '-';
      const entries = Object.entries(formData);
      if (entries.length === 0) return '-';
      return entries.map(([key, value]) => `${key}:${value}`).join('、');
    };

    const orderStatusText = (status: string) => {
      const map: Record<string, string> = {
        pending: '待支付',
        paid: '已支付',
        cancelled: '已取消',
      };
      return map[status] || status || '-';
    };

    // 返回列表
    const goBack = () => {
      window.location.hash = '#/members';
    };

    onMounted(() => {
      memberId.value = getMemberIdFromUrl();
      if (memberId.value) {
        loadMemberDetail();
      }
    });

    return {
      member,
      loading,
      showTagDialog,
      availableTags,
      selectedTags,
      handleAddTags,
      formatAmount,
      formatDate,
      genderText,
      formatFormData,
      orderStatusText,
      goBack,
    };
  },
  template: `
    <div class="member-detail-view">
      <div class="header">
        <button @click="goBack" class="btn btn-secondary">← 返回列表</button>
        <h2>会员详情</h2>
      </div>

      <div v-if="loading" class="loading">加载中...</div>

      <div v-else-if="member" class="detail-container">
        <!-- 基本信息 -->
        <div class="section">
          <h3>基本信息</h3>
          <div class="info-grid">
            <div class="info-item">
              <label>会员ID:</label>
              <span>{{ member.id }}</span>
            </div>
            <div class="info-item">
              <label>UnionID:</label>
              <span>{{ member.unionid }}</span>
            </div>
            <div class="info-item">
              <label>昵称:</label>
              <span>{{ member.nickname || '-' }}</span>
            </div>
            <div class="info-item">
              <label>手机号:</label>
              <span>{{ member.phone || '-' }}</span>
            </div>
            <div class="info-item">
              <label>性别:</label>
              <span>{{ genderText(member.gender) }}</span>
            </div>
            <div class="info-item">
              <label>来源:</label>
              <span>{{ member.source || '-' }}</span>
            </div>
            <div class="info-item">
              <label>状态:</label>
              <span :class="member.status === 'active' ? 'text-success' : 'text-danger'">
                {{ member.status === 'active' ? '正常' : '禁用' }}
              </span>
            </div>
            <div class="info-item">
              <label>注册时间:</label>
              <span>{{ formatDate(member.createdAt) }}</span>
            </div>
          </div>
        </div>

        <!-- 统计数据 -->
        <div class="section">
          <h3>统计数据</h3>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-label">累计订单</div>
              <div class="stat-value">{{ member.totalOrders }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-label">累计支付</div>
              <div class="stat-value">{{ formatAmount(member.totalPayment) }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-label">累计奖励</div>
              <div class="stat-value">{{ formatAmount(member.totalReward) }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-label">参与活动</div>
              <div class="stat-value">{{ member.participatedCampaigns }}</div>
            </div>
          </div>
          <div class="info-grid" style="margin-top: 20px;">
            <div class="info-item">
              <label>首次下单:</label>
              <span>{{ formatDate(member.firstOrderAt) }}</span>
            </div>
            <div class="info-item">
              <label>最后下单:</label>
              <span>{{ formatDate(member.lastOrderAt) }}</span>
            </div>
          </div>
        </div>

        <!-- 标签 -->
        <div class="section">
          <div class="section-header">
            <h3>会员标签</h3>
            <button @click="showTagDialog = true" class="btn btn-primary btn-sm">
              添加标签
            </button>
          </div>
          <div class="tags">
            <span 
              v-for="tag in member.tags" 
              :key="tag.id" 
              class="tag"
              :style="{ backgroundColor: tag.color || '#999' }"
            >
              {{ tag.name }}
            </span>
            <span v-if="!member.tags || member.tags.length === 0" class="text-muted">
              暂无标签
            </span>
          </div>
        </div>

        <!-- 关联品牌 -->
        <div class="section">
          <h3>关联品牌</h3>
          <div class="table-container">
            <table class="table">
              <thead>
                <tr>
                  <th>品牌名称</th>
                  <th>首次参与活动ID</th>
                  <th>关联时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="brand in member.brands" :key="brand.brandId">
                  <td>{{ brand.brandName }}</td>
                  <td>{{ brand.firstCampaignId }}</td>
                  <td>{{ formatDate(brand.createdAt) }}</td>
                </tr>
                <tr v-if="!member.brands || member.brands.length === 0">
                  <td colspan="3" class="text-center text-muted">暂无关联品牌</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 订单记录 -->
        <div class="section">
          <h3>订单记录</h3>
          <div class="table-container">
            <table class="table">
              <thead>
                <tr>
                  <th>订单ID</th>
                  <th>活动</th>
                  <th>金额</th>
                  <th>状态</th>
                  <th>表单</th>
                  <th>创建时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="order in member.orders" :key="order.id">
                  <td>{{ order.id }}</td>
                  <td>{{ order.campaignName || order.campaignId }}</td>
                  <td>{{ formatAmount(order.amount || 0) }}</td>
                  <td>{{ orderStatusText(order.status) }}</td>
                  <td>{{ formatFormData(order.formData) }}</td>
                  <td>{{ formatDate(order.createdAt) }}</td>
                </tr>
                <tr v-if="!member.orders || member.orders.length === 0">
                  <td colspan="6" class="text-center text-muted">暂无订单</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 添加标签对话框 -->
      <div v-if="showTagDialog" class="modal">
        <div class="modal-content">
          <div class="modal-header">
            <h3>添加标签</h3>
            <button @click="showTagDialog = false" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>选择标签:</label>
              <div class="checkbox-group">
                <label v-for="tag in availableTags" :key="tag.id">
                  <input 
                    type="checkbox" 
                    :value="tag.id"
                    v-model="selectedTags"
                  />
                  {{ tag.name }}
                </label>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button @click="showTagDialog = false" class="btn btn-secondary">
              取消
            </button>
            <button @click="handleAddTags" class="btn btn-primary">
              确定
            </button>
          </div>
        </div>
      </div>
    </div>
  `,
};
