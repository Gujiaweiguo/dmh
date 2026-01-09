// DMH H5 品牌管理端
let authToken = localStorage.getItem('h5_token');
let campaigns = [];
let currentCampaign = null;
let currentTab = 'home';

// 初始化应用
function init() {
    render();
    if (authToken) {
        showMainPage();
    }
}

// 渲染应用
function render() {
    document.getElementById('app').innerHTML = `
        <!-- 登录页面 -->
        <div class="login-page" id="loginPage">
            <div class="login-card">
                <div class="logo">
                    <h1>DMH 品牌管理</h1>
                    <p>数字营销中台 · 品牌管理端</p>
                </div>
                <form id="loginForm">
                    <div class="form-group">
                        <label>用户名</label>
                        <input type="text" id="username" value="brand_manager" required>
                    </div>
                    <div class="form-group">
                        <label>密码</label>
                        <input type="password" id="password" value="123456" required>
                    </div>
                    <div id="errorMsg"></div>
                    <button type="submit" class="btn" id="loginBtn">登录</button>
                </form>
                <div class="test-info">
                    <p><strong>测试账号</strong></p>
                    <p>用户名: brand_manager | 密码: 123456</p>
                </div>
            </div>
        </div>

        <!-- 主页面 -->
        <div class="main-page" id="mainPage">
            <div class="header">
                <h2>品牌管理中心</h2>
                <button class="logout-btn" onclick="logout()">退出</button>
            </div>
            <div class="stats">
                <div class="stat-card purple"><div class="number" id="totalCampaigns">0</div><div class="label">总活动</div></div>
                <div class="stat-card green"><div class="number" id="activeCampaigns">0</div><div class="label">进行中</div></div>
                <div class="stat-card orange"><div class="number" id="totalParticipants">0</div><div class="label">参与数</div></div>
                <div class="stat-card red"><div class="number" id="conversionRate">0%</div><div class="label">转化率</div></div>
            </div>
            <div class="section">
                <div class="section-header">
                    <span class="section-title">📋 我的活动</span>
                    <button class="btn btn-sm" onclick="openCreateModal()">+ 创建活动</button>
                </div>
                <div id="campaignList"><div class="empty-state">加载中...</div></div>
            </div>
            <div class="tab-bar">
                <div class="tab-item active" onclick="switchTab('home')"><div class="icon">🏠</div>首页</div>
                <div class="tab-item" onclick="switchTab('campaigns')"><div class="icon">📋</div>活动</div>
                <div class="tab-item" onclick="openCreateModal()"><div class="icon">➕</div>创建</div>
                <div class="tab-item" onclick="switchTab('profile')"><div class="icon">👤</div>我的</div>
            </div>
        </div>
        ${renderModals()}
    `;
    bindEvents();
}


// 渲染模态框
function renderModals() {
    return `
        <!-- 创建/编辑活动模态框 -->
        <div class="modal" id="campaignModal">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 id="modalTitle">创建活动</h3>
                    <button class="modal-close" onclick="closeModal('campaignModal')">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="campaignForm">
                        <div class="form-group">
                            <label>活动名称 *</label>
                            <input type="text" id="campaignName" required placeholder="请输入活动名称">
                        </div>
                        <div class="form-group">
                            <label>活动描述</label>
                            <textarea id="campaignDesc" rows="3" placeholder="请输入活动描述"></textarea>
                        </div>
                        <div class="form-group">
                            <label>开始时间 *</label>
                            <input type="date" id="startTime" required>
                        </div>
                        <div class="form-group">
                            <label>结束时间 *</label>
                            <input type="date" id="endTime" required>
                        </div>
                        <div class="form-group">
                            <label>奖励金额 (元)</label>
                            <input type="number" id="rewardRule" value="0" min="0" placeholder="每人奖励金额">
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button class="btn btn-sm btn-secondary" onclick="closeModal('campaignModal')">取消</button>
                    <button class="btn btn-sm" onclick="saveCampaign()">保存</button>
                </div>
            </div>
        </div>

        <!-- 查看活动详情模态框 -->
        <div class="modal" id="viewModal">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>活动详情</h3>
                    <button class="modal-close" onclick="closeModal('viewModal')">&times;</button>
                </div>
                <div class="modal-body" id="viewContent"></div>
                <div class="modal-footer">
                    <button class="btn btn-sm btn-secondary" onclick="closeModal('viewModal')">关闭</button>
                </div>
            </div>
        </div>
    `;
}

