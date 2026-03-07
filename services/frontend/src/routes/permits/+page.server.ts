import type { PageServerLoad } from './$types';
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
                'search', 'is_closed', 'custom_filters'
            ];

            // If is_closed is not provided in URL, default to displaying only active permits
            if (!url.searchParams.has('is_closed')) {
                filters['is_closed'] = 'false';
            } else if (url.searchParams.get('is_closed') === 'all') {
                filters['is_closed'] = 'all';
            }

            for (const key of keys) {
                const val = url.searchParams.get(key);
                if (val !== null && val !== undefined && key !== 'is_closed') {
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

            const usersRes = await locals.coreClient.listData<any>('users', 1, 100);
            if (usersRes) {
                users = usersRes.data.map((u: any) => ({
                    ID: u.ID,
                    name: `${u.first_name} ${u.last_name}`
                }));
            }

            const modesRes = await locals.coreClient.listData<any>('customs-modes', 1, 100);
            if (modesRes) {
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
