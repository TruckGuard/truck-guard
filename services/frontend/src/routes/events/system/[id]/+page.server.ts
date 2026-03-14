import type { PageServerLoad } from './$types';
import type { SystemEvent } from '$lib/types/events';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, locals }) => {
    if (!locals.coreClient) {
        throw error(401, 'Unauthorized');
    }

    const { id } = params;
    try {
        const event = await locals.coreClient.getSystemEvent<SystemEvent>(id);
        return { event };
    } catch (e: any) {
        if (e.status) {
            throw error(e.status, e.data?.error || 'Помилка при отриманні системної події');
        }
        console.error('Failed to fetch system event:', e);
        throw error(500, 'Не вдалося завантажити деталі системної події');
    }
};
