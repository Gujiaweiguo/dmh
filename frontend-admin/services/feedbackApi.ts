const API_BASE_URL = '/api/v1';

export interface FeedbackItem {
  id: number;
  userId: number;
  userName: string;
  userRole: string;
  category: string;
  subcategory: string;
  rating: number | null;
  title: string;
  content: string;
  featureUseCase: string;
  deviceInfo: string;
  browserInfo: string;
  priority: 'low' | 'medium' | 'high' | string;
  status: 'pending' | 'reviewing' | 'resolved' | 'closed' | string;
  assigneeId: number | null;
  response: string;
  resolvedAt: string | null;
  createdAt: string;
}

export interface FeedbackListResponse {
  total: number;
  feedbacks: FeedbackItem[];
}

export interface FeedbackStatsResponse {
  totalFeedbacks: number;
  byCategory: Record<string, number>;
  byStatus: Record<string, number>;
  byPriority: Record<string, number>;
  averageRating: number;
  resolutionRate: number;
  avgResolutionTime: number;
  byRating: Record<string, number>;
}

export interface UpdateFeedbackStatusRequest {
  id: number;
  status: string;
  assigneeId?: number;
  response?: string;
}

class FeedbackApiService {
  private getAuthHeaders() {
    const token = localStorage.getItem('dmh_token');
    return {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    };
  }

  async getFeedbackList(params: {
    page?: number;
    pageSize?: number;
    category?: string;
    status?: string;
    priority?: string;
  } = {}): Promise<FeedbackListResponse> {
    const query = new URLSearchParams();
    if (params.page) query.append('page', String(params.page));
    if (params.pageSize) query.append('pageSize', String(params.pageSize));
    if (params.category) query.append('category', params.category);
    if (params.status) query.append('status', params.status);
    if (params.priority) query.append('priority', params.priority);

    const response = await fetch(`${API_BASE_URL}/feedback/list?${query.toString()}`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch feedback list: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  async getFeedbackDetail(id: number): Promise<FeedbackItem> {
    const response = await fetch(`${API_BASE_URL}/feedback/detail?id=${id}`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch feedback detail: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  async updateFeedbackStatus(data: UpdateFeedbackStatusRequest): Promise<FeedbackItem> {
    const response = await fetch(`${API_BASE_URL}/feedback/status`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to update feedback status: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  async getFeedbackStatistics(params: {
    startDate?: string;
    endDate?: string;
    category?: string;
  } = {}): Promise<FeedbackStatsResponse> {
    const query = new URLSearchParams();
    if (params.startDate) query.append('startDate', params.startDate);
    if (params.endDate) query.append('endDate', params.endDate);
    if (params.category) query.append('category', params.category);

    const response = await fetch(`${API_BASE_URL}/feedback/statistics?${query.toString()}`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch feedback statistics: ${response.status} ${errorText}`);
    }
    return response.json();
  }
}

export const feedbackApi = new FeedbackApiService();
