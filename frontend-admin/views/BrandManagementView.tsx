import { defineComponent, h, ref, onMounted, reactive, computed } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { PermissionGuard, usePermission } from '../components/PermissionGuard';

// 品牌管理视图
export const BrandManagementView = defineComponent({
  setup() {
    const { hasPermission, canAccessBrand, user } = usePermission();
    const loading = ref(false);
    const activeTab = ref<'brands' | 'assets' | 'stats' | 'relations'>('brands');
    
    // 品牌数据
    const brands = ref<any[]>([]);
    const selectedBrand = ref<any>(null);
    const brandAssets = ref<any[]>([]);
    const brandStats = ref<any>(null);
    const brandRelations = ref<any[]>([]);
    
    // 编辑品牌弹窗
    const showBrandDialog = ref(false);
    const editingBrand = ref<any>(null);
    const brandForm = reactive({
      name: '',
      logo: '',
      description: '',
      status: 'active'
    });
    
    // 素材管理弹窗
    const showAssetDialog = ref(false);
    const editingAsset = ref<any>(null);
    const assetForm = reactive({
      name: '',
      type: 'image',
      category: '',
      tags: '',
      fileUrl: '',
      fileSize: 0,
      description: ''
    });
    
    // 品牌关系管理弹窗
    const showRelationDialog = ref(false);
    const availableUsers = ref<any[]>([]);
    const selectedUsers = ref<number[]>([]);

    // 筛选的品牌列表（根据用户权限）
    const filteredBrands = computed(() => {
      if (!user.value) return [];
      
      // 平台管理员可以看到所有品牌
      if (user.value.roles.includes('platform_admin')) {
        return brands.value;
      }
      
      // 其他角色无权限查看品牌
      return [];
      }
      
      return [];
    });

    // 加载品牌列表
    const loadBrands = async () => {
      loading.value = true;
      try {
        // TODO: 调用真实API
        // const response = await fetch('/api/v1/brands', {
        //   headers: { 'Authorization': `Bearer ${localStorage.getItem('dmh_token')}` }
        // });
        // brands.value = await response.json();
        
        // 模拟数据
        brands.value = [
          {
            id: 1,
            name: '品牌A',
            logo: 'https://via.placeholder.com/150',
            description: '这是品牌A的描述',
            status: 'active',
            createdAt: '2025-01-01 10:00:00',
            updatedAt: '2025-01-01 10:00:00'
          },
          {
            id: 2,
            name: '品牌B',
            logo: 'https://via.placeholder.com/150',
            description: '这是品牌B的描述',
            status: 'active',
            createdAt: '2025-01-01 10:00:00',
            updatedAt: '2025-01-01 10:00:00'
          },
          {
            id: 3,
            name: '品牌C',
            logo: 'https://via.placeholder.com/150',
            description: '这是品牌C的描述',
            status: 'disabled',
            createdAt: '2025-01-01 10:00:00',
            updatedAt: '2025-01-01 10:00:00'
          }
        ];
        
        // 默认选择第一个品牌
        if (filteredBrands.value.length > 0) {
          selectedBrand.value = filteredBrands.value[0];
          loadBrandAssets();
          loadBrandStats();
        }
      } catch (error) {
        console.error('加载品牌列表失败', error);
      } finally {
        loading.value = false;
      }
    };

    // 加载品牌素材
    const loadBrandAssets = async () => {
      if (!selectedBrand.value) return;
      
      try {
        // TODO: 调用真实API
        // const response = await fetch(`/api/v1/brands/${selectedBrand.value.id}/assets`, {
        //   headers: { 'Authorization': `Bearer ${localStorage.getItem('dmh_token')}` }
        // });
        // brandAssets.value = await response.json();
        
        // 模拟数据
        brandAssets.value = [
          {
            id: 1,
            brandId: selectedBrand.value.id,
            name: '宣传海报1',
            type: 'image',
            category: '海报',
            tags: '促销,限时',
            fileUrl: 'https://via.placeholder.com/400x600',
            fileSize: 1024000,
            description: '春季促销活动海报',
            createdAt: '2025-01-01 10:00:00',
            updatedAt: '2025-01-01 10:00:00'
          },
          {
            id: 2,
            brandId: selectedBrand.value.id,
            name: '产品介绍视频',
            type: 'video',
            category: '视频',
            tags: '产品,介绍',
            fileUrl: 'https://example.com/video.mp4',
            fileSize: 5120000,
            description: '产品功能介绍视频',
            createdAt: '2025-01-01 10:00:00',
            updatedAt: '2025-01-01 10:00:00'
          }
        ];
      } catch (error) {
        console.error('加载品牌素材失败', error);
      }
    };

    // 加载品牌统计
    const loadBrandStats = async () => {
      if (!selectedBrand.value) return;
      
      try {
        // TODO: 调用真实API
        // const response = await fetch(`/api/v1/brands/${selectedBrand.value.id}/stats`, {
        //   headers: { 'Authorization': `Bearer ${localStorage.getItem('dmh_token')}` }
        // });
        // brandStats.value = await response.json();
        
        // 模拟数据
        brandStats.value = {
          brandId: selectedBrand.value.id,
          totalCampaigns: 12,
          activeCampaigns: 8,
          totalOrders: 1284,
          totalRevenue: 42050.00,
          totalRewards: 8410.00,
          participantCount: 856,
          conversionRate: 68.5,
          lastUpdated: '2025-01-02 15:30:00'
        };
      } catch (error) {
        console.error('加载品牌统计失败', error);
      }
    };

    // 加载品牌关系
    const loadBrandRelations = async () => {
      try {
        // TODO: 调用真实API
        // 模拟数据
        brandRelations.value = [
          {
            id: 1,
            userId: 2,
            username: 'brand_manager',
            realName: '品牌经理',
            brandIds: [1, 2],
            createdAt: '2025-01-01 10:00:00'
          }
        ];
        
        availableUsers.value = [
          { id: 4, username: 'user002', realName: '李四', roles: ['participant'] }
        ];
      } catch (error) {
        console.error('加载品牌关系失败', error);
      }
    };

    // 打开品牌编辑对话框
    const openBrandDialog = (brand?: any) => {
      if (brand) {
        editingBrand.value = brand;
        brandForm.name = brand.name;
        brandForm.logo = brand.logo;
        brandForm.description = brand.description;
        brandForm.status = brand.status;
      } else {
        editingBrand.value = null;
        brandForm.name = '';
        brandForm.logo = '';
        brandForm.description = '';
        brandForm.status = 'active';
      }
      showBrandDialog.value = true;
    };

    // 保存品牌
    const saveBrand = async () => {
      try {
        if (editingBrand.value) {
          // 更新品牌
          // TODO: 调用真实API
          Object.assign(editingBrand.value, brandForm);
        } else {
          // 创建品牌
          // TODO: 调用真实API
          const newBrand = {
            id: Date.now(),
            ...brandForm,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
          };
          brands.value.push(newBrand);
        }
        
        showBrandDialog.value = false;
        alert('品牌保存成功');
      } catch (error) {
        console.error('保存品牌失败', error);
        alert('保存失败，请重试');
      }
    };

    // 打开素材编辑对话框
    const openAssetDialog = (asset?: any) => {
      if (asset) {
        editingAsset.value = asset;
        assetForm.name = asset.name;
        assetForm.type = asset.type;
        assetForm.category = asset.category;
        assetForm.tags = asset.tags;
        assetForm.fileUrl = asset.fileUrl;
        assetForm.fileSize = asset.fileSize;
        assetForm.description = asset.description;
      } else {
        editingAsset.value = null;
        assetForm.name = '';
        assetForm.type = 'image';
        assetForm.category = '';
        assetForm.tags = '';
        assetForm.fileUrl = '';
        assetForm.fileSize = 0;
        assetForm.description = '';
      }
      showAssetDialog.value = true;
    };

    // 保存素材
    const saveAsset = async () => {
      try {
        if (editingAsset.value) {
          // 更新素材
          Object.assign(editingAsset.value, assetForm);
        } else {
          // 创建素材
          const newAsset = {
            id: Date.now(),
            brandId: selectedBrand.value.id,
            ...assetForm,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
          };
          brandAssets.value.push(newAsset);
        }
        
        showAssetDialog.value = false;
        alert('素材保存成功');
      } catch (error) {
        console.error('保存素材失败', error);
        alert('保存失败，请重试');
      }
    };

    // 删除素材
    const deleteAsset = async (asset: any) => {
      if (!confirm(`确定要删除素材"${asset.name}"吗？`)) return;
      
      try {
        // TODO: 调用真实API
        brandAssets.value = brandAssets.value.filter(a => a.id !== asset.id);
        alert('素材删除成功');
      } catch (error) {
        console.error('删除素材失败', error);
        alert('删除失败，请重试');
      }
    };

    // 格式化文件大小
    const formatFileSize = (bytes: number) => {
      if (bytes === 0) return '0 B';
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    // 获取素材类型图标
    const getAssetTypeIcon = (type: string) => {
      const icons: Record<string, any> = {
        image: LucideIcons.Image,
        video: LucideIcons.Video,
        document: LucideIcons.FileText,
      };
      return icons[type] || LucideIcons.File;
    };

    onMounted(() => {
      loadBrands();
      loadBrandRelations();
    });

    return () => h(PermissionGuard, { roles: ['platform_admin'] }, () => [
      h('div', { class: 'space-y-8 animate-in fade-in' }, [
        // 页面标题
        h('div', { class: 'flex justify-between items-end' }, [
          h('div', [
            h('h2', { class: 'text-4xl font-black text-slate-900' }, '品牌管理'),
            h('p', { class: 'text-slate-400 mt-2' }, '管理品牌信息、素材库、数据统计和管理员关系')
          ]),
          h(PermissionGuard, { permission: 'brand:create' }, () => [
            h('button', {
              onClick: () => openBrandDialog(),
              class: 'bg-slate-900 text-white px-6 py-3 rounded-2xl font-bold shadow-lg flex items-center gap-2 hover:bg-slate-800 transition-colors'
            }, [
              h(LucideIcons.Plus, { size: 18 }),
              '创建品牌'
            ])
          ])
        ]),

        // 品牌选择器
        filteredBrands.value.length > 0 && h('div', { class: 'bg-white rounded-3xl p-6 border border-slate-100' }, [
          h('div', { class: 'flex items-center gap-4 mb-4' }, [
            h(LucideIcons.Building, { size: 20, class: 'text-slate-600' }),
            h('span', { class: 'font-bold text-slate-900' }, '选择品牌')
          ]),
          h('div', { class: 'grid grid-cols-1 md:grid-cols-4 gap-4' }, 
            filteredBrands.value.map(brand => 
              h('button', {
                onClick: () => {
                  selectedBrand.value = brand;
                  loadBrandAssets();
                  loadBrandStats();
                },
                class: `p-4 rounded-2xl border-2 transition-all text-left ${
                  selectedBrand.value?.id === brand.id
                    ? 'border-indigo-500 bg-indigo-50'
                    : 'border-slate-200 hover:border-indigo-200'
                }`
              }, [
                h('div', { class: 'flex items-center gap-3' }, [
                  h('img', { 
                    src: brand.logo, 
                    class: 'w-12 h-12 rounded-xl object-cover',
                    alt: brand.name
                  }),
                  h('div', [
                    h('div', { class: 'font-bold text-slate-900' }, brand.name),
                    h('div', { class: `text-xs px-2 py-1 rounded-full mt-1 ${
                      brand.status === 'active' 
                        ? 'bg-green-100 text-green-700' 
                        : 'bg-red-100 text-red-700'
                    }` }, brand.status === 'active' ? '正常' : '已禁用')
                  ])
                ])
              ])
            )
          )
        ]),

        // 标签页导航
        selectedBrand.value && h('div', { class: 'bg-white rounded-3xl p-2 border border-slate-100' }, [
          h('div', { class: 'flex gap-2' }, [
            h('button', {
              onClick: () => activeTab.value = 'brands',
              class: `flex-1 py-3 px-6 rounded-2xl font-bold text-sm transition-all ${
                activeTab.value === 'brands'
                  ? 'bg-indigo-600 text-white shadow-lg'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
              }`
            }, [
              h(LucideIcons.Building, { size: 18, class: 'inline mr-2' }),
              '品牌信息'
            ]),
            h('button', {
              onClick: () => activeTab.value = 'assets',
              class: `flex-1 py-3 px-6 rounded-2xl font-bold text-sm transition-all ${
                activeTab.value === 'assets'
                  ? 'bg-indigo-600 text-white shadow-lg'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
              }`
            }, [
              h(LucideIcons.Image, { size: 18, class: 'inline mr-2' }),
              '素材库'
            ]),
            h('button', {
              onClick: () => activeTab.value = 'stats',
              class: `flex-1 py-3 px-6 rounded-2xl font-bold text-sm transition-all ${
                activeTab.value === 'stats'
                  ? 'bg-indigo-600 text-white shadow-lg'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
              }`
            }, [
              h(LucideIcons.BarChart3, { size: 18, class: 'inline mr-2' }),
              '数据统计'
            ])
          ])
        ]),

        // 品牌信息标签页
        activeTab.value === 'brands' && selectedBrand.value && h('div', { class: 'bg-white rounded-3xl border border-slate-100 p-8' }, [
          h('div', { class: 'flex items-start justify-between mb-8' }, [
            h('div', { class: 'flex items-center gap-6' }, [
              h('img', { 
                src: selectedBrand.value.logo, 
                class: 'w-24 h-24 rounded-3xl object-cover shadow-lg',
                alt: selectedBrand.value.name
              }),
              h('div', [
                h('h3', { class: 'text-3xl font-black text-slate-900' }, selectedBrand.value.name),
                h('p', { class: 'text-slate-600 mt-2' }, selectedBrand.value.description),
                h('div', { class: 'flex items-center gap-4 mt-4' }, [
                  h('div', { class: `px-3 py-1 rounded-full text-sm font-bold ${
                    selectedBrand.value.status === 'active' 
                      ? 'bg-green-100 text-green-700' 
                      : 'bg-red-100 text-red-700'
                  }` }, selectedBrand.value.status === 'active' ? '正常运营' : '已禁用'),
                  h('span', { class: 'text-sm text-slate-500' }, `创建于 ${selectedBrand.value.createdAt}`)
                ])
              ])
            ]),
            h(PermissionGuard, { permission: 'brand:update', brandId: selectedBrand.value.id }, () => [
              h('button', {
                onClick: () => openBrandDialog(selectedBrand.value),
                class: 'px-6 py-3 rounded-2xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors flex items-center gap-2'
              }, [
                h(LucideIcons.Edit, { size: 18 }),
                '编辑品牌'
              ])
            ])
          ])
        ]),

        // 素材库标签页
        activeTab.value === 'assets' && selectedBrand.value && h('div', { class: 'space-y-6' }, [
          h('div', { class: 'flex justify-between items-center' }, [
            h('div', [
              h('h3', { class: 'text-2xl font-black text-slate-900' }, '品牌素材库'),
              h('p', { class: 'text-slate-500 mt-1' }, `${selectedBrand.value.name} 的素材资源管理`)
            ]),
            h(PermissionGuard, { permission: 'asset:create', brandId: selectedBrand.value.id }, () => [
              h('button', {
                onClick: () => openAssetDialog(),
                class: 'bg-indigo-600 text-white px-6 py-3 rounded-2xl font-bold shadow-lg flex items-center gap-2 hover:bg-indigo-500 transition-colors'
              }, [
                h(LucideIcons.Plus, { size: 18 }),
                '上传素材'
              ])
            ])
          ]),
          
          brandAssets.value.length === 0
            ? h('div', { class: 'text-center py-20 text-slate-400' }, '暂无素材，点击"上传素材"开始')
            : h('div', { class: 'grid grid-cols-1 md:grid-cols-3 gap-6' }, 
                brandAssets.value.map(asset => 
                  h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden hover:shadow-lg transition-all' }, [
                    h('div', { class: 'aspect-video bg-slate-100 relative overflow-hidden' }, [
                      asset.type === 'image' 
                        ? h('img', { src: asset.fileUrl, class: 'w-full h-full object-cover' })
                        : h('div', { class: 'w-full h-full flex items-center justify-center' }, [
                            h(getAssetTypeIcon(asset.type), { size: 48, class: 'text-slate-400' })
                          ]),
                      h('div', { class: 'absolute top-3 right-3 flex gap-2' }, [
                        h(PermissionGuard, { permission: 'asset:update', brandId: selectedBrand.value.id }, () => [
                          h('button', {
                            onClick: () => openAssetDialog(asset),
                            class: 'p-2 bg-white/90 backdrop-blur rounded-xl hover:bg-white transition-colors shadow-sm'
                          }, h(LucideIcons.Edit, { size: 16, class: 'text-slate-600' }))
                        ]),
                        h(PermissionGuard, { permission: 'asset:delete', brandId: selectedBrand.value.id }, () => [
                          h('button', {
                            onClick: () => deleteAsset(asset),
                            class: 'p-2 bg-white/90 backdrop-blur rounded-xl hover:bg-white transition-colors shadow-sm'
                          }, h(LucideIcons.Trash2, { size: 16, class: 'text-red-600' }))
                        ])
                      ])
                    ]),
                    h('div', { class: 'p-6' }, [
                      h('h4', { class: 'font-bold text-slate-900 mb-2' }, asset.name),
                      h('p', { class: 'text-sm text-slate-600 mb-4' }, asset.description),
                      h('div', { class: 'flex items-center justify-between text-xs text-slate-500' }, [
                        h('span', asset.category),
                        h('span', formatFileSize(asset.fileSize))
                      ]),
                      asset.tags && h('div', { class: 'flex flex-wrap gap-1 mt-3' }, 
                        asset.tags.split(',').map((tag: string) => 
                          h('span', { class: 'px-2 py-1 bg-slate-100 text-slate-600 rounded-lg text-xs' }, tag.trim())
                        )
                      )
                    ])
                  ])
                )
              )
        ]),

        // 数据统计标签页
        activeTab.value === 'stats' && selectedBrand.value && brandStats.value && h('div', { class: 'space-y-6' }, [
          h('div', [
            h('h3', { class: 'text-2xl font-black text-slate-900' }, '品牌数据统计'),
            h('p', { class: 'text-slate-500 mt-1' }, `${selectedBrand.value.name} 的运营数据概览`)
          ]),
          
          // 统计卡片
          h('div', { class: 'grid grid-cols-1 md:grid-cols-4 gap-6' }, [
            {
              label: '总活动数',
              value: brandStats.value.totalCampaigns,
              icon: LucideIcons.Target,
              color: 'bg-blue-600'
            },
            {
              label: '活跃活动',
              value: brandStats.value.activeCampaigns,
              icon: LucideIcons.Activity,
              color: 'bg-green-600'
            },
            {
              label: '总订单数',
              value: brandStats.value.totalOrders.toLocaleString(),
              icon: LucideIcons.ShoppingCart,
              color: 'bg-purple-600'
            },
            {
              label: '参与用户',
              value: brandStats.value.participantCount.toLocaleString(),
              icon: LucideIcons.Users,
              color: 'bg-indigo-600'
            }
          ].map(stat => 
            h('div', { class: 'bg-white p-8 rounded-3xl border border-slate-100 shadow-sm' }, [
              h('div', { class: `w-12 h-12 ${stat.color} text-white rounded-2xl flex items-center justify-center mb-6` }, 
                h(stat.icon, { size: 24 })
              ),
              h('p', { class: 'text-xs font-black text-slate-400 uppercase tracking-widest' }, stat.label),
              h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, stat.value)
            ])
          )),
          
          // 收入统计
          h('div', { class: 'grid grid-cols-1 md:grid-cols-3 gap-6' }, [
            h('div', { class: 'bg-white p-8 rounded-3xl border border-slate-100 shadow-sm' }, [
              h('div', { class: 'w-12 h-12 bg-emerald-600 text-white rounded-2xl flex items-center justify-center mb-6' }, 
                h(LucideIcons.DollarSign, { size: 24 })
              ),
              h('p', { class: 'text-xs font-black text-slate-400 uppercase tracking-widest' }, '总收入'),
              h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, `¥${brandStats.value.totalRevenue.toLocaleString()}`)
            ]),
            h('div', { class: 'bg-white p-8 rounded-3xl border border-slate-100 shadow-sm' }, [
              h('div', { class: 'w-12 h-12 bg-amber-600 text-white rounded-2xl flex items-center justify-center mb-6' }, 
                h(LucideIcons.Gift, { size: 24 })
              ),
              h('p', { class: 'text-xs font-black text-slate-400 uppercase tracking-widest' }, '奖励发放'),
              h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, `¥${brandStats.value.totalRewards.toLocaleString()}`)
            ]),
            h('div', { class: 'bg-white p-8 rounded-3xl border border-slate-100 shadow-sm' }, [
              h('div', { class: 'w-12 h-12 bg-rose-600 text-white rounded-2xl flex items-center justify-center mb-6' }, 
                h(LucideIcons.TrendingUp, { size: 24 })
              ),
              h('p', { class: 'text-xs font-black text-slate-400 uppercase tracking-widest' }, '转化率'),
              h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, `${brandStats.value.conversionRate}%`)
            ])
          ]),
          
          h('div', { class: 'bg-slate-50 rounded-2xl p-4 text-sm text-slate-600' }, [
            `数据更新时间：${brandStats.value.lastUpdated}`
          ])
        ]),

        // 品牌编辑对话框
        showBrandDialog.value && h('div', { 
          class: 'fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4',
          onClick: () => showBrandDialog.value = false
        }, [
          h('div', { 
            class: 'bg-white rounded-3xl p-8 max-w-2xl w-full max-h-[80vh] overflow-auto',
            onClick: (e: Event) => e.stopPropagation()
          }, [
            h('div', { class: 'flex items-center justify-between mb-6' }, [
              h('h3', { class: 'text-2xl font-black text-slate-900' }, 
                editingBrand.value ? '编辑品牌' : '创建品牌'
              ),
              h('button', {
                onClick: () => showBrandDialog.value = false,
                class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors'
              }, h(LucideIcons.X, { size: 20 }))
            ]),
            
            h('div', { class: 'space-y-6' }, [
              h('div', [
                h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '品牌名称'),
                h('input', {
                  type: 'text',
                  value: brandForm.name,
                  onInput: (e: any) => brandForm.name = e.target.value,
                  placeholder: '请输入品牌名称',
                  class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                })
              ]),
              
              h('div', [
                h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '品牌Logo'),
                h('input', {
                  type: 'url',
                  value: brandForm.logo,
                  onInput: (e: any) => brandForm.logo = e.target.value,
                  placeholder: '请输入Logo URL',
                  class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                })
              ]),
              
              h('div', [
                h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '品牌描述'),
                h('textarea', {
                  value: brandForm.description,
                  onInput: (e: any) => brandForm.description = e.target.value,
                  placeholder: '请输入品牌描述',
                  class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors h-32 resize-none'
                })
              ]),
              
              editingBrand.value && h('div', [
                h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '品牌状态'),
                h('select', {
                  value: brandForm.status,
                  onChange: (e: any) => brandForm.status = e.target.value,
                  class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                }, [
                  h('option', { value: 'active' }, '正常'),
                  h('option', { value: 'disabled' }, '禁用')
                ])
              ])
            ]),
            
            h('div', { class: 'flex gap-3 mt-8' }, [
              h('button', {
                onClick: () => showBrandDialog.value = false,
                class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
              }, '取消'),
              h('button', {
                onClick: saveBrand,
                class: 'flex-1 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-500 transition-colors'
              }, '保存')
            ])
          ])
        ]),

        // 素材编辑对话框
        showAssetDialog.value && h('div', { 
          class: 'fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4',
          onClick: () => showAssetDialog.value = false
        }, [
          h('div', { 
            class: 'bg-white rounded-3xl p-8 max-w-2xl w-full max-h-[80vh] overflow-auto',
            onClick: (e: Event) => e.stopPropagation()
          }, [
            h('div', { class: 'flex items-center justify-between mb-6' }, [
              h('h3', { class: 'text-2xl font-black text-slate-900' }, 
                editingAsset.value ? '编辑素材' : '上传素材'
              ),
              h('button', {
                onClick: () => showAssetDialog.value = false,
                class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors'
              }, h(LucideIcons.X, { size: 20 }))
            ]),
            
            h('div', { class: 'space-y-6' }, [
              h('div', { class: 'grid grid-cols-2 gap-4' }, [
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '素材名称'),
                  h('input', {
                    type: 'text',
                    value: assetForm.name,
                    onInput: (e: any) => assetForm.name = e.target.value,
                    placeholder: '请输入素材名称',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '素材类型'),
                  h('select', {
                    value: assetForm.type,
                    onChange: (e: any) => assetForm.type = e.target.value,
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  }, [
                    h('option', { value: 'image' }, '图片'),
                    h('option', { value: 'video' }, '视频'),
                    h('option', { value: 'document' }, '文档')
                  ])
                ])
              ]),
              
              h('div', { class: 'grid grid-cols-2 gap-4' }, [
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '分类'),
                  h('input', {
                    type: 'text',
                    value: assetForm.category,
                    onInput: (e: any) => assetForm.category = e.target.value,
                    placeholder: '如：海报、视频、文档',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '标签'),
                  h('input', {
                    type: 'text',
                    value: assetForm.tags,
                    onInput: (e: any) => assetForm.tags = e.target.value,
                    placeholder: '用逗号分隔，如：促销,限时',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ])
              ]),
              
              h('div', [
                h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '文件URL'),
                h('input', {
                  type: 'url',
                  value: assetForm.fileUrl,
                  onInput: (e: any) => assetForm.fileUrl = e.target.value,
                  placeholder: '请输入文件URL',
                  class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                })
              ]),
              
              h('div', [
                h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '描述'),
                h('textarea', {
                  value: assetForm.description,
                  onInput: (e: any) => assetForm.description = e.target.value,
                  placeholder: '请输入素材描述',
                  class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors h-32 resize-none'
                })
              ])
            ]),
            
            h('div', { class: 'flex gap-3 mt-8' }, [
              h('button', {
                onClick: () => showAssetDialog.value = false,
                class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
              }, '取消'),
              h('button', {
                onClick: saveAsset,
                class: 'flex-1 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-500 transition-colors'
              }, '保存')
            ])
          ])
        ])
      ])
    ]);
  }
});