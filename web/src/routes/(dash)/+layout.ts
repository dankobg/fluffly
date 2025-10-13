import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { config } from '$lib/kratos/config';
import type { LayoutLoad } from './$types';

export const prerender = false;

export const load: LayoutLoad = async ({ parent, url }) => {
	const data = await parent();
	const { auth } = data;

	if (!auth.session?.active) {
		if (browser) {
			goto(config.routes.login.path + `?return_to=${encodeURIComponent(url.toString())}`);
		}
		return;
	}
	return data;
};
