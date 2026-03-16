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


export const actions: Actions = {
    linkEntity: async ({ request, locals }) => {
        if (!locals.coreClient) return { type: 'error', error: { message: 'Unauthorized' } };
        const data = await request.formData();
        const permitId = Number(data.get('permit_id'));
        const eventId = Number(data.get('event_id'));
        const eventType = data.get('event_type') as 'plate' | 'weight';

        try {
            await locals.coreClient.linkPermit(permitId, eventId, eventType);
            return { type: 'success' };
        } catch (e: any) {
            console.error('Link entity error:', e);
            return { type: 'error', error: { message: e.message || 'Помилка зв\'язування' } };
        }
    }
};
