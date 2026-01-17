import { ref, onMounted, computed } from 'vue';
import { memberApi } from '../services/memberApi';

export default {
  name: 'MemberListView',
  setup() {
    const members = ref<any[]>([]);
    const loading = ref(false);
    const total = ref(0);
    const currentPage = ref(1);
    const pageSize = ref(20);
    
    // 筛选条件
    const filters = ref({
      keyword: '',
      brandId: null as number | null,
      source: '',
      status: '',
      startDate: '',
      endDate: '',
    });

    // 选中的会员
    const selectedMembers = ref<number[]>([]);

    // 加载会员列表
    const loadMembers = async () => {
      loading.value = true;
      try {
        const response = await memberApi.getMembers({
          page: currentPage.value,
          pageSize: pageSize.value,
          ...filters.value,
        });
        members.value = response.members || [];
        total.value = response.total || 0;
      } catch (error: any) {
        console.error('加载会员列表失败:', error);
        const errorMsg = error.message || '加载失败';
        // 如果是登录相关错误，提示用户刷新页面
        if (errorMsg.includes('登录')) {
          alert(errorMsg + '\n\n点击确定后将刷新页面');
          window.location.reload();
        } else {
          alert(errorMsg);
        }
      } finally {
        loading.value = false;
      }
    };

    // 搜索
    const handleSearch = () => {
      currentPage.value = 1;
      loadMembers();
    };

    // 重置筛选
    const handleReset = () => {
      filters.value = {
        keyword: '',
        brandId: null,
        source: '',
        status: '',
        startDate: '',
        endDate: '',
      };
      handleSearch();
    };

    // 分页变化
    const handlePageChange = (page: number) => {
      currentPage.value = page;
      loadMembers();
    };

    // 查看详情
    const viewDetail = (memberId: number) => {
      window.location.hash = `#/members/${memberId}`;
    };

    // 导出会员
    const handleExport = () => {
      window.location.hash = '#/members/export';
    };

    // 合并会员
    const handleMerge = () => {
      if (selectedMembers.value.length !== 2) {
        alert('请选择两个会员进行合并');
        return;
      }
      window.location.hash = `#/members/merge?source=${selectedMembers.value[0]}&target=${selectedMembers.value[1]}`;
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

    // 状态颜色
    const statusColor = (status: string) => {
      const colors: Record<string, string> = {
        active: 'green',
        disabled: 'red',
      };
      return colors[status] || 'gray';
    };

    onMounted(() => {
      loadMembers();
    });

    return {
      members,
      loading,
      total,
      currentPage,
      pageSize,
      filters,
      selectedMembers,
      loadMembers,
      handleSearch,
      handleReset,
      handlePageChange,
      viewDetail,
      handleExport,
      handleMerge,
      formatAmount,
      formatDate,
      genderText,
      statusColor,
    };
  },
  template: `
    <div class="member-list-view">
      <div class="header">
        <h2>会员管理</h2>
        <div class="actions">
          <button @click="handleExport" class="btn btn-primary">
            导出会员
          </button>
          <button 
            @click="handleMerge" 
            class="btn btn-secondary"
            :disabled="selectedMembers.length !== 2"
          >
            合并会员
          </button>
        </div>
      </div>

      <!-- 筛选条件 -->
      <div class="filters">
        <div class="filter-row">
          <input
            v-model="filters.keyword"
            type="text"
            placeholder="搜索昵称/手机号/UnionID"
            class="input"
            @keyup.enter="handleSearch"
          />
          <select v-model="filters.status" class="select">
            <option value="">全部状态</option>
            <option value="active">正常</option>
            <option value="disabled">禁用</option>
          </select>
          <input
            v-model="filters.startDate"
            type="date"
            placeholder="开始日期"
            class="input"
          />
          <input
            v-model="filters.endDate"
            type="date"
            placeholder="结束日期"
            class="input"
          />
          <button @click="handleSearch" class="btn btn-primary">搜索</button>
          <button @click="handleReset" class="btn btn-secondary">重置</button>
        </div>
      </div>

      <!-- 会员列表 -->
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th width="50">
                <input type="checkbox" />
              </th>
              <th>ID</th>
              <th>昵称</th>
              <th>手机号</th>
              <th>性别</th>
              <th>来源</th>
              <th>订单数</th>
              <th>累计支付</th>
              <th>累计奖励</th>
              <th>参与活动</th>
              <th>状态</th>
              <th>注册时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="member in members" :key="member.id">
              <td>
                <input 
                  type="checkbox" 
                  :value="member.id"
                  v-model="selectedMembers"
                />
              </td>
              <td>{{ member.id }}</td>
              <td>
                <div class="member-info">
                  <img 
                    v-if="member.avatar" 
                    :src="member.avatar" 
                    class="avatar"
                  />
                  <span>{{ member.nickname || '-' }}</span>
                </div>
              </td>
              <td>{{ member.phone || '-' }}</td>
              <td>{{ genderText(member.gender) }}</td>
              <td>{{ member.source || '-' }}</td>
              <td>{{ member.totalOrders }}</td>
              <td>{{ formatAmount(member.totalPayment) }}</td>
              <td>{{ formatAmount(member.totalReward) }}</td>
              <td>{{ member.participatedCampaigns }}</td>
              <td>
                <span 
                  class="status-badge" 
                  :style="{ color: statusColor(member.status) }"
                >
                  {{ member.status === 'active' ? '正常' : '禁用' }}
                </span>
              </td>
              <td>{{ formatDate(member.createdAt) }}</td>
              <td>
                <button 
                  @click="viewDetail(member.id)" 
                  class="btn-link"
                >
                  查看详情
                </button>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr>
              <td colspan="13" class="text-center">加载中...</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div class="pagination" v-if="total > 0">
        <button 
          @click="handlePageChange(currentPage - 1)"
          :disabled="currentPage === 1"
          class="btn btn-secondary"
        >
          上一页
        </button>
        <span class="page-info">
          第 {{ currentPage }} 页，共 {{ Math.ceil(total / pageSize) }} 页，
          总计 {{ total }} 条
        </span>
        <button 
          @click="handlePageChange(currentPage + 1)"
          :disabled="currentPage >= Math.ceil(total / pageSize)"
          class="btn btn-secondary"
        >
          下一页
        </button>
      </div>
    </div>
  `,
};
