import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import type { Permit } from '$lib/types/permits';

export const GET: RequestHandler = async ({ locals, url }) => {
    if (!locals.user || !locals.coreClient) {
        return json({ error: 'Unauthorized' }, { status: 401 });
    }

    try {
        const page = parseInt(url.searchParams.get('page') || '1');
        const limit = parseInt(url.searchParams.get('limit') || '5000');
        
        const filters: Record<string, string> = {};
        const keys = [
            'sort_field', 'sort_order', 'filter_from', 'filter_to',
            'filter_post_id', 'filter_vehicle_type', 'filter_payment_type', 'filter_payer',
            'search', 'is_closed', 'custom_filters'
        ];

        for (const key of keys) {
            const val = url.searchParams.get(key);
            if (val !== null && val !== undefined) {
                filters[key] = val;
            }
        }

        const permitsRes = await locals.coreClient.getPermits<Permit>(page, limit, filters);
        return json(permitsRes);
    } catch (e: any) {
        console.error('Failed to fetch permits for export:', e);
        return json({ error: e.message }, { status: 500 });
    }
};
