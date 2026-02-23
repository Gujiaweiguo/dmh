<template>
  <div class="brand-materials">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">素材库</h1>
      <button @click="showUploadModal = true" class="upload-btn">
        📤 上传
      </button>
    </div>

    <!-- 分类标签 -->
    <div class="category-tabs">
      <button
        v-for="category in categories"
        :key="category.value"
        @click="currentCategory = category.value"
        :class="['category-tab', { active: currentCategory === category.value }]"
      >
        {{ category.label }}
      </button>
    </div>

    <!-- 素材列表 -->
    <div class="materials-grid">
      <div
        v-for="material in filteredMaterials"
        :key="material.id"
        class="material-card"
        @click="viewMaterial(material)"
      >
        <div class="material-preview">
          <img
            v-if="material.type === 'image'"
            :src="material.url"
            :alt="material.name"
            class="material-image"
          >
          <div v-else class="material-text">
            <div class="text-icon">📝</div>
            <div class="text-preview">{{ material.content.substring(0, 50) }}...</div>
          </div>
        </div>
        
        <div class="material-info">
          <h3 class="material-name">{{ material.name }}</h3>
          <p class="material-desc">{{ material.description }}</p>
          <div class="material-meta">
            <span class="material-type">{{ getTypeText(material.type) }}</span>
            <span class="material-date">{{ formatDate(material.createdAt) }}</span>
          </div>
        </div>

        <div class="material-actions">
          <button @click.stop="useMaterial(material)" class="action-btn use">
            使用
          </button>
          <button @click.stop="editMaterial(material)" class="action-btn edit">
            编辑
          </button>
          <button @click.stop="deleteMaterial(material)" class="action-btn delete">
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- AI生成区域 -->
    <div class="ai-section">
      <h2 class="ai-title">🤖 AI智能生成</h2>
      <div class="ai-tools">
        <button @click="showAITextModal = true" class="ai-tool-btn">
          <div class="tool-icon">✍️</div>
          <div class="tool-text">生成文案</div>
        </button>
        <button @click="showAIImageModal = true" class="ai-tool-btn">
          <div class="tool-icon">🎨</div>
          <div class="tool-text">生成图片</div>
        </button>
      </div>
    </div>

    <!-- 上传模态框 -->
    <div v-if="showUploadModal" class="modal-overlay" @click="showUploadModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>上传素材</h3>
          <button @click="showUploadModal = false" class="close-btn">✕</button>
        </div>
        
        <div class="upload-area">
          <input
            ref="fileInput"
            type="file"
            accept="image/*"
            @change="handleFileUpload"
            class="file-input"
          >
          <div class="upload-placeholder" @click="$refs.fileInput.click()">
            <div class="upload-icon">📁</div>
            <p>点击选择文件或拖拽到此处</p>
            <p class="upload-hint">支持 JPG、PNG、GIF 格式</p>
          </div>
        </div>

        <div class="upload-form">
          <div class="form-group">
            <label>素材名称</label>
            <input v-model="uploadForm.name" type="text" class="form-input" placeholder="请输入素材名称">
          </div>
          <div class="form-group">
            <label>素材描述</label>
            <textarea v-model="uploadForm.description" class="form-textarea" placeholder="请输入素材描述"></textarea>
          </div>
          <div class="form-group">
            <label>分类</label>
            <select v-model="uploadForm.category" class="form-select">
              <option value="image">图片素材</option>
              <option value="text">文案素材</option>
            </select>
          </div>
        </div>

        <div class="modal-actions">
          <button @click="showUploadModal = false" class="cancel-btn">取消</button>
          <button @click="uploadMaterial" :disabled="uploading" class="confirm-btn">
            {{ uploading ? '上传中...' : '确认上传' }}
          </button>
        </div>
      </div>
    </div>

    <!-- AI文案生成模态框 -->
    <div v-if="showAITextModal" class="modal-overlay" @click="showAITextModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>AI文案生成</h3>
          <button @click="showAITextModal = false" class="close-btn">✕</button>
        </div>
        
        <div class="ai-form">
          <div class="form-group">
            <label>文案主题</label>
            <input v-model="aiTextForm.topic" type="text" class="form-input" placeholder="如：春节促销活动">
          </div>
          <div class="form-group">
            <label>文案风格</label>
            <select v-model="aiTextForm.style" class="form-select">
              <option value="professional">专业正式</option>
              <option value="casual">轻松活泼</option>
              <option value="urgent">紧迫感</option>
              <option value="emotional">情感化</option>
            </select>
          </div>
          <div class="form-group">
            <label>字数要求</label>
            <select v-model="aiTextForm.length" class="form-select">
              <option value="short">简短 (50字以内)</option>
              <option value="medium">中等 (100字左右)</option>
              <option value="long">详细 (200字以上)</option>
            </select>
          </div>
        </div>

        <div class="modal-actions">
          <button @click="showAITextModal = false" class="cancel-btn">取消</button>
          <button @click="generateAIText" :disabled="aiGenerating" class="confirm-btn">
            {{ aiGenerating ? '生成中...' : '生成文案' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 素材详情模态框 -->
    <div v-if="showDetailModal" class="modal-overlay" @click="showDetailModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>素材详情</h3>
          <button @click="showDetailModal = false" class="close-btn">✕</button>
        </div>
        
        <div class="detail-body">
          <div v-if="currentMaterial?.type === 'image'" class="detail-preview">
            <img :src="currentMaterial?.url" :alt="currentMaterial?.name" class="preview-image">
          </div>
          <div v-else-if="currentMaterial?.type === 'text'" class="detail-content">
            <p>{{ currentMaterial?.content }}</p>
          </div>
          
          <div class="detail-info">
            <div v-for="item in getMaterialDetailItems(currentMaterial)" :key="item.label" class="detail-row">
              <span class="detail-label">{{ item.label }}:</span>
              <span class="detail-value">{{ item.value }}</span>
            </div>
          </div>
        </div>

        <div class="modal-actions">
          <button @click="useMaterial(currentMaterial)" class="action-btn use">复制内容</button>
          <button @click="showDetailModal = false; editMaterial(currentMaterial)" class="action-btn edit">编辑</button>
          <button @click="showDetailModal = false" class="cancel-btn">关闭</button>
        </div>
      </div>
    </div>

    <!-- 素材编辑模态框 -->
    <div v-if="showEditModal" class="modal-overlay" @click="showEditModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑素材</h3>
          <button @click="showEditModal = false" class="close-btn">✕</button>
        </div>
        
        <div class="edit-form">
          <div class="form-group">
            <label>素材名称</label>
            <input v-model="editForm.name" type="text" class="form-input" placeholder="请输入素材名称">
          </div>
          <div class="form-group">
            <label>素材描述</label>
            <textarea v-model="editForm.description" class="form-textarea" placeholder="请输入素材描述"></textarea>
          </div>
          <div class="form-group">
            <label>分类</label>
            <select v-model="editForm.category" class="form-select">
              <option value="image">图片素材</option>
              <option value="text">文案素材</option>
            </select>
          </div>
        </div>

        <div class="modal-actions">
          <button @click="showEditModal = false" class="cancel-btn">取消</button>
          <button @click="saveMaterialEdit" class="confirm-btn">保存</button>
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <router-link to="/brand/dashboard" class="nav-item">
        <div class="nav-icon">🏠</div>
        <div class="nav-text">工作台</div>
      </router-link>
      <router-link to="/brand/campaigns" class="nav-item">
        <div class="nav-icon">🎯</div>
        <div class="nav-text">活动</div>
      </router-link>
      <router-link to="/brand/orders" class="nav-item">
        <div class="nav-icon">📋</div>
        <div class="nav-text">订单</div>
      </router-link>
      <router-link to="/brand/distributors" class="nav-item">
        <div class="nav-icon">🧭</div>
        <div class="nav-text">分销</div>
      </router-link>
      <router-link to="/brand/members" class="nav-item">
        <div class="nav-icon">👤</div>
        <div class="nav-text">会员</div>
      </router-link>
      <router-link to="/brand/promoters" class="nav-item">
        <div class="nav-icon">👥</div>
        <div class="nav-text">推广员</div>
      </router-link>
      <router-link to="/brand/settings" class="nav-item">
        <div class="nav-icon">⚙️</div>
        <div class="nav-text">设置</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import {
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
  getMaterialDetailItems,
  copyMaterialToClipboard,
  getDefaultEditForm,
  validateMaterialEdit,
  buildAIGeneratePayload,
  createMaterialFromAI,
} from './materials.logic.js'
import { aiApi } from '../../services/brandApi.js'

const materials = ref([])
const currentCategory = ref('all')
const showUploadModal = ref(false)
const showAITextModal = ref(false)
const showAIImageModal = ref(false)
const showDetailModal = ref(false)
const showEditModal = ref(false)
const currentMaterial = ref(null)
const editForm = reactive(getDefaultEditForm())
const uploading = ref(false)
const aiGenerating = ref(false)

const categories = getMaterialCategories()

const uploadForm = reactive(getDefaultUploadForm())

const aiTextForm = reactive(getDefaultAITextForm())

const filteredMaterials = computed(() => {
  return filterMaterialsByCategory(materials.value, currentCategory.value)
})

const getTypeText = (type) => {
  return getMaterialTypeText(type)
}

const formatDate = (dateString) => {
  return formatMaterialDate(dateString)
}

const loadMaterials = async () => {
  try {
    const response = await materialApi.getMaterials()
    const data = response.data || response
    materials.value = Array.isArray(data) ? data : (data.list || [])
  } catch (error) {
    console.error('加载素材失败:', error)
  }
}

const handleFileUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    uploadForm.file = file
    if (!uploadForm.name) {
      uploadForm.name = file.name.split('.')[0]
    }
  }
}

