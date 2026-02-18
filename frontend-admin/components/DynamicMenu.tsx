import { defineComponent, h, ref, onMounted, computed, type DefineComponent } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { PermissionGuard, usePermission } from './PermissionGuard';
import { menuApi } from '../services/menuApi';

// 组件元数据类型
export interface ComponentWithProps {
  props: Record<string, any>;
  emits?: string[];
}

// 菜单项接口
interface MenuItem {
  id: number;
  name: string;
  code: string;
  path: string;
  icon?: string;
  parentId?: number;
  sort: number;
  type: 'menu' | 'button';
  platform: 'admin' | 'h5';
  status: 'active' | 'disabled';
  permission?: string;
  role?: string;
  children?: MenuItem[];
}

const fallbackMenus: MenuItem[] = [
  { id: 1, name: '仪表板', code: 'dashboard', path: '/dashboard', icon: 'LayoutDashboard', parentId: undefined, sort: 1, type: 'menu', platform: 'admin', status: 'active' },
  { id: 2, name: '用户管理', code: 'user-management', path: '/users', icon: 'Users', parentId: undefined, sort: 2, type: 'menu', platform: 'admin', status: 'active', permission: 'user:read' },
  { id: 3, name: '品牌管理', code: 'brand-management', path: '/brands', icon: 'Building', parentId: undefined, sort: 3, type: 'menu', platform: 'admin', status: 'active', permission: 'brand:read' },
  { id: 4, name: '活动管理', code: 'campaign-management', path: '/campaigns', icon: 'Target', parentId: undefined, sort: 4, type: 'menu', platform: 'admin', status: 'active', permission: 'campaign:read' },
  { id: 5, name: '订单管理', code: 'order-management', path: '/orders', icon: 'ShoppingCart', parentId: undefined, sort: 5, type: 'menu', platform: 'admin', status: 'active', permission: 'order:read' },
  { id: 6, name: '奖励管理', code: 'reward-management', path: '/rewards', icon: 'Gift', parentId: undefined, sort: 6, type: 'menu', platform: 'admin', status: 'active', permission: 'reward:read' },
  { id: 7, name: '提现管理', code: 'withdrawal-management', path: '/withdrawals', icon: 'Wallet', parentId: undefined, sort: 7, type: 'menu', platform: 'admin', status: 'active', permission: 'withdrawal:approve' },
  { id: 8, name: '系统管理', code: 'system-management', path: '/system', icon: 'Settings', parentId: undefined, sort: 8, type: 'menu', platform: 'admin', status: 'active', role: 'platform_admin' },
  { id: 9, name: '用户列表', code: 'user-list', path: '/users/list', parentId: 2, sort: 1, type: 'menu', platform: 'admin', status: 'active', permission: 'user:read' },
  { id: 10, name: '创建用户', code: 'user-create', path: '/users/create', parentId: 2, sort: 2, type: 'menu', platform: 'admin', status: 'active', permission: 'user:create' },
  { id: 11, name: '品牌列表', code: 'brand-list', path: '/brands/list', parentId: 3, sort: 1, type: 'menu', platform: 'admin', status: 'active', permission: 'brand:read' },
  { id: 12, name: '创建品牌', code: 'brand-create', path: '/brands/create', parentId: 3, sort: 2, type: 'menu', platform: 'admin', status: 'active', permission: 'brand:create' },
  { id: 13, name: '品牌关系管理', code: 'brand-relation', path: '/brands/relations', parentId: 3, sort: 3, type: 'menu', platform: 'admin', status: 'active', permission: 'user:update' },
  { id: 14, name: '活动列表', code: 'campaign-list', path: '/campaigns/list', parentId: 4, sort: 1, type: 'menu', platform: 'admin', status: 'active', permission: 'campaign:read' },
  { id: 15, name: '创建活动', code: 'campaign-create', path: '/campaigns/create', parentId: 4, sort: 2, type: 'menu', platform: 'admin', status: 'active', permission: 'campaign:create' },
  { id: 16, name: '页面配置', code: 'campaign-config', path: '/campaigns/config', parentId: 4, sort: 3, type: 'menu', platform: 'admin', status: 'active', permission: 'campaign:update' },
  { id: 17, name: '角色管理', code: 'role-management', path: '/system/roles', parentId: 8, sort: 1, type: 'menu', platform: 'admin', status: 'active', role: 'platform_admin' },
  { id: 18, name: '菜单管理', code: 'menu-management', path: '/system/menus', parentId: 8, sort: 2, type: 'menu', platform: 'admin', status: 'active', role: 'platform_admin' },
  { id: 19, name: '权限配置', code: 'permission-config', path: '/system/permissions', parentId: 8, sort: 3, type: 'menu', platform: 'admin', status: 'active', role: 'platform_admin' },
];

