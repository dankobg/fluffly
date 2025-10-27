import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	try {
		const countriesResult = await fluffly.GET('/countries', {
			fetch,
			params: {
				query: { page_size: 500 }
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
