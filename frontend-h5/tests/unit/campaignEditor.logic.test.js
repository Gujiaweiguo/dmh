import { describe, it, expect } from 'vitest'
import {
  STATUS_OPTIONS,
  PAYMENT_TYPE_OPTIONS,
  REQUIRE_PAYMENT_OPTIONS,
  FIELD_TYPE_MAP,
  getFieldTypeText,
  isEditMode,
  getPaymentAmountLabel,
  getPaymentAmountPlaceholder,
  getDefaultForm,
  resolveCampaignBrandId,
  validateCampaignForm,
  canShowQrcodePreview,
  getSaveButtonText,
  getTitleText
} from '../../src/views/brand/campaignEditor.logic.js'

describe('campaignEditor.logic', () => {
  describe('Constants', () => {
    it('should have correct status options', () => {
      expect(STATUS_OPTIONS).toHaveLength(3)
      expect(STATUS_OPTIONS[0].value).toBe('active')
    })

    it('should have correct payment type options', () => {
      expect(PAYMENT_TYPE_OPTIONS).toHaveLength(2)
    })

    it('should have correct field type map', () => {
      expect(FIELD_TYPE_MAP['text']).toBe('文本')
      expect(FIELD_TYPE_MAP['tel']).toBe('手机号')
    })
  })

  describe('getFieldTypeText', () => {
    it('should return correct text for known type', () => {
      expect(getFieldTypeText('text')).toBe('文本')
      expect(getFieldTypeText('email')).toBe('邮箱')
    })

    it('should return type for unknown type', () => {
      expect(getFieldTypeText('unknown')).toBe('unknown')
    })
  })

  describe('isEditMode', () => {
    it('should return true for valid id', () => {
      expect(isEditMode('123')).toBe(true)
      expect(isEditMode(123)).toBe(true)
    })

    it('should return false for null or undefined', () => {
      expect(isEditMode(null)).toBe(false)
      expect(isEditMode(undefined)).toBe(false)
    })
  })

  describe('getPaymentAmountLabel', () => {
    it('should return deposit label for deposit type', () => {
      expect(getPaymentAmountLabel('deposit')).toBe('订金金额')
    })

    it('should return full label for full type', () => {
      expect(getPaymentAmountLabel('full')).toBe('支付金额')
    })
  })

  describe('getPaymentAmountPlaceholder', () => {
    it('should return deposit placeholder for deposit type', () => {
      expect(getPaymentAmountPlaceholder('deposit')).toBe('请输入订金金额')
    })

    it('should return full placeholder for full type', () => {
      expect(getPaymentAmountPlaceholder('full')).toBe('请输入支付金额')
    })
  })

  describe('getDefaultForm', () => {
    it('should return default form object', () => {
      const form = getDefaultForm()
      expect(form.brandId).toBe(0)
      expect(form.name).toBe('')
      expect(form.status).toBe('active')
      expect(form.formFields).toEqual([])
    })
  })

  describe('resolveCampaignBrandId', () => {
    it('should prefer brand id from storage', () => {
      expect(resolveCampaignBrandId('12', '')).toBe(12)
    })

    it('should fallback to first brand id from user info', () => {
      expect(resolveCampaignBrandId('', JSON.stringify({ brandIds: [33, 34] }))).toBe(33)
    })

    it('should return 0 when both sources invalid', () => {
      expect(resolveCampaignBrandId('', 'invalid-json')).toBe(0)
      expect(resolveCampaignBrandId('', JSON.stringify({ brandIds: [] }))).toBe(0)
    })
  })

  describe('validateCampaignForm', () => {
    it('should return errors for empty form', () => {
      const result = validateCampaignForm(getDefaultForm())
      expect(result.valid).toBe(false)
      expect(result.errors.length).toBeGreaterThan(0)
    })

    it('should return error for short name', () => {
      const result = validateCampaignForm({ ...getDefaultForm(), name: 'a' })
      expect(result.errors).toContain('活动名称至少需要2个字符')
    })

    it('should return error for invalid time range', () => {
      const result = validateCampaignForm({
        ...getDefaultForm(),
        name: '测试活动',
        startTime: '2024-12-31',
        endTime: '2024-01-01',
        rewardRule: 10
      })
      expect(result.errors).toContain('结束时间必须晚于开始时间')
    })

    it('should return valid for correct form', () => {
      const result = validateCampaignForm({
        ...getDefaultForm(),
        name: '测试活动',
        startTime: '2024-01-01',
        endTime: '2024-12-31',
        rewardRule: 10
      })
      expect(result.valid).toBe(true)
    })
  })

  describe('canShowQrcodePreview', () => {
    it('should return true for edit mode with payment', () => {
      expect(canShowQrcodePreview(true, 'true')).toBe(true)
    })

    it('should return false for new campaign', () => {
      expect(canShowQrcodePreview(false, 'true')).toBe(false)
    })

    it('should return false for free campaign', () => {
      expect(canShowQrcodePreview(true, 'false')).toBe(false)
    })
  })

  describe('getSaveButtonText', () => {
    it('should return saving text when saving', () => {
      expect(getSaveButtonText(true)).toBe('保存中...')
    })

    it('should return normal text when not saving', () => {
      expect(getSaveButtonText(false)).toBe('保存')
    })
  })

  describe('getTitleText', () => {
    it('should return edit text for edit mode', () => {
      expect(getTitleText(true)).toBe('编辑活动')
    })

    it('should return create text for new mode', () => {
      expect(getTitleText(false)).toBe('创建活动')
    })
  })
})
