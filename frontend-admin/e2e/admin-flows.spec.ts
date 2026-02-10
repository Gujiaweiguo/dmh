import { test, expect } from '@playwright/test';

test.describe('Login Flow', () => {
  test('admin can login successfully', async ({ page }) => {
    await page.goto('/');
    
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', '123456');
    await page.click('button[type="submit"]');
    
    await expect(page).toHaveURL(/.*dashboard.*/);
    await expect(page.locator('text=Welcome')).toBeVisible();
  });

  test('shows error on invalid credentials', async ({ page }) => {
    await page.goto('/');
    
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');
    
    await expect(page.locator('text=用户名或密码错误')).toBeVisible();
  });
});

test.describe('Campaign Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', '123456');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/.*dashboard.*/);
  });

  test('can create new campaign', async ({ page }) => {
    await page.click('text=活动管理');
    await page.click('text=新建活动');
    
    await page.fill('input[name="name"]', 'E2E Test Campaign');
    await page.fill('textarea[name="description"]', 'This is a test campaign');
    await page.fill('input[name="rewardRule"]', '100');
    
    await page.click('button[type="submit"]');
    
    await expect(page.locator('text=创建成功')).toBeVisible();
    await expect(page.locator('text=E2E Test Campaign')).toBeVisible();
  });

  test('can edit existing campaign', async ({ page }) => {
    await page.click('text=活动管理');
    await page.locator('text=E2E Test Campaign').first().click();
    await page.click('text=编辑');
    
    await page.fill('input[name="name"]', 'Updated Campaign Name');
    await page.click('button[type="submit"]');
    
    await expect(page.locator('text=更新成功')).toBeVisible();
    await expect(page.locator('text=Updated Campaign Name')).toBeVisible();
  });

  test('can delete campaign', async ({ page }) => {
    await page.click('text=活动管理');
    await page.locator('.campaign-item').first().hover();
    await page.click('text=删除');
    await page.click('text=确认');
    
    await expect(page.locator('text=删除成功')).toBeVisible();
  });
});

test.describe('User Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', '123456');
    await page.click('button[type="submit"]');
  });

  test('can view user list', async ({ page }) => {
    await page.click('text=用户管理');
    
    await expect(page.locator('.user-table')).toBeVisible();
    await expect(page.locator('text=用户名')).toBeVisible();
  });

  test('can create new user', async ({ page }) => {
    await page.click('text=用户管理');
    await page.click('text=新建用户');
    
    const timestamp = Date.now();
    await page.fill('input[name="username"]', `testuser_${timestamp}`);
    await page.fill('input[name="password"]', 'password123');
    await page.fill('input[name="phone"]', `138${timestamp.toString().slice(-8)}`);
    await page.selectOption('select[name="role"]', 'participant');
    
    await page.click('button[type="submit"]');
    
    await expect(page.locator('text=创建成功')).toBeVisible();
  });
});

test.describe('Distributor Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', '123456');
    await page.click('button[type="submit"]');
  });

  test('can view distributor list', async ({ page }) => {
    await page.click('text=分销管理');
    
    await expect(page.locator('.distributor-table')).toBeVisible();
    await expect(page.locator('text=分销商列表')).toBeVisible();
  });

  test('can filter distributors by status', async ({ page }) => {
    await page.click('text=分销管理');
    await page.selectOption('select[name="status"]', 'active');
    
    await expect(page.locator('.distributor-row')).toHaveCount.greaterThan(0);
  });
});
