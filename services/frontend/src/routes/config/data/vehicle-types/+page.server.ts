import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { VehicleType } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 10;
	const name = url.searchParams.get('name') || undefined;

	try {
		const { data, metadata } = await locals.coreClient.listData<VehicleType>('vehicle-types', page, limit, { name });
		return {
			vehicleTypes: data,
			pagination: metadata,
			error: null
		};
	} catch (e: any) {
		return {
			vehicleTypes: [],
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
		const entry_price = Number(formData.get('entry_price'));
		const daily_price = Number(formData.get('daily_price'));
		const color = formData.get('color') as string;

		if (!name || !code) {
			return fail(400, { message: 'Назва та Код обов\'язкові' });
		}

		try {
			await locals.coreClient.createData('vehicle-types', { 
				name, code, description, entry_price, daily_price, color 
			});
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка створення' });
		}
	},
	update: async ({ request, locals }) => {
		const formData = await request.formData();
		const id = formData.get('id') as string;
		const name = formData.get('name') as string;
		const code = formData.get('code') as string;
		const description = formData.get('description') as string;
		const entry_price = Number(formData.get('entry_price'));
		const daily_price = Number(formData.get('daily_price'));
		const color = formData.get('color') as string;

		if (!id || !name || !code) {
			return fail(400, { message: 'ID, Назва та Код обов\'язкові' });
		}

		try {
			await locals.coreClient.updateData('vehicle-types', id, { 
				name, code, description, entry_price, daily_price, color 
			});
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
			await locals.coreClient.deleteData('vehicle-types', id);
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка видалення' });
		}
	}
};
