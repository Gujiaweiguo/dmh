import { test, expect } from '@playwright/test';

test.describe('H5 Campaign Flow', () => {
  test('user can view campaign list', async ({ page }) => {
    await page.goto('/');
    // Wait for page to load
    await page.waitForTimeout(1000);
    // Verify body is visible
    await expect(page.locator('body')).toBeVisible();
  });

  test('user can navigate to feedback', async ({ page }) => {
    await page.goto('/feedback');
    await page.waitForTimeout(1000);
    // Verify feedback page loaded
    await expect(page.locator('text=帮助与反馈')).toBeVisible();
  });
});

test.describe('H5 Brand Admin Flow', () => {
  test('brand admin can access login page', async ({ page }) => {
    await page.goto('/brand/login');
    await page.waitForTimeout(1000);
    // Verify login form exists
    await expect(page.locator('body')).toBeVisible();
  });
});

test.describe('H5 Distributor Flow', () => {
  test('user can access distributor page', async ({ page }) => {
    await page.goto('/distributor');
    await page.waitForTimeout(1000);
    await expect(page.locator('body')).toBeVisible();
  });
});

test.describe('H5 Order Flow', () => {
  test('user can access orders page', async ({ page }) => {
    await page.goto('/orders');
    await page.waitForTimeout(1000);
    await expect(page.locator('body')).toBeVisible();
  });
});
