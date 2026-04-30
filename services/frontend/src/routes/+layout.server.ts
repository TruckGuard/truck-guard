import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, cookies }) => {
	let user = locals.user;
	const token = cookies.get('session');	
	let notifications: any[] = [];

	if (user && token) {
		try {
			const coreUser = await locals.coreClient.getUser(user.id);
			if (coreUser) {
				user = {
					...user,
					...coreUser
				};
			} 
			
			const notifsRes = await locals.coreClient.getNotifications(50);
			notifications = notifsRes.data || [];
		} catch (e) {
			console.error('Failed to fetch core user or notifications', e);
		}
	}

	return {
		user,
		notifications
	};
};