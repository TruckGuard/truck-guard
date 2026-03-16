import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import type { Permit } from '$lib/types/permits';
import type { Company, CustomsMode, PaymentType, VehicleType } from '$lib/types/data';

export const load: PageServerLoad = async ({ params, url, locals }) => {
    const id = params.id;
    if (id === 'new') {
        throw redirect(303, '/permits');
    }

    if (!locals.coreClient) {
        throw redirect(303, '/login');
    }

    let permit: Partial<Permit> = {};

    try {
        const fetched = await locals.coreClient.getPermit<Permit>(id);
        if (fetched) {
            permit = fetched;
        } else {
            throw error(404, 'Перепустку не знайдено');
        }
    } catch (e: any) {
        if (e.status === 303 || e.status === 404 || e.status === 500) {
            if (e.status === 404) throw error(404, 'Перепустку не знайдено');
            if (e.status === 303) throw e;
            throw error(e.status, e.message || 'Помилка завантаження');
        }
        throw error(e.status || 500, e.message || 'Не вдалося завантажити перепустку');
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


    console.log('permit', permit);
    return {
        permit,
        isNew: false,
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
    },

    restore: async ({ locals, params }) => {
        if (!locals.coreClient || params.id === 'new') return { type: 'error', error: { message: 'Invalid operation' } };

        try {
            const result = await locals.coreClient.restorePermit(params.id);
            return { type: 'success', data: result };
        } catch (e: any) {
            console.error('Restore permit error:', e);
            return { type: 'error', error: { message: e.message || 'Error restoring permit' } };
        }
    },

    void: async ({ locals, params }) => {
        if (!locals.coreClient || params.id === 'new') return { type: 'error', error: { message: 'Invalid operation' } };

        try {
            const result = await locals.coreClient.voidPermit(params.id);
            return { type: 'success', data: result };
        } catch (e: any) {
            console.error('Void permit error:', e);
            return { type: 'error', error: { message: e.message || 'Error voiding permit' } };
        }
    },

    delete: async ({ locals, params }) => {
        if (!locals.coreClient || params.id === 'new') return { type: 'error', error: { message: 'Invalid operation' } };

        try {
            await locals.coreClient.deletePermit(params.id);
            return { type: 'success' };
        } catch (e: any) {
            console.error('Delete permit error:', e);
            return { type: 'error', error: { message: e.message || 'Error deleting permit' } };
        }
    },

    createCompany: async ({ request, locals }) => {
        if (!locals.coreClient) return { type: 'error', error: { message: 'Unauthorized' } };

        const data = await request.formData();
        const name = data.get('name') as string;
        const edrpou = data.get('edrpou') as string;

        if (!name || !edrpou) {
            return { type: 'error', error: { message: 'Назва та ЄДРПОУ обов’язкові' } };
        }

        try {
            const result = await locals.coreClient.createData<Company>('companies', {
                name,
                edrpou
            });
            return { type: 'success', data: result };
        } catch (e: any) {
            console.error('Create company error:', e);
            return { type: 'error', error: { message: e.message || 'Error creating company' } };
        }
    },

    searchCompanies: async ({ request, locals }) => {
        if (!locals.coreClient) return { type: 'error', error: { message: 'Unauthorized' } };

        const data = await request.formData();
        const q = data.get('q') as string || '';

        try {
            let filters: Record<string, string> = {};
            if (q.trim()) {
                filters.q = q;
            }
            

            const res = await locals.coreClient.listData<Company>('companies', 1, 15, filters);
            return { type: 'success', data: res.data };
        } catch (e: any) {
            console.error('Search companies error:', e);
            return { type: 'error', error: { message: e.message || 'Error searching companies' } };
        }
    },

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
    },
    unlinkEntity: async ({ request, locals }) => {
        if (!locals.coreClient) return { type: 'error', error: { message: 'Unauthorized' } };
        const data = await request.formData();
        const permitId = Number(data.get('permit_id'));
        const eventId = Number(data.get('event_id'));
        const eventType = data.get('event_type') as 'plate' | 'weight';

        try {
            await locals.coreClient.unlinkPermit(permitId, eventId, eventType);
            return { type: 'success' };
        } catch (e: any) {
            console.error('Unlink entity error:', e);
            return { type: 'error', error: { message: e.message || 'Помилка відв\'язування' } };
        }
    }
};
