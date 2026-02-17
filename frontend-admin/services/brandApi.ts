const API_BASE_URL = '/api/v1';

export interface BrandItem {
  id: number;
  name: string;
  logo: string;
  description: string;
  status: 'active' | 'disabled' | string;
  createdAt: string;
  updatedAt?: string;
}

export interface BrandAssetItem {
  id: number;
  brandId: number;
  name: string;
  type: 'image' | 'video' | 'document' | string;
  category: string;
  tags: string;
  fileUrl: string;
  fileSize: number;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface BrandStats {
  brandId: number;
  totalCampaigns: number;
  activeCampaigns: number;
  totalOrders: number;
  totalRevenue: number;
  totalRewards: number;
  participantCount: number;
  conversionRate: number;
  lastUpdated: string;
}

class BrandApiService {
  private getToken() {
    return localStorage.getItem('dmh_token');
  }

  private async request<T>(url: string, options: RequestInit = {}): Promise<T> {
    const token = this.getToken();
    if (!token) {
      throw new Error('请先登录');
    }

    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
        ...options.headers,
      },
    });

    if (response.status === 401) {
      localStorage.removeItem('dmh_token');
      throw new Error('登录已过期，请重新登录');
    }

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: '请求失败' }));
      throw new Error(error.message || '请求失败');
    }

    return response.json() as Promise<T>;
  }

  async getBrands(): Promise<BrandItem[]> {
    const data = await this.request<{ brands?: BrandItem[]; total?: number }>('/brands');
    return Array.isArray(data.brands) ? data.brands : [];
  }

  async createBrand(payload: { name: string; logo?: string; description?: string }): Promise<BrandItem> {
    return this.request<BrandItem>('/brands', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  }

  async updateBrand(id: number, payload: {
    name?: string;
    logo?: string;
    description?: string;
    status?: string;
  }): Promise<BrandItem> {
    return this.request<BrandItem>(`/brands/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
  }

  async getBrandAssets(brandId: number): Promise<BrandAssetItem[]> {
    const data = await this.request<{ assets?: BrandAssetItem[]; total?: number }>(`/brands/${brandId}/assets`);
    return Array.isArray(data.assets) ? data.assets : [];
  }

  async createBrandAsset(brandId: number, payload: {
    name: string;
    type: string;
    category?: string;
    tags?: string;
    fileUrl: string;
    fileSize: number;
    description?: string;
  }): Promise<BrandAssetItem> {
    return this.request<BrandAssetItem>(`/brands/${brandId}/assets`, {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  }

  async updateBrandAsset(brandId: number, assetId: number, payload: {
    name: string;
    type: string;
    category?: string;
    tags?: string;
    fileUrl: string;
    fileSize: number;
    description?: string;
  }): Promise<BrandAssetItem> {
    return this.request<BrandAssetItem>(`/brands/${brandId}/assets/${assetId}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
  }

  async deleteBrandAsset(brandId: number, assetId: number): Promise<void> {
    await this.request<{ message: string }>(`/brands/${brandId}/assets/${assetId}`, {
      method: 'DELETE',
    });
  }

  async getBrandStats(brandId: number): Promise<BrandStats> {
    return this.request<BrandStats>(`/brands/${brandId}/stats`);
  }
}

export const brandApi = new BrandApiService();
