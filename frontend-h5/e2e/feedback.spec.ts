import { test, expect } from '@playwright/test';

test.describe('反馈中心 - H5 端测试', () => {

  test('场景1: 页面加载和表单显示', async ({ page }) => {
    await page.goto('/feedback');

    // 等待页面加载
    await expect(page.locator('text=帮助与反馈')).toBeVisible({ timeout: 10000 });

    // 验证表单元素存在 - 使用 heading 角色选择器
    await expect(page.getByRole('heading', { name: '提交反馈' })).toBeVisible();
    await expect(page.getByRole('heading', { name: '我的反馈' })).toBeVisible();
    await expect(page.getByRole('heading', { name: '常见问题' })).toBeVisible();
  });

  test('场景2: 返回按钮存在', async ({ page }) => {
    await page.goto('/feedback');

    await expect(page.locator('text=帮助与反馈')).toBeVisible();

    // 验证返回按钮
    await expect(page.locator('button:has-text("返回")')).toBeVisible();
  });
});
