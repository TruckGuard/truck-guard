import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { Company } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, params }) => {
    if (!locals.coreClient) {
        throw error(401, 'Unauthorized');
    }

	try {
		const company = await locals.coreClient.getData<Company>('companies', params.id);
		return { company };
	} catch (e: any) {
        if (e.status) {
            throw error(e.status, e.data?.error || 'Помилка при отриманні даних компанії');
        }
		console.error('Failed to load company details:', e);
        throw error(500, 'Не вдалося завантажити деталі компанії');
	}
};

export const actions: Actions = {
	update: async ({ request, locals, params }) => {
		const formData = await request.formData();
		const name = formData.get('name') as string;
		const edrpou = formData.get('edrpou') as string;
		const detailsStr = formData.get('details') as string;

		if (!name || !edrpou) {
			return fail(400, { message: 'Назва та ЄДРПОУ обов\'язкові' });
		}

		let details = {};
		try {
			details = JSON.parse(detailsStr || '{}');
		} catch (e) {
			return fail(400, { message: 'Некоректний формат JSON деталей' });
		}

		try {
			await locals.coreClient.updateData('companies', params.id, { 
				name, edrpou, details 
			});
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка оновлення' });
		}
	}
};
