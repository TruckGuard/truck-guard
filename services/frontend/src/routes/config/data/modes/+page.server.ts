import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { CustomsMode } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 10;
	const name = url.searchParams.get('name') || undefined;
	const code = url.searchParams.get('code') || undefined;

	try {
		const { data, metadata } = await locals.coreClient.listData<CustomsMode>('modes', page, limit, { name, code });
		return {
			modes: data,
			pagination: metadata,
			error: null
		};
	} catch (e: any) {
		return {
			modes: [],
			pagination: null,
			error: e.message || 'Помилка завантаження даних'
		};
	}
};

export const actions: Actions = {
	create: async ({ request, locals }) => {
		const formData = await request.formData();
		const name = formData.get('name') as string;
		const code = formData.get('code') as string;
		const description = formData.get('description') as string;
		const requiredFieldsRaw = formData.get('required_fields') as string;
		let required_fields: string[] = [];
		try { required_fields = JSON.parse(requiredFieldsRaw || '[]'); } catch { }

		if (!name || !code) {
			return fail(400, { message: 'Назва та Код обов\'язкові' });
		}

		try {
			let test = await locals.coreClient.createData('modes', { name, code, description, required_fields });
			console.log('Create mode success:', test);
			return { success: true };
		} catch (e: any) {
			console.error('Create mode failed:', e);
			return fail(500, { message: e.message || 'Помилка створення' });
		}
	},
	update: async ({ request, locals }) => {
		const formData = await request.formData();
		const id = formData.get('id') as string;
		const name = formData.get('name') as string;
		const code = formData.get('code') as string;
		const description = formData.get('description') as string;
		const requiredFieldsRaw = formData.get('required_fields') as string;
		let required_fields: string[] = [];
		try { required_fields = JSON.parse(requiredFieldsRaw || '[]'); } catch { }
		console.log('Update mode:', { id, name, code, description, required_fields });
		if (!id || !name || !code) {
			return fail(400, { message: 'ID, Назва та Код обов\'язкові' });
		}

		try {
			await locals.coreClient.updateData('modes', id, { name, code, description, required_fields });
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
			await locals.coreClient.deleteData('modes', id);
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка видалення' });
		}
	}
};
