import { test, expect, type Page } from '@playwright/test';

async function loginAsAdmin(page: Page) {
  await page.goto('/');
  await page.fill('input[placeholder="请输入用户名"]', 'admin');
  await page.fill('input[placeholder="请输入密码"]', '123456');
  await page.click('button[type="submit"]');
  await expect(page.getByRole('heading', { name: '控制面板' })).toBeVisible({ timeout: 10000 });
}

test.describe('安全管理 E2E', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
  });

  test('可以进入安全管理并返回非空结构数据', async ({ page }) => {
    const policyRespPromise = page.waitForResponse((resp) =>
      resp.url().includes('/api/v1/security/password-policy') && resp.request().method() === 'GET'
    );
    const sessionsRespPromise = page.waitForResponse((resp) =>
      resp.url().includes('/api/v1/security/sessions') && resp.request().method() === 'GET'
    );
    const eventsRespPromise = page.waitForResponse((resp) =>
      resp.url().includes('/api/v1/security/events') && resp.request().method() === 'GET'
    );

    await page.getByRole('button', { name: '安全管理' }).click();

    await expect(page.locator('text=密码策略').first()).toBeVisible();
    await expect(page.locator('text=活跃会话').first()).toBeVisible();
    await expect(page.locator('text=安全事件').first()).toBeVisible();

    const [policyResp, sessionsResp, eventsResp] = await Promise.all([
      policyRespPromise,
      sessionsRespPromise,
      eventsRespPromise,
    ]);

    expect(policyResp.status()).toBe(200);
    expect(sessionsResp.status()).toBe(200);
    expect(eventsResp.status()).toBe(200);

    const policyBody = await policyResp.json();
    expect(policyBody).toBeTruthy();
    expect(policyBody).toEqual(
      expect.objectContaining({
        id: expect.any(Number),
        minLength: expect.any(Number),
      })
    );

    const sessionsBody = await sessionsResp.json();
    expect(sessionsBody).toEqual(
      expect.objectContaining({
        total: expect.any(Number),
        sessions: expect.any(Array),
      })
    );

    const eventsBody = await eventsResp.json();
    expect(eventsBody).toEqual(
      expect.objectContaining({
        total: expect.any(Number),
        events: expect.any(Array),
      })
    );
  });

  test('保存策略后返回成功提示与对象响应', async ({ page }) => {
    await page.getByRole('button', { name: '安全管理' }).click();
    await expect(page.locator('text=密码策略').first()).toBeVisible();

    const saveRespPromise = page.waitForResponse((resp) =>
      resp.url().includes('/api/v1/security/password-policy') &&
      resp.request().method() === 'PUT' &&
      resp.status() === 200
    );

    await page.click('button:has-text("保存策略")');

    const saveResp = await saveRespPromise;
    const saveBody = await saveResp.json();
    expect(saveBody).toBeTruthy();
    expect(saveBody).toEqual(
      expect.objectContaining({
        id: expect.any(Number),
        minLength: expect.any(Number),
      })
    );

    await expect(page.locator('text=密码策略已更新')).toBeVisible();
  });

  test('会话操作按钮可以触发确认弹窗并取消', async ({ page }) => {
    await page.getByRole('button', { name: '安全管理' }).click();
    await expect(page.locator('text=活跃会话').first()).toBeVisible();

    const revokeButton = page.locator('button:has-text("撤销会话")').first();
    const revokeCount = await revokeButton.count();

    if (revokeCount === 0) {
      await expect(page.locator('text=暂无会话数据')).toBeVisible();
      return;
    }

    let revokeConfirmMessage = '';
    page.once('dialog', async (dialog) => {
      revokeConfirmMessage = dialog.message();
      await dialog.dismiss();
    });

    await revokeButton.click();
    await expect.poll(() => revokeConfirmMessage).toContain('确认撤销会话');

    const forceButton = page.locator('button:has-text("强制下线")').first();
    let forceConfirmMessage = '';
    page.once('dialog', async (dialog) => {
      forceConfirmMessage = dialog.message();
      await dialog.dismiss();
    });

    await forceButton.click();
    await expect.poll(() => forceConfirmMessage).toContain('确认强制用户');
  });
});
