// 前端性能监控工具
const PerformanceMonitor = {
  // 记录 API 请求时间
  logApiRequest: (endpoint, duration, success = true) => {
    const threshold = 500; // 500ms 阈值
    if (duration > threshold) {
      console.warn(`[SLOW API] ${endpoint} took ${duration}ms`);
    }
    
    // 可以发送到监控服务
    if (window.navigator.sendBeacon) {
      window.navigator.sendBeacon('/api/metrics/frontend', JSON.stringify({
        type: 'api',
        endpoint,
        duration,
        success,
        timestamp: new Date().toISOString(),
      }));
    }
  },

  // 监控组件渲染时间
  measureRender: (componentName, fn) => {
    const start = performance.now();
    const result = fn();
    const duration = performance.now() - start;
    
    if (duration > 100) {
      console.warn(`[SLOW RENDER] ${componentName} took ${duration.toFixed(2)}ms`);
    }
    
    return result;
  },

  // 防抖函数
  debounce: (fn, delay = 300) => {
    let timer = null;
    return function(...args) {
      if (timer) clearTimeout(timer);
      timer = setTimeout(() => {
        fn.apply(this, args);
      }, delay);
    };
  },

  // 节流函数
  throttle: (fn, limit = 300) => {
    let inThrottle = false;
    return function(...args) {
      if (!inThrottle) {
        fn.apply(this, args);
        inThrottle = true;
        setTimeout(() => { inThrottle = false; }, limit);
      }
    };
  },

  // 监控页面加载性能
  init: () => {
    if (typeof window !== 'undefined') {
      window.addEventListener('load', () => {
        setTimeout(() => {
          const perfData = performance.timing;
          const pageLoadTime = perfData.loadEventEnd - perfData.navigationStart;
          
          if (pageLoadTime > 3000) {
            console.warn(`[SLOW PAGE LOAD] ${pageLoadTime}ms`);
          }
          
          console.log('[PERFORMANCE] Page load time:', pageLoadTime, 'ms');
        }, 0);
      });
    }
  },
};

export default PerformanceMonitor;
