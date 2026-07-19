import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { OrganizationStatus } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch }) => {
	depends('data:dashboard-animals-submit');

	try {
		const myOrganizationsResult = await fluffly.GET('/me/organizations', {
			fetch,
			params: {
				query: { status: [OrganizationStatus.approved] }
			}
		});

		if (myOrganizationsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

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
			myOrganizationsResult,
			animalTypesResult,
			animalSpeciesResult,
			animalBreedsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
