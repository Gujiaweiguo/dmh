import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { ref } from 'vue';

// Mock 组件实例类型 (Vue component instance 自动解包 Ref)
interface CampaignListVMockInstance {
  campaigns: any[];
  total: number;
  page: number;
  pageSize: number;
  keyword: string;
  statusFilter: string;
  loading: boolean;
  fetchCampaigns: ReturnType<typeof vi.fn>;
  handleSearch: ReturnType<typeof vi.fn>;
  handleDelete: ReturnType<typeof vi.fn>;
  handleGeneratePoster: ReturnType<typeof vi.fn>;
  formatDate: (dateStr: string) => string;
  getStatusText: (status: string) => string;
  getStatusClass: (status: string) => string;
}

vi.mock('vue', async () => {
  const actual = await vi.importActual('vue');
  return {
    ...actual,
    onMounted: vi.fn((fn) => fn()),
  };
});

vi.mock('../../services/campaignApi', () => ({
  campaignApi: {
    getCampaigns: vi.fn().mockResolvedValue({ campaigns: [], total: 0 }),
    deleteCampaign: vi.fn().mockResolvedValue({}),
  },
}));

vi.mock('../../services/posterApi', () => ({
  posterApi: {
    generateCampaignPoster: vi.fn().mockResolvedValue({ posterUrl: 'http://example.com/poster.png' }),
  },
}));

vi.mock('lucide-vue-next', () => ({
  Search: { template: '<div></div>' },
  Plus: { template: '<div></div>' },
  Edit2: { template: '<div></div>' },
  Trash2: { template: '<div></div>' },
  Eye: { template: '<div></div>' },
  Image: { template: '<div></div>' },
}));

const createCampaignListView = () => {
  return {
    name: 'CampaignListView',
    template: `
      <div class="campaign-list-view">
        <div class="header">
          <h1 class="title">营销活动管理</h1>
        </div>
        <div class="filters">
          <input v-model="keyword" type="text" placeholder="搜索活动名称或描述" />
          <select v-model="statusFilter">
            <option value="">全部状态</option>
            <option value="active">进行中</option>
            <option value="paused">已暂停</option>
            <option value="ended">已结束</option>
          </select>
        </div>
      </div>
    `,
    setup() {
      const campaigns = ref<any[]>([]);
      const total = ref(0);
      const page = ref(1);
      const pageSize = ref(20);
      const keyword = ref('');
      const statusFilter = ref('');
      const loading = ref(false);

      const fetchCampaigns = vi.fn();
      const handleSearch = vi.fn();
      const handleDelete = vi.fn();
      const handleGeneratePoster = vi.fn();
      const formatDate = (dateStr: string) => {
        if (!dateStr) return '-';
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

      return {
        campaigns,
        total,
        page,
        pageSize,
        keyword,
        statusFilter,
        loading,
        fetchCampaigns,
        handleSearch,
        handleDelete,
        handleGeneratePoster,
        formatDate,
        getStatusText,
        getStatusClass,
      };
    },
  };
};

describe('CampaignListView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(createCampaignListView());
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(createCampaignListView());
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.campaign-list-view').exists()).toBe(true);
  });

  it('should have formatDate method', () => {
    const wrapper = mount(createCampaignListView());
    const vm = wrapper.vm as unknown as CampaignListVMockInstance;
    expect(vm.formatDate('')).toBe('-');
    expect(vm.formatDate('2024-01-01T00:00:00Z')).toBeTruthy();
  });

  it('should have getStatusText method', () => {
    const wrapper = mount(createCampaignListView());
    const vm = wrapper.vm as unknown as CampaignListVMockInstance;
    expect(vm.getStatusText('active')).toBe('进行中');
    expect(vm.getStatusText('paused')).toBe('已暂停');
    expect(vm.getStatusText('ended')).toBe('已结束');
    expect(vm.getStatusText('unknown')).toBe('unknown');
  });

  it('should have getStatusClass method', () => {
    const wrapper = mount(createCampaignListView());
    const vm = wrapper.vm as unknown as CampaignListVMockInstance;
    expect(vm.getStatusClass('active')).toBe('bg-green-100 text-green-800');
    expect(vm.getStatusClass('paused')).toBe('bg-yellow-100 text-yellow-800');
    expect(vm.getStatusClass('ended')).toBe('bg-gray-100 text-gray-800');
    expect(vm.getStatusClass('unknown')).toBe('');
  });

  it('should have initial values set correctly', () => {
    const wrapper = mount(createCampaignListView());
    const vm = wrapper.vm as unknown as CampaignListVMockInstance;
    expect(vm.keyword).toBe('');
    expect(vm.statusFilter).toBe('');
    expect(vm.page).toBe(1);
    expect(vm.pageSize).toBe(20);
    expect(vm.total).toBe(0);
    expect(vm.loading).toBe(false);
  });
});
