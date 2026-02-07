import { describe, expect, it } from 'vitest';
import {
  computeDistributorStats,
  filterDistributors,
  withTimeout,
  type DistributorRecord,
} from '../../views/DistributorManagementView';

const mockDistributors: DistributorRecord[] = [
  {
    id: 1,
    username: 'alice',
    brandName: 'Brand A',
    level: 1,
    status: 'active',
  },
  {
    id: 2,
    username: 'bob',
    brandName: 'Brand B',
    level: 2,
    status: 'suspended',
  },
  {
    id: 3,
    username: 'carol',
    brandName: 'Alpha',
    level: 3,
    status: 'pending',
  },
];

describe('DistributorManagement helpers', () => {
  it('computes distributor stats', () => {
    expect(computeDistributorStats(mockDistributors)).toEqual({
      total: 3,
      active: 1,
      suspended: 1,
      pending: 1,
    });
  });

  it('filters by status and level', () => {
    const filtered = filterDistributors(mockDistributors, 'active', 1, '');
    expect(filtered).toHaveLength(1);
    expect(filtered[0].username).toBe('alice');
  });

  it('filters keyword case-insensitively and trims whitespace', () => {
    const filtered = filterDistributors(mockDistributors, '', 0, '  brAnd b  ');
    expect(filtered).toHaveLength(1);
    expect(filtered[0].username).toBe('bob');
  });

  it('resolves promise before timeout', async () => {
    await expect(withTimeout(Promise.resolve('ok'), 50)).resolves.toBe('ok');
  });

  it('rejects with timeout error when promise is too slow', async () => {
    const slowPromise = new Promise<string>((resolve) => {
      setTimeout(() => resolve('late'), 30);
    });

    await expect(withTimeout(slowPromise, 5)).rejects.toMatchObject({
      name: 'TimeoutError',
      message: 'REQUEST_TIMEOUT',
    });
  });
});
