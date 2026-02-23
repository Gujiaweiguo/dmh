<template>
  <div class="campaign-editor">
    <!-- 顶部导航栏 -->
    <van-nav-bar
      :title="isEditMode ? '编辑活动' : '创建活动'"
      left-text="返回"
      right-text="保存"
      left-arrow
      @click-left="goBack"
      @click-right="saveCampaign"
    />

    <!-- 表单内容 -->
    <van-form @submit="saveCampaign" class="campaign-form">
      <!-- 基本信息 -->
      <van-cell-group inset title="基本信息">
        <van-field
          v-model="form.name"
          name="name"
          label="活动名称"
          placeholder="请输入活动名称（2-100字符）"
          :rules="[{ required: true, message: '请输入活动名称' }]"
          maxlength="100"
          show-word-limit
        />
        
        <van-field
          v-model="form.description"
          name="description"
          label="活动描述"
          type="textarea"
          placeholder="请输入活动描述（最多500字符）"
          maxlength="500"
          rows="3"
          autosize
          show-word-limit
        />
        
        <van-field
          v-model="form.startTime"
          name="startTime"
          label="开始时间"
          placeholder="请选择开始时间"
          readonly
          @click="showStartTimePicker = true"
          :rules="[{ required: true, message: '请选择开始时间' }]"
        />
        
        <van-field
          v-model="form.endTime"
          name="endTime"
          label="结束时间"
          placeholder="请选择结束时间"
          readonly
          @click="showEndTimePicker = true"
          :rules="[{ required: true, message: '请选择结束时间' }]"
        />
        
        <van-field
          v-model="form.rewardRule"
          name="rewardRule"
          label="奖励金额"
          type="number"
          placeholder="请输入奖励金额"
          :rules="[{ required: true, message: '请输入奖励金额' }]"
        >
          <template #right-icon>
            <span class="currency-unit">元</span>
          </template>
        </van-field>
        
        <van-field
          v-if="isEditMode"
          v-model="form.status"
          name="status"
          label="活动状态"
          placeholder="请选择活动状态"
          readonly
          @click="showStatusPicker = true"
        />
      </van-cell-group>

      <!-- 动态表单字段 -->
      <van-cell-group inset title="表单字段配置">
        <van-cell
          v-for="(field, index) in form.formFields"
          :key="index"
          :title="field.label"
          :label="`类型: ${getFieldTypeText(field.type)} ${field.required ? '(必填)' : ''}`"
          is-link
          @click="editField(index)"
        >
          <template #right-icon>
            <van-icon name="edit" @click.stop="editField(index)" />
            <van-icon name="delete-o" @click.stop="removeField(index)" class="delete-icon" />
          </template>
        </van-cell>
        
        <van-cell
          v-if="form.formFields.length === 0"
          title="暂无表单字段"
          label="请添加表单字段"
        />
        
        <van-cell
          title="添加字段"
          is-link
          @click="showFieldTypeSheet = true"
        >
          <template #icon>
            <van-icon name="plus" />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 页面设计 -->
      <van-cell-group inset title="页面设计" v-if="isEditMode">
        <van-cell
          title="设计活动页面"
          is-link
          @click="goToPageDesigner"
        >
          <template #icon>
            <van-icon name="edit" />
          </template>
          <template #label>
            自定义活动展示页面的布局和内容
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 表单预览 -->
      <van-cell-group inset title="表单预览">
        <div class="form-preview">
          <div class="preview-header">
            <h3>{{ form.name || '活动名称' }}</h3>
            <p>{{ form.description || '活动描述' }}</p>
          </div>
          
          <van-form class="preview-form">
            <van-field
              v-for="field in form.formFields"
              :key="field.name"
              :label="field.label"
              :placeholder="field.placeholder"
              :required="field.required"
              readonly
            />
            
            <div v-if="form.formFields.length === 0" class="preview-empty">
              请先添加表单字段
            </div>
            
            <van-button
              type="primary"
              block
              disabled
              class="preview-submit"
            >
              立即参与（奖励 ¥{{ form.rewardRule || 0 }}）
            </van-button>
          </van-form>
        </div>
      </van-cell-group>
    </van-form>

    <!-- 字段类型选择 -->
    <van-action-sheet
      v-model:show="showFieldTypeSheet"
      :actions="fieldTypeActions"
      @select="onFieldTypeSelect"
      cancel-text="取消"
      title="选择字段类型"
    />

    <!-- 状态选择 -->
    <van-action-sheet
      v-model:show="showStatusPicker"
      :actions="statusActions"
      @select="onStatusSelect"
      cancel-text="取消"
      title="选择活动状态"
    />

    <!-- 字段编辑弹窗 -->
    <van-popup
      v-model:show="showFieldModal"
      position="bottom"
      :style="{ height: '70%' }"
      round
    >
      <div class="field-modal">
        <van-nav-bar
          :title="editingFieldIndex >= 0 ? '编辑字段' : '添加字段'"
          left-text="取消"
          right-text="保存"
          @click-left="closeFieldModal"
          @click-right="saveField"
        />
        
        <van-form class="field-form">
          <van-cell-group>
            <van-field
              v-model="fieldForm.type"
              name="type"
              label="字段类型"
              placeholder="请选择字段类型"
              readonly
              :disabled="editingFieldIndex >= 0"
              @click="showFieldTypeSheet = true"
            />
            
            <van-field
              v-model="fieldForm.label"
              name="label"
              label="字段标签"
              placeholder="如：姓名、联系方式"
              :rules="[{ required: true, message: '请输入字段标签' }]"
            />
            
            <van-field
              v-model="fieldForm.name"
              name="name"
              label="字段名称"
              placeholder="如：name、phone（英文）"
              :rules="[{ required: true, message: '请输入字段名称' }]"
            />
            
            <van-field
              v-if="fieldForm.type !== 'phone'"
              v-model="fieldForm.placeholder"
              name="placeholder"
              label="占位符"
              placeholder="输入框的提示文字"
            />
            
            <van-cell title="必填字段">
              <template #right-icon>
                <van-switch v-model="fieldForm.required" />
              </template>
            </van-cell>
          </van-cell-group>
          
          <!-- 下拉选择选项 -->
          <van-cell-group v-if="fieldForm.type === 'select'" title="选项列表">
            <van-field
              v-for="(option, index) in fieldForm.options"
              :key="index"
              v-model="fieldForm.options[index]"
              :label="`选项${index + 1}`"
              placeholder="选项内容"
            >
              <template #right-icon>
                <van-icon name="delete-o" @click="removeOption(index)" />
              </template>
            </van-field>
            
            <van-cell
              title="添加选项"
              is-link
              @click="addOption"
            >
              <template #icon>
                <van-icon name="plus" />
              </template>
            </van-cell>
          </van-cell-group>
        </van-form>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { campaignApi } from '../../services/brandApi.js'
