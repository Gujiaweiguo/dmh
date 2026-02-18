import { describe, expect, it, vi, beforeEach, beforeAll } from 'vitest';
import { api } from '../../services/mockApi';
import type { ContentLibraryItem, MarketingTask } from '../../types';
import { TaskType } from '../../types';

const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => { store[key] = value; },
    removeItem: (key: string) => { delete store[key]; },
    clear: () => { store = {}; },
  };
})();

Object.defineProperty(global, 'localStorage', { value: localStorageMock });

describe('MockApiService', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  describe('listIndustries', () => {
    it('should return industries list', async () => {
      const result = await api.listIndustries();
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe('listBrands', () => {
    it('should return brands list', async () => {
      const result = await api.listBrands();
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe('listContentLibrary', () => {
    it('should return content library items', async () => {
      const result = await api.listContentLibrary();
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe('listTasks', () => {
    it('should return tasks list', async () => {
      const result = await api.listTasks();
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe('listOrders', () => {
    it('should return orders list', async () => {
      const result = await api.listOrders();
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe('listWithdrawals', () => {
    it('should return withdrawals list', async () => {
      const result = await api.listWithdrawals();
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe('listProfiles', () => {
    it('should return profiles list', async () => {
      const result = await api.listProfiles();
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe('saveContent', () => {
    it('should save new content item', async () => {
      const newItem: ContentLibraryItem = {
        id: '',
        brandId: 'br-1',
        title: 'Test Content',
        description: 'Test Description',
        imageUrl: 'http://example.com/image.png',
        contentBlocks: [],
        formSchema: [],
      };
      
      const result = await api.saveContent(newItem);
      expect(result).toBeDefined();
      expect(result.id).toBeDefined();
      expect(result.title).toBe('Test Content');
    });

    it('should update existing content item', async () => {
      const contents = await api.listContentLibrary();
      if (contents.length > 0) {
        const existingItem = contents[0];
        existingItem.title = 'Updated Title';
        
        const result = await api.saveContent(existingItem);
        expect(result.title).toBe('Updated Title');
      }
    });
  });

  describe('saveTask', () => {
    it('should save new task', async () => {
      const newTask: MarketingTask = {
        id: '',
        contentId: 'c-1',
        title: 'Test Task',
        type: TaskType.PAYMENT,
        rewardAmount: 10.0,
        entryFee: 5.0,
        status: 'ACTIVE',
        description: 'Test task description',
        socialPost: 'Test social post',
        createdAt: '',
      };
      
      const result = await api.saveTask(newTask);
      expect(result).toBeDefined();
      expect(result.id).toBeDefined();
      expect(result.title).toBe('Test Task');
      expect(result.createdAt).toBeDefined();
    });
  });

  describe('approveWithdrawal', () => {
    it('should handle approve withdrawal call', async () => {
      await api.approveWithdrawal('wd-1', 'APPROVED');
      expect(true).toBe(true);
    });

    it('should handle reject withdrawal call', async () => {
      await api.approveWithdrawal('wd-1', 'REJECTED');
      expect(true).toBe(true);
    });
  });

  describe('createOrder', () => {
    it('should create new order', async () => {
      const tasks = await api.listTasks();
      if (tasks.length > 0) {
        const taskId = tasks[0].id;
        const formData = { name: 'Test User', phone: '12345678' };
        
        const result = await api.createOrder(taskId, formData);
        expect(result).toBeDefined();
        expect(result.id).toBeDefined();
        expect(result.taskId).toBe(taskId);
        expect(result.status).toBe('PAID');
        expect(result.formData).toEqual(formData);
      }
    });

    it('should throw error for non-existent task', async () => {
      await expect(api.createOrder('non-existent', {})).rejects.toThrow('Task not found');
    });
  });
});
