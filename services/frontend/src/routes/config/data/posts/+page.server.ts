import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { CustomsPost } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 10;
	const name = url.searchParams.get('name') || undefined;

	try {
		const { data, metadata } = await locals.coreClient.listData<CustomsPost>('posts', page, limit, { name });
		return {
			posts: data,
			pagination: metadata,
			error: null
		};
	} catch (e: any) {
		return {
			posts: [],
			pagination: null,
			error: e.message || 'Помилка завантаження даних'
		};
	}
};

export const actions: Actions = {
	create: async ({ request, locals }) => {
		const formData = await request.formData();
		const name = formData.get('name') as string;
		const description = formData.get('description') as string;

		if (!name) {
			return fail(400, { message: 'Назва обов\'язкова' });
		}

		try {
			await locals.coreClient.createData('posts', { name, description });
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка створення' });
		}
	},
	update: async ({ request, locals }) => {
		const formData = await request.formData();
		const id = formData.get('id') as string;
		const name = formData.get('name') as string;
		const description = formData.get('description') as string;

		if (!id || !name) {
			return fail(400, { message: 'ID та Назва обов\'язкові' });
		}

		try {
			await locals.coreClient.updateData('posts', id, { name, description });
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка оновлення' });
		}
	},
	delete: async ({ request, locals }) => {
		const formData = await request.formData();
		const id = formData.get('id') as string;

		if (!id) {
			return fail(400, { message: 'ID обов\'язковий' });
		}

		try {
			await locals.coreClient.deleteData('posts', id);
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка видалення' });
		}
	}
};
