export const getPromoterStatusText = (status) => {
  const statusMap = {
    active: '活跃',
    inactive: '不活跃',
    blocked: '已封禁',
  }
  return statusMap[status] || status
}

export const formatPromoterTime = (timeString) => {
  const date = new Date(timeString)
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export const filterAndSortPromoters = (promoters, currentFilter, searchKeyword, now = new Date()) => {
  let filtered = promoters

  if (currentFilter !== 'all') {
    switch (currentFilter) {
      case 'active':
        filtered = filtered.filter((promoter) => promoter.status === 'active')
        break
      case 'top':
        filtered = filtered.filter((promoter) => promoter.level === 'VIP' || promoter.totalRewards > 1000)
        break
      case 'new': {
        const oneWeekAgo = new Date(now)
        oneWeekAgo.setDate(oneWeekAgo.getDate() - 7)
        filtered = filtered.filter((promoter) => new Date(promoter.joinDate) > oneWeekAgo)
        break
      }
    }
  }

  if (searchKeyword) {
    const keyword = searchKeyword.toLowerCase()
    filtered = filtered.filter(
      (promoter) => promoter.name.toLowerCase().includes(keyword) || promoter.phone.includes(keyword),
    )
  }

  return [...filtered].sort((a, b) => b.totalRewards - a.totalRewards)
}

export const calculatePromoterStats = (promoters) => {
  const active = promoters.filter((promoter) => promoter.status === 'active').length
  const totalRewards = promoters.reduce((sum, promoter) => sum + promoter.totalRewards, 0)
  const todayOrders = promoters.reduce((sum, promoter) => sum + (promoter.todayOrders || 0), 0)
  const conversionRate =
    promoters.length === 0
      ? 0
      : Math.round(promoters.reduce((sum, promoter) => sum + promoter.conversionRate, 0) / promoters.length)

  return {
    active,
    totalRewards,
    todayOrders,
    conversionRate,
  }
}

export const buildPromoterLinkForm = (promoter) => ({
  promoterId: promoter.id,
  promoterName: promoter.name,
  campaignId: '',
})

export const buildPromoterLink = (baseUrl, campaignId, promoterId) => {
  if (!campaignId || !promoterId) return ''
  return `${baseUrl}/campaign/${campaignId}?ref=${promoterId}`
}

export const buildPhoneLink = (phone) => {
  if (!phone) return ''
  const cleanPhone = phone.replace(/[^\d+]/g, '')
  return `tel:${cleanPhone}`
}

export const buildSmsLink = (phone, message = '') => {
  if (!phone) return ''
  const cleanPhone = phone.replace(/[^\d+]/g, '')
  if (message) {
    return `sms:${cleanPhone}?body=${encodeURIComponent(message)}`
  }
  return `sms:${cleanPhone}`
}

export const getDefaultContactForm = () => ({
  method: 'phone',
  message: '',
})

export const validateContactForm = (form, promoter) => {
  if (!promoter?.phone) {
    return '推广员没有联系电话'
  }
  if (form?.method === 'sms' && !form?.message?.trim()) {
    return '请输入短信内容'
  }
  return ''
}

export const buildContactAction = (form, promoter) => {
  if (!promoter?.phone || !form) return null
  
  if (form.method === 'phone') {
    return {
      type: 'phone',
      link: buildPhoneLink(promoter.phone),
    }
  }
  
  return {
    type: 'sms',
    link: buildSmsLink(promoter.phone, form.message),
  }
}
