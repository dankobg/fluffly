import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-countries-${params.country_id}-update`);

	try {
		const countryResult = await fluffly.GET('/countries/{id}', {
			fetch,
			params: {
				path: { id: Number(params.country_id) }
			}
		});

		if (countryResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			countryResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
