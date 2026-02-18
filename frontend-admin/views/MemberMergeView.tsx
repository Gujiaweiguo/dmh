import { ref, onMounted } from 'vue';
import { memberApi } from '../services/memberApi';

// 公开接口类型定义 (Vue component instance 自动解包 Ref)
export interface MemberMergeViewInstance {
  loading: boolean;
  sourceMemberId: number;
  targetMemberId: number;
  reason: string;
  preview: any;
  showPreview: boolean;
  handlePreview: () => Promise<void>;
  handleMerge: () => Promise<void>;
  goBack: () => void;
  formatAmount: (amount: number) => string;
}

export default {
  name: 'MemberMergeView',
  setup() {
    const loading = ref(false);
    const sourceMemberId = ref<number>(0);
    const targetMemberId = ref<number>(0);
    const reason = ref('');
    const preview = ref<any>(null);
    const showPreview = ref(false);

    // 从 URL 获取参数
    const getParamsFromUrl = () => {
      const hash = window.location.hash;
      const params = new URLSearchParams(hash.split('?')[1] || '');
      sourceMemberId.value = parseInt(params.get('source') || '0');
      targetMemberId.value = parseInt(params.get('target') || '0');
    };

    // 预览合并
    const handlePreview = async () => {
      if (!sourceMemberId.value || !targetMemberId.value) {
        alert('请输入源会员ID和目标会员ID');
        return;
      }

      if (sourceMemberId.value === targetMemberId.value) {
        alert('源会员和目标会员不能相同');
        return;
      }

      loading.value = true;
      try {
        const response = await memberApi.previewMerge({
          sourceMemberId: sourceMemberId.value,
          targetMemberId: targetMemberId.value,
          reason: reason.value,
        });
        preview.value = response;
        showPreview.value = true;
      } catch (error: any) {
        console.error('预览合并失败:', error);
        alert(error.message || '预览失败');
      } finally {
        loading.value = false;
      }
    };

    // 执行合并
    const handleMerge = async () => {
      if (!preview.value?.canMerge) {
        alert('存在冲突，无法合并');
        return;
      }

      if (!confirm('确定要合并这两个会员吗？此操作不可逆！')) {
        return;
      }

      loading.value = true;
      try {
        await memberApi.mergeMember({
          sourceMemberId: sourceMemberId.value,
          targetMemberId: targetMemberId.value,
          reason: reason.value,
        });
        alert('合并成功');
        window.location.hash = '#/members';
      } catch (error: any) {
        console.error('合并失败:', error);
        alert(error.message || '合并失败');
      } finally {
        loading.value = false;
      }
    };

    // 返回列表
    const goBack = () => {
      window.location.hash = '#/members';
    };

    // 格式化金额
    const formatAmount = (amount: number) => {
      return `¥${amount.toFixed(2)}`;
    };

    onMounted(() => {
      getParamsFromUrl();
    });

    return {
      loading,
      sourceMemberId,
      targetMemberId,
      reason,
      preview,
      showPreview,
      handlePreview,
      handleMerge,
      goBack,
      formatAmount,
    };
  },
  template: `
    <div class="member-merge-view">
      <div class="header">
        <button @click="goBack" class="btn btn-secondary">← 返回列表</button>
        <h2>会员合并</h2>
      </div>

      <div class="merge-form">
        <div class="form-section">
          <h3>合并信息</h3>
          <div class="form-group">
            <label>源会员ID（将被合并）:</label>
            <input 
              v-model.number="sourceMemberId" 
              type="number" 
              class="input"
              placeholder="输入源会员ID"
            />
          </div>
          <div class="form-group">
            <label>目标会员ID（保留）:</label>
            <input 
              v-model.number="targetMemberId" 
              type="number" 
              class="input"
              placeholder="输入目标会员ID"
            />
          </div>
          <div class="form-group">
            <label>合并原因:</label>
            <textarea 
              v-model="reason" 
              class="textarea"
              placeholder="请说明合并原因"
              rows="3"
            ></textarea>
          </div>
          <div class="form-actions">
            <button 
              @click="handlePreview" 
              class="btn btn-primary"
              :disabled="loading"
            >
              {{ loading ? '加载中...' : '预览合并' }}
            </button>
          </div>
        </div>

        <!-- 预览结果 -->
        <div v-if="showPreview && preview" class="preview-section">
          <h3>合并预览</h3>
          
          <!-- 冲突提示 -->
          <div v-if="!preview.canMerge" class="alert alert-danger">
            <strong>存在冲突，无法合并！</strong>
            <ul>
              <li v-for="(conflict, index) in preview.conflicts" :key="index">
                {{ conflict }}
              </li>
            </ul>
          </div>

          <div v-else class="alert alert-success">
            <strong>可以合并</strong>
            <ul v-if="preview.conflicts && preview.conflicts.length > 0">
              <li v-for="(conflict, index) in preview.conflicts" :key="index">
                {{ conflict }}
              </li>
            </ul>
          </div>

          <!-- 会员对比 -->
          <div class="member-compare">
            <div class="member-card">
              <h4>源会员（将被合并）</h4>
              <div class="member-info">
                <p><strong>ID:</strong> {{ preview.sourceMember.id }}</p>
                <p><strong>昵称:</strong> {{ preview.sourceMember.nickname }}</p>
                <p><strong>手机号:</strong> {{ preview.sourceMember.phone }}</p>
                <p><strong>UnionID:</strong> {{ preview.sourceMember.unionid }}</p>
                <p><strong>订单数:</strong> {{ preview.sourceMember.totalOrders }}</p>
                <p><strong>累计支付:</strong> {{ formatAmount(preview.sourceMember.totalPayment) }}</p>
                <p><strong>累计奖励:</strong> {{ formatAmount(preview.sourceMember.totalReward) }}</p>
              </div>
            </div>

            <div class="merge-arrow">→</div>

            <div class="member-card">
              <h4>目标会员（保留）</h4>
              <div class="member-info">
                <p><strong>ID:</strong> {{ preview.targetMember.id }}</p>
                <p><strong>昵称:</strong> {{ preview.targetMember.nickname }}</p>
                <p><strong>手机号:</strong> {{ preview.targetMember.phone }}</p>
                <p><strong>UnionID:</strong> {{ preview.targetMember.unionid }}</p>
                <p><strong>订单数:</strong> {{ preview.targetMember.totalOrders }}</p>
                <p><strong>累计支付:</strong> {{ formatAmount(preview.targetMember.totalPayment) }}</p>
                <p><strong>累计奖励:</strong> {{ formatAmount(preview.targetMember.totalReward) }}</p>
              </div>
            </div>
          </div>

          <!-- 合并按钮 -->
          <div class="merge-actions" v-if="preview.canMerge">
            <button 
              @click="handleMerge" 
              class="btn btn-danger"
              :disabled="loading"
            >
              {{ loading ? '合并中...' : '确认合并' }}
            </button>
            <p class="warning-text">
              ⚠️ 警告：合并操作不可逆，请确认后再执行！
            </p>
          </div>
        </div>
      </div>
    </div>
  `,
};
