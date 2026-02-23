import { describe, expect, it } from 'vitest'
import {
  getPromoterStatusText,
  formatPromoterDetail,
  buildRewardRecordsPath,
  buildGenerateLinkPath
} from '../../src/views/brand/promoterDetail.logic.js'

describe('promoterDetail logic', () => {
  it('maps status to display text', () => {
    expect(getPromoterStatusText('active')).toBe('活跃')
    expect(getPromoterStatusText('inactive')).toBe('不活跃')
    expect(getPromoterStatusText('blocked')).toBe('已禁用')
    expect(getPromoterStatusText('unknown')).toBe('unknown')
    expect(getPromoterStatusText(null)).toBe('-')
    expect(getPromoterStatusText(undefined)).toBe('-')
  })

  it('formats promoter detail with defaults', () => {
    const result = formatPromoterDetail({
      id: 1,
      name: '测试推广员',
      status: 'active',
      totalOrders: 10
    })

    expect(result.displayStatus).toBe('活跃')
    expect(result.displayLevel).toBe('普通')
    expect(result.displayTotalOrders).toBe(10)
    expect(result.displayTotalRewards).toBe(0)
    expect(result.displayConversionRate).toBe(0)
    expect(result.displayCampaignCount).toBe(0)
  })

  it('formats promoter detail with all fields', () => {
    const result = formatPromoterDetail({
      id: 2,
      name: '高级推广员',
      status: 'active',
      level: 'VIP',
      totalOrders: 50,
      totalRewards: 5000,
      conversionRate: 25,
      campaignCount: 8
    })

    expect(result.displayStatus).toBe('活跃')
    expect(result.displayLevel).toBe('VIP')
    expect(result.displayTotalOrders).toBe(50)
    expect(result.displayTotalRewards).toBe(5000)
    expect(result.displayConversionRate).toBe(25)
    expect(result.displayCampaignCount).toBe(8)
  })

  it('returns empty object for null input', () => {
    expect(formatPromoterDetail(null)).toEqual({})
    expect(formatPromoterDetail(undefined)).toEqual({})
  })

  it('builds reward records path', () => {
    expect(buildRewardRecordsPath(123)).toBe('/brand/reward-records/123')
    expect(buildRewardRecordsPath('abc')).toBe('/brand/reward-records/abc')
    expect(buildRewardRecordsPath(null)).toBe('/brand/reward-records')
    expect(buildRewardRecordsPath(undefined)).toBe('/brand/reward-records')
    expect(buildRewardRecordsPath('')).toBe('/brand/reward-records')
  })

  it('builds generate link path', () => {
    expect(buildGenerateLinkPath(456)).toBe('/brand/promoters?generateLink=456')
    expect(buildGenerateLinkPath('xyz')).toBe('/brand/promoters?generateLink=xyz')
    expect(buildGenerateLinkPath(null)).toBe('/brand/promoters')
    expect(buildGenerateLinkPath(undefined)).toBe('/brand/promoters')
    expect(buildGenerateLinkPath('')).toBe('/brand/promoters')
  })
})
