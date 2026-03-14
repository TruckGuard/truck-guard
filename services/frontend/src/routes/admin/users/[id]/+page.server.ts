import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { can } from '$lib/auth';

export const load: PageServerLoad = async ({ params, locals, cookies }) => {
     if (!can(locals.user, 'update:users')) {
        throw redirect(303, '/admin/users');
    }

    const { id } = params;

    const roles = await locals.authClient.getRoles();

    const allAuthUsers = await locals.authClient.listUsers();
    const authUser = allAuthUsers.find((u: { id: string; }) => u.id == id);

    if (!authUser) {
        throw error(404, 'Користувача не знайдено в сервісі автентифікації');
    }

    let coreUser = null;
    try {
        coreUser = await locals.coreClient.getUser(id);
    } catch (e: any) {
        if (e.status && e.status !== 404) {
            throw error(e.status, e.data?.error || 'Помилка при отриманні профілю користувача');
        }
        // If 404 in core, we might still want to show the auth user, but maybe not?
        // For now, let's just propagate the error if it's important.
        if (e.status === 404) {
             throw error(404, 'Профіль користувача не знайдено');
        }
    }
    const postsResponse = await locals.coreClient.listData<any>('posts', 1, 100);

    return {
        user: {
            ...authUser,
            profile: coreUser || {}
        },
        roles,
        posts: postsResponse.data
    };
};

export const actions: Actions = {
    default: async ({ request, params, locals, cookies }) => {
        if (!can(locals.user, 'update:users')) {
            return fail(403, { error: 'You do not have permission to update users' });
        }

        const { id } = params;
        const data = await request.formData();
        
        const roleId = data.get('role_id');
        const firstName = data.get('first_name') as string;
        const lastName = data.get('last_name') as string;
        const thirdName = data.get('third_name') as string;
        const phone = data.get('phone_number') as string;
        const email = data.get('email') as string;
        const notes = data.get('notes') as string;
        const customsPostId = data.get('customs_post_id');

        if (roleId) {
             const success = await locals.authClient.updateUserRole(id, Number(roleId));
             if (!success) {
                 return fail(500, { error: 'Failed to update user role' });
             }
        }

        const result = await locals.coreClient.updateUser(id, {
            first_name: firstName,
            last_name: lastName,
            third_name: thirdName,
            phone_number: phone,
            email,
            notes,
            customs_post_id: customsPostId ? Number(customsPostId) : null
        });

        if (!result) {
            return fail(500, { error: 'Failed to update user profile' });
        }

        throw redirect(303, '/admin/users');
    }
};
