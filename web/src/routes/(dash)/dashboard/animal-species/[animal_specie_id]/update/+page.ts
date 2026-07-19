import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-animal-species-${params.animal_specie_id}-update`);

	try {
		const animalSpeciesResult = await fluffly.GET('/animal_species/{id}', {
			fetch,
			params: {
				path: { id: Number(params.animal_specie_id) }
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