// 绑定事件
function bindEvents() {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
}


// 登录处理
async function handleLogin(e) {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const loginBtn = document.getElementById('loginBtn');
    const errorMsg = document.getElementById('errorMsg');
    
    loginBtn.disabled = true;
    loginBtn.textContent = '登录中...';
    errorMsg.innerHTML = '';
    
    try {
        const response = await fetch('/api/v1/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await response.json();
        if (response.ok && data.token) {
            authToken = data.token;
            localStorage.setItem('h5_token', authToken);
            showMainPage();
        } else {
            throw new Error(data.message || '登录失败');
        }
    } catch (error) {
        errorMsg.innerHTML = `<div class="error-msg">登录失败: ${error.message}</div>`;
    } finally {
        loginBtn.disabled = false;
        loginBtn.textContent = '登录';
    }
}

// 退出登录
function logout() {
    localStorage.removeItem('h5_token');
    authToken = null;
    document.getElementById('loginPage').classList.remove('hidden');
    document.getElementById('mainPage').classList.remove('active');
}

// 显示主页面
function showMainPage() {
    document.getElementById('loginPage').classList.add('hidden');
    document.getElementById('mainPage').classList.add('active');
    loadCampaigns();
}

// 切换标签
function switchTab(tab) {
    currentTab = tab;
    document.querySelectorAll('.tab-item').forEach((el, i) => {
        el.classList.toggle('active', ['home', 'campaigns', '', 'profile'][i] === tab);
    });
}


// 加载活动列表
async function loadCampaigns() {
    try {
        const response = await fetch('/api/v1/campaigns?page=1&pageSize=100', {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (response.ok) {
            const data = await response.json();
            campaigns = data.campaigns || data.list || [];
            updateStats();
            renderCampaignList();
        }
    } catch (error) {
        document.getElementById('campaignList').innerHTML = `<div class="empty-state">加载失败: ${error.message}</div>`;
    }
}

// 更新统计数据
function updateStats() {
    const total = campaigns.length;
    const active = campaigns.filter(c => c.status === 'ACTIVE' || c.status === 'active').length;
    const participants = campaigns.reduce((sum, c) => sum + (c.orderCount || 0), 0);
    document.getElementById('totalCampaigns').textContent = total;
    document.getElementById('activeCampaigns').textContent = active;
    document.getElementById('totalParticipants').textContent = participants;
    document.getElementById('conversionRate').textContent = participants > 0 ? '15%' : '0%';
}

// 渲染活动列表
function renderCampaignList() {
    const listEl = document.getElementById('campaignList');
    if (campaigns.length === 0) {
        listEl.innerHTML = '<div class="empty-state">暂无活动，点击上方按钮创建</div>';
        return;
    }
    listEl.innerHTML = campaigns.map(c => `
        <div class="campaign-card">
            <h4>${c.name}</h4>
            <div class="campaign-meta">
                <span class="status ${getStatusClass(c.status)}">${getStatusText(c.status)}</span>
                <span>👥 ${c.orderCount || 0}人</span>
                <span>📅 ${(c.startTime || '').substring(0, 10)}</span>
            </div>
            <div class="campaign-actions">
                <button class="btn-view" onclick="viewCampaign(${c.id})">查看</button>
                <button class="btn-edit" onclick="editCampaign(${c.id})">编辑</button>
                <button style="background: #e0e7ff; color: #4f46e5;" onclick="openPageDesign(${c.id})">📐 页面设计</button>
                ${c.status === 'draft' ? `<button class="btn-publish" onclick="publishCampaign(${c.id})">发布</button>` : ''}
                ${c.status === 'ACTIVE' || c.status === 'active' ? `<button class="btn-publish" onclick="pauseCampaign(${c.id})">暂停</button>` : ''}
                ${c.status === 'PAUSED' || c.status === 'paused' ? `<button class="btn-success" onclick="resumeCampaign(${c.id})">恢复</button>` : ''}
                <button class="btn-delete" onclick="deleteCampaign(${c.id})">删除</button>
            </div>
        </div>
    `).join('');
}

function getStatusClass(status) {
    if (status === 'ACTIVE' || status === 'active') return 'active';
    if (status === 'PAUSED' || status === 'paused') return 'paused';
    return 'draft';
}

function getStatusText(status) {
    if (status === 'ACTIVE' || status === 'active') return '进行中';
    if (status === 'PAUSED' || status === 'paused') return '已暂停';
    return '草稿';
}


// 打开创建模态框
function openCreateModal() {
    currentCampaign = null;
    document.getElementById('modalTitle').textContent = '创建活动';
    document.getElementById('campaignName').value = '';
    document.getElementById('campaignDesc').value = '';
    document.getElementById('startTime').value = '';
    document.getElementById('endTime').value = '';
    document.getElementById('rewardRule').value = '0';
    openModal('campaignModal');
}

// 编辑活动
function editCampaign(id) {
    currentCampaign = campaigns.find(c => c.id === id);
    if (!currentCampaign) return;
    document.getElementById('modalTitle').textContent = '编辑活动';
    document.getElementById('campaignName').value = currentCampaign.name || '';
    document.getElementById('campaignDesc').value = currentCampaign.description || '';
    document.getElementById('startTime').value = (currentCampaign.startTime || '').substring(0, 10);
    document.getElementById('endTime').value = (currentCampaign.endTime || '').substring(0, 10);
    document.getElementById('rewardRule').value = currentCampaign.rewardRule || 0;
    openModal('campaignModal');
}

// 查看活动详情
function viewCampaign(id) {
    const c = campaigns.find(c => c.id === id);
    if (!c) return;
    document.getElementById('viewContent').innerHTML = `
        <div style="space-y: 15px;">
            <p><strong>活动名称：</strong>${c.name}</p>
            <p><strong>活动状态：</strong><span class="status ${getStatusClass(c.status)}">${getStatusText(c.status)}</span></p>
            <p><strong>活动描述：</strong>${c.description || '暂无描述'}</p>
            <p><strong>开始时间：</strong>${(c.startTime || '').substring(0, 10)}</p>
            <p><strong>结束时间：</strong>${(c.endTime || '').substring(0, 10)}</p>
            <p><strong>参与人数：</strong>${c.orderCount || 0}人</p>
            <p><strong>奖励金额：</strong>¥${c.rewardRule || 0}</p>
            <hr style="margin: 15px 0; border: none; border-top: 1px solid #eee;">
            <p><strong>数据统计</strong></p>
            <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 10px; margin-top: 10px;">
                <div style="background: #f0f0f0; padding: 15px; border-radius: 8px; text-align: center;">
                    <div style="font-size: 24px; font-weight: bold; color: #667eea;">${c.orderCount || 0}</div>
                    <div style="font-size: 12px; color: #666;">总参与</div>
                </div>
                <div style="background: #f0f0f0; padding: 15px; border-radius: 8px; text-align: center;">
                    <div style="font-size: 24px; font-weight: bold; color: #10b981;">${Math.floor((c.orderCount || 0) * 0.8)}</div>
                    <div style="font-size: 12px; color: #666;">有效报名</div>
                </div>
            </div>
        </div>
    `;
    openModal('viewModal');
}


// 保存活动
async function saveCampaign() {
    const name = document.getElementById('campaignName').value;
    const description = document.getElementById('campaignDesc').value;
    const startTime = document.getElementById('startTime').value;
    const endTime = document.getElementById('endTime').value;
    const rewardRule = parseInt(document.getElementById('rewardRule').value) || 0;
    
    if (!name || !startTime || !endTime) {
        alert('请填写必填字段');
        return;
    }
    
    try {
        const url = currentCampaign ? `/api/v1/campaigns/${currentCampaign.id}` : '/api/v1/campaigns';
        const method = currentCampaign ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method,
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({
                name, description, startTime, endTime, rewardRule, brandId: 1,
                status: currentCampaign?.status || 'active'
            })
        });
        
        if (response.ok) {
            alert(currentCampaign ? '活动更新成功' : '活动创建成功');
            closeModal('campaignModal');
            loadCampaigns();
        } else {
            const data = await response.json();
            throw new Error(data.message || '操作失败');
        }
    } catch (error) {
        alert('操作失败: ' + error.message);
    }
}

// 发布活动
async function publishCampaign(id) {
    if (!confirm('确定要发布此活动吗？')) return;
    await updateCampaignStatus(id, 'active');
}

// 暂停活动
async function pauseCampaign(id) {
    if (!confirm('确定要暂停此活动吗？')) return;
    await updateCampaignStatus(id, 'paused');
}

// 恢复活动
async function resumeCampaign(id) {
    if (!confirm('确定要恢复此活动吗？')) return;
    await updateCampaignStatus(id, 'active');
}

// 更新活动状态
async function updateCampaignStatus(id, status) {
    try {
        const response = await fetch(`/api/v1/campaigns/${id}/status`, {
            method: 'PUT',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({ status })
        });
        if (response.ok) {
            alert('操作成功');
            loadCampaigns();
        } else {
            throw new Error('操作失败');
        }
    } catch (error) {
        alert('操作失败: ' + error.message);
    }
}

// 删除活动
async function deleteCampaign(id) {
    if (!confirm('确定要删除此活动吗？此操作不可恢复！')) return;
    try {
        const response = await fetch(`/api/v1/campaigns/${id}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (response.ok) {
            alert('删除成功');
            loadCampaigns();
        } else {
            throw new Error('删除失败');
        }
    } catch (error) {
        alert('删除失败: ' + error.message);
    }
}

// 模态框操作
function openModal(id) { document.getElementById(id).classList.add('active'); }
function closeModal(id) { document.getElementById(id).classList.remove('active'); }

// 初始化
init();


// ==================== 页面设计功能 ====================

// 组件类型定义
const componentTypes = [
    { type: 'banner', name: '横幅图片', icon: '🖼️', desc: '添加横幅图片' },
    { type: 'text', name: '文本内容', icon: '📝', desc: '添加文字说明' },
    { type: 'video', name: '视频播放', icon: '🎬', desc: '嵌入视频' },
    { type: 'countdown', name: '倒计时', icon: '⏰', desc: '活动倒计时' },
    { type: 'button', name: '按钮', icon: '🔘', desc: '行动按钮' },
    { type: 'divider', name: '分割线', icon: '➖', desc: '内容分隔' }
];

let pageComponents = [];
let pageSettings = {
    title: '',
    description: '',
    backgroundColor: '#ffffff',
    primaryColor: '#667eea',
    buttonColor: '#667eea'
};

// 打开页面设计
function openPageDesign(id) {
    currentCampaign = campaigns.find(c => c.id === id);
    if (!currentCampaign) return;
    
    // 加载已保存的页面配置
    pageComponents = currentCampaign.pageComponents || [];
    pageSettings = currentCampaign.pageSettings || {
        title: currentCampaign.name || '',
        description: currentCampaign.description || '',
        backgroundColor: '#ffffff',
        primaryColor: '#667eea',
        buttonColor: '#667eea'
    };
    
    renderPageDesignModal();
    openModal('pageDesignModal');
}

// 渲染页面设计模态框
function renderPageDesignModal() {
    let modal = document.getElementById('pageDesignModal');
    if (!modal) {
        modal = document.createElement('div');
        modal.id = 'pageDesignModal';
        modal.className = 'modal';
        document.body.appendChild(modal);
    }
    
    modal.innerHTML = `
        <div class="modal-content" style="max-width: 100%; height: 100%; max-height: 100%; border-radius: 0;">
            <div class="modal-header" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white;">
                <h3>📐 页面设计 - ${currentCampaign?.name || ''}</h3>
                <button class="modal-close" onclick="closeModal('pageDesignModal')" style="color: white;">&times;</button>
            </div>
            <div style="display: flex; height: calc(100% - 130px); overflow: hidden;">
                <!-- 左侧：组件库 -->
                <div style="width: 200px; background: #f8f9fa; padding: 15px; overflow-y: auto; border-right: 1px solid #eee;">
                    <h4 style="margin-bottom: 15px; font-size: 14px; color: #333;">📦 组件库</h4>
                    ${componentTypes.map(c => `
                        <div onclick="addComponent('${c.type}')" style="background: white; padding: 12px; border-radius: 8px; margin-bottom: 10px; cursor: pointer; border: 1px solid #eee; transition: all 0.2s;">
                            <div style="font-size: 20px; margin-bottom: 5px;">${c.icon}</div>
                            <div style="font-size: 13px; font-weight: 600; color: #333;">${c.name}</div>
                            <div style="font-size: 11px; color: #999;">${c.desc}</div>
                        </div>
                    `).join('')}
                    
                    <h4 style="margin: 20px 0 15px; font-size: 14px; color: #333;">🎨 页面设置</h4>
                    <div style="background: white; padding: 12px; border-radius: 8px; border: 1px solid #eee;">
                        <div style="margin-bottom: 10px;">
                            <label style="font-size: 12px; color: #666;">页面标题</label>
                            <input type="text" id="pageTitle" value="${pageSettings.title}" onchange="updatePageSetting('title', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 12px;">
                        </div>
                        <div style="margin-bottom: 10px;">
                            <label style="font-size: 12px; color: #666;">背景色</label>
                            <input type="color" id="pageBgColor" value="${pageSettings.backgroundColor}" onchange="updatePageSetting('backgroundColor', this.value)" style="width: 100%; height: 30px; border: none; cursor: pointer;">
                        </div>
                        <div style="margin-bottom: 10px;">
                            <label style="font-size: 12px; color: #666;">主题色</label>
                            <input type="color" id="pagePrimaryColor" value="${pageSettings.primaryColor}" onchange="updatePageSetting('primaryColor', this.value)" style="width: 100%; height: 30px; border: none; cursor: pointer;">
                        </div>
                        <div>
                            <label style="font-size: 12px; color: #666;">按钮色</label>
                            <input type="color" id="pageButtonColor" value="${pageSettings.buttonColor}" onchange="updatePageSetting('buttonColor', this.value)" style="width: 100%; height: 30px; border: none; cursor: pointer;">
                        </div>
                    </div>
                </div>
                
                <!-- 中间：组件配置 -->
                <div style="flex: 1; padding: 15px; overflow-y: auto; background: #fff;">
                    <h4 style="margin-bottom: 15px; font-size: 14px; color: #333;">⚙️ 已添加组件</h4>
                    <div id="componentList">
                        ${renderComponentList()}
                    </div>
                </div>
                
                <!-- 右侧：预览 -->
                <div style="width: 320px; background: #f0f0f0; padding: 15px; overflow-y: auto;">
                    <h4 style="margin-bottom: 15px; font-size: 14px; color: #333;">👁️ 实时预览</h4>
                    <div id="pagePreview" style="background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 20px rgba(0,0,0,0.1);">
                        ${renderPagePreview()}
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button class="btn btn-sm btn-secondary" onclick="closeModal('pageDesignModal')">取消</button>
                <button class="btn btn-sm btn-success" onclick="previewPage()">预览</button>
                <button class="btn btn-sm" onclick="savePageDesign()">保存设计</button>
            </div>
        </div>
    `;
}


// 渲染组件列表
function renderComponentList() {
    if (pageComponents.length === 0) {
        return '<div style="text-align: center; padding: 40px; color: #999;">从左侧拖拽组件到这里</div>';
    }
    return pageComponents.map((comp, index) => `
        <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 10px; border: 1px solid #eee;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                <span style="font-weight: 600; font-size: 14px;">${componentTypes.find(c => c.type === comp.type)?.icon || '📦'} ${componentTypes.find(c => c.type === comp.type)?.name || comp.type}</span>
                <div>
                    ${index > 0 ? `<button onclick="moveComponent(${index}, -1)" style="padding: 4px 8px; border: none; background: #e0e0e0; border-radius: 4px; cursor: pointer; margin-right: 5px;">↑</button>` : ''}
                    ${index < pageComponents.length - 1 ? `<button onclick="moveComponent(${index}, 1)" style="padding: 4px 8px; border: none; background: #e0e0e0; border-radius: 4px; cursor: pointer; margin-right: 5px;">↓</button>` : ''}
                    <button onclick="removeComponent(${index})" style="padding: 4px 8px; border: none; background: #fee; color: #c33; border-radius: 4px; cursor: pointer;">删除</button>
                </div>
            </div>
            ${renderComponentConfig(comp, index)}
        </div>
    `).join('');
}

// 渲染组件配置
function renderComponentConfig(comp, index) {
    switch (comp.type) {
        case 'banner':
            return `<input type="text" value="${comp.config?.imageUrl || ''}" placeholder="图片URL" onchange="updateComponentConfig(${index}, 'imageUrl', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">`;
        case 'text':
            return `<textarea rows="3" placeholder="输入文本内容" onchange="updateComponentConfig(${index}, 'content', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">${comp.config?.content || ''}</textarea>
                    <select onchange="updateComponentConfig(${index}, 'align', this.value)" style="margin-top: 8px; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">
                        <option value="left" ${comp.config?.align === 'left' ? 'selected' : ''}>左对齐</option>
                        <option value="center" ${comp.config?.align === 'center' ? 'selected' : ''}>居中</option>
                        <option value="right" ${comp.config?.align === 'right' ? 'selected' : ''}>右对齐</option>
                    </select>`;
        case 'video':
            return `<input type="text" value="${comp.config?.videoUrl || ''}" placeholder="视频URL" onchange="updateComponentConfig(${index}, 'videoUrl', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">`;
        case 'countdown':
            return `<input type="datetime-local" value="${comp.config?.endTime || ''}" onchange="updateComponentConfig(${index}, 'endTime', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">`;
        case 'button':
            return `<input type="text" value="${comp.config?.text || '立即参与'}" placeholder="按钮文字" onchange="updateComponentConfig(${index}, 'text', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; margin-bottom: 8px;">
                    <input type="text" value="${comp.config?.link || ''}" placeholder="跳转链接" onchange="updateComponentConfig(${index}, 'link', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">`;
        case 'divider':
            return `<select onchange="updateComponentConfig(${index}, 'style', this.value)" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">
                        <option value="solid" ${comp.config?.style === 'solid' ? 'selected' : ''}>实线</option>
                        <option value="dashed" ${comp.config?.style === 'dashed' ? 'selected' : ''}>虚线</option>
                        <option value="dotted" ${comp.config?.style === 'dotted' ? 'selected' : ''}>点线</option>
                    </select>`;
        default:
            return '';
    }
}


// 渲染页面预览
function renderPagePreview() {
    const bgColor = pageSettings.backgroundColor;
    const primaryColor = pageSettings.primaryColor;
    const buttonColor = pageSettings.buttonColor;
    
    let content = `
        <div style="background: ${bgColor}; min-height: 400px; padding: 20px;">
            <h2 style="color: ${primaryColor}; text-align: center; margin-bottom: 10px; font-size: 18px;">${pageSettings.title || '活动标题'}</h2>
            <p style="color: ${primaryColor}; opacity: 0.8; text-align: center; font-size: 13px; margin-bottom: 20px;">${pageSettings.description || '活动描述'}</p>
    `;
    
    pageComponents.forEach(comp => {
        switch (comp.type) {
            case 'banner':
                content += `<div style="margin-bottom: 15px;"><img src="${comp.config?.imageUrl || 'https://via.placeholder.com/300x120?text=横幅图片'}" style="width: 100%; border-radius: 8px;" onerror="this.src='https://via.placeholder.com/300x120?text=图片加载失败'"></div>`;
                break;
            case 'text':
                content += `<div style="margin-bottom: 15px; padding: 10px; background: rgba(255,255,255,0.5); border-radius: 8px; text-align: ${comp.config?.align || 'left'}; color: ${primaryColor}; font-size: 14px;">${comp.config?.content || '文本内容'}</div>`;
                break;
            case 'video':
                content += `<div style="margin-bottom: 15px; background: #000; border-radius: 8px; height: 150px; display: flex; align-items: center; justify-content: center; color: white;">🎬 视频播放器</div>`;
                break;
            case 'countdown':
                content += `<div style="margin-bottom: 15px; padding: 15px; background: rgba(255,255,255,0.5); border-radius: 8px; text-align: center;">
                    <div style="font-size: 12px; color: #666; margin-bottom: 8px;">活动倒计时</div>
                    <div style="display: flex; justify-content: center; gap: 10px;">
                        <span style="background: ${primaryColor}; color: white; padding: 8px 12px; border-radius: 6px; font-weight: bold;">10天</span>
                        <span style="background: ${primaryColor}; color: white; padding: 8px 12px; border-radius: 6px; font-weight: bold;">12时</span>
                        <span style="background: ${primaryColor}; color: white; padding: 8px 12px; border-radius: 6px; font-weight: bold;">30分</span>
                    </div>
                </div>`;
                break;
            case 'button':
                content += `<div style="margin-bottom: 15px; text-align: center;"><button style="background: ${buttonColor}; color: white; border: none; padding: 12px 40px; border-radius: 25px; font-size: 16px; font-weight: 600;">${comp.config?.text || '立即参与'}</button></div>`;
                break;
            case 'divider':
                content += `<hr style="margin: 15px 0; border: none; border-top: 1px ${comp.config?.style || 'solid'} #ddd;">`;
                break;
        }
    });
    
    // 默认报名表单
    content += `
        <div style="margin-top: 20px; padding: 15px; background: rgba(255,255,255,0.8); border-radius: 12px;">
            <div style="margin-bottom: 12px;">
                <label style="font-size: 13px; color: #333; display: block; margin-bottom: 5px;">姓名 *</label>
                <input type="text" placeholder="请输入姓名" style="width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 6px;">
            </div>
            <div style="margin-bottom: 15px;">
                <label style="font-size: 13px; color: #333; display: block; margin-bottom: 5px;">手机号 *</label>
                <input type="tel" placeholder="请输入手机号" style="width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 6px;">
            </div>
            <button style="width: 100%; background: ${buttonColor}; color: white; border: none; padding: 14px; border-radius: 8px; font-size: 16px; font-weight: 600;">立即报名</button>
        </div>
    `;
    
    content += '</div>';
    return content;
}


// 添加组件
function addComponent(type) {
    const defaultConfigs = {
        banner: { imageUrl: '' },
        text: { content: '', align: 'left' },
        video: { videoUrl: '' },
        countdown: { endTime: '' },
        button: { text: '立即参与', link: '' },
        divider: { style: 'solid' }
    };
    
    pageComponents.push({
        id: Date.now(),
        type: type,
        config: defaultConfigs[type] || {}
    });
    
    refreshPageDesign();
}

// 移动组件
function moveComponent(index, direction) {
    const newIndex = index + direction;
    if (newIndex >= 0 && newIndex < pageComponents.length) {
        const temp = pageComponents[index];
        pageComponents[index] = pageComponents[newIndex];
        pageComponents[newIndex] = temp;
        refreshPageDesign();
    }
}

// 删除组件
function removeComponent(index) {
    pageComponents.splice(index, 1);
    refreshPageDesign();
}

// 更新组件配置
function updateComponentConfig(index, key, value) {
    if (!pageComponents[index].config) {
        pageComponents[index].config = {};
    }
    pageComponents[index].config[key] = value;
    refreshPageDesign();
}

// 更新页面设置
function updatePageSetting(key, value) {
    pageSettings[key] = value;
    refreshPageDesign();
}

// 刷新页面设计
function refreshPageDesign() {
    document.getElementById('componentList').innerHTML = renderComponentList();
    document.getElementById('pagePreview').innerHTML = renderPagePreview();
}

// 预览页面
function previewPage() {
    const previewWindow = window.open('', '_blank');
    previewWindow.document.write(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>${pageSettings.title || '活动页面'}</title>
            <style>* { margin: 0; padding: 0; box-sizing: border-box; } body { font-family: -apple-system, sans-serif; }</style>
        </head>
        <body>${renderPagePreview()}</body>
        </html>
    `);
}

// 保存页面设计
async function savePageDesign() {
    try {
        // 这里应该调用API保存页面配置
        // 暂时保存到本地
        if (currentCampaign) {
            currentCampaign.pageComponents = pageComponents;
            currentCampaign.pageSettings = pageSettings;
        }
        
        alert('页面设计保存成功！');
        closeModal('pageDesignModal');
    } catch (error) {
        alert('保存失败: ' + error.message);
    }
}