const uploadMaterial = async () => {
  const errorMsg = validateUploadInput(uploadForm)
  if (errorMsg) {
    alert(errorMsg)
    return
  }

  uploading.value = true
  try {
    // 调用API上传素材 - API expects (file, type)
    const fileType = uploadForm.category || 'image'
    const response = await materialApi.uploadMaterial(uploadForm.file, fileType)
    const newMaterial = response.data || response
    
    // 更新素材元数据（名称和描述）
    if (newMaterial.id) {
      // 如果API返回了新素材，添加额外信息
      newMaterial.name = uploadForm.name
      newMaterial.description = uploadForm.description
      newMaterial.category = uploadForm.category
    }
    
    // 添加到列表
    materials.value.unshift(newMaterial)
    showUploadModal.value = false
    // Reset form
    Object.assign(uploadForm, getDefaultUploadForm())
    
    materials.value.unshift(newMaterial)
    
    // 重置表单
    Object.assign(uploadForm, getDefaultUploadForm())
    
    showUploadModal.value = false
    alert('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    alert('上传失败')
  } finally {
    uploading.value = false
  }
}

const generateAIText = async () => {
  const errorMsg = validateAITextInput(aiTextForm.topic)
  if (errorMsg) {
    alert(errorMsg)
    return
  }

  aiGenerating.value = true
  try {
    const payload = buildAIGeneratePayload(aiTextForm)
    let content = ''
    
    try {
      const response = await aiApi.generateCopywriting(payload)
      const result = response?.data || response
      content = result?.content || result?.text || ''
    } catch (apiError) {
      console.warn('AI API 调用失败，使用本地生成:', apiError)
    }
    
    const newMaterial = createMaterialFromAI(
      aiTextForm.topic,
      content,
      new Date().toISOString().split('T')[0]
    )
    
    materials.value.unshift(newMaterial)
    
    Object.assign(aiTextForm, getDefaultAITextForm())
    showAITextModal.value = false
    alert('文案生成成功')
  } catch (error) {
    console.error('AI生成失败:', error)
    alert('生成失败，请重试')
  } finally {
    aiGenerating.value = false
  }
}

const viewMaterial = (material) => {
  currentMaterial.value = material
  showDetailModal.value = true
}

const useMaterial = async (material) => {
  const success = await copyMaterialToClipboard(material)
  if (success) {
    alert('素材内容已复制到剪贴板')
  } else {
    alert('复制失败，请手动复制')
  }
}

const editMaterial = (material) => {
  currentMaterial.value = material
  Object.assign(editForm, getDefaultEditForm(material))
  showEditModal.value = true
}

const saveMaterialEdit = async () => {
  const errorMsg = validateMaterialEdit(editForm)
  if (errorMsg) {
    alert(errorMsg)
    return
  }

  try {
    const index = materials.value.findIndex(m => m.id === currentMaterial.value.id)
    if (index > -1) {
      materials.value[index] = {
        ...materials.value[index],
        name: editForm.name,
        description: editForm.description,
        type: editForm.category,
      }
      alert('保存成功')
      showEditModal.value = false
    }
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败')
  }
}

const deleteMaterial = (material) => {
  if (confirm(`确定要删除素材"${material.name}"吗？`)) {
    const index = materials.value.findIndex(m => m.id === material.id)
    if (index > -1) {
      materials.value.splice(index, 1)
      alert('删除成功')
    }
  }
}

onMounted(() => {
  loadMaterials()
})
</script>

<style scoped>
.brand-materials {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 80px;
}

.top-nav {
  background: white;
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 10;
}

.nav-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
  color: #333;
}

