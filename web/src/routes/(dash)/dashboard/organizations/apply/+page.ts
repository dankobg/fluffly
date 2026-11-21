import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { fluffly } from '#lib/fluffly/client.ts';
import { config } from '#lib/kratos/config.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch }) => {
	depends('data:dashboard-organizations-create');

	try {
		const countriesResult = await fluffly.GET('/countries', {
			fetch,
			params: {
				query: { page_size: 1_000 }
			}
		});

		if (countriesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			countriesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
