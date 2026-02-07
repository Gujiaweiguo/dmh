import { Campaign } from '../types';

const API_BASE_URL = '/api/v1';

export interface CampaignListResponse {
  total: number;
  campaigns: Campaign[];
}

export interface CreateCampaignRequest {
  name: string;
  description: string;
  formFields: string[];
  rewardRule: number;
  startTime: string;
  endTime: string;
  paymentConfig?: string;
  enableDistribution?: boolean;
  distributionLevel?: number;
  distributionRewards?: string;
  posterTemplateId?: number;
}

export interface UpdateCampaignRequest extends CreateCampaignRequest {
  status: 'active' | 'paused' | 'ended';
}

class CampaignApiService {
  private getAuthHeaders() {
    const token = localStorage.getItem('dmh_token');
    return {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
    };
  }

  // 创建活动
  async createCampaign(data: CreateCampaignRequest): Promise<Campaign> {
    const response = await fetch(`${API_BASE_URL}/campaigns`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to create campaign: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  // 获取活动列表
  async getCampaigns(page = 1, pageSize = 20, status?: string, keyword?: string): Promise<CampaignListResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      pageSize: pageSize.toString(),
    });
    if (status) params.append('status', status);
    if (keyword) params.append('keyword', keyword);

    const response = await fetch(`${API_BASE_URL}/campaigns?${params}`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch campaigns: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  // 获取活动详情
  async getCampaign(id: number): Promise<Campaign> {
    const response = await fetch(`${API_BASE_URL}/campaigns/${id}`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch campaign: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  // 更新活动
  async updateCampaign(id: number, data: UpdateCampaignRequest): Promise<Campaign> {
    const response = await fetch(`${API_BASE_URL}/campaigns/${id}`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to update campaign: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  // 删除活动
  async deleteCampaign(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/campaigns/${id}`, {
      method: 'DELETE',
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to delete campaign: ${response.status} ${errorText}`);
    }
  }
}

export const campaignApi = new CampaignApiService();