// 动态菜单组件
export const DynamicMenu = defineComponent({
  props: {
    platform: {
      type: String as () => 'admin' | 'h5',
      default: 'admin'
    },
    currentPath: {
      type: String,
      default: ''
    }
  },
  emits: ['navigate'],
  setup(props, { emit }) {
    const { user, hasPermission, hasRole } = usePermission();
    const menuItems = ref<MenuItem[]>([]);
    const loading = ref(false);

    // 加载菜单数据
    const loadMenus = async () => {
      loading.value = true;
      try {
        const menus = await menuApi.getUserMenus(props.platform);
        if (menus.length > 0) {
          menuItems.value = menus;
        } else {
          menuItems.value = fallbackMenus.filter(menu =>
            menu.platform === props.platform && menu.status === 'active'
          );
        }
      } catch (error) {
        console.error('加载菜单失败', error);
        menuItems.value = fallbackMenus.filter(menu =>
          menu.platform === props.platform && menu.status === 'active'
        );
      } finally {
        loading.value = false;
      }
    };

    // 检查菜单项权限
    const hasMenuAccess = (menu: MenuItem): boolean => {
      // 检查权限
      if (menu.permission && !hasPermission(menu.permission)) {
        return false;
      }
      
      // 检查角色
      if (menu.role && !hasRole(menu.role as any)) {
        return false;
      }
      
      return true;
    };

    // 构建菜单树结构
    const menuTree = computed(() => {
      const filteredMenus = menuItems.value.filter(hasMenuAccess);
      
      // 构建树形结构
      const buildTree = (parentId?: number): MenuItem[] => {
        return filteredMenus
          .filter(menu => menu.parentId === parentId)
          .sort((a, b) => a.sort - b.sort)
          .map(menu => ({
            ...menu,
            children: buildTree(menu.id)
          }));
      };
      
      return buildTree();
    });

    // 获取图标组件
    const getIcon = (iconName?: string) => {
      if (!iconName) return null;
      const IconComponent = (LucideIcons as any)[iconName];
      return IconComponent ? h(IconComponent, { size: 20 }) : null;
    };

    // 处理菜单点击
    const handleMenuClick = (menu: MenuItem) => {
      if (menu.type === 'menu' && menu.path) {
        emit('navigate', menu.path);
      }
    };

    // 检查菜单是否激活
    const isMenuActive = (menu: MenuItem): boolean => {
      if (menu.path === props.currentPath) {
        return true;
      }
      
      // 检查子菜单是否激活
      if (menu.children) {
        return menu.children.some(child => isMenuActive(child));
      }
      
      return false;
    };

    // 渲染菜单项
    const renderMenuItem = (menu: MenuItem, level = 0) => {
      const isActive = isMenuActive(menu);
      const hasChildren = menu.children && menu.children.length > 0;
      
      return h('div', { key: menu.id, class: 'menu-item' }, [
        h('div', {
          onClick: () => handleMenuClick(menu),
          class: `flex items-center gap-3 px-4 py-3 rounded-2xl cursor-pointer transition-all ${
            level === 0 ? 'mb-2' : 'mb-1 ml-4'
          } ${
            isActive 
              ? 'bg-indigo-600 text-white shadow-lg' 
              : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'
          }`
        }, [
          level === 0 && getIcon(menu.icon),
          h('span', { class: `font-medium ${level === 0 ? 'text-sm' : 'text-xs'}` }, menu.name),
          hasChildren && h('div', { class: 'ml-auto' }, [
            h(LucideIcons.ChevronRight, { 
              size: 16, 
              class: `transition-transform ${isActive ? 'rotate-90' : ''}` 
            })
          ])
        ]),
        
        // 子菜单
        hasChildren && isActive && h('div', { class: 'mt-2 space-y-1' }, 
          menu.children!.map(child => renderMenuItem(child, level + 1))
        )
      ]);
    };

    onMounted(() => {
      loadMenus();
    });

    return () => h('div', { class: 'dynamic-menu' }, [
      loading.value 
        ? h('div', { class: 'p-4 text-center text-slate-400' }, '加载菜单中...')
        : h('nav', { class: 'space-y-2' }, 
            menuTree.value.map(menu => renderMenuItem(menu))
          )
    ]);
  }
});

// 面包屑导航组件
export const Breadcrumb = defineComponent({
  props: {
    items: {
      type: Array as () => { name: string; path?: string }[],
      required: true
    }
  },
  emits: ['navigate'],
  setup(props, { emit }) {
    return () => h('nav', { class: 'flex items-center space-x-2 text-sm' }, [
      props.items.map((item, index) => [
        index > 0 && h(LucideIcons.ChevronRight, { size: 16, class: 'text-slate-400' }),
        item.path 
          ? h('button', {
              onClick: () => emit('navigate', item.path),
              class: 'text-indigo-600 hover:text-indigo-500 font-medium'
            }, item.name)
          : h('span', { class: 'text-slate-600 font-medium' }, item.name)
      ]).flat()
    ]);
  }
});

// 页面标题组件
export const PageHeader = defineComponent({
  props: {
    title: String,
    description: String,
    breadcrumb: Array as () => { name: string; path?: string }[]
  },
  emits: ['navigate'],
  setup(props, { emit, slots }) {
    return () => h('div', { class: 'mb-8' }, [
      // 面包屑
      props.breadcrumb && h('div', { class: 'mb-4' }, [
        h(Breadcrumb, { 
          items: props.breadcrumb,
          onNavigate: (path: string) => emit('navigate', path)
        })
      ]),
      
      // 标题区域
      h('div', { class: 'flex items-end justify-between' }, [
        h('div', [
          h('h1', { class: 'text-4xl font-black text-slate-900' }, props.title),
          props.description && h('p', { class: 'text-slate-400 mt-2' }, props.description)
        ]),
        
        // 操作按钮区域
        slots.actions && h('div', { class: 'flex items-center gap-3' }, slots.actions())
      ])
    ]);
  }
});
