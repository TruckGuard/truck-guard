import type { Actions, PageServerLoad } from './$types';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
    let profile = {};
    try {
        profile = (await locals.coreClient.getMyProfile()) || {};
    } catch (e) {
        console.error("Failed to fetch profile", e);
    }

    const postsResponse = await locals.coreClient.listData<any>('posts', 1, 100);
    
    return {
        user: locals.user,
        profile,
        posts: postsResponse.data
    };
};

export const actions: Actions = {
    default: async ({ request, locals }) => {
        const data = await request.formData();
        
        const firstName = data.get('first_name') as string;
        const lastName = data.get('last_name') as string;
        const thirdName = data.get('third_name') as string;
        const phone = data.get('phone_number') as string;
        const email = data.get('email') as string;
        const notes = data.get('notes') as string;

        const customsPostId = data.get('customs_post_id');
        
        const updatedProfile = await locals.coreClient.updateMyProfile({
            first_name: firstName,
            last_name: lastName,
            third_name: thirdName,
            phone_number: phone,
            email,
            notes,
            customs_post_id: customsPostId ? Number(customsPostId) : null
        });

        if (!updatedProfile) {
            return fail(500, { error: 'Failed to update profile' });
        }

        return { success: true };
    }
};
