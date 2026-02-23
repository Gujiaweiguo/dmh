import { describe, expect, it } from 'vitest'
import {
  buildPromoterLink,
  buildPromoterLinkForm,
  calculatePromoterStats,
  filterAndSortPromoters,
  formatPromoterTime,
  getPromoterStatusText,
} from '../../src/views/brand/promoters.logic.js'

describe('promoters logic', () => {
  const promoters = [
    {
      id: 1,
      name: '张推广',
      phone: '138****1234',
      status: 'active',
      level: 'VIP',
      totalRewards: 1200,
      conversionRate: 16,
      joinDate: '2026-02-12',
      todayOrders: 3,
    },
    {
      id: 2,
      name: '李推广',
      phone: '139****5678',
      status: 'inactive',
      level: '银牌',
      totalRewards: 400,
      conversionRate: 8,
      joinDate: '2025-12-01',
      todayOrders: 0,
    },
    {
      id: 3,
      name: '王推广',
      phone: '137****9999',
      status: 'active',
      level: '金牌',
      totalRewards: 1600,
      conversionRate: 12,
      joinDate: '2026-02-09',
      todayOrders: 1,
    },
  ]

  it('maps promoter status text', () => {
    expect(getPromoterStatusText('active')).toBe('活跃')
    expect(getPromoterStatusText('inactive')).toBe('不活跃')
    expect(getPromoterStatusText('blocked')).toBe('已封禁')
    expect(getPromoterStatusText('other')).toBe('other')
  })

  it('formats promoter activity time', () => {
    expect(formatPromoterTime('2026-02-13 10:00:00')).toContain('2')
  })

  it('filters by active and sorts by rewards desc', () => {
    const result = filterAndSortPromoters(promoters, 'active', '')
    expect(result.map((item) => item.id)).toEqual([3, 1])
  })

  it('filters top promoters and supports keyword search', () => {
    const top = filterAndSortPromoters(promoters, 'top', '')
    expect(top.map((item) => item.id)).toEqual([3, 1])

    const byKeyword = filterAndSortPromoters(promoters, 'all', '139****')
    expect(byKeyword).toHaveLength(1)
    expect(byKeyword[0].id).toBe(2)
  })

  it('filters new promoters by one week window', () => {
    const now = new Date('2026-02-13T00:00:00Z')
    const result = filterAndSortPromoters(promoters, 'new', '', now)
    expect(result.map((item) => item.id)).toEqual([3, 1])
  })

  it('calculates promoter stats with safe zero handling', () => {
    const stats = calculatePromoterStats(promoters)
    expect(stats).toEqual({
      active: 2,
      totalRewards: 3200,
      todayOrders: 4,
      conversionRate: 12,
    })

    const emptyStats = calculatePromoterStats([])
    expect(emptyStats.conversionRate).toBe(0)
  })

  it('builds link form and promo url', () => {
    const form = buildPromoterLinkForm({ id: 9, name: '赵推广' })
    expect(form).toEqual({ promoterId: 9, promoterName: '赵推广', campaignId: '' })

    expect(buildPromoterLink('https://dmh.test', 11, 9)).toBe('https://dmh.test/campaign/11?ref=9')
    expect(buildPromoterLink('https://dmh.test', '', 9)).toBe('')
  })
})

import {
  buildPhoneLink,
  buildSmsLink,
  getDefaultContactForm,
  validateContactForm,
  buildContactAction,
} from '../../src/views/brand/promoters.logic.js'

describe('contact functions', () => {
  describe('buildPhoneLink', () => {
    it('builds tel link from phone number', () => {
      expect(buildPhoneLink('13800138000')).toBe('tel:13800138000')
    })

    it('cleans special characters from phone', () => {
      expect(buildPhoneLink('138-0013-8000')).toBe('tel:13800138000')
      expect(buildPhoneLink('138****8000')).toBe('tel:1388000')
    })

    it('returns empty string for null/undefined', () => {
      expect(buildPhoneLink(null)).toBe('')
      expect(buildPhoneLink(undefined)).toBe('')
      expect(buildPhoneLink('')).toBe('')
    })
  })

  describe('buildSmsLink', () => {
    it('builds sms link without message', () => {
      expect(buildSmsLink('13800138000')).toBe('sms:13800138000')
    })

    it('builds sms link with message', () => {
      const link = buildSmsLink('13800138000', '测试短信')
      expect(link).toContain('sms:13800138000')
      expect(link).toContain('body=')
    })

    it('returns empty string for null phone', () => {
      expect(buildSmsLink(null)).toBe('')
    })
  })

  describe('getDefaultContactForm', () => {
    it('returns default form values', () => {
      const form = getDefaultContactForm()
      expect(form).toEqual({
        method: 'phone',
        message: '',
      })
    })
  })

  describe('validateContactForm', () => {
    it('returns error if promoter has no phone', () => {
      expect(validateContactForm({ method: 'phone' }, { name: 'Test' })).toBe('推广员没有联系电话')
      expect(validateContactForm({ method: 'phone' }, null)).toBe('推广员没有联系电话')
    })

    it('returns error for sms without message', () => {
      expect(
        validateContactForm({ method: 'sms', message: '' }, { phone: '13800000000' })
      ).toBe('请输入短信内容')
      expect(
        validateContactForm({ method: 'sms', message: '   ' }, { phone: '13800000000' })
      ).toBe('请输入短信内容')
    })

    it('returns empty string for valid phone contact', () => {
      expect(
        validateContactForm({ method: 'phone' }, { phone: '13800000000' })
      ).toBe('')
    })

    it('returns empty string for valid sms contact', () => {
      expect(
        validateContactForm({ method: 'sms', message: 'Hello' }, { phone: '13800000000' })
      ).toBe('')
    })
  })

  describe('buildContactAction', () => {
    it('returns null for missing phone', () => {
      expect(buildContactAction({ method: 'phone' }, { name: 'Test' })).toBeNull()
    })

    it('returns phone action', () => {
      const action = buildContactAction({ method: 'phone' }, { phone: '13800000000' })
      expect(action).toEqual({
        type: 'phone',
        link: 'tel:13800000000',
      })
    })

    it('returns sms action', () => {
      const action = buildContactAction({ method: 'sms', message: 'Hello' }, { phone: '13800000000' })
      expect(action.type).toBe('sms')
      expect(action.link).toContain('sms:13800000000')
    })
  })
})
