import { describe, expect, it, beforeEach } from 'vitest'
import {
  buildAIGeneratedText,
  createAIGeneratedMaterial,
  createUploadedMaterial,
  filterMaterialsByCategory,
  formatMaterialDate,
  getDefaultAITextForm,
  getDefaultUploadForm,
  getMaterialCategories,
  getMaterialTypeText,
  getMockMaterials,
  validateAITextInput,
  validateUploadInput,
} from '../../src/views/brand/materials.logic.js'

describe('materials logic', () => {
  it('provides default categories and forms', () => {
    expect(getMaterialCategories()).toHaveLength(4)
    expect(getDefaultUploadForm()).toEqual({
      name: '',
      description: '',
      category: 'image',
      file: null,
    })
    expect(getDefaultAITextForm()).toEqual({
      topic: '',
      style: 'professional',
      length: 'medium',
    })
  })

  it('filters materials by category', () => {
    const list = [
      { id: 1, type: 'image' },
      { id: 2, type: 'text' },
      { id: 3, type: 'image' },
    ]
    expect(filterMaterialsByCategory(list, 'all')).toHaveLength(3)
    expect(filterMaterialsByCategory(list, 'image').map((item) => item.id)).toEqual([1, 3])
  })

  it('maps type text and formats date', () => {
    expect(getMaterialTypeText('image')).toBe('图片')
    expect(getMaterialTypeText('text')).toBe('文案')
    expect(getMaterialTypeText('video')).toBe('视频')
    expect(getMaterialTypeText('other')).toBe('other')
    expect(formatMaterialDate('2026-01-01')).toContain('2026')
  })

  it('validates upload and ai text input', () => {
    expect(validateUploadInput({ file: null, name: '' })).toBe('请选择文件并填写素材名称')
    expect(validateUploadInput({ file: { name: 'a.png' }, name: '素材A' })).toBe('')

    expect(validateAITextInput('')).toBe('请输入文案主题')
    expect(validateAITextInput('春节促销')).toBe('')
  })

  it('creates uploaded and ai generated materials', () => {
    const uploaded = createUploadedMaterial({
      id: 1,
      name: '海报',
      description: 'desc',
      category: 'image',
      url: 'blob:test',
      createdAt: '2026-02-13',
    })
    expect(uploaded).toMatchObject({
      id: 1,
      type: 'image',
      url: 'blob:test',
    })

    const text = buildAIGeneratedText('春节活动')
    expect(text).toContain('春节活动')

    const aiMaterial = createAIGeneratedMaterial({
      id: 2,
      topic: '春节活动',
      createdAt: '2026-02-13',
    })
    expect(aiMaterial).toMatchObject({
      id: 2,
      type: 'text',
      name: 'AI生成-春节活动',
    })
    expect(aiMaterial.content).toContain('春节活动')
  })

  it('provides mock materials', () => {
    const mock = getMockMaterials()
    expect(mock).toHaveLength(3)
    expect(mock[0]).toMatchObject({ type: 'image' })
  })
})

import {
  getMaterialDetailItems,
  copyMaterialToClipboard,
  getDefaultEditForm,
  validateMaterialEdit,
} from '../../src/views/brand/materials.logic.js'

describe('getMaterialDetailItems', () => {
  it('returns empty array for null material', () => {
    expect(getMaterialDetailItems(null)).toEqual([])
  })

  it('returns basic items for image material', () => {
    const material = {
      name: '测试图片',
      type: 'image',
      description: '测试描述',
      url: 'https://example.com/image.png',
      createdAt: '2026-02-20',
    }
    const items = getMaterialDetailItems(material)

    expect(items).toHaveLength(5)
    expect(items.find(i => i.label === '素材名称')?.value).toBe('测试图片')
    expect(items.find(i => i.label === '素材类型')?.value).toBe('图片')
    expect(items.find(i => i.label === '图片链接')?.value).toBe('https://example.com/image.png')
  })

  it('returns content for text material', () => {
    const material = {
      name: '测试文案',
      type: 'text',
      description: '文案描述',
      content: '这是测试文案内容',
      createdAt: '2026-02-20',
    }
    const items = getMaterialDetailItems(material)

    expect(items.find(i => i.label === '文案内容')?.value).toBe('这是测试文案内容')
  })
})

