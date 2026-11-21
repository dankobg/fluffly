import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { fluffly } from '#lib/fluffly/client.ts';
import { config } from '#lib/kratos/config.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch }) => {
	depends('data:dashboard-animal-species-create');

	try {
		const animalTypesResult = await fluffly.GET('/animal_types', {
			fetch,
			params: {
				query: { page_size: 1_000 }
			}
		});

		if (animalTypesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalTypesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
