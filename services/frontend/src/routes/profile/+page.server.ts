import type { Actions, PageServerLoad } from './$types';
import { fail, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
    let profile: any = null;
    try {
        profile = (await locals.coreClient.getMyProfile());
    } catch (e) {
        console.error("Failed to fetch profile", e);
    }

    const postsResponse = await locals.coreClient.listData<any>('posts', 1, 100);

    // Fetch active sessions
    let sessions: any[] = [];
    try {
        sessions = await locals.authClient.listSessions();
    } catch (e) {
        console.error("Failed to fetch sessions", e);
    }

    return {
        user: locals.user,
        profile,
        posts: postsResponse.data,
        sessions
    };
};

export const actions: Actions = {
    updateProfile: async ({ request, locals }) => {
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
    },
    revokeAllSessions: async ({ locals, cookies }) => {
        try {
            await locals.authClient.revokeAllSessions();
        } catch (e) {
            console.error("Failed to revoke all sessions", e);
            return fail(500, { error: 'Failed to revoke all sessions' });
        }

        cookies.delete('session', { path: '/' });
        throw redirect(303, '/login');
    },
    changePassword: async ({ request, locals }) => {
        const data = await request.formData();
        const currentPass = data.get('current_password') as string;
        const newPass = data.get('new_password') as string;
        const confirmPass = data.get('confirm_password') as string;

        if (newPass !== confirmPass) {
            return fail(400, { error: 'Паролі не збігаються' });
        }

        try {
            await locals.authClient.changePassword(currentPass, newPass);
            return { success: true };
        } catch (e: any) {
            console.error("Change password error", e);
            return fail(e.status || 500, { error: e.data?.error || 'Помилка зміни пароля' });
        }
    }
};
