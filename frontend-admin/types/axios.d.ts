declare module 'axios' {
  export interface AxiosInstance {
    get<T = any>(url: string, config?: any): Promise<{ data: T }>;
    post<T = any>(url: string, data?: any, config?: any): Promise<{ data: T }>;
    put<T = any>(url: string, data?: any, config?: any): Promise<{ data: T }>;
    delete<T = any>(url: string, config?: any): Promise<{ data: T }>;
    create(config?: AxiosRequestConfig): AxiosInstance;
    interceptors: {
      request: { use: (fn: (config: any) => any) => number };
      response: { use: (onFulfilled: (response: any) => any, onRejected?: (error: any) => any) => number };
    };
  }

  export interface AxiosRequestConfig {
    baseURL?: string;
    timeout?: number;
    headers?: Record<string, string>;
  }

  const axios: AxiosInstance;
  export default axios;
}
