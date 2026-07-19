import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import type { operations } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-breeds');

	const listBreedsParams: operations['listBreeds']['parameters'] = {
		query: {}
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listBreedsParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listBreedsParams.query!.page_size = pageSize;
		}
	}
	const animalSpecieIdParams = url.searchParams.getAll('animal_specie_id');
	if (animalSpecieIdParams.length > 0) {
		const animalSpecieIds = animalSpecieIdParams.map(Number).filter(Number.isFinite);
		if (animalSpecieIds.length > 0) {
			listBreedsParams.query!.animal_specie_id = animalSpecieIds;
		}
	}
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listBreedsParams.query!.name = nameParam;
	}
	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listBreedsParams.query!.sort = sortParams;
	}

	try {
		const animalSpeciesResults = await fluffly.GET('/animal_species', {
			fetch,
			params: { query: { page_size: 1_000 } }
		});

		if (animalSpeciesResults.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const breedsResult = await fluffly.GET('/animal_breeds', {
			fetch,
			params: listBreedsParams
		});

		if (breedsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalSpeciesResults,
			breedsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