import { resolveCampaignBrandId } from './campaignEditor.logic.js'

const router = useRouter()
const route = useRoute()

const form = reactive({
  brandId: 0,
  name: '',
  description: '',
  formFields: [],
  rewardRule: 0,
  startTime: '',
  endTime: '',
  status: 'active'
})

const fieldForm = reactive({
  type: 'text',
  name: '',
  label: '',
  required: false,
  placeholder: '',
  options: []
})

const saving = ref(false)
const loading = ref(false)
const showFieldModal = ref(false)
const showFieldTypeSheet = ref(false)
const showStatusPicker = ref(false)
const showStartTimePicker = ref(false)
const showEndTimePicker = ref(false)
const editingFieldIndex = ref(-1)

const isEditMode = computed(() => !!route.params.id)

// 字段类型选项
const fieldTypeActions = [
  { name: '文本框', value: 'text' },
  { name: '手机号', value: 'phone' },
  { name: '下拉选择', value: 'select' }
]

// 状态选项
const statusActions = [
  { name: '进行中', value: 'active' },
  { name: '已暂停', value: 'paused' },
  { name: '已结束', value: 'ended' }
]

const goBack = () => {
  router.back()
}

const goToPageDesigner = () => {
  router.push(`/brand/campaigns/${route.params.id}/page-design`)
}

