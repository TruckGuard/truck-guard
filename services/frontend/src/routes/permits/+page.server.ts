import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import type { ApiResponse } from '$lib/types/events';
import type { CustomsPost, VehicleType } from '$lib/types/data';
import type { Permit } from '$lib/types/permits';

export const load: PageServerLoad = async ({ locals, url }) => {
    let activePermits: ApiResponse<Permit> = {
        data: [],
        metadata: {
            total_items: 0,
            total_pages: 0,
            current_page: 1,
            limit: 50
        }
    };

    let customsPosts: CustomsPost[] = [];
    let vehicleTypes: VehicleType[] = [];
    let paymentTypes: any[] = [];
    let customsModes: any[] = [];
    let users: any[] = [];
    let hasAllPermitsAccess = false;

    if (locals.user && locals.user.permissions) {
        hasAllPermitsAccess = locals.user.permissions.some((p: string) =>
            p === 'admin' || p === '*:*' || p.includes(':permits:all') || p.includes('*:permits:all')
        );
    }

    if (locals.coreClient) {
        try {
            // Extract pagination
            const page = parseInt(url.searchParams.get('page') || '1');
            const limit = parseInt(url.searchParams.get('limit') || '50');

            // Build filters record from searchParams
            const filters: Record<string, string> = {};
            const keys = [
                'sort_field', 'sort_order', 'filter_from', 'filter_to',
                'filter_post_id', 'filter_vehicle_type', 'filter_payment_type', 'filter_payer',
                'search', 'is_closed', 'is_void', 'custom_filters'
            ];

            // Pass is_closed filter directly
            const isClosedVal = url.searchParams.get('is_closed');
  
            filters['is_closed'] = isClosedVal || "";

            const isVoidVal = url.searchParams.get('is_void');
            if (isVoidVal === 'show_all') {
                filters['is_void'] = ""; // Show all
            } else {
                filters['is_void'] = isVoidVal ?? "false"; // Default to false
            }
            

            for (const key of keys) {
                const val = url.searchParams.get(key);
                if (val !== null && val !== undefined && key !== 'is_closed' && key !== 'is_void') {
                    filters[key] = val;
                }
            }

            // Fetch active permits
            const permitsRes = await locals.coreClient.getPermits<Permit>(page, limit, filters);
            if (permitsRes) {
                activePermits = permitsRes as unknown as ApiResponse<Permit>;
            }

            // Fetch dictionaries for filters
            // Posts are needed if user has 'all' access, but we fetch them just in case
            const postsRes = await locals.coreClient.listData<CustomsPost>('posts', 1, 100);
            if (postsRes) {
                customsPosts = postsRes.data;
            }

            const vehRes = await locals.coreClient.listData<VehicleType>('vehicle-types', 1, 100);
            if (vehRes) {
                vehicleTypes = vehRes.data;
            }

            const ptRes = await locals.coreClient.listData<any>('payment-types', 1, 100);
            if (ptRes) {
                paymentTypes = ptRes.data;
            }

            const usersList = await locals.coreClient.listUsers();
            if (usersList) {
                users = usersList.map((u: any) => ({
                    ID: u.ID,
                    name: `${u.first_name} ${u.last_name}`
                }));
            }

            const modesRes = await locals.coreClient.listData<any>('modes', 1, 100);
            if (modesRes) {
                console.log(modesRes);
                customsModes = modesRes.data;
            }
        } catch (e) {
            console.error('Failed to fetch dashboard data:', e);
        }
    }

    return {
        activePermits,
        customsPosts,
        vehicleTypes,
        paymentTypes,
        customsModes,
        users,
        hasAllPermitsAccess
    };
};

export const actions: Actions = {
    create: async ({ locals, request }) => {
        if (!locals.coreClient) {
            throw error(401, 'Unauthorized');
        }

        const formData = await request.formData();
        const camera_event_id = formData.get('camera_event_id');
        const scale_event_id = formData.get('scale_event_id');

        try {
            const payload: any = {};
            if (camera_event_id) payload.camera_event_id = Number(camera_event_id);
            if (scale_event_id) payload.scale_event_id = Number(scale_event_id);

            // Create an empty permit. Backend will generate code and set CreatedBy.
            const result = await locals.coreClient.createPermit<any>(payload);
            if (result && result.ID) {
                throw redirect(303, `/permits/${result.ID}`);
            }
            return { type: 'error', error: 'Failed to create permit' };
        } catch (e: any) {
            if (e.status === 303) throw e;
            console.error('Action create error:', e);
            return { type: 'error', error: e.message || 'Error creating permit' };
        }
    }
};
