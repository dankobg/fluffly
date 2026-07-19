import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { animalAgeValues, animalGenderValues, animalSizeValues } from '$lib/enum-values';
import { fluffly } from '$lib/fluffly/client';
import {
	AnimalAge,
	AnimalGender,
	AnimalSize,
	PathsMeAnimalsGetParametersQueryEmbed,
	type operations
} from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-posted-animals');

	const listyMyAnimalsParams: operations['listMyAnimals']['parameters'] = {
		query: {
			embed: [PathsMeAnimalsGetParametersQueryEmbed.breeds, PathsMeAnimalsGetParametersQueryEmbed.tags]
		}
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listyMyAnimalsParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listyMyAnimalsParams.query!.page_size = pageSize;
		}
	}
	const animalTypeIdParam = url.searchParams.get('animal_type_id');
	if (animalTypeIdParam) {
		const animalTypeId = Number.parseInt(animalTypeIdParam);
		if (!Number.isNaN(animalTypeId)) {
			listyMyAnimalsParams.query!.animal_type_id = [animalTypeId];
		}
	}
	const animalSpecieIdParam = url.searchParams.get('animal_specie_id');
	if (animalSpecieIdParam) {
		const animalSpecieId = Number.parseInt(animalSpecieIdParam);
		if (!Number.isNaN(animalSpecieId)) {
			listyMyAnimalsParams.query!.animal_specie_id = [animalSpecieId];
		}
	}
	const animalBreedsIdParams = url.searchParams.getAll('animal_breed_id');
	if (animalBreedsIdParams.length > 0) {
		const animalBreedIds = animalBreedsIdParams.map(Number).filter(Number.isFinite);
		if (animalBreedIds.length > 0) {
			listyMyAnimalsParams.query!.animal_breed_id = animalBreedIds;
		}
	}
	const ageParams = url.searchParams.getAll('age');
	if (ageParams?.length > 0) {
		const ages = ageParams.filter(x => animalAgeValues.includes(x as AnimalAge)) as AnimalAge[];
		if (ages.length > 0) {
			listyMyAnimalsParams.query!.age = ages;
		}
	}
	const sizeParams = url.searchParams.getAll('size');
	if (sizeParams?.length > 0) {
		const sizes = sizeParams.filter(x => animalSizeValues.includes(x as AnimalSize)) as AnimalSize[];
		if (sizes.length > 0) {
			listyMyAnimalsParams.query!.size = sizes;
		}
	}
	const genderParams = url.searchParams.getAll('gender');
	if (genderParams?.length > 0) {
		const genders = genderParams.filter(x => animalGenderValues.includes(x as AnimalGender)) as AnimalGender[];
		if (genders.length > 0) {
			listyMyAnimalsParams.query!.gender = genders;
		}
	}
	const hermaphroditeParam = url.searchParams.get('hermaphrodite');
	if (hermaphroditeParam === 'true' || hermaphroditeParam === 'false') {
		listyMyAnimalsParams.query!.hermaphrodite = hermaphroditeParam === 'true';
	}
	const microchipParam = url.searchParams.get('microchip');
	if (microchipParam === 'true' || microchipParam === 'false') {
		listyMyAnimalsParams.query!.microchip = microchipParam === 'true';
	}
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listyMyAnimalsParams.query!.name = nameParam;
	}
	const tagParams = url.searchParams.getAll('tag');
	if (tagParams.length > 0) {
		listyMyAnimalsParams.query!.tag = tagParams;
	}
	const daysLtParam = url.searchParams.get('days_lt');
	if (daysLtParam) {
		const daysLt = Number.parseInt(daysLtParam);
		if (!Number.isNaN(daysLt)) {
			listyMyAnimalsParams.query!.days_lt = daysLt;
		}
	}
	const daysGtParam = url.searchParams.get('days_gt');
	if (daysGtParam) {
		const daysGt = Number.parseInt(daysGtParam);
		if (!Number.isNaN(daysGt)) {
			listyMyAnimalsParams.query!.days_gt = daysGt;
		}
	}
	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listyMyAnimalsParams.query!.sort = sortParams;
	}

	const properties: Record<string, string[]> = {};
	url.searchParams.entries().forEach(([key, value]) => {
		const match = key.match(/^properties\[(\w+)\]\[(\d+)\]$/);
		if (match) {
			const [, name, indexStr] = match;
			const index = Number(indexStr);
			if (name) {
				properties[name] ||= [];
				properties[name][index] = value;
			}
		}
	});
	if (Object.keys(properties).length > 0) {
		listyMyAnimalsParams.query!.properties = properties;
	}

	try {
		const animalsResult = await fluffly.GET('/me/animals', {
			fetch,
			params: listyMyAnimalsParams,
			querySerializer(queryParams) {
				const sp = new URLSearchParams();

				Object.entries(queryParams).forEach(([key, val]) => {
					if (Array.isArray(val)) {
						for (const x of val) {
							sp.append(key, `${x}`);
						}
					} else {
						if (typeof val === 'object') {
							if (key === 'properties') {
								Object.entries(val).forEach(([propKey, propVal]) => {
									propVal.forEach((x, i) => {
										sp.append(`properties[${propKey}][${i}]`, `${x}`);
									});
								});
							}
						} else {
							sp.set(key, `${val}`);
						}
					}
				});

				return sp.size > 0 ? `?${sp}` : '';
			}
		});

		if (animalsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
