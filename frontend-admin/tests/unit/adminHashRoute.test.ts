import { describe, expect, it } from 'vitest';
import { resolveAdminHashRoute } from '../../utils/adminHashRoute';

describe('resolveAdminHashRoute', () => {
  it('normalizes legacy distributor approval route', () => {
    const route = resolveAdminHashRoute('#/distributor-approval');

    expect(route.activeTab).toBe('distributor-management');
    expect(route.normalizedHash).toBe('#/distributor-management');
  });

  it('resolves member sub routes correctly', () => {
    expect(resolveAdminHashRoute('#/members').memberRoute).toBe('list');
    expect(resolveAdminHashRoute('#/members/merge').memberRoute).toBe('merge');
    expect(resolveAdminHashRoute('#/members/export').memberRoute).toBe('export');
    expect(resolveAdminHashRoute('#/members/123').memberRoute).toBe('detail');
  });

  it('resolves known top-level tab routes', () => {
    const route = resolveAdminHashRoute('#/campaigns');
    expect(route.activeTab).toBe('campaigns');
  });

  it('returns empty resolution for unknown route', () => {
    expect(resolveAdminHashRoute('#/unknown')).toEqual({});
  });
});
