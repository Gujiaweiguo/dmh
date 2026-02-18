import { defineConfig, devices } from '@playwright/test';

const playwrightPort = process.env.PW_WEB_PORT || '4173';
const playwrightBaseUrl = process.env.PLAYWRIGHT_BASE_URL || `http://127.0.0.1:${playwrightPort}`;

export default defineConfig({
  testDir: './e2e',
  timeout: 30 * 1000,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: playwrightBaseUrl,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    // CI 环境下使用 headless，本地开发保持非 headless
    headless: process.env.CI === 'true' ? true : false,
  },

  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        channel: 'chrome', // Use system chrome browser
      },
    },
  ],

  webServer: {
    command: `npm run dev -- --host 127.0.0.1 --port ${playwrightPort} --strictPort`,
    url: playwrightBaseUrl,
    reuseExistingServer: process.env.PW_REUSE_SERVER === 'true',
  },
});
