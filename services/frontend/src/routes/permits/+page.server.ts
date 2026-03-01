import type { PageServerLoad } from './$types';
import type { ApiResponse } from '$lib/types/events';
import type { Permit } from '$lib/types/permits';

export const load: PageServerLoad = async ({ locals }) => {
    let activePermits: ApiResponse<Permit> = {
        data: [],
        metadata: {
            total_items: 0,
            total_pages: 0,
            current_page: 1,
            limit: 50
        }
    };

    let recentPlateEvents: ApiResponse<any> = {
        data: [],
        metadata: {
            total_items: 0,
            total_pages: 0,
            current_page: 1,
            limit: 10
        }
    };

    let recentWeightEvents: ApiResponse<any> = {
        data: [],
        metadata: {
            total_items: 0,
            total_pages: 0,
            current_page: 1,
            limit: 10
        }
    };

    if (locals.coreClient) {
        try {
            // Fetch active permits
            const permitsRes = await locals.coreClient.getPermits<Permit>(1, 50, { is_closed: 'false' });
            if (permitsRes) {
                // The backend API might return different casing, map to standard format if needed
                activePermits = permitsRes as unknown as ApiResponse<Permit>;
            }

            // Fetch recent plate events
            const plateRes = await locals.coreClient.getEvents<ApiResponse<any>>('plate', 1, 15);
            if (plateRes) recentPlateEvents = plateRes;

            // Fetch recent weight events
            const weightRes = await locals.coreClient.getEvents<ApiResponse<any>>('weight', 1, 15);
            if (weightRes) recentWeightEvents = weightRes;

        } catch (e) {
            console.error('Failed to fetch dashboard data:', e);
        }
    }

    return {
        activePermits,
        recentPlateEvents,
        recentWeightEvents
    };
};
