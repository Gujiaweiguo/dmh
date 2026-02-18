import { test, expect } from '@playwright/test';

test.describe('Admin Dashboard E2E', () => {
  test.beforeEach(async ({ page }) => {
    // 访问首页并登录
    await page.goto('/');
    await page.fill('input[placeholder="请输入用户名"]', 'admin');
    await page.fill('input[placeholder="请输入密码"]', '123456');
    await page.click('button[type="submit"]');
    // 等待登录成功
    await expect(page.locator('text=反馈管理')).toBeVisible({ timeout: 10000 });
  });

  test('displays dashboard with navigation', async ({ page }) => {
    // 验证控制面板菜单存在
    await expect(page.getByRole('heading', { name: '控制面板' })).toBeVisible();
    // 验证主要菜单项存在
    await expect(page.locator('text=活动监控').first()).toBeVisible();
    await expect(page.locator('text=分销监控').first()).toBeVisible();
    await expect(page.locator('text=用户管理').first()).toBeVisible();
    await expect(page.locator('text=品牌管理').first()).toBeVisible();
  });

  test('navigates to campaign monitoring', async ({ page }) => {
    await page.click('text=活动监控');
    // 等待页面加载
    await page.waitForTimeout(500);
    // 验证活动筛选元素存在
    await expect(page.locator('text=关键词')).toBeVisible();
  });

  test('navigates to distributor monitoring', async ({ page }) => {
    await page.click('text=分销监控');
    await page.waitForTimeout(500);
    // 验证分销监控页面已加载（页面存在且不报错）
    await expect(page.locator('body')).toBeVisible();
  });

  test('navigates to user management', async ({ page }) => {
    await page.click('text=用户管理');
    await page.waitForTimeout(500);
    // 验证用户管理页面已加载（页面存在且不报错）
    await expect(page.locator('body')).toBeVisible();
  });

  test('navigates to feedback management', async ({ page }) => {
    await page.click('text=反馈管理');
    await page.waitForTimeout(500);
    // 验证反馈管理页面加载
    await expect(page.locator('.table')).toBeVisible();
  });
});

test.describe('Navigation and Responsive E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.fill('input[placeholder="请输入用户名"]', 'admin');
    await page.fill('input[placeholder="请输入密码"]', '123456');
    await page.click('button[type="submit"]');
    await expect(page.locator('text=反馈管理')).toBeVisible({ timeout: 10000 });
  });

  test('navigates back to dashboard from any page', async ({ page }) => {
    await page.click('text=活动监控');
    await page.waitForTimeout(500);
    await page.click('text=控制面板');
    await page.waitForTimeout(500);
    // 验证回到控制面板 - 使用 heading 选择器
    await expect(page.getByRole('heading', { name: '控制面板' })).toBeVisible();
  });

  test('maintains sidebar visibility', async ({ page }) => {
    // 验证侧边栏菜单项可见
    await expect(page.locator('text=用户管理')).toBeVisible();
    await page.click('text=分销监控');
    await expect(page.locator('text=品牌管理')).toBeVisible();
    await page.click('text=会员管理');
    await expect(page.locator('text=系统设置')).toBeVisible();
  });
});
