import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { Company } from '$lib/types/data';

export const load: PageServerLoad = async ({ locals, params }) => {
    if (!locals.coreClient) throw error(401, 'Unauthorized');
    try {
        const company = await locals.coreClient.getData<Company>('companies', params.id);
        return { company };
    } catch (e: any) {
        if (e.status) throw error(e.status, e.data?.error || 'Помилка при отриманні даних компанії');
        throw error(500, 'Не вдалося завантажити деталі компанії');
    }
};

export const actions: Actions = {
    update: async ({ request, locals, params }) => {
        const fd = await request.formData();
        const name = fd.get('name') as string;
        const edrpou = fd.get('edrpou') as string;
        const notes = (fd.get('notes') as string) || '';
        const detailsStr = fd.get('details') as string;
        const discountPct = parseFloat((fd.get('discount_percentage') as string) || '0');
        const discountFixed = parseFloat((fd.get('discount_fixed') as string) || '0');

        if (!name || !edrpou) return fail(400, { message: "Назва та ЄДРПОУ обов'язкові" });

        let details = {};
        try { details = JSON.parse(detailsStr || '{}'); } catch {
            return fail(400, { message: 'Некоректний формат JSON деталей' });
        }

        try {
            const updated = await locals.coreClient.updateData<Company>('companies', params.id, {
                name, edrpou, notes, details,
                discount_percentage: isNaN(discountPct) ? 0 : discountPct,
                discount_fixed: isNaN(discountFixed) ? 0 : discountFixed,
            });
            return { success: true, company: updated };
        } catch (e: any) {
            return fail(500, { message: e.message || 'Помилка оновлення' });
        }
    },

    syncEdrpou: async ({ locals, params }) => {
        const company = await locals.coreClient.getData<Company>('companies', params.id).catch(() => null);
        if (!company) return fail(404, { message: 'Компанію не знайдено' });

        let edrData: Awaited<ReturnType<typeof locals.coreClient.getCompanyByEdrpou>> | null = null;
        try {
            edrData = await locals.coreClient.getCompanyByEdrpou(company.edrpou);
        } catch (e: any) {
            return fail(502, { message: `Не вдалося отримати дані з ЄДР: ${e.message}` });
        }

        if (!edrData) return fail(502, { message: 'Порожня відповідь від сервісу ЄДР' });

        // Map customs-parser field names → Ukrainian labels used for display
        const fieldMap: Record<string, string> = {
            name: 'Повна назва', short_name: 'Скорочена назва',
            boss: 'Директор', address: 'Адреса',
            registration_date: 'Дата реєстрації', kved: 'КВЕД',
            state: 'Статус', phone: 'Телефон', email: 'Email',
            iban: 'IBAN', bank: 'Банк', mfo: 'МФО',
        };

        const merged: Record<string, any> = { ...(company.details || {}) };
        for (const [apiKey, label] of Object.entries(fieldMap)) {
            const val = (edrData as any)[apiKey];
            if (val !== undefined && val !== null && val !== '') merged[label] = val;
        }

        const now = new Date().toISOString();
        try {
            await locals.coreClient.updateData('companies', params.id, {
                ...company,
                details: merged,
                last_synced_at: now,
            });
            return { success: true, synced: true, details: merged, last_synced_at: now };
        } catch (e: any) {
            return fail(500, { message: e.message || 'Помилка збереження після синхронізації' });
        }
    },
};
