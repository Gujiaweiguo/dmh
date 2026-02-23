import { describe, expect, it } from 'vitest'
import {
  getRewardTypeText,
  getRewardStatusText,
  calculateRewardSummary,
  filterRewardRecords,
  formatRewardAmount,
  getAmountClass
} from '../../src/views/brand/rewardRecords.logic.js'

describe('rewardRecords logic', () => {
  const sampleRecords = [
    { id: 1, type: 'commission', status: 'paid', amount: 100, description: '订单佣金' },
    { id: 2, type: 'bonus', status: 'pending', amount: 50, description: '活动奖金' },
    { id: 3, type: 'commission', status: 'pending', amount: 80, description: '推广佣金' },
    { id: 4, type: 'penalty', status: 'paid', amount: -20, description: '违规扣款' }
  ]

  it('maps reward type to display text', () => {
    expect(getRewardTypeText('commission')).toBe('佣金')
    expect(getRewardTypeText('bonus')).toBe('奖金')
    expect(getRewardTypeText('penalty')).toBe('扣款')
    expect(getRewardTypeText('other')).toBe('other')
  })

  it('maps reward status to display text', () => {
    expect(getRewardStatusText('pending')).toBe('待发放')
    expect(getRewardStatusText('paid')).toBe('已发放')
    expect(getRewardStatusText('cancelled')).toBe('已取消')
    expect(getRewardStatusText('other')).toBe('other')
  })

  it('calculates reward summary correctly', () => {
    const summary = calculateRewardSummary(sampleRecords)
    expect(summary.totalAmount).toBe('210.00')
    expect(summary.pendingAmount).toBe('130.00')
    expect(summary.totalCount).toBe(4)
  })

  it('calculates empty summary for empty array', () => {
    const summary = calculateRewardSummary([])
    expect(summary).toEqual({
      totalAmount: '0.00',
      pendingAmount: '0.00',
      totalCount: 0
    })
  })

  it('handles null and undefined input', () => {
    expect(calculateRewardSummary(null)).toEqual({
      totalAmount: '0.00',
      pendingAmount: '0.00',
      totalCount: 0
    })
    expect(calculateRewardSummary(undefined)).toEqual({
      totalAmount: '0.00',
      pendingAmount: '0.00',
      totalCount: 0
    })
  })

  it('filters by type', () => {
    const result = filterRewardRecords(sampleRecords, 'commission', 'all')
    expect(result).toHaveLength(2)
    expect(result.map(r => r.id)).toEqual([1, 3])
  })

  it('filters by status', () => {
    const result = filterRewardRecords(sampleRecords, 'all', 'pending')
    expect(result).toHaveLength(2)
    expect(result.map(r => r.id)).toEqual([2, 3])
  })

  it('filters by both type and status', () => {
    const result = filterRewardRecords(sampleRecords, 'commission', 'pending')
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe(3)
  })

  it('returns all for all filters', () => {
    const result = filterRewardRecords(sampleRecords, 'all', 'all')
    expect(result).toHaveLength(4)
  })

  it('handles empty array input', () => {
    expect(filterRewardRecords([], 'all', 'all')).toEqual([])
    expect(filterRewardRecords(null, 'all', 'all')).toEqual([])
  })

  it('formats reward amount with sign', () => {
    expect(formatRewardAmount(100)).toBe('+¥100')
    expect(formatRewardAmount(0)).toBe('+¥0')
    expect(formatRewardAmount(-50)).toBe('¥-50')
    expect(formatRewardAmount(null)).toBe('¥0')
    expect(formatRewardAmount(undefined)).toBe('¥0')
  })

  it('determines amount class', () => {
    expect(getAmountClass(100)).toBe('positive')
    expect(getAmountClass(0)).toBe('positive')
    expect(getAmountClass(-50)).toBe('negative')
  })
})
