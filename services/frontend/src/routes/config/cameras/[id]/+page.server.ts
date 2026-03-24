import { error, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ params, url, locals }) => {
    const { id } = params;
    const page = Number(url.searchParams.get('page')) || 1;
    const limit = 10;

    if (!locals.coreClient) throw error(401, 'Unauthorized');

    try {
        const camera = await locals.coreClient.getCamera(id);
        
        // Fetch events for this camera
        const events = await locals.coreClient.getEvents('plate', page, limit, { camera_id: camera.camera_id });
        const postsRes = await locals.coreClient.listData<any>('posts', 1, 100);

        return {
            camera,
            events,
            posts: postsRes.data || [],
            page,
            limit
        };
    } catch (e: any) {
        console.error('Failed to load camera details:', e);
        throw error(e.status || 500, e.message || 'Failed to load camera details');
    }
};

export const actions: Actions = {
    update: async ({ request, params, locals }) => {
        const { id } = params;
        const formData = await request.formData();
        const data = Object.fromEntries(formData);

        const payload = {
            ...data,
            match_permit: data.match_permit === 'on',
            run_anpr: data.run_anpr === 'on',
            customs_post_id: data.customs_post_id && data.customs_post_id !== 'null' ? parseInt(data.customs_post_id as string) : null
        };

        if (!locals.coreClient) return fail(401, { error: 'Unauthorized' });

        try {
            await locals.coreClient.updateCamera(id, payload);
            return { success: true };
        } catch (e: any) {
            return fail(e.status || 500, { error: e.message || 'Failed to update camera' });
        }
    }
};
