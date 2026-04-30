import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals }) => {
    const type = url.searchParams.get('type') || 'permits';
    const page = Number(url.searchParams.get('page')) || 1;
    const limit = Number(url.searchParams.get('limit')) || 20;

    const filters: Record<string, string | undefined> = {
        action: url.searchParams.get('action') || undefined,
        user_id: url.searchParams.get('user_id') || undefined,
        filter_from: url.searchParams.get('filter_from') || undefined,
        filter_to: url.searchParams.get('filter_to') || undefined,
    };

    let auditEvents = {
        data: [] as any[],
        metadata: {
            total_items: 0,
            total_pages: 0,
            current_page: page,
            limit,
        }
    };

    try {
        if (type === 'system') {
            const res = await locals.coreClient.getSystemAuditEvents(page, limit);
            if (res) auditEvents = res;
        } else {
            const res = await locals.coreClient.getAuditEvents(page, limit, filters);
            if (res) auditEvents = res;
        }
    } catch (e) {
        console.error('Failed to fetch audit events:', e);
    }

    return { auditEvents, page, limit, type };
};
