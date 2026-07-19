import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends }) => {
	depends(`data:home`);

	try {
		const animalTypesResult = await fluffly.GET('/animal_types', {
			fetch
		});

		if (animalTypesResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		const animalSpeciesResult = await fluffly.GET('/animal_species', {
			fetch
		});

		if (animalSpeciesResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
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
