import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { fluffly } from '#lib/fluffly/client.ts';
import { config } from '#lib/kratos/config.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends }) => {
	depends(`data:dashboard-my-organizations`);

	try {
		const myOrganizationsResult = await fluffly.GET('/me/organizations', {
			fetch
		});

		if (myOrganizationsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			myOrganizationsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
