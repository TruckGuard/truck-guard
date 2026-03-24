import { error, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ params, url, locals }) => {
    const { id } = params;
    const page = Number(url.searchParams.get('page')) || 1;
    const limit = 10;

    if (!locals.coreClient) throw error(401, 'Unauthorized');

    try {
        const scale = await locals.coreClient.getScale(id);
        
        // Fetch events for this scale
        const events = await locals.coreClient.getEvents('weight', page, limit, { scale_id: scale.scale_id });
        const postsRes = await locals.coreClient.listData<any>('posts', 1, 100);

        return {
            scale,
            events,
            posts: postsRes.data || [],
            page,
            limit
        };
    } catch (e: any) {
        console.error('Failed to load scale details:', e);
        throw error(e.status || 500, e.message || 'Failed to load scale details');
    }
};

export const actions: Actions = {
    update: async ({ request, params, locals }) => {
        const { id } = params;
        const formData = await request.formData();
        const data = Object.fromEntries(formData);
        
        // Convert checkbox to boolean and identifiers to numeric
        const payload = {
            ...data,
            match_permit: data.match_permit === 'on',
            customs_post_id: data.customs_post_id && data.customs_post_id !== 'null' ? parseInt(data.customs_post_id as string) : null
        };

        if (!locals.coreClient) return fail(401, { error: 'Unauthorized' });

        try {
            await locals.coreClient.updateScale(id, payload);
            return { success: true };
        } catch (e: any) {
            return fail(e.status || 500, { error: e.message || 'Failed to update scale' });
        }
    }
};
