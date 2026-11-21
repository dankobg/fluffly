import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { fluffly } from '#lib/fluffly/client.ts';
import { config } from '#lib/kratos/config.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch, params }) => {
	depends(`data:dashboard-breeds-${params.breed_id}`);

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
