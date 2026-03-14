import type { Actions, PageServerLoad } from './$types';
import { error, fail } from '@sveltejs/kit';
import { can } from '$lib/auth';

export const load: PageServerLoad = async ({ params, locals }) => {
    if (!locals.coreClient) {
        throw error(401, 'Unauthorized');
    }

    const { id } = params;
    try {
        const event = await locals.coreClient.getPlateEvent(id);
        return { event };
    } catch (e: any) {
        if (e.status) {
            throw error(e.status, e.data?.error || 'Помилка при отриманні події');
        }
        console.error('Failed to load plate event:', e);
        throw error(500, 'Не вдалося завантажити деталі події');
    }
};

