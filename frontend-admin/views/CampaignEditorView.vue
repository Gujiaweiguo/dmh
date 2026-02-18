<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { Save, ArrowLeft, Plus, Trash2, Settings } from 'lucide-vue-next';
import { campaignApi, type CreateCampaignRequest, type UpdateCampaignRequest } from '../services/campaignApi';
import { posterApi, type PosterTemplate } from '../services/posterApi';
import type { Campaign } from '../types';

const props = defineProps<{
  campaignId?: number;
}>();

const emit = defineEmits(['back', 'saved']);

const form = ref<CreateCampaignRequest & { 
  status?: 'active' | 'paused' | 'ended';
  paymentConfig?: string;
  posterTemplateId?: number;
  enableDistribution?: boolean;
  distributionLevel?: number;
  distributionRewards?: string;
}>({
  name: '',
  description: '',
  formFields: [],
  rewardRule: 0,
  startTime: '',
  endTime: '',
  status: 'active',
  paymentConfig: '',
  posterTemplateId: 1,
  enableDistribution: false,
  distributionLevel: 1,
  distributionRewards: ''
});

const posterTemplates = ref<PosterTemplate[]>([]);
const posterTemplatesLoading = ref(false);

const newField = ref('');
const loading = ref(false);
const saving = ref(false);

const showPaymentConfigDialog = ref(false);

const isEditMode = computed(() => !!props.campaignId);

const paymentConfigForm = ref({
  paymentType: 'order_amount' as 'order_amount' | 'full_amount',
  merchantId: '',
  appId: '',
  apiKey: '',
  minAmount: undefined as number | undefined,
  maxAmount: undefined as number | undefined
});

const distributionRewardsForm = ref({
  level1: 10,
  level2: 5,
  level3: 3
});

const parseDistributionRewards = (value?: string) => {
  if (!value) return;
  try {
    const parsed = JSON.parse(value);
    if (parsed && typeof parsed === 'object') {
      distributionRewardsForm.value = {
        level1: Number(parsed.level1 ?? distributionRewardsForm.value.level1),
        level2: Number(parsed.level2 ?? distributionRewardsForm.value.level2),
        level3: Number(parsed.level3 ?? distributionRewardsForm.value.level3)
      };
    }
  } catch (error) {
    console.warn('Failed to parse distribution rewards:', error);
  }
};

const loadCampaign = async () => {
  if (!props.campaignId) return;

  loading.value = true;
  try {
    const campaign = await campaignApi.getCampaign(props.campaignId);
    form.value = {
      name: campaign.name,
      description: campaign.description,
      formFields: [...campaign.formFields],
      rewardRule: campaign.rewardRule,
      startTime: campaign.startTime,
      endTime: campaign.endTime,
      status: campaign.status,
      paymentConfig: campaign.paymentConfig || '',
      posterTemplateId: campaign.posterTemplateId || 1,
      enableDistribution: campaign.enableDistribution ?? false,
      distributionLevel: campaign.distributionLevel || 1,
      distributionRewards: campaign.distributionRewards || ''
    };
    parseDistributionRewards(campaign.distributionRewards);
  } catch (error) {
    console.error('Failed to load campaign:', error);
    alert('加载活动失败');
  } finally {
    loading.value = false;
  }
};

const loadPosterTemplates = async () => {
  posterTemplatesLoading.value = true;
  try {
    const response = await posterApi.getPosterTemplates();
    posterTemplates.value = response.templates.filter((t: PosterTemplate) => t.status === 'active');
  } catch (error) {
    console.error('Failed to load poster templates:', error);
    alert('加载海报模板失败');
  } finally {
    posterTemplatesLoading.value = false;
  }
};

const addFormField = () => {
  if (!newField.value.trim()) return;
  form.value.formFields.push(newField.value.trim());
  newField.value = '';
};

const removeFormField = (index: number) => {
  form.value.formFields.splice(index, 1);
};

const validateForm = (): boolean => {
  if (!form.value.name.trim()) {
    alert('请输入活动名称');
    return false;
  }
  if (form.value.name.length < 2 || form.value.name.length > 100) {
    alert('活动名称长度必须在2-100个字符之间');
    return false;
  }
  if (form.value.rewardRule < 0) {
    alert('奖励金额不能为负数');
    return false;
  }
  if (!form.value.startTime || !form.value.endTime) {
    alert('请选择活动时间范围');
    return false;
  }
  if (new Date(form.value.startTime) > new Date(form.value.endTime)) {
    alert('开始时间不能晚于结束时间');
    return false;
  }

  if (form.value.enableDistribution) {
    const level = Number(form.value.distributionLevel || 1);
    if (level < 1 || level > 3) {
      alert('分销层级必须是 1 到 3');
      return false;
    }
    const rewards = distributionRewardsForm.value;
    const requiredLevels = ['level1', 'level2', 'level3'].slice(0, level);
    for (const key of requiredLevels) {
      const value = Number((rewards as any)[key]);
      if (!Number.isFinite(value) || value < 0 || value > 100) {
        alert('分销奖励比例必须是 0-100 的数字');
        return false;
      }
    }
  }
  return true;
};

