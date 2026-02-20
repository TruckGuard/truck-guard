import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	try {
		const settings = await locals.coreClient.listSettings();
		return { 
			settings,
			error: null 
		};
	} catch (e: any) {
		return { 
			settings: [],
			error: e.message || 'Помилка завантаження налаштувань'
		};
	}
};

export const actions: Actions = {
	update: async ({ request, locals }) => {
		const formData = await request.formData();
		const key = formData.get('key') as string;
		const value = formData.get('value') as string;

		if (!key || value === null) {
			return fail(400, { message: 'Ключ та значення обов\'язкові' });
		}

		try {
			await locals.coreClient.updateSetting(key, value);
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка оновлення' });
		}
	}
};