.upload-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
}

.category-tabs {
  background: white;
  padding: 16px;
  display: flex;
  gap: 8px;
  border-bottom: 1px solid #eee;
}

.category-tab {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.category-tab.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.materials-grid {
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.material-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s;
}

.material-card:hover {
  transform: translateY(-2px);
}

.material-preview {
  height: 160px;
  overflow: hidden;
}

.material-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.material-text {
  height: 100%;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.text-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.text-preview {
  font-size: 14px;
  color: #666;
  text-align: center;
  line-height: 1.4;
}

.material-info {
  padding: 16px;
}

.material-name {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 8px 0;
  color: #333;
}

.material-desc {
  font-size: 14px;
  color: #666;
  margin: 0 0 12px 0;
  line-height: 1.4;
}

.material-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.material-type {
  font-size: 12px;
  background: #e3f2fd;
  color: #1976d2;
  padding: 2px 8px;
  border-radius: 12px;
}

.material-date {
  font-size: 12px;
  color: #999;
}

.material-actions {
  padding: 12px 16px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 8px;
}

.action-btn {
  flex: 1;
  padding: 6px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn.use {
  border-color: #4caf50;
  color: #4caf50;
}

.action-btn.edit {
  border-color: #2196f3;
  color: #2196f3;
}

.action-btn.delete {
  border-color: #f44336;
  color: #f44336;
}

.ai-section {
  padding: 16px;
  margin-top: 20px;
}

.ai-title {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 16px 0;
  color: #333;
}

.ai-tools {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.ai-tool-btn {
  background: white;
  border: 2px solid #667eea;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
}

.ai-tool-btn:hover {
  background: #667eea;
  color: white;
}

.tool-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.tool-text {
  font-size: 14px;
  font-weight: 500;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 20px;
  width: 100%;
  max-width: 400px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  padding: 20px 20px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #999;
}

.upload-area {
  padding: 20px;
}

.file-input {
  display: none;
}

.upload-placeholder {
  border: 2px dashed #ddd;
  border-radius: 12px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.3s;
}

.upload-placeholder:hover {
  border-color: #667eea;
}

.upload-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.upload-hint {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}

.upload-form,
.ai-form {
  padding: 0 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-input,
.form-textarea,
.form-select {
  width: 100%;
  padding: 12px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus {
  outline: none;
  border-color: #667eea;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.modal-actions {
  padding: 20px;
  display: flex;
  gap: 12px;
}

.cancel-btn,
.confirm-btn {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: opacity 0.3s;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.confirm-btn {
  background: #667eea;
  color: white;
}

.confirm-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  border-top: 1px solid #eee;
  padding: 8px 0;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #999;
  padding: 8px;
}

.nav-item.active {
  color: #667eea;
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 4px;
}

.nav-text {
  font-size: 12px;
}
</style>
