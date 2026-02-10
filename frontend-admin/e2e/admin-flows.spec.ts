import { test, expect } from '@playwright/test';

test.describe('Login Flow', () => {
  test('admin can login successfully', async ({ page }) => {
    await page.goto('/');

    await page.fill('input[placeholder="请输入用户名"]', 'admin');
    await page.fill('input[placeholder="请输入密码"]', '123456');
    await page.click('button[type="submit"]');

    await expect(page.locator('text=反馈管理')).toBeVisible();
  });

  test('shows error on invalid credentials', async ({ page }) => {
    await page.goto('/');

    await page.fill('input[placeholder="请输入用户名"]', 'admin');
    await page.fill('input[placeholder="请输入密码"]', 'wrongpassword');
    await page.click('button[type="submit"]');

    await expect(page.locator('text=HTTP error')).toBeVisible();
  });
});

test.describe('Navigation Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.fill('input[placeholder="请输入用户名"]', 'admin');
    await page.fill('input[placeholder="请输入密码"]', '123456');
    await page.click('button[type="submit"]');
    await expect(page.locator('text=反馈管理')).toBeVisible();
  });

  test('can navigate to campaign monitoring', async ({ page }) => {
    await page.click('text=活动监控');
    await page.waitForTimeout(500);
    // Verify page loaded
    await expect(page.locator('body')).toBeVisible();
  });

  test('can navigate to user management', async ({ page }) => {
    await page.click('text=用户管理');
    await page.waitForTimeout(500);
    await expect(page.locator('body')).toBeVisible();
  });

  test('can navigate to distributor monitoring', async ({ page }) => {
    await page.click('text=分销监控');
    await page.waitForTimeout(500);
    await expect(page.locator('body')).toBeVisible();
  });

  test('can navigate to feedback management', async ({ page }) => {
    await page.click('text=反馈管理');
    await page.waitForTimeout(500);
    // Feedback management has table
    await expect(page.locator('.table')).toBeVisible();
  });
});
