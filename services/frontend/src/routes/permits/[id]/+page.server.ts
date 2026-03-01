import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import type { Permit } from '$lib/types/permits';
import type { Company, CustomsMode, PaymentType, VehicleType } from '$lib/types/data';

export const load: PageServerLoad = async ({ params, url, locals }) => {
    const id = params.id;
    const isNew = id === 'new';

    if (!locals.coreClient) {
        throw redirect(303, '/login');
    }

    let permit: Partial<Permit> = {
        is_closed: false,
        is_void: false
    };

    // Pre-fill from URL params if creating new
    if (isNew) {
        const plate = url.searchParams.get('plate');
        const plate_back = url.searchParams.get('plate_back');
        const weight = url.searchParams.get('weight');
        const camera_event_id = url.searchParams.get('camera_event_id');
        const scale_event_id = url.searchParams.get('scale_event_id');

        if (plate) permit.plate_front = plate;
        if (plate_back) permit.plate_back = plate_back;
        if (weight) permit.total_weight = Number(weight);

        // At this point we just set frontend states. Actual linking 
        // will require backend support (e.g passing event IDs on create)
        // For the sake of UI we pass them along.
        (permit as any)._initial_camera_event = camera_event_id;
        (permit as any)._initial_scale_event = scale_event_id;
    } else {
        try {
            const fetched = await locals.coreClient.getPermit<Permit>(id);
            if (fetched) {
                permit = fetched;
            } else {
                throw error(404, 'Перепустку не знайдено');
            }
        } catch (e: any) {
            throw error(e.status || 500, 'Не вдалося завантажити перепустку');
        }
    }

    // Load reference data for selects
    let paramsData = {
        vehicleTypes: [] as VehicleType[],
        customsModes: [] as CustomsMode[],
        companies: [] as Company[],
        paymentTypes: [] as PaymentType[]
    };

    try {
        const [vtReq, cmReq, cReq, ptReq] = await Promise.all([
            locals.coreClient.listData<VehicleType>('vehicle-types', 1, 100),
            locals.coreClient.listData<CustomsMode>('modes', 1, 100),
            locals.coreClient.listData<Company>('companies', 1, 100),
            locals.coreClient.listData<PaymentType>('payment-types', 1, 100)
        ]);

        paramsData.vehicleTypes = vtReq.data;
        paramsData.customsModes = cmReq.data;
        paramsData.companies = cReq.data;
        paramsData.paymentTypes = ptReq.data;
    } catch (e) {
        console.error('Failed to load dictionary data', e);
    }

    // Pass user permissions for UI toggling
    const userRole = locals.user?.role || '';
    const canValidate = typeof userRole === 'string' ? userRole === 'admin' || userRole === 'manager' : (userRole as any)?.name === 'admin' || (userRole as any)?.name === 'manager';
    const canManageCustoms = true; // Based on earlier seed, operator/manager/admin all have read:customs

    return {
        permit,
        isNew,
        userRole,
        canValidate,
        canManageCustoms,
        ...paramsData
    };
};

export const actions: Actions = {
    save: async ({ request, locals, params }) => {
        if (!locals.coreClient) return { type: 'error', error: { message: 'Unauthorized' } };

        const data = await request.formData();
        const payloadStr = data.get('data') as string;
        if (!payloadStr) return { type: 'error', error: { message: 'Empty payload' } };

        try {
            const payload = JSON.parse(payloadStr);
            let result;

            if (params.id === 'new') {
                result = await locals.coreClient.createPermit(payload);
            } else {
                console.log('Updating permit', params.id, payload);
                result = await locals.coreClient.updatePermit(params.id, payload);
            }

            return { type: 'success', data: result };
        } catch (e: any) {
            console.error('Save permit error:', e);
            return { type: 'error', error: { message: e.message || 'Error saving permit' } };
        }
    },

    close: async ({ locals, params }) => {
        if (!locals.coreClient || params.id === 'new') return { type: 'error', error: { message: 'Invalid operation' } };

        try {
            // To close a permit, we send a PUT with is_closed: true
            // Ensure this aligns with backend logic. 
            const result = await locals.coreClient.updatePermit(params.id, { is_closed: true });
            return { type: 'success', data: result };
        } catch (e: any) {
            console.error('Close permit error:', e);
            return { type: 'error', error: { message: e.message || 'Error closing permit' } };
        }
    },

    validate: async ({ locals, params }) => {
        if (!locals.coreClient || params.id === 'new') return { type: 'error', error: { message: 'Invalid operation' } };

        try {
            const result = await locals.coreClient.validatePermit(params.id);
            return { type: 'success', data: result };
        } catch (e: any) {
            console.error('Validate permit error:', e);
            return { type: 'error', error: { message: e.message || 'Error validating permit' } };
        }
    }
};
