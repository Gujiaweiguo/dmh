import { describe, expect, it } from 'vitest'
import {
  getCurrentBrandId,
  getDefaultBrandInfo,
  getDefaultNotificationSettings,
  getDefaultPasswordForm,
  getDefaultRewardSettings,
  getDefaultSyncSettings,
  getSyncStatusText,
  resolveUploadedLogoUrl,
  resolveSyncStatus,
  unwrapApiResponse,
  validateLogoFile,
  validatePasswordForm,
} from '../../src/views/brand/settings.logic.js'

describe('settings logic', () => {
  it('provides default settings objects', () => {
    expect(getDefaultBrandInfo()).toMatchObject({
      name: '示例品牌',
      phone: '400-123-4567',
    })
    expect(getDefaultRewardSettings()).toEqual({
      defaultRate: 20,
      minWithdraw: 100,
      settlementType: 'instant',
    })
    expect(getDefaultNotificationSettings()).toEqual({
      newOrder: true,
      newPromoter: true,
      dailyReport: false,
      email: 'admin@brand.com',
    })
    expect(getDefaultSyncSettings()).toEqual({
      status: 'connected',
      frequency: 'realtime',
      dataTypes: ['orders', 'users'],
    })
    expect(getDefaultPasswordForm()).toEqual({
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    })
  })

  it('maps sync status text', () => {
    expect(getSyncStatusText('connected')).toBe('已连接')
    expect(getSyncStatusText('disconnected')).toBe('未连接')
    expect(getSyncStatusText('error')).toBe('连接错误')
    expect(getSyncStatusText('custom')).toBe('custom')
  })

  it('validates password form fields', () => {
    expect(
      validatePasswordForm({ oldPassword: '', newPassword: '123456', confirmPassword: '123456' }),
    ).toBe('请填写完整的密码信息')

    expect(
      validatePasswordForm({ oldPassword: 'old', newPassword: '123456', confirmPassword: '654321' }),
    ).toBe('两次输入的新密码不一致')

    expect(
      validatePasswordForm({ oldPassword: 'old', newPassword: '123456', confirmPassword: '123456' }),
    ).toBe('')
  })

  it('resolves current brand id from storage and user info', () => {
    expect(getCurrentBrandId('12', '')).toBe(12)
    expect(getCurrentBrandId('invalid', JSON.stringify({ brandIds: [22] }))).toBe(22)
    expect(getCurrentBrandId('', JSON.stringify({ brandIds: [] }))).toBe(0)
    expect(getCurrentBrandId('', 'invalid-json')).toBe(0)
  })

  it('unwraps api payload', () => {
    expect(unwrapApiResponse({ data: { id: 1, name: 'A' } })).toEqual({ id: 1, name: 'A' })
    expect(unwrapApiResponse({ id: 2 })).toEqual({ id: 2 })
    expect(unwrapApiResponse(null)).toEqual({})
  })

  it('maps sync health status to ui status', () => {
    expect(resolveSyncStatus({ status: 'healthy' })).toBe('connected')
    expect(resolveSyncStatus({ status: 'ok' })).toBe('connected')
    expect(resolveSyncStatus({ status: 'error' })).toBe('error')
    expect(resolveSyncStatus({})).toBe('error')
  })

  it('validates logo file', () => {
    expect(validateLogoFile(null)).toBe('请选择图片文件')
    expect(validateLogoFile({ type: 'text/plain', size: 1 })).toBe('仅支持图片文件上传')
    expect(validateLogoFile({ type: 'image/png', size: 6 * 1024 * 1024 })).toBe('图片大小不能超过 5MB')
    expect(validateLogoFile({ type: 'image/png', size: 1024 })).toBe('')
  })

  it('resolves uploaded logo url from api payload', () => {
    expect(resolveUploadedLogoUrl({ data: { url: 'https://a/logo.png' } })).toBe('https://a/logo.png')
    expect(resolveUploadedLogoUrl({ data: { fileUrl: 'https://a/file.png' } })).toBe('https://a/file.png')
    expect(resolveUploadedLogoUrl({ data: {} })).toBe('')
  })
})
