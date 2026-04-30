import { fail } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
	markAllRead: async ({ locals }) => {
		if (!locals.coreClient) return fail(401, { error: "Unauthorized" });

		try {
			await locals.coreClient.markAllNotificationsRead();
			return { success: true };
		} catch (e: any) {
			console.error("markAllRead error:", e);
			return fail(500, { error: e.message || "Internal Error" });
		}
	},
	markRead: async ({ request, locals }) => {
		if (!locals.coreClient) return fail(401, { error: "Unauthorized" });

		const data = await request.formData();
		const id = data.get("id")?.toString();

		if (!id) return fail(400, { error: "Missing notification ID" });

		try {
			await locals.core.markNotificationRead(id);
			return { success: true };
		} catch (e: any) {
			console.error("markRead error:", e);
			return fail(500, { error: e.message || "Internal Error" });
		}
	}
};
