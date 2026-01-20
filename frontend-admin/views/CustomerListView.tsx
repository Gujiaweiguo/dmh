import { ref, onMounted, computed } from 'vue';
import axios from '../services/axios';

export default {
  name: 'CustomerListView',
  setup() {
    const customers = ref<any[]>([]);
    const loading = ref(false);
    const total = ref(0);
    const currentPage = ref(1);
    const pageSize = ref(20);
    
    // 筛选条件
    const filters = ref({
      keyword: '',
      brandId: null as number | null,
      campaignId: null as number | null,
      status: '',
      startDate: '',
      endDate: '',
    });

    const brands = ref<any[]>([]);
    const campaigns = ref<any[]>([]);

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

    // 加载活动列表
    const loadCampaigns = async () => {
      try {
        const { data } = await axios.get('/api/v1/campaigns', {
          params: { page: 1, pageSize: 100, status: 'active' }
        });
        if (data.code === 200) {
          campaigns.value = data.data?.campaigns || [];
        }
      } catch (error) {
        console.error('加载活动列表失败:', error);
      }
    };

    // 加载顾客列表
    const loadCustomers = async () => {
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
        if (filters.value.campaignId) {
          params.campaignId = filters.value.campaignId;
        }
        if (filters.value.status) {
          params.status = filters.value.status;
        }
        if (filters.value.startDate) {
          params.startDate = filters.value.startDate;
        }
        if (filters.value.endDate) {
          params.endDate = filters.value.endDate;
        }

        const { data } = await axios.get('/api/v1/brands/:brandId/customers', { params });
        if (data.code === 200) {
          customers.value = data.data?.customers || [];
          total.value = data.data?.total || 0;
        }
      } catch (error: any) {
        console.error('加载顾客列表失败:', error);
        alert(error.response?.data?.message || '加载失败');
      } finally {
        loading.value = false;
      }
    };

    // 搜索
    const handleSearch = () => {
      currentPage.value = 1;
      loadCustomers();
    };

    // 重置筛选
    const handleReset = () => {
      filters.value = {
        keyword: '',
        brandId: null,
        campaignId: null,
        status: '',
        startDate: '',
        endDate: '',
      };
      handleSearch();
    };

    // 分页变化
    const handlePageChange = (page: number) => {
      currentPage.value = page;
      loadCustomers();
    };

    // 查看详情
    const viewDetail = (customerId: number) => {
      window.location.hash = `#/customers/${customerId}`;
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

    // 支付状态文本
    const paymentStatusText = (status: string) => {
      const map: Record<string, string> = {
        unpaid: '未支付',
        paid: '已支付',
        refunded: '已退款',
      };
      return map[status] || status;
    };

    // 支付状态颜色
    const paymentStatusColor = (status: string) => {
      const colors: Record<string, string> = {
        unpaid: 'orange',
        paid: 'green',
        refunded: 'gray',
      };
      return colors[status] || 'gray';
    };

    onMounted(() => {
      loadBrands();
      loadCampaigns();
      loadCustomers();
    });

    return {
      customers,
      loading,
      total,
      currentPage,
      pageSize,
      filters,
      brands,
      campaigns,
      loadCustomers,
      handleSearch,
      handleReset,
      handlePageChange,
      viewDetail,
      formatAmount,
      formatDate,
      paymentStatusText,
      paymentStatusColor,
    };
  },
  template: `
    <div class="customer-list-view">
      <div class="header">
        <h2>顾客管理</h2>
        <div class="actions">
          <button @click="loadCustomers" class="btn btn-primary">
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
            placeholder="搜索姓名/手机号"
            class="input"
            @keyup.enter="handleSearch"
          />
          <select v-model="filters.brandId" class="select">
            <option :value="null">全部品牌</option>
            <option v-for="brand in brands" :key="brand.id" :value="brand.id">
              {{ brand.name }}
            </option>
          </select>
          <select v-model="filters.campaignId" class="select">
            <option :value="null">全部活动</option>
            <option v-for="campaign in campaigns" :key="campaign.id" :value="campaign.id">
              {{ campaign.name }}
            </option>
          </select>
          <select v-model="filters.status" class="select">
            <option value="">全部状态</option>
            <option value="paid">已支付</option>
            <option value="unpaid">未支付</option>
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
 
      <!-- 顾客列表 -->
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>姓名</th>
              <th>手机号</th>
              <th>品牌</th>
              <th>活动</th>
              <th>订单数</th>
              <th>累计金额</th>
              <th>首单时间</th>
              <th>末单时间</th>
              <th>注册时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="customer in customers" :key="customer.id">
              <td>{{ customer.id }}</td>
              <td>{{ customer.username || '-' }}</td>
              <td>{{ customer.phone }}</td>
              <td>{{ customer.brandName || '-' }}</td>
              <td>{{ customer.campaignName || '-' }}</td>
              <td>{{ customer.orderCount }}</td>
              <td>{{ formatAmount(customer.totalAmount) }}</td>
              <td>{{ formatDate(customer.firstOrderAt) }}</td>
              <td>{{ formatDate(customer.lastOrderAt) }}</td>
              <td>{{ formatDate(customer.createdAt) }}</td>
              <td>
                <button 
                  @click="viewDetail(customer.id)" 
                  class="btn-link"
                >
                  查看详情
                </button>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr>
              <td colspan="11" class="text-center">加载中...</td>
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
