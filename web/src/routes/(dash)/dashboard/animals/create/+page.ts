import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch }) => {
	depends('data:dashboard-animals-create');

	// SELECT o.*
	// FROM organization o
	// JOIN organization_membership om
	//   ON om.organization_id = o.id
	// WHERE om.user_id = 'd3d4fc3d-ab70-46c3-a823-2bc5611fd193'
	//   AND om.status = 'approved';

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

		const animalSpeciesResult = await fluffly.GET('/animal_species', {
			fetch,
			params: {
				query: { page_size: 1_000, sort: ['name'] }
			}
		});

		if (animalSpeciesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const animalBreedsResult = await fluffly.GET('/animal_breeds', {
			fetch,
			params: {
				query: { page_size: 1_000, sort: ['name'] }
			}
		});

		if (animalBreedsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalTypesResult,
			animalSpeciesResult,
			animalBreedsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
