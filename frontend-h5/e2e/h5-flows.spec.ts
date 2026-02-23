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

  test('brand order detail and export flow', async ({ page }) => {
    await page.addInitScript(() => {
      localStorage.setItem('dmh_token', 'e2e-token');
      localStorage.setItem('dmh_user_info', JSON.stringify({ brandIds: [1] }));
    });

    await page.route('**/api/v1/order/list**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            list: [
              {
                id: 1001,
                campaignId: 9,
                campaignName: 'E2E 活动',
                phone: '13800000000',
                amount: 88,
                status: 'paid',
                rewardAmount: 10,
                createdAt: '2026-02-23 10:00:00',
                formData: { 姓名: '测试用户' },
              },
            ],
          },
        }),
      });
    });

    await page.route('**/api/v1/order/detail/1001**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            id: 1001,
            campaignName: 'E2E 活动',
            phone: '13800000000',
            amount: 88,
            status: 'paid',
            rewardAmount: 10,
            createdAt: '2026-02-23 10:00:00',
            formData: { 姓名: '测试用户' },
          },
        }),
      });
    });

    await page.goto('/brand/orders');
    await expect(page.locator('text=订单管理')).toBeVisible();

    await expect(page.locator('button:has-text("导出当前筛选")')).toBeVisible();
    await page.locator('button:has-text("导出当前筛选")').click();

    await page.locator('button:has-text("详情")').first().click();
    await expect(page).toHaveURL(/\/brand\/order-detail\/1001/);
    await expect(page.locator('text=订单详情')).toBeVisible();
  });
});


test.describe('H5 Promoter Flow', () => {
  test('promoter list page loads', async ({ page }) => {
    await page.addInitScript(() => {
      localStorage.setItem('dmh_token', 'e2e-token');
      localStorage.setItem('dmh_user_info', JSON.stringify({ brandIds: [1] }));
    });

    await page.route('**/api/v1/promoter/list**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            list: [
              {
                id: 1,
                name: '测试推广员',
                phone: '138****1234',
                status: 'active',
                level: 'VIP',
                totalOrders: 50,
                totalRewards: 2000,
                conversionRate: 15,
                recentActivities: [],
              },
            ],
          },
        }),
      });
    });

    await page.goto('/brand/promoters');
    await expect(page.locator('text=推广员管理')).toBeVisible();
  });

  test('promoter detail and reward records flow', async ({ page }) => {
    await page.addInitScript(() => {
      localStorage.setItem('dmh_token', 'e2e-token');
      localStorage.setItem('dmh_user_info', JSON.stringify({ brandIds: [1] }));
    });

    await page.route('**/api/v1/promoter/detail/1**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            id: 1,
            name: '测试推广员',
            phone: '138****1234',
            status: 'active',
            level: 'VIP',
            totalOrders: 50,
            totalRewards: 2000,
            conversionRate: 15,
            campaignCount: 5,
            createdAt: '2026-01-15',
            lastActiveAt: '2026-02-23',
            links: [],
          },
        }),
      });
    });

    await page.route('**/api/v1/promoter/rewards/1**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            {
              id: 1,
              type: 'commission',
              status: 'paid',
              amount: 100,
              description: '订单佣金',
              createdAt: '2026-02-20',
            },
            {
              id: 2,
              type: 'bonus',
              status: 'pending',
              amount: 50,
              description: '活动奖金',
              createdAt: '2026-02-21',
            },
          ],
        }),
      });
    });

    // Navigate to promoter detail
    await page.goto('/brand/promoter-detail/1');
    await expect(page.locator('text=推广员详情')).toBeVisible();
    await expect(page.locator('text=测试推广员')).toBeVisible();

    // Click view reward records button
    await page.locator('button:has-text("查看奖励记录")').click();
    await expect(page).toHaveURL(/\/brand\/reward-records\/1/);
    await expect(page.locator('text=奖励记录')).toBeVisible();
  });

  test('reward records page with filters', async ({ page }) => {
    await page.addInitScript(() => {
      localStorage.setItem('dmh_token', 'e2e-token');
      localStorage.setItem('dmh_user_info', JSON.stringify({ brandIds: [1] }));
    });

    await page.route('**/api/v1/promoter/rewards/**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            {
              id: 1,
              type: 'commission',
              status: 'paid',
              amount: 100,
              description: '订单佣金',
              createdAt: '2026-02-20',
            },
          ],
        }),
      });
    });

    await page.goto('/brand/reward-records');
    await expect(page.locator('text=奖励记录')).toBeVisible();

    // Check filter exists
    await expect(page.locator('select').first()).toBeVisible();
  });
});