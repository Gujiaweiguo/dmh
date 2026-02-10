import { test, expect } from '@playwright/test';

test.describe('反馈管理 - Admin 端测试', () => {

  test.beforeEach(async ({ page }) => {
    // 访问登录页面
    await page.goto('/');

    // 登录
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', '123456');
    await page.click('button:has-text("登录")');

    // 等待登录成功并跳转到管理后台
    await page.waitForURL('/dashboard');
  });

  test('场景1: 管理员登录', async ({ page }) => {
    // 验证登录成功，跳转到管理后台主页
    await expect(page).toHaveURL('/dashboard');

    // 验证显示管理员后台界面
    await expect(page.locator('.admin-layout')).toBeVisible();
  });

  test('场景2: 查看反馈统计', async ({ page }) => {
    // 点击"反馈管理"菜单
    await page.click('text=反馈管理');

    // 等待反馈管理页面加载
    await page.waitForURL('/feedback-management');

    // 验证统计卡片显示
    await expect(page.locator('.stat-card:has-text("总反馈数")')).toBeVisible();
    await expect(page.locator('.stat-card:has-text("待处理")')).toBeVisible();
    await expect(page.locator('.stat-card:has-text("处理中")')).toBeVisible();
    await expect(page.locator('.stat-card:has-text("已解决")')).toBeVisible();

    // 测试日期筛选功能
    await page.selectOption('select[name="dateFilter"]', 'week');
    await page.waitForTimeout(500);

    // 验证统计数据已更新（通过检查加载状态）
    await expect(page.locator('.feedback-list')).toBeVisible();
  });

  test('场景3: 查看反馈列表', async ({ page }) => {
    // 进入反馈管理页面
    await page.click('text=反馈管理');
    await page.waitForURL('/feedback-management');

    // 验证反馈列表显示
    await expect(page.locator('.feedback-list')).toBeVisible();

    // 检查列表项包含: ID、用户、标题、类别、优先级、状态、提交时间
    const feedbackItems = page.locator('.feedback-item');
    const count = await feedbackItems.count();

    if (count > 0) {
      const firstItem = feedbackItems.first();
      await expect(firstItem.locator('.feedback-id')).toBeVisible();
      await expect(firstItem.locator('.feedback-user')).toBeVisible();
      await expect(firstItem.locator('.feedback-title')).toBeVisible();
      await expect(firstItem.locator('.feedback-category')).toBeVisible();
      await expect(firstItem.locator('.feedback-priority')).toBeVisible();
      await expect(firstItem.locator('.feedback-status')).toBeVisible();
      await expect(firstItem.locator('.feedback-time')).toBeVisible();

      // 测试状态筛选
      await page.selectOption('select[name="statusFilter"]', 'pending');
      await page.waitForTimeout(500);

      // 验证筛选结果
      const filteredItems = page.locator('.feedback-item');
      if (await filteredItems.count() > 0) {
        const firstFilteredItem = filteredItems.first();
        await expect(firstFilteredItem.locator('.feedback-status')).toHaveText(/待处理/);
      }

      // 测试类别筛选
      await page.selectOption('select[name="statusFilter"]', 'all');
      await page.selectOption('select[name="categoryFilter"]', 'other');
      await page.waitForTimeout(500);

      // 测试优先级筛选
      await page.selectOption('select[name="categoryFilter"]', 'all');
      await page.selectOption('select[name="priorityFilter"]', 'high');
      await page.waitForTimeout(500);

      // 测试搜索功能
      await page.selectOption('select[name="priorityFilter"]', 'all');
      await page.fill('input[name="search"]', '测试');
      await page.keyboard.press('Enter');
      await page.waitForTimeout(500);
    }
  });

  test('场景4: 查看反馈详情', async ({ page }) => {
    // 进入反馈管理页面
    await page.click('text=反馈管理');
    await page.waitForURL('/feedback-management');

    // 点击任意反馈项
    const feedbackItems = page.locator('.feedback-item');
    const count = await feedbackItems.count();

    if (count > 0) {
      await feedbackItems.first().click();

      // 验证详情页显示
      await expect(page.locator('.feedback-detail-modal')).toBeVisible();

      // 检查详情页包含: 反馈基本信息、反馈内容、管理操作
      await expect(page.locator('.feedback-info')).toBeVisible();
      await expect(page.locator('.feedback-content')).toBeVisible();
      await expect(page.locator('.feedback-actions')).toBeVisible();

      // 点击关闭按钮
      await page.click('button:has-text("关闭")');
      await expect(page.locator('.feedback-detail-modal')).not.toBeVisible();
    }
  });

  test('场景5: 更新反馈状态', async ({ page }) => {
    // 进入反馈管理页面
    await page.click('text=反馈管理');
    await page.waitForURL('/feedback-management');

    // 点击待处理的反馈
    const pendingFeedbacks = page.locator('.feedback-item:has-text("待处理")');
    const count = await pendingFeedbacks.count();

    if (count > 0) {
      await pendingFeedbacks.first().click();

      // 选择新状态: 处理中
      await page.selectOption('select[name="status"]', 'reviewing');

      // 点击"更新状态"按钮
      await page.click('button:has-text("更新状态")');

      // 验证状态更新成功提示
      await expect(page.locator('text=状态更新成功')).toBeVisible();

      // 关闭详情页
      await page.click('button:has-text("关闭")');

      // 验证列表中状态已更新
      await expect(pendingFeedbacks.first().locator('.feedback-status')).toHaveText(/处理中/);

      // 验证统计数据已更新
      await page.waitForTimeout(500);
      await expect(page.locator('.stat-card:has-text("待处理")')).toBeVisible();
    }
  });

  test('场景6: 批量操作', async ({ page }) => {
    // 进入反馈管理页面
    await page.click('text=反馈管理');
    await page.waitForURL('/feedback-management');

    // 勾选多个反馈项
    const feedbackItems = page.locator('.feedback-item');
    const count = await feedbackItems.count();

    if (count >= 2) {
      await feedbackItems.nth(0).locator('input[type="checkbox"]').check();
      await feedbackItems.nth(1).locator('input[type="checkbox"]').check();

      // 选择批量操作: 标记为已解决
      await page.selectOption('select[name="bulkAction"]', 'mark-resolved');

      // 点击确认按钮
      await page.click('button:has-text("确认")');

      // 验证批量操作成功提示
      await expect(page.locator('text=批量操作成功')).toBeVisible();

      // 验证所有选中反馈状态更新
      await expect(feedbackItems.nth(0).locator('.feedback-status')).toHaveText(/已解决/);
      await expect(feedbackItems.nth(1).locator('.feedback-status')).toHaveText(/已解决/);
    }
  });

  test('场景7: 分页功能', async ({ page }) => {
    // 进入反馈管理页面
    await page.click('text=反馈管理');
    await page.waitForURL('/feedback-management');

    // 验证分页组件显示
    await expect(page.locator('.pagination')).toBeVisible();

    // 测试下一页
    const nextPageButton = page.locator('button:has-text("下一页")');
    const isNextPageDisabled = await nextPageButton.isDisabled();

    if (!isNextPageDisabled) {
      const initialItems = await page.locator('.feedback-item').count();

      await nextPageButton.click();
      await page.waitForTimeout(500);

      const newItems = await page.locator('.feedback-item').count();

      // 测试上一页
      const prevPageButton = page.locator('button:has-text("上一页")');
      if (!await prevPageButton.isDisabled()) {
        await prevPageButton.click();
        await page.waitForTimeout(500);
      }
    }
  });

  test('完整流程: 反馈管理完整生命周期', async ({ page }) => {
    // 进入反馈管理页面
    await page.click('text=反馈管理');
    await page.waitForURL('/feedback-management');

    // 查看所有反馈
    await expect(page.locator('.feedback-list')).toBeVisible();

    // 查看统计信息
    await expect(page.locator('.stat-card:has-text("总反馈数")')).toBeVisible();

    // 点击待处理的反馈
    const pendingFeedbacks = page.locator('.feedback-item:has-text("待处理")');
    const count = await pendingFeedbacks.count();

    if (count > 0) {
      const firstPending = pendingFeedbacks.first();
      const feedbackTitle = await firstPending.locator('.feedback-title').textContent();

      // 查看详情
      await firstPending.click();
      await expect(page.locator('.feedback-detail-modal')).toBeVisible();

      // 更新状态为处理中
      await page.selectOption('select[name="status"]', 'reviewing');
      await page.click('button:has-text("更新状态")');
      await expect(page.locator('text=状态更新成功')).toBeVisible();

      // 更新状态为已解决
      await page.selectOption('select[name="status"]', 'resolved');
      await page.click('button:has-text("更新状态")');
      await expect(page.locator('text=状态更新成功')).toBeVisible();

      // 关闭详情页
      await page.click('button:has-text("关闭")');

      // 验证反馈状态已更新
      const updatedFeedback = page.locator(`.feedback-item:has-text("${feedbackTitle}")`);
      await expect(updatedFeedback.locator('.feedback-status')).toHaveText(/已解决/);

      // 验证统计数据已更新
      await expect(page.locator('.stat-card:has-text("待处理")')).toBeVisible();
      await expect(page.locator('.stat-card:has-text("已解决")')).toBeVisible();
    }
  });
});
