import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async (event) => {
        const { request, cookies, locals, getClientAddress } = event;
        const data = await request.formData();
        const username = data.get('username') as string;
        const password = data.get('password') as string;

        if (!username || !password) {
            return fail(400, { message: 'Username and password are required' });
        }

        let ip = getClientAddress();
        const userAgent = request.headers.get('user-agent') || undefined;

        const result = await locals.authClient.login(username, password, ip, userAgent);

        if (!result) {
            console.error('AuthClient.login failed', { result, username, password });
            return fail(401, { message: 'Invalid username or password', status: 401 });
        }

        cookies.set('session', result.session_id, {
            path: '/',
            httpOnly: true,
            sameSite: 'strict',
            maxAge: 60 * 60 * 12,
        });

        throw redirect(303, '/');
    }
};