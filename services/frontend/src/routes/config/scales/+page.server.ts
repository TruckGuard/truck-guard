import type { PageServerLoad, Actions } from './$types';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ url, locals }) => {
    const page = Number(url.searchParams.get('page')) || 1;
    const limit = 10;
    const search = url.searchParams.get('search') || '';
    
    try {
        const [scalesRes, posts] = await Promise.all([
            locals.coreClient.listScales(page, limit, search),
            locals.coreClient.listData('posts', 1, 100)
        ]);
        
        return {
            scales: scalesRes.data || [],
            pagination: scalesRes.metadata,
            posts: posts.data || []
        };
    } catch (error: any) {
        return {
            scales: [],
            posts: [],
            error: error.message
        };
    }
};

export const actions: Actions = {
    create: async ({ request, locals }) => {
        const formData = await request.formData();
        const name = formData.get('name') as string;
        const description = formData.get('description') as string;
        const match_permit = formData.get('match_permit') === 'on';
        const customs_post_id = formData.get('customs_post_id') ? parseInt(formData.get('customs_post_id') as string) : undefined;

        try {
            const result = await locals.coreClient.createScale({
                name,
                description,
                match_permit,
                customs_post_id
            });
            return { success: true, api_key: result.api_key, scale: result.scale };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    },
    update: async ({ request, locals }) => {
        const formData = await request.formData();
        const id = formData.get('id') as string;
        const name = formData.get('name') as string;
        const description = formData.get('description') as string;
        const match_permit = formData.get('match_permit') === 'on';
        const customs_post_id = formData.get('customs_post_id') ? parseInt(formData.get('customs_post_id') as string) : undefined;

        try {
            await locals.coreClient.updateScale(id, {
                name,
                description,
                match_permit,
                customs_post_id
            });
            return { success: true };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    },
    delete: async ({ request, locals }) => {
        const formData = await request.formData();
        const id = formData.get('id') as string;

        try {
            await locals.coreClient.deleteScale(id);
            return { success: true };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    },
    regenerate: async ({ request, locals }) => {
        const formData = await request.formData();
        const id = formData.get('id') as string;

        try {
            const result = await locals.coreClient.regenerateScaleKey(id);
            return { success: true, api_key: result.api_key, scale: result.scale };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    }
};
