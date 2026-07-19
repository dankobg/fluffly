import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-breeds-${params.breed_id}-update`);

	try {
		const breedResult = await fluffly.GET('/animal_breeds/{id}', {
			fetch,
			params: {
				path: { id: Number(params.breed_id) }
			}
		});

		if (breedResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			breedResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
