import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { ExcludedPlate } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 10;
	const plate = url.searchParams.get('plate') || undefined;

	try {
		const { data, metadata } = await locals.coreClient.listData<ExcludedPlate>('excluded-plates', page, limit, { plate });
		return {
			plates: data,
			pagination: metadata,
			error: null
		};
	} catch (e: any) {
		return {
			plates: [],
			pagination: null,
			error: e.message || 'Помилка завантаження даних'
		};
	}
};

export const actions: Actions = {
	create: async ({ request, locals }) => {
		const formData = await request.formData();
		const plate = formData.get('plate') as string;
		const comment = formData.get('comment') as string;

		if (!plate) {
			return fail(400, { message: 'Номер обов\'язковий' });
		}

		try {
			await locals.coreClient.createData('excluded-plates', { 
				plate, comment: comment || ''
			});
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка створення' });
		}
	},
	delete: async ({ request, locals }) => {
		const formData = await request.formData();
		const id = formData.get('id') as string;

		if (!id) {
			return fail(400, { message: 'ID обов\'язковий' });
		}

		try {
			await locals.coreClient.deleteData('excluded-plates', id);
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка видалення' });
		}
	}
};
