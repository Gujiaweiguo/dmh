import { describe, expect, it } from 'vitest'
import {
  applyOrderStatus,
  buildOrdersCsv,
  buildExportOrderData,
  calculateOrderStats,
  filterAndSortOrders,
  formatOrderDateTime,
  getOrderDetailItems,
  getOrderStatusText,
} from '../../src/views/brand/orders.logic.js'

describe('orders logic', () => {
  const orders = [
    { id: 1, status: 'pending', amount: 10, rewardAmount: 0, referrerId: null, createdAt: '2026-02-13 10:00:00', campaignName: 'A', phone: '138', formData: { 姓名: '甲' } },
    { id: 2, status: 'paid', amount: 20, rewardAmount: 4, referrerId: 8, createdAt: '2026-02-12 10:00:00', campaignName: 'B', phone: '139', formData: {} },
    { id: 3, status: 'cancelled', amount: 30, rewardAmount: 0, referrerId: 9, createdAt: '2026-02-14 09:00:00', campaignName: 'C', phone: '137', formData: {} },
  ]

  it('maps order status text', () => {
    expect(getOrderStatusText('pending')).toBe('待支付')
    expect(getOrderStatusText('paid')).toBe('已支付')
    expect(getOrderStatusText('cancelled')).toBe('已取消')
    expect(getOrderStatusText('custom')).toBe('custom')
  })

  it('formats order datetime', () => {
    expect(formatOrderDateTime('2026-02-13 10:00:00')).toContain('2')
  })

  it('filters by status, date and sorts desc', () => {
    const byStatus = filterAndSortOrders(orders, 'paid', { start: '', end: '' })
    expect(byStatus).toHaveLength(1)
    expect(byStatus[0].id).toBe(2)

    const byDate = filterAndSortOrders(orders, 'all', { start: '2026-02-13', end: '2026-02-14' })
    expect(byDate.map((item) => item.id)).toEqual([3, 1])
  })

  it('calculates order stats', () => {
    const stats = calculateOrderStats(orders, 'Fri Feb 13 2026')
    expect(stats).toEqual({
      total: 3,
      totalAmount: 60,
      totalRewards: 4,
      todayOrders: 1,
    })
  })

  it('applies paid status and reward for referrer', () => {
    const next = applyOrderStatus({ id: 9, amount: 100, referrerId: 20, status: 'pending', rewardAmount: 0 }, 'paid')
    expect(next.status).toBe('paid')
    expect(next.rewardAmount).toBe(20)
  })

  it('builds export order data with mapped status', () => {
    const row = buildExportOrderData({
      id: 1,
      campaignName: '活动A',
      phone: '138****0000',
      amount: 88,
      status: 'paid',
      createdAt: '2026-02-13 10:00:00',
      formData: { 姓名: '李雷' },
    })

    expect(row).toMatchObject({
      订单号: 1,
      活动名称: '活动A',
      用户手机: '138****0000',
      订单金额: 88,
      订单状态: '已支付',
      姓名: '李雷',
    })
  })

  it('builds csv content for multiple orders', () => {
    const csv = buildOrdersCsv(orders)
    expect(csv).toContain('订单号')
    expect(csv).toContain('活动名称')
    expect(csv).toContain('A')
    expect(csv.split('\n').length).toBeGreaterThan(1)
  })

  it('returns empty csv for empty list', () => {
    expect(buildOrdersCsv([])).toBe('')
  })


  it('returns empty csv for null/undefined list', () => {
    expect(buildOrdersCsv(null)).toBe('')
    expect(buildOrdersCsv(undefined)).toBe('')
  })

  it('escapes csv cells with special characters', () => {
    const ordersWithSpecialChars = [
      {
        id: 1,
        status: 'paid',
        amount: 100,
        rewardAmount: 0,
        referrerId: null,
        createdAt: '2026-02-13',
        campaignName: '活动,A', // comma
        phone: '138',
        formData: { name: '张"三' }, // quote
      },
      {
        id: 2,
        status: 'paid',
        amount: 100,
        rewardAmount: 0,
        referrerId: null,
        createdAt: '2026-02-13',
        campaignName: '活动\nB', // newline
        phone: '139',
        formData: {},
      },
    ]
    const csv = buildOrdersCsv(ordersWithSpecialChars)
    // Cells with commas, quotes, or newlines should be quoted
    expect(csv).toContain('"活动,A"')
    expect(csv).toContain('"张""三"') // escaped quote
  })

  it('builds order detail items', () => {
    const items = getOrderDetailItems({
      id: 1,
      campaignName: '活动A',
      phone: '13800000000',
      status: 'paid',
      amount: 100,
      rewardAmount: 20,
      referrerName: '张三',
      createdAt: '2026-02-13 10:00:00',
    })
    expect(items.find((item) => item.label === '订单号')?.value).toBe(1)
    expect(items.find((item) => item.label === '订单状态')?.value).toBe('已支付')
	expect(items.find((item) => item.label === '订单状态')?.value).toBe('已支付')
  })

  it('builds order detail items with missing fields', () => {
    const items = getOrderDetailItems({
      id: null,
      campaignName: null,
      phone: null,
      status: 'pending',
      amount: null,
      rewardAmount: null,
      referrerName: null,
      referrerId: null,
      createdAt: null,
    })
    expect(items.find((item) => item.label === '订单号')?.value).toBe('-')
    expect(items.find((item) => item.label === '活动名称')?.value).toBe('-')
    expect(items.find((item) => item.label === '手机号')?.value).toBe('-')
    expect(items.find((item) => item.label === '订单金额')?.value).toBe('¥0.00')
    expect(items.find((item) => item.label === '奖励金额')?.value).toBe('¥0.00')
    expect(items.find((item) => item.label === '推荐人')?.value).toBe('-')
    expect(items.find((item) => item.label === '创建时间')?.value).toBe('-')
  })

  it('builds order detail items with referrerId only', () => {
    const items = getOrderDetailItems({
      id: 1,
      campaignName: '活动A',
      phone: '13800000000',
      status: 'paid',
      amount: 100,
      rewardAmount: 20,
      referrerName: null,
      referrerId: 123,
      createdAt: '2026-02-13',
    })
    expect(items.find((item) => item.label === '推荐人')?.value).toBe(123)
  })
})
