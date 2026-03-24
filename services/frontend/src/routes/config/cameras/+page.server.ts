import type { PageServerLoad, Actions } from './$types';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ url, locals }) => {
    const page = Number(url.searchParams.get('page')) || 1;
    const limit = 10;
    const search = url.searchParams.get('search') || '';
    
    try {
        const [camerasRes, posts] = await Promise.all([
            locals.coreClient.listCameras(page, limit, search),
            locals.coreClient.listData('posts', 1, 100)
        ]);
        
        return {
            cameras: camerasRes.data || [],
            pagination: camerasRes.metadata,
            posts: posts.data || []
        };
    } catch (error: any) {
        return {
            cameras: [],
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
        const type = formData.get('type') as 'front' | 'back';
        const match_permit = formData.get('match_permit') === 'on';
        const customs_post_id = formData.get('customs_post_id') ? parseInt(formData.get('customs_post_id') as string) : undefined;

        try {
            const result = await locals.coreClient.createCamera({
                name,
                description,
                type,
                match_permit,
                customs_post_id
            });
            return { success: true, api_key: result.api_key, camera: result.camera };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    },
    update: async ({ request, locals }) => {
        const formData = await request.formData();
        const id = formData.get('id') as string;
        const name = formData.get('name') as string;
        const description = formData.get('description') as string;
        const type = formData.get('type') as 'front' | 'back';
        const match_permit = formData.get('match_permit') === 'on';
        const customs_post_id = formData.get('customs_post_id') ? parseInt(formData.get('customs_post_id') as string) : undefined;

        try {
            await locals.coreClient.updateCamera(id, {
                name,
                description,
                type,
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
            await locals.coreClient.deleteCamera(id);
            return { success: true };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    },
    regenerate: async ({ request, locals }) => {
        const formData = await request.formData();
        const id = formData.get('id') as string;

        try {
            const result = await locals.coreClient.regenerateCameraKey(id);
            return { success: true, api_key: result.api_key, camera: result.camera };
        } catch (error: any) {
            return fail(500, { error: error.message });
        }
    }
};
