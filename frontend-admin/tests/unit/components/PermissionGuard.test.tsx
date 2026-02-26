import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'
import { PermissionProvider, PermissionGuard, PermissionButton, usePermission } from '../../../components/PermissionGuard'
import type { UserRole, CurrentUser } from '../../../types'

const createTestUser = (roles: UserRole[], username: string = 'test'): CurrentUser => ({
  id: 1,
  username,
  phone: '13800138000',
  email: 'test@test.com',
  realName: 'Test User',
  avatar: '',
  status: 'active',
  roles
})

describe('PermissionGuard', () => {
  const mockUser = createTestUser(['platform_admin'], 'admin')

  describe('PermissionProvider', () => {
    it('provides permission context to children', () => {
      const Child = defineComponent({
        setup() {
          return () => h('div', 'child')
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.exists()).toBe(true)
    })

    it('hasPermission returns true for platform_admin', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasPermission('any:permission')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })

    it('hasPermission returns false for anonymous user', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasPermission('any:permission')))
        }
      })

      const anonymousUser = createTestUser(['anonymous'], 'guest')

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: anonymousUser },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })

    it('hasRole returns correct value', () => {
      const userWithRoles = createTestUser(['participant'])

      const Child = defineComponent({
        setup() {
          const { hasRole } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasRole('participant')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userWithRoles },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })

    it('hasRole returns false for non-matching role', () => {
      const userWithDifferentRole = createTestUser(['participant'])

      const Child = defineComponent({
        setup() {
          const { hasRole } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasRole('platform_admin')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userWithDifferentRole },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })

    it('canAccessBrand returns true for platform_admin', () => {
      const Child = defineComponent({
        setup() {
          const { canAccessBrand } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(canAccessBrand(123)))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })
  })

  describe('PermissionGuard', () => {
    it('renders children when has permission', () => {
      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionGuard,
          { permission: 'campaign:read' },
          () => h('div', { 'data-testid': 'content' }, 'Protected Content')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(true)
    })

    it('does not render children when no permission', () => {
      const userNoPerm = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionGuard,
          { permission: 'admin:delete' },
          () => h('div', { 'data-testid': 'content' }, 'Protected Content')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })

    it('renders fallback when provided and no permission', () => {
      const userNoPerm = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionGuard,
          { permission: 'admin:delete', fallback: 'No Access' },
          () => h('div', { 'data-testid': 'content' }, 'Protected Content')
        )
      ))

      expect(wrapper.text()).toContain('No Access')
      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })

    it('checks role correctly', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionGuard,
          { role: 'participant' },
          () => h('div', { 'data-testid': 'content' }, 'Has Role')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(true)
    })

    it('checks roles array correctly', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionGuard,
          { roles: ['participant', 'platform_admin'] as UserRole[] },
          () => h('div', { 'data-testid': 'content' }, 'Has Role')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(true)
    })
  })

  describe('PermissionButton', () => {
    it('renders button when has permission', () => {
      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionButton,
          { permission: 'campaign:read' },
          () => h('span', 'Click Me')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(true)
      expect(wrapper.text()).toContain('Click Me')
    })

    it('does not render button when no permission', () => {
      const userNoPerm = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionButton,
          { permission: 'admin:delete' },
          () => h('span', 'Click Me')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(false)
    })

    it('handles click when enabled', async () => {
      const handleClick = vi.fn()

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionButton,
          { permission: 'campaign:read', onClick: handleClick },
          () => h('span', 'Click Me')
        )
      ))

      await wrapper.find('button').trigger('click')
      expect(handleClick).toHaveBeenCalled()
    })

    it('button does not render when no permission', () => {
      const userNoPerm = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionButton,
          { permission: 'admin:delete' },
          () => h('span', 'Click Me')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(false)
    })
  })

  describe('usePermission', () => {
    it('returns context when used within provider', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', String(hasPermission('test')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.text()).toBe('true')
    })

    it('returns fallback when used outside provider', () => {
      const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})

      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', String(hasPermission('test')))
        }
      })

      const wrapper = mount(() => h(Child))

      expect(wrapper.text()).toBe('false')
      expect(warnSpy).toHaveBeenCalledWith('usePermission must be used within PermissionProvider')
      warnSpy.mockRestore()
    })
  })

  describe('PermissionGuard additional tests', () => {
    it('warns when used outside PermissionProvider', () => {
      const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})

      const wrapper = mount(() => h(
        PermissionGuard,
        { permission: 'test:permission' },
        () => h('div', 'Content')
      ))

      expect(warnSpy).toHaveBeenCalledWith('PermissionGuard must be used within PermissionProvider')
      expect(wrapper.html()).toBe('')
      warnSpy.mockRestore()
    })

    it('checks brandId correctly - no access for non-platform-admin', () => {
      const userBrandAdmin = createTestUser(['brand_admin'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userBrandAdmin },
        () => h(
          PermissionGuard,
          { brandId: 123 },
          () => h('div', { 'data-testid': 'content' }, 'Brand Content')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })

    it('renders object fallback when no permission', () => {
      const userNoPerm = createTestUser(['participant'])
      const FallbackComponent = defineComponent({
        setup() {
          return () => h('div', { 'data-testid': 'fallback' }, 'Fallback Component')
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionGuard,
          { permission: 'admin:delete', fallback: FallbackComponent },
          () => h('div', { 'data-testid': 'content' }, 'Protected Content')
        )
      ))

      expect(wrapper.find('[data-testid="fallback"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })

    it('denies access when role does not match', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionGuard,
          { role: 'platform_admin' },
          () => h('div', { 'data-testid': 'content' }, 'Has Role')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })

    it('denies access when none of roles match', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionGuard,
          { roles: ['platform_admin', 'brand_admin'] as UserRole[] },
          () => h('div', { 'data-testid': 'content' }, 'Has Role')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })
  })

  describe('PermissionButton additional tests', () => {
    it('renders disabled button when disabled prop is true', async () => {
      const handleClick = vi.fn()

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionButton,
          { permission: 'campaign:read', disabled: true, onClick: handleClick },
          () => h('span', 'Disabled Button')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(true)
      expect(wrapper.find('button').attributes('disabled')).toBeDefined()
      await wrapper.find('button').trigger('click')
      expect(handleClick).not.toHaveBeenCalled()
    })

    it('renders button with custom class', () => {
      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionButton,
          { permission: 'campaign:read', class: 'btn-primary' },
          () => h('span', 'Styled Button')
        )
      ))

      expect(wrapper.find('button').classes()).toContain('btn-primary')
    })

    it('does not render button when no context', () => {
      const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})

      const wrapper = mount(() => h(
        PermissionButton,
        { permission: 'campaign:read' },
        () => h('span', 'No Context')
      ))

      expect(wrapper.find('button').exists()).toBe(false)
      warnSpy.mockRestore()
    })

    it('checks role correctly', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionButton,
          { role: 'participant' },
          () => h('span', 'Has Role')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('denies when role does not match', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionButton,
          { role: 'platform_admin' },
          () => h('span', 'No Role')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(false)
    })

    it('checks roles array correctly', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionButton,
          { roles: ['participant', 'brand_admin'] as UserRole[] },
          () => h('span', 'Has Role')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('denies when none of roles match', () => {
      const userParticipant = createTestUser(['participant'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionButton,
          { roles: ['platform_admin', 'brand_admin'] as UserRole[] },
          () => h('span', 'No Role')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(false)
    })

    it('checks brandId correctly - no access for non-platform-admin', () => {
      const userBrandAdmin = createTestUser(['brand_admin'])

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userBrandAdmin },
        () => h(
          PermissionButton,
          { brandId: 123 },
          () => h('span', 'Brand Access')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(false)
    })

    it('grants brandId access for platform_admin', () => {
      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionButton,
          { brandId: 123 },
          () => h('span', 'Brand Access')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(true)
    })
  })

  describe('PermissionProvider additional tests', () => {
    it('hasPermission returns true for brand_admin with campaign permissions', () => {
      const userBrandAdmin = createTestUser(['brand_admin'])

      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasPermission('campaign:read')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userBrandAdmin },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })

    it('hasPermission returns false for brand_admin without admin permissions', () => {
      const userBrandAdmin = createTestUser(['brand_admin'])

      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasPermission('admin:delete')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userBrandAdmin },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })

    it('canAccessBrand returns false for non-platform-admin', () => {
      const userBrandAdmin = createTestUser(['brand_admin'])

      const Child = defineComponent({
        setup() {
          const { canAccessBrand } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(canAccessBrand(123)))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userBrandAdmin },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })

    it('hasAnyRole returns true when any role matches', () => {
      const user = createTestUser(['participant', 'brand_admin'])

      const Child = defineComponent({
        setup() {
          const { hasAnyRole } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasAnyRole(['platform_admin', 'brand_admin'])))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })

    it('hasPermission returns false when user is null', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasPermission('any:permission')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: null as unknown as CurrentUser },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })

    it('canAccessBrand returns false when user is null', () => {
      const Child = defineComponent({
        setup() {
          const { canAccessBrand } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(canAccessBrand(123)))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: null as unknown as CurrentUser },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })
  })
})