describe('copyMaterialToClipboard', () => {
  beforeEach(() => {
    Object.assign(navigator, {
      clipboard: {
        writeText: async (text) => true,
      },
    })
  })

  it('returns false for null material', async () => {
    const result = await copyMaterialToClipboard(null)
    expect(result).toBe(false)
  })

  it('copies content for text material', async () => {
    const material = {
      name: '文案',
      type: 'text',
      content: '测试文案内容',
    }
    const result = await copyMaterialToClipboard(material)
    expect(result).toBe(true)
  })

  it('copies url for image material', async () => {
    const material = {
      name: '图片',
      type: 'image',
      url: 'https://example.com/image.png',
    }
    const result = await copyMaterialToClipboard(material)
    expect(result).toBe(true)
  })

  it('copies name as fallback', async () => {
    const material = {
      name: '素材名称',
      type: 'other',
    }
    const result = await copyMaterialToClipboard(material)
    expect(result).toBe(true)
  })
})

describe('getDefaultEditForm', () => {
  it('returns empty form without material', () => {
    const form = getDefaultEditForm()
    expect(form).toEqual({
      name: '',
      description: '',
      category: 'image',
    })
  })

  it('populates form from material', () => {
    const material = {
      name: '现有素材',
      description: '现有描述',
      type: 'text',
    }
    const form = getDefaultEditForm(material)
    expect(form).toEqual({
      name: '现有素材',
      description: '现有描述',
      category: 'text',
    })
  })
})

describe('validateMaterialEdit', () => {
  it('returns error for null form', () => {
    expect(validateMaterialEdit(null)).toBe('请填写素材名称')
  })

  it('returns error for empty name', () => {
    expect(validateMaterialEdit({ name: '' })).toBe('请填写素材名称')
  })

  it('returns error for whitespace-only name', () => {
    expect(validateMaterialEdit({ name: '   ' })).toBe('请填写素材名称')
  })

  it('returns empty string for valid form', () => {
    expect(validateMaterialEdit({ name: '素材名称' })).toBe('')
  })
})

import {
  buildAIGeneratePayload,
  parseAIGenerateResponse,
  createMaterialFromAI,
} from '../../src/views/brand/materials.logic.js'

describe('buildAIGeneratePayload', () => {
  it('builds payload from form', () => {
    const form = {
      topic: '春节促销',
      style: 'casual',
      length: 'short',
    }
    const payload = buildAIGeneratePayload(form)
    
    expect(payload).toEqual({
      topic: '春节促销',
      style: 'casual',
      length: 'short',
    })
  })

  it('provides defaults for missing fields', () => {
    const payload = buildAIGeneratePayload({})
    
    expect(payload).toEqual({
      topic: '',
      style: 'professional',
      length: 'medium',
    })
  })

  it('handles null form', () => {
    const payload = buildAIGeneratePayload(null)
    
    expect(payload.topic).toBe('')
    expect(payload.style).toBe('professional')
  })
})

describe('parseAIGenerateResponse', () => {
  it('parses response with content field', () => {
    const response = { data: { content: '生成的文案内容' } }
    const result = parseAIGenerateResponse(response)
    
    expect(result.content).toBe('生成的文案内容')
  })

  it('parses response with text field', () => {
    const response = { data: { text: '备选文案' } }
    const result = parseAIGenerateResponse(response)
    
    expect(result.content).toBe('备选文案')
  })

  it('handles response without data wrapper', () => {
    const response = { content: '直接内容' }
    const result = parseAIGenerateResponse(response)
    
    expect(result.content).toBe('直接内容')
  })

  it('returns empty content for null response', () => {
    const result = parseAIGenerateResponse(null)
    
    expect(result.content).toBe('')
    expect(result.name).toBeDefined()
  })
})

describe('createMaterialFromAI', () => {
  it('creates material with provided content', () => {
    const material = createMaterialFromAI('春节活动', '这是AI生成的文案', '2026-02-23')
    
    expect(material.name).toBe('AI生成-春节活动')
    expect(material.content).toBe('这是AI生成的文案')
    expect(material.type).toBe('text')
    expect(material.createdAt).toBe('2026-02-23')
    expect(material.id).toBeDefined()
  })

  it('generates content when not provided', () => {
    const material = createMaterialFromAI('春节活动', '', '2026-02-23')
    
    expect(material.content).toContain('春节活动')
  })

  it('uses current date when not provided', () => {
    const material = createMaterialFromAI('测试', '内容')
    
    expect(material.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}$/)
  })
})
