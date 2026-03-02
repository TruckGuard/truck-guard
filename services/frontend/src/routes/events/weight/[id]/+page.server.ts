import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, locals }) => {
    if (!locals.coreClient) {
        throw error(401, 'Unauthorized');
    }

    const { id } = params;
    try {
        const event = await locals.coreClient.getWeightEvent(id);
        if (!event) {
            throw error(404, 'Event not found');
        }
        return { event };
    } catch (e) {
        console.error('Failed to load weight event:', e);
        throw error(500, 'Failed to load event details');
    }
};
