export const getOrderStatusText = (status) => {
  const statusMap = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
  }
  return statusMap[status] || status
}

export const formatOrderDateTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export const filterAndSortOrders = (orders, currentStatus, dateRange) => {
  let filtered = orders

  if (currentStatus !== 'all') {
    filtered = filtered.filter((order) => order.status === currentStatus)
  }

  if (dateRange?.start) {
    filtered = filtered.filter((order) => new Date(order.createdAt) >= new Date(dateRange.start))
  }

  if (dateRange?.end) {
    filtered = filtered.filter((order) => new Date(order.createdAt) <= new Date(`${dateRange.end} 23:59:59`))
  }

  return [...filtered].sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
}

export const calculateOrderStats = (orders, today = new Date().toDateString()) => {
  return {
    total: orders.length,
    totalAmount: orders.reduce((sum, order) => sum + order.amount, 0),
    totalRewards: orders.reduce((sum, order) => sum + (order.rewardAmount || 0), 0),
    todayOrders: orders.filter((order) => new Date(order.createdAt).toDateString() === today).length,
  }
}

export const applyOrderStatus = (order, newStatus, rewardRate = 0.2) => {
  const next = { ...order, status: newStatus }
  if (newStatus === 'paid' && order.referrerId) {
    next.rewardAmount = order.amount * rewardRate
  }
  return next
}

export const buildExportOrderData = (order) => {
  return {
    订单号: order.id,
    活动名称: order.campaignName,
    用户手机: order.phone,
    订单金额: order.amount,
    订单状态: getOrderStatusText(order.status),
    创建时间: order.createdAt,
    ...order.formData,
  }
}

const escapeCsvCell = (value) => {
  const text = String(value ?? '')
  if (text.includes(',') || text.includes('"') || text.includes('\n')) {
    return `"${text.replace(/"/g, '""')}"`
  }
  return text
}

export const buildOrdersCsv = (orders) => {
  const rows = (orders || []).map(buildExportOrderData)
  if (rows.length === 0) return ''

  const headers = Object.keys(rows[0])
  const lines = [headers.map(escapeCsvCell).join(',')]

  rows.forEach((row) => {
    lines.push(headers.map((key) => escapeCsvCell(row[key])).join(','))
  })

  return lines.join('\n')
}

export const getOrderDetailItems = (order) => {
  return [
    { label: '订单号', value: order.id || '-' },
    { label: '活动名称', value: order.campaignName || '-' },
    { label: '手机号', value: order.phone || '-' },
    { label: '订单状态', value: getOrderStatusText(order.status) },
    { label: '订单金额', value: `¥${Number(order.amount || 0).toFixed(2)}` },
    { label: '奖励金额', value: `¥${Number(order.rewardAmount || 0).toFixed(2)}` },
    { label: '推荐人', value: order.referrerName || order.referrerId || '-' },
    { label: '创建时间', value: order.createdAt || '-' },
  ]
}