const getFieldTypeText = (type) => {
  const typeMap = {
    text: '文本框',
    phone: '手机号',
    select: '下拉选择'
  }
  return typeMap[type] || type
}

const onFieldTypeSelect = (action) => {
  if (showFieldModal.value) {
    fieldForm.type = action.value
  } else {
    addField(action.value)
  }
  showFieldTypeSheet.value = false
}

const onStatusSelect = (action) => {
  form.status = action.value
  showStatusPicker.value = false
}

const addField = (type) => {
  resetFieldForm()
  fieldForm.type = type
  
  switch (type) {
    case 'text':
      fieldForm.label = '姓名'
      fieldForm.name = 'name'
      fieldForm.placeholder = '请输入姓名'
      break
    case 'phone':
      fieldForm.label = '手机号'
      fieldForm.name = 'phone'
      fieldForm.placeholder = '请输入手机号'
      fieldForm.required = true
      break
    case 'select':
      fieldForm.label = '意向课程'
      fieldForm.name = 'course'
      fieldForm.options = ['前端开发', '后端开发']
      break
  }
  
  showFieldModal.value = true
}

const editField = (index) => {
  const field = form.formFields[index]
  Object.assign(fieldForm, {
    type: field.type,
    name: field.name,
    label: field.label,
    required: field.required || false,
    placeholder: field.placeholder || '',
    options: field.options ? [...field.options] : []
  })
  editingFieldIndex.value = index
  showFieldModal.value = true
}

const removeField = (index) => {
  form.formFields.splice(index, 1)
  showToast({ type: 'success', message: '字段已删除' })
}

const resetFieldForm = () => {
  Object.assign(fieldForm, {
    type: 'text',
    name: '',
    label: '',
    required: false,
    placeholder: '',
    options: []
  })
  editingFieldIndex.value = -1
}

const closeFieldModal = () => {
  showFieldModal.value = false
  resetFieldForm()
}

const saveField = () => {
  if (!fieldForm.label.trim()) {
    showToast({ type: 'fail', message: '请输入字段标签' })
    return
  }

  if (!fieldForm.name.trim()) {
    showToast({ type: 'fail', message: '请输入字段名称' })
    return
  }

  if (fieldForm.type === 'select' && fieldForm.options.length === 0) {
    showToast({ type: 'fail', message: '下拉选择字段至少需要一个选项' })
    return
  }

  const existingIndex = form.formFields.findIndex((field, index) =>
    field.name === fieldForm.name && index !== editingFieldIndex.value
  )
  if (existingIndex >= 0) {
    showToast({ type: 'fail', message: '字段名称已存在，请使用其他名称' })
    return
  }

  const fieldData = {
    type: fieldForm.type,
    name: fieldForm.name,
    label: fieldForm.label,
    required: fieldForm.required
  }

  if (fieldForm.placeholder) {
    fieldData.placeholder = fieldForm.placeholder
  }

  if (fieldForm.type === 'select') {
    fieldData.options = fieldForm.options.filter(option => option.trim())
  }

  if (editingFieldIndex.value >= 0) {
    form.formFields[editingFieldIndex.value] = fieldData
    showToast({ type: 'success', message: '字段已更新' })
  } else {
    form.formFields.push(fieldData)
    showToast({ type: 'success', message: '字段已添加' })
  }

  closeFieldModal()
}

const addOption = () => {
  fieldForm.options.push('')
}

const removeOption = (index) => {
  fieldForm.options.splice(index, 1)
}

