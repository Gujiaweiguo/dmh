
export enum TaskType {
  PAYMENT = 'PAYMENT',
  SHARE = 'SHARE',
  EXAM = 'EXAM'
}

// ============================================
// 权限管理相关类型
// ============================================

// 用户角色
export type UserRole = 'platform_admin' | 'participant' | 'anonymous';

// 登录请求
export interface LoginRequest {
  username: string;
  password: string;
}

// 注册请求
export interface RegisterRequest {
  username: string;
  password: string;
  phone: string;
  email?: string;
  realName?: string;
}

// 登录响应
export interface LoginResponse {
  token: string;
  userId: number;
  username: string;
  phone: string;
  realName: string;
  roles: UserRole[];
  brandIds?: number[];
}

// 当前用户信息
export interface CurrentUser {
  id: number;
  username: string;
  phone: string;
  email: string;
  realName: string;
  avatar: string;
  status: 'active' | 'disabled';
  roles: UserRole[];
  brandIds?: number[];
}

// 品牌信息
export interface BrandInfo {
  id: number;
  name: string;
  logo: string;
  description: string;
  status: 'active' | 'disabled';
  createdAt: string;
}

// 提现申请
export interface WithdrawalRequest {
  amount: number;
  bankName: string;
  bankAccount: string;
  accountName: string;
}

// 提现记录
export interface WithdrawalRecord {
  id: number;
  userId: number;
  username?: string;
  realName?: string;
  amount: number;
  bankName: string;
  bankAccount: string;
  accountName: string;
  status: 'pending' | 'approved' | 'rejected';
  remark?: string;
  approvedBy?: number;
  approvedAt?: string;
  createdAt: string;
}

// ============================================
// 业务相关类型
// ============================================

// Campaign 营销活动
export interface Campaign {
  id: number;
  name: string;
  description: string;
  formFields: string[];
  rewardRule: number;
  startTime: string;
  endTime: string;
  status: 'active' | 'paused' | 'ended';
  createdAt: string;
}

export interface Industry {
  id: string;
  name: string;
  icon: string;
  brandCount: number;
}

export interface Brand {
  id: string;
  name: string;
  industryId: string;
  logo: string;
  campaignCount: number;
}

export interface ContentBlock {
  id: string;
  type: 'text' | 'image' | 'heading';
  value: string;
}

export interface FormField {
  id: string;
  label: string;
  type: 'string' | 'number' | 'select';
  required: boolean;
  options?: string[];
}

// Content Library Item: Purely UI/UX structure and data capture definition
export interface ContentLibraryItem {
  id: string;
  brandId: string;
  title: string;
  description: string;
  imageUrl: string;
  videoUrl?: string;
  contentBlocks: ContentBlock[];
  formSchema: FormField[];
  successImageUrl?: string;
}

// Marketing Task: Orchestrates content with business logic (rewards, fees)
export interface MarketingTask {
  id: string;
  contentId: string; // References ContentLibraryItem
  title: string;
  type: TaskType;
  rewardAmount: number;
  entryFee: number;
  status: 'ACTIVE' | 'PAUSED';
  description: string;
  socialPost?: string;
  createdAt: string;
}

export interface Order {
  id: string;
  taskId: string;
  taskTitle?: string;
  userId: string;
  referrerId?: string;
  status: 'PAID' | 'PENDING';
  amount: number;
  formData: Record<string, any>;
  createdAt: string;
}

export interface Withdrawal {
  id: string;
  userId: string;
  userName: string;
  amount: number;
  status: 'PENDING' | 'APPROVED' | 'REJECTED';
  timestamp: string;
  paymentMethod: string;
}

export interface UserProfile {
  id: string;
  name: string;
  email: string;
  phone: string;
  avatarUrl: string;
  accumulatedRewards: number;
  referralCount: number;
  totalShares: number; 
  totalSpent: number;
  history: any[];
  joinDate: string;
}

export interface GenAIConfig {
  endpoint: string;
  apiKey: string;
  model: string;
}
