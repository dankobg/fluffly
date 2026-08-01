import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { PathsAnimalsIdGetParametersQueryEmbed } from '$lib/gen/fluffly_openapi';
import type { PageLoad } from './$types';

export const prerender = false;

export const load: PageLoad = async ({ fetch, depends, params }) => {
	depends(`data:animal:${params.animal_id}`);

	try {
		const animalResult = await fluffly.GET('/animals/{id}', {
			fetch,
			params: {
				path: {
					id: Number(params.animal_id)
				},
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
				goto('/');
			}
		}

		const animalSpecieResult = await fluffly.GET('/animal_species/{id}', {
			fetch,
			params: {
				path: {
					id: animalResult?.data?.specie?.id ?? 0
				}
			}
		});
		if (animalSpecieResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		return {
			animalResult,
			animalSpecieResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
