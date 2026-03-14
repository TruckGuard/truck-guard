import type { PageServerLoad, Actions } from './$types';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
    try {
        const settings = await locals.coreClient.listSettings();
        return {
            settings
        };
    } catch (error: any) {
        return {
            settings: [],
            error: error.message
        };
    }
};

export const actions: Actions = {
    default: async ({ request, locals }) => {
        const formData = await request.formData();
        const key = formData.get('key') as string;
        const value = formData.get('value') as string;

        if (!key || value === null) {
            return fail(400, { error: 'Key and Value are required' });
        }

        try {
            await locals.coreClient.updateSetting(key, value);
            return { success: true };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    }
};
