export const getMaterialCategories = () => [
  { value: 'all', label: '全部' },
  { value: 'image', label: '图片' },
  { value: 'text', label: '文案' },
  { value: 'video', label: '视频' },
]

export const getDefaultUploadForm = () => ({
  name: '',
  description: '',
  category: 'image',
  file: null,
})

export const getDefaultAITextForm = () => ({
  topic: '',
  style: 'professional',
  length: 'medium',
})

export const filterMaterialsByCategory = (materials, category) => {
  if (category === 'all') {
    return materials
  }
  return materials.filter((material) => material.type === category)
}

export const getMaterialTypeText = (type) => {
  const typeMap = {
    image: '图片',
    text: '文案',
    video: '视频',
  }
  return typeMap[type] || type
}

export const formatMaterialDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

export const validateUploadInput = (uploadForm) => {
  if (!uploadForm.file || !uploadForm.name) {
    return '请选择文件并填写素材名称'
  }
  return ''
}

export const validateAITextInput = (topic) => {
  if (!topic) {
    return '请输入文案主题'
  }
  return ''
}

export const createUploadedMaterial = ({ id, name, description, category, url, createdAt }) => ({
  id,
  name,
  description,
  type: category,
  url,
  createdAt,
})

export const buildAIGeneratedText = (topic) => {
  return `🎉 ${topic}火热进行中！限时优惠，机不可失！立即参与，享受超值福利，让您的生活更加精彩！赶快行动吧，名额有限，先到先得！`
}

export const createAIGeneratedMaterial = ({ id, topic, createdAt }) => ({
  id,
  name: `AI生成-${topic}`,
  description: 'AI智能生成的营销文案',
  type: 'text',
  content: buildAIGeneratedText(topic),
  createdAt,
})

export const getMockMaterials = () => [
  {
    id: 1,
    name: '春节促销海报',
    description: '2026年春节特惠活动主视觉海报',
    type: 'image',
    url: 'https://images.unsplash.com/photo-1607344645866-009c7d0f2e8d?w=400',
    createdAt: '2026-01-01',
  },
  {
    id: 2,
    name: '新年祝福文案',
    description: '温馨的新年祝福营销文案',
    type: 'text',
    content: '新年新气象，好运连连来！参与我们的春节特惠活动，让这个新年更加精彩...',
    createdAt: '2026-01-01',
  },
  {
    id: 3,
    name: '产品展示图',
    description: '主打产品的精美展示图片',
    type: 'image',
    url: 'https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=400',
    createdAt: '2025-12-28',
  },
]

export const getMaterialDetailItems = (material) => {
  if (!material) return []
  const items = [
    { label: '素材名称', value: material.name || '-' },
    { label: '素材类型', value: getMaterialTypeText(material.type) },
    { label: '素材描述', value: material.description || '-' },
    { label: '创建时间', value: material.createdAt || '-' },
  ]
  if (material.type === 'image' && material.url) {
    items.push({ label: '图片链接', value: material.url })
  }
  if (material.type === 'text' && material.content) {
    items.push({ label: '文案内容', value: material.content })
  }
  return items
}

export const copyMaterialToClipboard = async (material) => {
  if (!material) return false
  
  let textToCopy = ''
  if (material.type === 'text' && material.content) {
    textToCopy = material.content
  } else if (material.type === 'image' && material.url) {
    textToCopy = material.url
  } else {
    textToCopy = material.name || ''
  }
  
  try {
    await navigator.clipboard.writeText(textToCopy)
    return true
  } catch {
    return false
  }
}

export const buildMaterialUpdatePayload = (material, updates) => {
  return {
    id: material?.id,
    name: updates?.name || material?.name || '',
    description: updates?.description || material?.description || '',
    category: updates?.category || material?.type || 'image',
  }
}

export const validateMaterialEdit = (form) => {
  if (!form || !form.name || !form.name.trim()) {
    return '请填写素材名称'
  }
  return ''
}

export const getDefaultEditForm = (material = null) => ({
  name: material?.name || '',
  description: material?.description || '',
  category: material?.type || 'image',
})
