import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { Company } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 10;
	const name = url.searchParams.get('name') || undefined;
	const edrpou = url.searchParams.get('edrpou') || undefined;

	try {
		const { data, metadata } = await locals.coreClient.listData<Company>('companies', page, limit, { name, edrpou });
		return {
			companies: data,
			pagination: metadata,
			error: null
		};
	} catch (e: any) {
		return {
			companies: [],
			pagination: null,
			error: e.message || 'Помилка завантаження даних'
		};
	}
};

export const actions: Actions = {
	create: async ({ request, locals }) => {
		const formData = await request.formData();
		const name = formData.get('name') as string;
		const edrpou = formData.get('edrpou') as string;

		if (!name || !edrpou) {
			return fail(400, { message: 'Назва та ЄДРПОУ обов\'язкові' });
		}

		try {
			await locals.coreClient.createData('companies', { 
				name, edrpou, details: {} 
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
			await locals.coreClient.deleteData('companies', id);
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка видалення' });
		}
	}
};
