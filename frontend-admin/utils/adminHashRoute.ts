export type MemberRoute = 'list' | 'detail' | 'merge' | 'export';

export interface AdminRouteResolution {
  activeTab?: string;
  memberRoute?: MemberRoute;
  normalizedHash?: string;
}

const VALID_TABS = new Set([
  'dashboard',
  'users',
  'brands',
  'campaigns',
  'system',
  'distributor-management',
  'verification-records',
  'poster-records',
]);

export const resolveAdminHashRoute = (hash: string): AdminRouteResolution => {
  const normalizedHash = hash || '';

  if (normalizedHash.startsWith('#/distributor-approval')) {
    return {
      activeTab: 'distributor-management',
      normalizedHash: '#/distributor-management',
    };
  }

  if (normalizedHash.startsWith('#/members')) {
    let memberRoute: MemberRoute = 'list';
    if (normalizedHash.startsWith('#/members/merge')) {
      memberRoute = 'merge';
    } else if (normalizedHash.startsWith('#/members/export')) {
      memberRoute = 'export';
    } else if (/^#\/members\/\d+/.test(normalizedHash)) {
      memberRoute = 'detail';
    }

    return {
      activeTab: 'members',
      memberRoute,
    };
  }

  const tabFromHash = normalizedHash.replace('#/', '');
  if (VALID_TABS.has(tabFromHash)) {
    return {
      activeTab: tabFromHash,
    };
  }

  return {};
};
