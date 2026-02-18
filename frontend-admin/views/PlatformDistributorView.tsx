import { ref, onMounted } from 'vue';
import axios from '../services/axios';

// 公开接口类型定义 (Vue component instance 自动解包 Ref)
export interface PlatformDistributorViewInstance {
  distributors: any[];
  loading: boolean;
  total: number;
  currentPage: number;
  pageSize: number;
  filters: {
    keyword: string;
    brandId: number | null;
    level: number | null;
    status: string;
  };
  brands: any[];
  loadDistributors: () => Promise<void>;
  loadBrands: () => Promise<void>;
  handleSearch: () => void;
  handleReset: () => void;
  handlePageChange: (page: number) => void;
  viewDetail: (distributorId: number) => void;
  adjustLevel: (distributorId: number, currentLevel: number) => Promise<void>;
  toggleStatus: (distributorId: number, currentStatus: string) => Promise<void>;
  formatAmount: (amount: number) => string;
  formatDate: (date: string) => string;
  statusText: (status: string) => string;
  levelText: (level: number) => string;
  statusColor: (status: string) => string;
}

export default {
  name: 'PlatformDistributorView',
  setup() {
    const distributors = ref<any[]>([]);
    const loading = ref(false);
    const total = ref(0);
    const currentPage = ref(1);
    const pageSize = ref(20);
    
    // 筛选条件
    const filters = ref({
      keyword: '',
      brandId: null as number | null,
      level: null as number | null,
      status: '',
    });

    const brands = ref<any[]>([]);

    // 加载品牌列表
    const loadBrands = async () => {
      try {
        const { data } = await axios.get('/api/v1/brands', {
          params: { page: 1, pageSize: 100 }
        });
        if (data.code === 200) {
          brands.value = data.data?.brands || [];
        }
      } catch (error) {
        console.error('加载品牌列表失败:', error);
      }
    };

    // 加载分销商列表
    const loadDistributors = async () => {
      loading.value = true;
      try {
        const params: any = {
          page: currentPage.value,
          pageSize: pageSize.value,
        };
        
        if (filters.value.keyword) {
          params.keyword = filters.value.keyword;
        }
        if (filters.value.brandId) {
          params.brandId = filters.value.brandId;
        }
        if (filters.value.level) {
          params.level = filters.value.level;
        }
        if (filters.value.status) {
          params.status = filters.value.status;
        }

        const { data } = await axios.get('/api/v1/platform/distributors', { params });
        if (data.code === 200) {
          distributors.value = data.data?.distributors || [];
          total.value = data.data?.total || 0;
        }
      } catch (error: any) {
        console.error('加载分销商列表失败:', error);
        alert(error.response?.data?.message || '加载失败');
      } finally {
        loading.value = false;
      }
    };

    // 搜索
    const handleSearch = () => {
      currentPage.value = 1;
      loadDistributors();
    };

    // 重置筛选
    const handleReset = () => {
      filters.value = {
        keyword: '',
        brandId: null,
        level: null,
        status: '',
      };
      handleSearch();
    };

    // 分页变化
    const handlePageChange = (page: number) => {
      currentPage.value = page;
      loadDistributors();
    };

    // 查看详情
    const viewDetail = (distributorId: number) => {
      window.location.hash = `#/distributors/${distributorId}`;
    };

    // 调整级别
    const adjustLevel = async (distributorId: number, currentLevel: number) => {
      const newLevel = prompt('请输入新的级别 (1-3):', String(currentLevel));
      if (newLevel && !isNaN(Number(newLevel))) {
        const level = Number(newLevel);
        if (level < 1 || level > 3) {
          alert('级别必须在1-3之间');
          return;
        }
        try {
          const { data } = await axios.put(`/api/v1/distributors/${distributorId}/level`, { level });
          if (data.code === 200) {
            alert('级别更新成功');
            loadDistributors();
          }
        } catch (error: any) {
          alert(error.response?.data?.message || '更新失败');
        }
      }
    };

    // 调整状态
    const toggleStatus = async (distributorId: number, currentStatus: string) => {
      const newStatus = currentStatus === 'active' ? 'suspended' : 'active';
      try {
        const { data } = await axios.put(`/api/v1/distributors/${distributorId}/status`, { status: newStatus });
        if (data.code === 200) {
          alert('状态更新成功');
          loadDistributors();
        }
      } catch (error: any) {
        alert(error.response?.data?.message || '更新失败');
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

    // 状态文本
    const statusText = (status: string) => {
      const map: Record<string, string> = {
        active: '激活',
        suspended: '暂停',
        pending: '待审核',
      };
      return map[status] || status;
    };

    // 级别文本
    const levelText = (level: number) => {
      const map: Record<number, string> = {
        1: '一级',
        2: '二级',
        3: '三级',
      };
      return map[level] || `${level}级`;
    };

    // 状态颜色
    const statusColor = (status: string) => {
      const colors: Record<string, string> = {
        active: 'green',
        suspended: 'orange',
        pending: 'gray',
      };
      return colors[status] || 'gray';
    };

    onMounted(() => {
      loadBrands();
      loadDistributors();
    });

    return {
      distributors,
      loading,
      total,
      currentPage,
      pageSize,
      filters,
      brands,
      loadDistributors,
      loadBrands,
      handleSearch,
      handleReset,
      handlePageChange,
      viewDetail,
      adjustLevel,
      toggleStatus,
      formatAmount,
      formatDate,
      statusText,
      levelText,
      statusColor,
    };
  },
  template: `
    <div class="platform-distributor-view">
      <div class="header">
        <h2>全局分销商管理</h2>
        <div class="actions">
          <button @click="loadDistributors" class="btn btn-primary">
            刷新列表
          </button>
        </div>
      </div>
 
      <!-- 筛选条件 -->
      <div class="filters">
        <div class="filter-row">
          <input
            v-model="filters.keyword"
            type="text"
            placeholder="搜索用户名/手机号"
            class="input"
            @keyup.enter="handleSearch"
          />
          <select v-model="filters.brandId" class="select">
            <option :value="null">全部品牌</option>
            <option v-for="brand in brands" :key="brand.id" :value="brand.id">
              {{ brand.name }}
            </option>
          </select>
          <select v-model="filters.level" class="select">
            <option :value="null">全部级别</option>
            <option :value="1">一级</option>
            <option :value="2">二级</option>
            <option :value="3">三级</option>
          </select>
          <select v-model="filters.status" class="select">
            <option value="">全部状态</option>
            <option value="active">激活</option>
            <option value="suspended">暂停</option>
            <option value="pending">待审核</option>
          </select>
          <button @click="handleSearch" class="btn btn-primary">搜索</button>
          <button @click="handleReset" class="btn btn-secondary">重置</button>
        </div>
      </div>
 
      <!-- 分销商列表 -->
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>用户</th>
              <th>品牌</th>
              <th>级别</th>
              <th>状态</th>
              <th>累计收益</th>
              <th>下级数</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="distributor in distributors" :key="distributor.id">
              <td>{{ distributor.id }}</td>
              <td>{{ distributor.username || '-' }}</td>
              <td>{{ distributor.brandName || '-' }}</td>
              <td>{{ levelText(distributor.level) }}</td>
              <td>
                <span 
                  class="status-badge" 
                  :style="{ color: statusColor(distributor.status) }"
                >
                  {{ statusText(distributor.status) }}
                </span>
              </td>
              <td>{{ formatAmount(distributor.totalEarnings) }}</td>
              <td>{{ distributor.subordinatesCount }}</td>
              <td>{{ formatDate(distributor.createdAt) }}</td>
              <td>
                <button 
                  @click="viewDetail(distributor.id)" 
                  class="btn-link"
                >
                  详情
                </button>
                <button 
                  @click="adjustLevel(distributor.id, distributor.level)" 
                  class="btn-link"
                >
                  调级
                </button>
                <button 
                  @click="toggleStatus(distributor.id, distributor.status)" 
                  class="btn-link"
                >
                  {{ distributor.status === 'active' ? '暂停' : '激活' }}
                </button>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr>
              <td colspan="10" class="text-center">加载中...</td>
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