const buildDistributionRewards = () => {
  const level = Number(form.value.distributionLevel || 1);
  const rewards = distributionRewardsForm.value;
  const payload: Record<string, number> = {
    level1: Number(rewards.level1) || 0
  };
  if (level >= 2) {
    payload.level2 = Number(rewards.level2) || 0;
  }
  if (level >= 3) {
    payload.level3 = Number(rewards.level3) || 0;
  }
  return JSON.stringify(payload);
};

const handleSave = async () => {
  if (!validateForm()) return;

  saving.value = true;
  try {
    form.value.distributionLevel = Number(form.value.distributionLevel || 1);
    if (form.value.enableDistribution) {
      form.value.distributionRewards = buildDistributionRewards();
    } else {
      form.value.distributionRewards = '';
    }

    if (isEditMode.value && props.campaignId) {
      await campaignApi.updateCampaign(props.campaignId, form.value as UpdateCampaignRequest);
      alert('更新成功');
    } else {
      await campaignApi.createCampaign(form.value);
      alert('创建成功');
    }
    emit('saved');
  } catch (error) {
    console.error('Failed to save campaign:', error);
    const errorMessage = error instanceof Error ? error.message : '保存失败';
    alert(`保存失败: ${errorMessage}`);
  } finally {
    saving.value = false;
  }
};

const handleSavePaymentConfig = () => {
  form.value.paymentConfig = JSON.stringify(paymentConfigForm.value);
  showPaymentConfigDialog.value = false;
  alert('支付配置已保存');
};

onMounted(() => {
  loadCampaign();
  loadPosterTemplates();
});
</script>

