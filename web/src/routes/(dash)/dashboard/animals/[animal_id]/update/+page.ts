import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { PathsAnimalsIdGetParametersQueryEmbed } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-animals-${params.animal_id}-update`);

	try {
		const animalResult = await fluffly.GET('/animals/{id}', {
			fetch,
			params: {
				path: { id: Number(params.animal_id) },
				query: {
					embed: [
						PathsAnimalsIdGetParametersQueryEmbed.breeds,
						PathsAnimalsIdGetParametersQueryEmbed.microchip,
						PathsAnimalsIdGetParametersQueryEmbed.organization,
						PathsAnimalsIdGetParametersQueryEmbed.tags,
						PathsAnimalsIdGetParametersQueryEmbed.photos,
						PathsAnimalsIdGetParametersQueryEmbed.videos
					]
				}
			}
		});

		if (animalResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const organizationsResult = await fluffly.GET('/organizations', {
			fetch,
			params: {
				query: { page_size: 1_000 }
			}
		});

		if (organizationsResult.error?.status_code === 403) {
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
			animalResult,
			organizationsResult,
			animalTypesResult,
			animalSpeciesResult,
			animalBreedsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
