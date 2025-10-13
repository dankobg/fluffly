import { kratos } from '$lib/kratos/client';
import { createSessionService } from '$lib/kratos/service';
import type { LayoutLoad } from './$types';

export const prerender = true;

export const load: LayoutLoad = async () => {
	try {
		const session = await kratos.toSession();
		const svc = createSessionService(session);
		const logoutFlow = await kratos.createBrowserLogoutFlow();
		return {
			auth: svc,
			logoutToken: logoutFlow.logout_token,
			logoutUrl: logoutFlow.logout_url
		};
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
	} catch (_error: unknown) {
		return {
			auth: createSessionService(null)
		};
	}
};
