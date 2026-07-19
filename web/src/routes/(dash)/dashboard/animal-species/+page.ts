import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import type { operations } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-animal-species');

	const listAnimalSpeciesParams: operations['listAnimalSpecies']['parameters'] = {
		query: {}
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listAnimalSpeciesParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listAnimalSpeciesParams.query!.page_size = pageSize;
		}
	}
	const animalTypeIdParams = url.searchParams.getAll('animal_type_id');
	if (animalTypeIdParams.length > 0) {
		const animalTypeIds = animalTypeIdParams.map(Number).filter(Number.isFinite);
		if (animalTypeIds.length > 0) {
			listAnimalSpeciesParams.query!.animal_type_id = animalTypeIds;
		}
	}
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listAnimalSpeciesParams.query!.name = nameParam;
	}
	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listAnimalSpeciesParams.query!.sort = sortParams;
	}

	try {
		const animalTypesResult = await fluffly.GET('/animal_types', {
			fetch
		});

		if (animalTypesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const animalSpeciesResult = await fluffly.GET('/animal_species', {
			fetch,
			params: listAnimalSpeciesParams
		});

		if (animalSpeciesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalTypesResult,
			animalSpeciesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