const saveCampaign = async () => {
  if (!form.brandId) {
    showToast({ type: 'fail', message: '未找到品牌信息，请重新登录' })
    return
  }

  if (!form.name.trim()) {
    showToast({ type: 'fail', message: '请输入活动名称' })
    return
  }

  if (form.name.length < 2 || form.name.length > 100) {
    showToast({ type: 'fail', message: '活动名称长度必须在2-100个字符之间' })
    return
  }

  if (form.rewardRule < 0.01 || form.rewardRule > 999.99) {
    showToast({ type: 'fail', message: '奖励金额必须在0.01-999.99元之间' })
    return
  }

  if (!form.startTime || !form.endTime) {
    showToast({ type: 'fail', message: '请设置活动时间' })
    return
  }

  if (new Date(form.startTime) >= new Date(form.endTime)) {
    showToast({ type: 'fail', message: '结束时间必须晚于开始时间' })
    return
  }

  if (form.formFields.length === 0) {
    showToast({ type: 'fail', message: '请至少添加一个表单字段' })
    return
  }

  saving.value = true
  showLoadingToast({ message: '保存中...', forbidClick: true })

  try {
    const campaignData = {
      ...form,
      startTime: formatDateTime(form.startTime),
      endTime: formatDateTime(form.endTime)
    }

    if (isEditMode.value) {
      await campaignApi.updateCampaign(route.params.id, campaignData)
      closeToast()
      showToast({ type: 'success', message: '更新成功' })
    } else {
      await campaignApi.createCampaign(campaignData)
      closeToast()
      showToast({ type: 'success', message: '创建成功' })
    }

    router.push('/brand/campaigns')
  } catch (error) {
    console.error('保存活动失败:', error)
    closeToast()
    showToast({ type: 'fail', message: `保存失败: ${error.message || '未知错误'}` })
  } finally {
    saving.value = false
  }
}

const loadCampaign = async () => {
  if (!isEditMode.value) return

  loading.value = true
  showLoadingToast({ message: '加载中...', forbidClick: true })

  try {
    const campaign = await campaignApi.getCampaign(route.params.id)

    Object.assign(form, {
      brandId: campaign.brandId,
      name: campaign.name,
      description: campaign.description,
      formFields: [...campaign.formFields],
      rewardRule: campaign.rewardRule,
      startTime: formatDateTimeLocal(campaign.startTime),
      endTime: formatDateTimeLocal(campaign.endTime),
      status: campaign.status
    })

    closeToast()
  } catch (error) {
    console.error('加载活动失败:', error)
    closeToast()
    showToast({ type: 'fail', message: '加载活动信息失败' })
    goBack()
  } finally {
    loading.value = false
  }
}

const formatDateTimeLocal = (dateString) => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

const formatDateTime = (dateTimeLocal) => {
  const date = new Date(dateTimeLocal)
  return date.toISOString().slice(0, 19).replace('T', ' ')
}

onMounted(() => {
  form.brandId = resolveCampaignBrandId(
    localStorage.getItem('brandId') || localStorage.getItem('dmh_current_brand_id'),
    localStorage.getItem('dmh_user_info'),
  )

  if (!form.brandId && !isEditMode.value) {
    showToast({ type: 'fail', message: '未找到品牌信息，请重新登录' })
  }

  loadCampaign()
})
</script>

<style scoped>
.campaign-editor {
  min-height: 100vh;
  background: #f7f8fa;
}

.campaign-form {
  padding-bottom: 20px;
}

.currency-unit {
  color: #969799;
  font-size: 14px;
}

.delete-icon {
  margin-left: 12px;
  color: #ee0a24;
}

.form-preview {
  padding: 16px;
  background: white;
  margin: 16px;
  border-radius: 8px;
}

.preview-header {
  text-align: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebedf0;
}

.preview-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #323233;
  margin: 0 0 8px 0;
}

.preview-header p {
  font-size: 14px;
  color: #646566;
  margin: 0;
}

.preview-form {
  max-width: 100%;
}

.preview-empty {
  text-align: center;
  padding: 32px;
  color: #c8c9cc;
  font-size: 14px;
}

.preview-submit {
  margin-top: 20px;
  opacity: 0.7;
}

.field-modal {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.field-form {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

@media (max-width: 768px) {
  .form-preview {
    margin: 12px;
  }
}
</style>
