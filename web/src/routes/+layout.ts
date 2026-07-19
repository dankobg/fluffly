import { kratos } from '$lib/kratos/client';
import { createSessionService } from '$lib/kratos/service';
import type { LayoutLoad } from './$types';

export const ssr = true;
export const prerender = true;
export const trailingSlash = 'never';

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
	} catch (_: unknown) {
		return {
			auth: createSessionService(null)
		};
	}
};
