import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch, params }) => {
	depends(`data:dashboard-animal-types-${params.animal_type_id}`);

	try {
		const animalTypeResult = await fluffly.GET('/animal_types/{id}', {
			fetch,
			params: {
				path: { id: Number(params.animal_type_id) }
			}
		});

		if (animalTypeResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalTypeResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
