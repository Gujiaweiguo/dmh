import { test, expect } from '@playwright/test';

test.describe('Admin Dashboard E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:3000');
  });

  test('displays dashboard with stats', async ({ page }) => {
    await expect(page).toHaveTitle(/DMH/);
    await expect(page.locator('h1:has-text("控制面板")')).toBeVisible();
    await expect(page.locator('text=累计报名')).toBeVisible();
    await expect(page.locator('text=推广总收益')).toBeVisible();
    await expect(page.locator('text=活跃活动')).toBeVisible();
    await expect(page.locator('text=待处理提现')).toBeVisible();
  });

  test('navigates to campaign monitoring', async ({ page }) => {
    await page.click('text=活动监控');
    await expect(page).toHaveURL(/.*campaigns/);
    await expect(page.locator('h1:has-text("活动监控")')).toBeVisible();
    await expect(page.locator('text=关键词')).toBeVisible();
    await expect(page.locator('text=品牌')).toBeVisible();
    await expect(page.locator('text=状态')).toBeVisible();
  });

  test('navigates to distributor monitoring', async ({ page }) => {
    await page.click('text=分销监控');
    await expect(page).toHaveURL(/.*distributor-management/);
    await expect(page.locator('h1:has-text("分销监控")')).toBeVisible();
    await expect(page.locator('text=总数')).toBeVisible();
    await expect(page.locator('text=正常')).toBeVisible();
    await expect(page.locator('text=暂停')).toBeVisible();
    await expect(page.locator('text=待审核')).toBeVisible();
    await expect(page.locator('text=全部级别')).toBeVisible();
  });

  test('navigates to user management', async ({ page }) => {
    await page.click('text=用户管理');
    await expect(page.locator('h1:has-text("用户管理")')).toBeVisible();
  });

  test('navigates to brand management', async ({ page }) => {
    await page.click('text=品牌管理');
    await expect(page.locator('h1:has-text("品牌管理")')).toBeVisible();
  });

  test('navigates to member management', async ({ page }) => {
    await page.click('text=会员管理');
    await expect(page.locator('h1:has-text("会员管理")')).toBeVisible();
  });

  test('navigates to feedback management', async ({ page }) => {
    await page.click('text=反馈管理');
    await expect(page.locator('h1:has-text("反馈管理")')).toBeVisible();
  });

  test('navigates to verification records', async ({ page }) => {
    await page.click('text=核销记录');
    await expect(page.locator('h1:has-text("核销记录")')).toBeVisible();
  });

  test('navigates to poster records', async ({ page }) => {
    await page.click('text=海报记录');
    await expect(page.locator('h1:has-text("海报记录")')).toBeVisible();
  });

  test('navigates to system settings', async ({ page }) => {
    await page.click('text=系统设置');
    await expect(page.locator('h1:has-text("系统设置")')).toBeVisible();
  });
});

test.describe('Campaign Monitoring E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:3000/#/campaigns');
  });

  test('displays campaign filters', async ({ page }) => {
    await expect(page.locator('text=实时筛选 自动应用')).toBeVisible();
    await expect(page.locator('text=关键词')).toBeVisible();
    await expect(page.locator('text=品牌')).toBeVisible();
    await expect(page.locator('text=状态')).toBeVisible();
    await expect(page.locator('text=开始时间')).toBeVisible();
    await expect(page.locator('text=结束时间')).toBeVisible();
  });

  test('filters campaigns by status', async ({ page }) => {
    await page.selectOption('select:has-text("全部状态")', '进行中');
    await page.waitForTimeout(500);
    await expect(page.locator('text=实时筛选 自动应用')).toBeVisible();
  });

  test('resets filters', async ({ page }) => {
    await page.selectOption('select:has-text("全部状态")', '进行中');
    await page.click('text=重置筛选');
    await expect(page.locator('text=共找到')).toBeVisible();
  });
});

test.describe('Distributor Monitoring E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:3000/#/distributor-management');
  });

  test('displays distributor stats', async ({ page }) => {
    await expect(page.locator('text=总数')).toBeVisible();
    await expect(page.locator('text=正常')).toBeVisible();
    await expect(page.locator('text=暂停')).toBeVisible();
    await expect(page.locator('text=待审核')).toBeVisible();
  });

  test('displays distributor filters', async ({ page }) => {
    await expect(page.locator('input[placeholder*="搜索用户名或品牌"]').first()).toBeVisible();
    await expect(page.locator('text=全部状态')).toBeVisible();
    await expect(page.locator('text=全部级别')).toBeVisible();
    await expect(page.locator('text=一级')).toBeVisible();
    await expect(page.locator('text=二级')).toBeVisible();
    await expect(page.locator('text=三级')).toBeVisible();
  });

  test('refreshes distributor data', async ({ page }) => {
    await page.click('text=刷新');
    await page.waitForTimeout(1000);
    await expect(page.locator('h1:has-text("分销监控")')).toBeVisible();
  });

  test('filters by level', async ({ page }) => {
    await page.selectOption('select:has-text("全部级别")', '一级');
    await page.waitForTimeout(500);
    await expect(page.locator('text=第 1 页')).toBeVisible();
  });
});

test.describe('Navigation and Responsive E2E', () => {
  test('navigates back to dashboard from any page', async ({ page }) => {
    await page.goto('http://localhost:3000/#/campaigns');
    await expect(page.locator('h1:has-text("活动监控")')).toBeVisible();
    await page.click('text=控制面板');
    await expect(page.locator('h1:has-text("控制面板")')).toBeVisible();
  });

  test('maintains sidebar visibility', async ({ page }) => {
    await expect(page.locator('nav')).toBeVisible();
    await expect(page.locator('text=退出登录')).toBeVisible();
    await page.click('text=用户管理');
    await expect(page.locator('nav')).toBeVisible();
    await page.click('text=分销监控');
    await expect(page.locator('nav')).toBeVisible();
  });
});
