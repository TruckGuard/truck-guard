import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
    const [stats, recentPermits] = await Promise.all([
        locals.coreClient.getStats().catch(() => null),
        locals.coreClient.getPermits(1, 5, { is_void: 'false' }).catch(() => ({ data: [], metadata: null })),
    ]);

    return {
        stats,
        recentPermits: recentPermits.data || [],
    };
};
