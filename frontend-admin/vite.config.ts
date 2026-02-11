import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';


export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, '.', '');
    return {
      cacheDir: '/tmp/dmh-admin-vite-cache',
      server: {
        port: 3000,
        host: '0.0.0.0',
        proxy: {
          '/api': {
            target: 'http://localhost:8889',
            changeOrigin: true,
          },
        },
      },
      plugins: [vue()],
      build: {
        // 代码分割优化
        rollupOptions: {
          output: {
            manualChunks: {
              // 将第三方库单独打包
              'vendor': ['vue'],
              // 将路由组件按需加载
              'feedback': ['./views/FeedbackManagementView.vue'],
            },
          },
        },
        // 压缩选项（避免依赖可选的 terser）
        minify: 'esbuild',
        esbuild: {
          drop: ['console', 'debugger'],
        },
        // 资源内联阈值
        assetsInlineLimit: 4096,
      },
      define: {
        'process.env.API_KEY': JSON.stringify(env.GEMINI_API_KEY),
        'process.env.GEMINI_API_KEY': JSON.stringify(env.GEMINI_API_KEY),
        __VUE_OPTIONS_API__: true,
        __VUE_PROD_DEVTOOLS__: false
      },
      resolve: {
        alias: {
          '@': path.resolve(__dirname, '.'),
          vue: 'vue/dist/vue.esm-bundler.js',
        }
      }
    };
});
