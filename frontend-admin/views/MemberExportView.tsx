import { ref, onMounted } from 'vue';
import { memberApi } from '../services/memberApi';

// 公开接口类型定义 (Vue component instance 自动解包 Ref)
export interface MemberExportViewInstance {
  loading: boolean;
  exportRequests: any[];
  total: number;
  currentPage: number;
  pageSize: number;
  showCreateDialog: boolean;
  createForm: {
    brandId: number | null;
    reason: string;
    filters: string;
  };
  showApproveDialog: boolean;
  approveForm: {
    requestId: number;
    approve: boolean;
    reason: string;
  };
  loadExportRequests: () => Promise<void>;
  handleCreate: () => Promise<void>;
  openApproveDialog: (request: any, approve: boolean) => void;
  handleApprove: () => Promise<void>;
  handleDownload: (fileUrl: string) => void;
  handlePageChange: (page: number) => void;
  formatDate: (date: string) => string;
  statusColor: (status: string) => string;
  statusText: (status: string) => string;
  goBack: () => void;
}

export default {
  name: 'MemberExportView',
  setup() {
    const loading = ref(false);
    const exportRequests = ref<any[]>([]);
    const total = ref(0);
    const currentPage = ref(1);
    const pageSize = ref(20);

    // 新建导出申请
    const showCreateDialog = ref(false);
    const createForm = ref({
      brandId: null as number | null,
      reason: '',
      filters: '',
    });

    // 审批对话框
    const showApproveDialog = ref(false);
    const approveForm = ref({
      requestId: 0,
      approve: true,
      reason: '',
    });

    // 加载导出申请列表
    const loadExportRequests = async () => {
      loading.value = true;
      try {
        const response = await memberApi.getExportRequests({
          page: currentPage.value,
          pageSize: pageSize.value,
        });
        exportRequests.value = response.requests || [];
        total.value = response.total || 0;
      } catch (error: any) {
        console.error('加载导出申请失败:', error);
        alert(error.message || '加载失败');
      } finally {
        loading.value = false;
      }
    };

    // 创建导出申请
    const handleCreate = async () => {
      if (!createForm.value.brandId) {
        alert('请选择品牌');
        return;
      }
      if (!createForm.value.reason) {
        alert('请填写导出原因');
        return;
      }

      loading.value = true;
      try {
        await memberApi.createExportRequest(createForm.value as any);
        alert('申请提交成功');
        showCreateDialog.value = false;
        createForm.value = {
          brandId: null,
          reason: '',
          filters: '',
        };
        loadExportRequests();
      } catch (error: any) {
        console.error('创建导出申请失败:', error);
        alert(error.message || '创建失败');
      } finally {
        loading.value = false;
      }
    };

    // 打开审批对话框
    const openApproveDialog = (request: any, approve: boolean) => {
      approveForm.value = {
        requestId: request.id,
        approve,
        reason: '',
      };
      showApproveDialog.value = true;
    };

    // 审批导出申请
    const handleApprove = async () => {
      if (!approveForm.value.approve && !approveForm.value.reason) {
        alert('驳回时请填写原因');
        return;
      }

      loading.value = true;
      try {
        await memberApi.approveExportRequest(
          approveForm.value.requestId,
          {
            approve: approveForm.value.approve,
            reason: approveForm.value.reason,
          }
        );
        alert(approveForm.value.approve ? '审批通过' : '已驳回');
        showApproveDialog.value = false;
        loadExportRequests();
      } catch (error: any) {
        console.error('审批失败:', error);
        alert(error.message || '审批失败');
      } finally {
        loading.value = false;
      }
    };

    // 下载导出文件
    const handleDownload = (fileUrl: string) => {
      window.open(fileUrl, '_blank');
    };

    // 分页变化
    const handlePageChange = (page: number) => {
      currentPage.value = page;
      loadExportRequests();
    };

    // 格式化日期
    const formatDate = (date: string) => {
      if (!date) return '-';
      return new Date(date).toLocaleString('zh-CN');
    };

    // 状态颜色
    const statusColor = (status: string) => {
      const colors: Record<string, string> = {
        pending: 'orange',
        approved: 'blue',
        rejected: 'red',
        completed: 'green',
      };
      return colors[status] || 'gray';
    };

    // 状态文本
    const statusText = (status: string) => {
      const texts: Record<string, string> = {
        pending: '待审批',
        approved: '已批准',
        rejected: '已驳回',
        completed: '已完成',
      };
      return texts[status] || status;
    };

    // 返回列表
    const goBack = () => {
      window.location.hash = '#/members';
    };

    onMounted(() => {
      loadExportRequests();
    });

    return {
      loading,
      exportRequests,
      total,
      currentPage,
      pageSize,
      showCreateDialog,
      createForm,
      showApproveDialog,
      approveForm,
      loadExportRequests,
      handleCreate,
      openApproveDialog,
      handleApprove,
      handleDownload,
      handlePageChange,
      formatDate,
      statusColor,
      statusText,
      goBack,
    };
  },
  template: `
    <div class="member-export-view">
      <div class="header">
        <button @click="goBack" class="btn btn-secondary">← 返回列表</button>
        <h2>会员导出管理</h2>
        <button @click="showCreateDialog = true" class="btn btn-primary">
          新建导出申请
        </button>
      </div>

      <!-- 导出申请列表 -->
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>申请ID</th>
              <th>品牌</th>
              <th>申请人</th>
              <th>申请原因</th>
              <th>状态</th>
              <th>记录数</th>
              <th>审批人</th>
              <th>申请时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="request in exportRequests" :key="request.id">
              <td>{{ request.id }}</td>
              <td>{{ request.brandName }}</td>
              <td>{{ request.requestedByName }}</td>
              <td>{{ request.reason }}</td>
              <td>
                <span 
                  class="status-badge"
                  :style="{ color: statusColor(request.status) }"
                >
                  {{ statusText(request.status) }}
                </span>
              </td>
              <td>{{ request.recordCount || '-' }}</td>
              <td>{{ request.approvedByName || '-' }}</td>
              <td>{{ formatDate(request.createdAt) }}</td>
              <td>
                <div class="action-buttons">
                  <button 
                    v-if="request.status === 'pending'"
                    @click="openApproveDialog(request, true)"
                    class="btn-link text-success"
                  >
                    批准
                  </button>
                  <button 
                    v-if="request.status === 'pending'"
                    @click="openApproveDialog(request, false)"
                    class="btn-link text-danger"
                  >
                    驳回
                  </button>
                  <button 
                    v-if="request.status === 'completed' && request.fileUrl"
                    @click="handleDownload(request.fileUrl)"
                    class="btn-link"
                  >
                    下载
                  </button>
                  <span v-if="request.status === 'rejected'" class="text-muted">
                    {{ request.rejectReason }}
                  </span>
                </div>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr>
              <td colspan="9" class="text-center">加载中...</td>
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
          第 {{ currentPage }} 页，共 {{ Math.ceil(total / pageSize) }} 页
        </span>
        <button 
          @click="handlePageChange(currentPage + 1)"
          :disabled="currentPage >= Math.ceil(total / pageSize)"
          class="btn btn-secondary"
        >
          下一页
        </button>
      </div>

      <!-- 新建导出申请对话框 -->
      <div v-if="showCreateDialog" class="modal">
        <div class="modal-content">
          <div class="modal-header">
            <h3>新建导出申请</h3>
            <button @click="showCreateDialog = false" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>品牌ID:</label>
              <input 
                v-model.number="createForm.brandId" 
                type="number"
                class="input"
                placeholder="输入品牌ID"
              />
            </div>
            <div class="form-group">
              <label>导出原因:</label>
              <textarea 
                v-model="createForm.reason" 
                class="textarea"
                placeholder="请说明导出原因"
                rows="3"
              ></textarea>
            </div>
            <div class="form-group">
              <label>筛选条件（可选）:</label>
              <textarea 
                v-model="createForm.filters" 
                class="textarea"
                placeholder='JSON格式，例如: {"status":"active"}'
                rows="2"
              ></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button @click="showCreateDialog = false" class="btn btn-secondary">
              取消
            </button>
            <button @click="handleCreate" class="btn btn-primary" :disabled="loading">
              {{ loading ? '提交中...' : '提交申请' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 审批对话框 -->
      <div v-if="showApproveDialog" class="modal">
        <div class="modal-content">
          <div class="modal-header">
            <h3>{{ approveForm.approve ? '批准导出' : '驳回申请' }}</h3>
            <button @click="showApproveDialog = false" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <div class="form-group" v-if="!approveForm.approve">
              <label>驳回原因:</label>
              <textarea 
                v-model="approveForm.reason" 
                class="textarea"
                placeholder="请说明驳回原因"
                rows="3"
              ></textarea>
            </div>
            <p v-else>
              确认批准此导出申请吗？系统将生成导出文件。
            </p>
          </div>
          <div class="modal-footer">
            <button @click="showApproveDialog = false" class="btn btn-secondary">
              取消
            </button>
            <button 
              @click="handleApprove" 
              :class="approveForm.approve ? 'btn btn-primary' : 'btn btn-danger'"
              :disabled="loading"
            >
              {{ loading ? '处理中...' : '确认' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  `,
};
