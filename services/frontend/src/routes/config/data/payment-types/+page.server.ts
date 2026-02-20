import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { PaymentType } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 10;
	const name = url.searchParams.get('name') || undefined;

	try {
		const { data, metadata } = await locals.coreClient.listData<PaymentType>('payment-types', page, limit, { name });
		return {
			paymentTypes: data,
			pagination: metadata,
			error: null
		};
	} catch (e: any) {
		return {
			paymentTypes: [],
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
		const is_active = formData.get('is_active') === 'on';
		const icon = formData.get('icon') as string;

		if (!name || !code) {
			return fail(400, { message: 'Назва та Код обов\'язкові' });
		}

		try {
			await locals.coreClient.createData('payment-types', { 
				name, code, description, is_active, icon 
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
		const is_active = formData.get('is_active') === 'on';
		const icon = formData.get('icon') as string;

		if (!id || !name || !code) {
			return fail(400, { message: 'ID, Назва та Код обов\'язкові' });
		}

		try {
			await locals.coreClient.updateData('payment-types', id, { 
				name, code, description, is_active, icon 
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
			await locals.coreClient.deleteData('payment-types', id);
			return { success: true };
		} catch (e: any) {
			return fail(500, { message: e.message || 'Помилка видалення' });
		}
	}
};
