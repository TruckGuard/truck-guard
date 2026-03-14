import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, locals }) => {
    if (!locals.coreClient) {
        throw error(401, 'Unauthorized');
    }

    const { id } = params;
    try {
        const event = await locals.coreClient.getWeightEvent(id);
        return { event };
    } catch (e: any) {
        if (e.status) {
            throw error(e.status, e.data?.error || 'Помилка при отриманні події зважування');
        }
        console.error('Failed to load weight event:', e);
        throw error(500, 'Не вдалося завантажити деталі події зважування');
    }
};
