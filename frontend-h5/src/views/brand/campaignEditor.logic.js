/**
 * CampaignEditor business logic
 */

export const STATUS_OPTIONS = [
  { value: 'active', label: '进行中' },
  { value: 'paused', label: '已暂停' },
  { value: 'ended', label: '已结束' }
]

export const PAYMENT_TYPE_OPTIONS = [
  { value: 'deposit', label: '订金支付' },
  { value: 'full', label: '全款支付' }
]

export const REQUIRE_PAYMENT_OPTIONS = [
  { value: 'true', label: '需要支付' },
  { value: 'false', label: '免费活动' }
]

export const FIELD_TYPE_MAP = {
  'text': '文本',
  'tel': '手机号',
  'email': '邮箱',
  'select': '下拉选择',
  'textarea': '多行文本',
  'number': '数字'
}

export const getFieldTypeText = (type) => FIELD_TYPE_MAP[type] || type

export const isEditMode = (campaignId) => !!campaignId

export const getPaymentAmountLabel = (paymentType) => {
  return paymentType === 'deposit' ? '订金金额' : '支付金额'
}

export const getPaymentAmountPlaceholder = (paymentType) => {
  return paymentType === 'deposit' ? '请输入订金金额' : '请输入支付金额'
}

export const getDefaultForm = () => ({
  brandId: 0,
  name: '',
  description: '',
  startTime: '',
  endTime: '',
  rewardRule: 0,
  status: 'active',
  requirePayment: 'false',
  paymentType: 'deposit',
  paymentAmount: 0,
  formFields: []
})

export const resolveCampaignBrandId = (storedBrandId, userInfoRaw) => {
  if (storedBrandId) {
    const parsed = Number.parseInt(storedBrandId, 10)
    if (!Number.isNaN(parsed) && parsed > 0) {
      return parsed
    }
  }

  if (!userInfoRaw) return 0

  try {
    const userInfo = JSON.parse(userInfoRaw)
    const firstBrandId = Array.isArray(userInfo.brandIds) ? Number(userInfo.brandIds[0]) : 0
    return Number.isFinite(firstBrandId) && firstBrandId > 0 ? firstBrandId : 0
  } catch {
    return 0
  }
}

export const validateCampaignForm = (form) => {
  const errors = []
  
  if (!form.name || form.name.trim().length < 2) {
    errors.push('活动名称至少需要2个字符')
  }
  
  if (!form.startTime) {
    errors.push('请选择开始时间')
  }
  
  if (!form.endTime) {
    errors.push('请选择结束时间')
  }
  
  if (form.startTime && form.endTime && new Date(form.startTime) >= new Date(form.endTime)) {
    errors.push('结束时间必须晚于开始时间')
  }
  
  if (!form.rewardRule || form.rewardRule <= 0) {
    errors.push('请输入有效的奖励金额')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

export const canShowQrcodePreview = (isEdit, requirePayment) => {
  return isEdit && (requirePayment === 'true' || requirePayment === true)
}

export const getSaveButtonText = (saving) => {
  return saving ? '保存中...' : '保存'
}

export const getTitleText = (isEdit) => {
  return isEdit ? '编辑活动' : '创建活动'
}
