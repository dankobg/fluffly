import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch }) => {
	depends('data:dashboard-breeds-create');

	try {
		const animalSpeciesResult = await fluffly.GET('/animal_species', {
			fetch,
			params: {
				query: { page_size: 1_000 }
			}
		});

		if (animalSpeciesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalSpeciesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
