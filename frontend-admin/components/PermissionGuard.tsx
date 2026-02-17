import { defineComponent, h, computed, inject, provide, ref, onMounted } from 'vue';
import { CurrentUser, UserRole } from '../types';
import { authApi } from '../services/authApi';

// 权限上下文
export const PermissionContext = Symbol('PermissionContext');

// 权限提供者组件
export const PermissionProvider = defineComponent({
  props: {
    user: {
      type: Object as () => CurrentUser | null,
      required: true
    }
  },
  setup(props, { slots }) {
    const hasPermission = (permission: string): boolean => {
      if (!props.user) return false;
      
      // 平台管理员拥有所有权限
      if (props.user.roles.includes('platform_admin')) {
        return true;
      }
      
      // 根据角色检查权限
      const rolePermissions: Record<UserRole, string[]> = {
        platform_admin: ['*'], // 所有权限
        participant: [
          'campaign:read',
          'order:create',
          'reward:read',
          'withdrawal:apply'
        ],
        anonymous: []
      };
      
      for (const role of props.user.roles) {
        const permissions = rolePermissions[role] || [];
        if (permissions.includes('*') || permissions.includes(permission)) {
          return true;
        }
      }
      
      return false;
    };

    const hasRole = (role: UserRole): boolean => {
      return props.user?.roles.includes(role) || false;
    };

    const hasAnyRole = (roles: UserRole[]): boolean => {
      return roles.some(role => hasRole(role));
    };

    const canAccessBrand = (brandId: number): boolean => {
      if (!props.user) return false;
      
      // 平台管理员可以访问所有品牌
      if (hasRole('platform_admin')) {
        return true;
      }
      
      // 其他角色无品牌访问权限
      return false;
    };

    const permissionContext = {
      user: computed(() => props.user),
      hasPermission,
      hasRole,
      hasAnyRole,
      canAccessBrand
    };

    provide(PermissionContext, permissionContext);

    return () => slots.default?.();
  }
});

// 权限守卫组件
export const PermissionGuard = defineComponent({
  props: {
    permission: String,
    role: String as () => UserRole,
    roles: Array as () => UserRole[],
    brandId: Number,
    fallback: {
      type: [String, Object],
      default: null
    }
  },
  setup(props, { slots }) {
    const context = inject(PermissionContext) as any;
    
    if (!context) {
      console.warn('PermissionGuard must be used within PermissionProvider');
      return () => null;
    }

    const hasAccess = computed(() => {
      // 检查权限
      if (props.permission && !context.hasPermission(props.permission)) {
        return false;
      }
      
      // 检查单个角色
      if (props.role && !context.hasRole(props.role)) {
        return false;
      }
      
      // 检查多个角色（任一匹配即可）
      if (props.roles && !context.hasAnyRole(props.roles)) {
        return false;
      }
      
      // 检查品牌访问权限
      if (props.brandId && !context.canAccessBrand(props.brandId)) {
        return false;
      }
      
      return true;
    });

    return () => {
      if (hasAccess.value) {
        return slots.default?.();
      }
      
      if (props.fallback) {
        if (typeof props.fallback === 'string') {
          return h('div', { class: 'text-slate-400 text-sm' }, props.fallback);
        }
        return h(props.fallback);
      }
      
      return null;
    };
  }
});

// 权限按钮组件
export const PermissionButton = defineComponent({
  props: {
    permission: String,
    role: String as () => UserRole,
    roles: Array as () => UserRole[],
    brandId: Number,
    disabled: Boolean,
    class: String,
    onClick: Function
  },
  setup(props, { slots }) {
    const context = inject(PermissionContext) as any;
    
    const hasAccess = computed(() => {
      if (!context) return false;
      
      if (props.permission && !context.hasPermission(props.permission)) {
        return false;
      }
      
      if (props.role && !context.hasRole(props.role)) {
        return false;
      }
      
      if (props.roles && !context.hasAnyRole(props.roles)) {
        return false;
      }
      
      if (props.brandId && !context.canAccessBrand(props.brandId)) {
        return false;
      }
      
      return true;
    });

    const isDisabled = computed(() => {
      return props.disabled || !hasAccess.value;
    });

    return () => {
      if (!hasAccess.value) {
        return null; // 没有权限时不显示按钮
      }
      
      return h('button', {
        class: `${props.class} ${isDisabled.value ? 'opacity-50 cursor-not-allowed' : ''}`,
        disabled: isDisabled.value,
        onClick: isDisabled.value ? undefined : props.onClick
      }, slots.default?.());
    };
  }
});

// 使用权限的Hook
export const usePermission = () => {
  const context = inject(PermissionContext) as any;
  
  if (!context) {
    console.warn('usePermission must be used within PermissionProvider');
    return {
      user: ref(null),
      hasPermission: () => false,
      hasRole: () => false,
      hasAnyRole: () => false,
      canAccessBrand: () => false
    };
  }
  
  return context;
};

// 路由权限守卫
export const createRouteGuard = () => {
  return {
    beforeEnter: async (to: any, from: any, next: any) => {
      const token = authApi.getToken();
      
      if (!token) {
        // 未登录，重定向到登录页
        next('/login');
        return;
      }
      
      // 检查路由权限
      const requiredRole = to.meta?.role as string;
      const requiredPermission = to.meta?.permission as string;
      
      if (!requiredRole && !requiredPermission) {
        next();
        return;
      }

      try {
        const userResp = await fetch('/api/v1/auth/userinfo', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!userResp.ok) {
          throw new Error('获取用户信息失败');
        }

        const user = await userResp.json() as { id: number; roles?: string[] };
        const roles = Array.isArray(user.roles) ? user.roles : [];
        const isPlatformAdmin = roles.includes('platform_admin');

        if (requiredRole && !roles.includes(requiredRole)) {
          next('/dashboard');
          return;
        }

        if (!requiredPermission || isPlatformAdmin) {
          next();
          return;
        }

        const permResp = await fetch(`/api/v1/users/${user.id}/permissions`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!permResp.ok) {
          throw new Error('获取用户权限失败');
        }

        const permissionData = await permResp.json() as { permissions?: string[] };
        const permissions = Array.isArray(permissionData.permissions)
          ? permissionData.permissions
          : [];

        if (!permissions.includes(requiredPermission)) {
          next('/dashboard');
          return;
        }
      } catch (error) {
        console.error('路由权限校验失败', error);
        authApi.logout();
        next('/login');
        return;
      }
      
      next();
    }
  };
};