<template>
  <div class="campaign-editor-view">
    <div class="header">
      <button class="btn-back" @click="emit('back')">
        <ArrowLeft :size="20" />
        <span>返回</span>
      </button>
      <h1 class="title">{{ isEditMode ? '编辑活动' : '创建活动' }}</h1>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="form-container">
      <div class="form-section">
        <h2 class="section-title">基本信息</h2>
        
        <div class="form-group">
          <label>活动名称 *</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="请输入活动名称（2-100字符）"
            maxlength="100"
          />
        </div>

        <div class="form-group">
          <label>活动描述</label>
          <textarea
            v-model="form.description"
            placeholder="请输入活动描述"
            rows="4"
          ></textarea>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>开始时间 *</label>
            <input v-model="form.startTime" type="datetime-local" />
          </div>
          <div class="form-group">
            <label>结束时间 *</label>
            <input v-model="form.endTime" type="datetime-local" />
          </div>
        </div>

        <div class="form-group">
          <label>奖励金额 *</label>
          <input
            v-model.number="form.rewardRule"
            type="number"
            step="0.01"
            min="0"
            placeholder="请输入奖励金额"
          />
        </div>

        <div class="form-group" v-if="isEditMode">
          <label>活动状态</label>
          <select v-model="form.status">
            <option value="active">进行中</option>
            <option value="paused">已暂停</option>
            <option value="ended">已结束</option>
          </select>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>海报模板</label>
            <select v-model="form.posterTemplateId" :disabled="posterTemplatesLoading">
              <option value="">暂不选择</option>
              <option v-for="template in posterTemplates" :key="template.id" :value="template.id">
                {{ template.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>支付配置</label>
            <button class="btn-config" @click="showPaymentConfigDialog = true">
          <Settings :size="20" />
              <span>{{ form.paymentConfig ? '已配置' : '未配置' }}</span>
            </button>
          </div>
        </div>
      </div>

      <div class="form-section">
        <h2 class="section-title">分销规则</h2>

        <div class="form-group form-inline">
          <label>启用分销</label>
          <input type="checkbox" v-model="form.enableDistribution" />
          <span class="form-tip">开启后可设置多级分销奖励</span>
        </div>

        <div class="form-row" v-if="form.enableDistribution">
          <div class="form-group">
            <label>分销层级</label>
            <select v-model.number="form.distributionLevel">
              <option :value="1">一级</option>
              <option :value="2">二级</option>
              <option :value="3">三级</option>
            </select>
          </div>
          <div class="form-group">
            <label>一级奖励比例 (%)</label>
            <input v-model.number="distributionRewardsForm.level1" type="number" min="0" max="100" />
          </div>
        </div>

        <div class="form-row" v-if="form.enableDistribution && form.distributionLevel >= 2">
          <div class="form-group">
            <label>二级奖励比例 (%)</label>
            <input v-model.number="distributionRewardsForm.level2" type="number" min="0" max="100" />
          </div>
          <div class="form-group">
            <label>三级奖励比例 (%)</label>
            <input v-model.number="distributionRewardsForm.level3" type="number" min="0" max="100" :disabled="form.distributionLevel < 3" />
          </div>
        </div>
      </div>

      <div class="form-section">
        <h2 class="section-title">动态表单字段</h2>
        
        <div class="field-list">
          <div v-for="(field, index) in form.formFields" :key="index" class="field-item">
            <span class="field-name">{{ field }}</span>
            <button class="btn-remove" @click="removeFormField(index)">
              <Trash2 :size="16" />
            </button>
          </div>
          <div v-if="form.formFields.length === 0" class="empty-state">
            暂无表单字段，请添加
          </div>
        </div>

        <div class="add-field">
          <input
            v-model="newField"
            type="text"
            placeholder="输入字段名称（如：姓名、手机号）"
            @keyup.enter="addFormField"
          />
          <button class="btn-add" @click="addFormField">
            <Plus :size="18" />
            <span>添加</span>
          </button>
        </div>
      </div>

      <div class="form-actions">
        <button class="btn-cancel" @click="emit('back')">取消</button>
        <button class="btn-save" @click="handleSave" :disabled="saving">
          <Save :size="20" />
          <span>{{ saving ? '保存中...' : '保存' }}</span>
        </button>
      </div>
    </div>

    <div v-if="showPaymentConfigDialog" class="dialog-overlay">
      <div class="dialog-content">
        <div class="dialog-header">
          <h3>支付配置</h3>
          <button class="btn-close" @click="showPaymentConfigDialog = false">
            <ArrowLeft :size="20" />
          </button>
        </div>

        <div class="dialog-body">
          <div class="form-group">
            <label>支付方式</label>
            <select v-model="paymentConfigForm.paymentType">
              <option value="order_amount">按订单金额支付</option>
              <option value="full_amount">固定全款</option>
            </select>
          </div>

          <div class="form-group" v-if="paymentConfigForm.paymentType === 'order_amount'">
            <label>最低金额（元）</label>
            <input v-model.number="paymentConfigForm.minAmount" type="number" step="0.01" min="0" />
          </div>

          <div class="form-group" v-if="paymentConfigForm.paymentType === 'order_amount'">
            <label>最高金额（元）</label>
            <input v-model.number="paymentConfigForm.maxAmount" type="number" step="0.01" min="0" />
          </div>

          <div class="form-group">
            <label>商户号</label>
            <input v-model="paymentConfigForm.merchantId" type="text" placeholder="请输入微信商户号" />
          </div>

          <div class="form-group">
            <label>应用ID</label>
            <input v-model="paymentConfigForm.appId" type="text" placeholder="请输入微信应用ID" />
          </div>

          <div class="form-group">
            <label>API密钥</label>
            <input v-model="paymentConfigForm.apiKey" type="password" placeholder="请输入微信API密钥" />
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn-cancel" @click="showPaymentConfigDialog = false">取消</button>
          <button class="btn-save" @click="handleSavePaymentConfig">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.campaign-editor-view {
  padding: 24px;
  max-width: 1000px;
  margin: 0 auto;
}

.header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.btn-back:hover {
  background: #f9fafb;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #6b7280;
  font-size: 16px;
}

.form-container {
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.form-section {
  margin-bottom: 32px;
  padding: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.form-inline {
  display: flex;
  align-items: center;
  gap: 12px;
}

.form-inline label {
  margin-bottom: 0;
}

.form-tip {
  font-size: 12px;
  color: #6b7280;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.btn-config {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  width: 100%;
}

.btn-config:hover {
  background: #f9fafb;
}

.field-list {
  background: #f9fafb;
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 16px;
  min-height: 100px;
}

.field-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  margin-bottom: 8px;
}

.field-name {
  font-size: 14px;
  color: #374151;
}

.btn-remove {
  padding: 6px;
  border: none;
  background: transparent;
  color: #ef4444;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn-remove:hover {
  background: #fee2e2;
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: #9ca3af;
  font-size: 14px;
}

.add-field {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.add-field input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
}

.btn-add {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-add:hover {
  background: #2563eb;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 32px;
  border-top: 1px solid #e5e7eb;
}

.btn-cancel {
  padding: 10px 24px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-cancel:hover {
  background: #f9fafb;
}

.btn-save {
  display: flex;
  align-items:.ts-center;
  gap: 8px;
  padding: 10px 24px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-save:hover:not(:disabled) {
  background: #2563eb;
}

.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog-content {
  background: white;
  border-radius: 12px;
  padding: 24px;
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.dialog-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
}

.btn-close {
  padding: 8px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: #6b7280;
}

.dialog-body {
  margin-bottom: 24px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.dialog-footer .btn-cancel,
.dialog-footer .btn-save {
  padding: 8px 20px;
}
</style>
