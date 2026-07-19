import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { animalAgeValues, animalGenderValues, animalSizeValues } from '$lib/enum-values';
import { fluffly } from '$lib/fluffly/client';
import {
	AnimalAge,
	AnimalGender,
	AnimalSize,
	AnimalStatus,
	PathsAnimalsGetParametersQueryEmbed,
	type operations
} from '$lib/gen/fluffly_openapi';
import type { PageLoad } from './$types';

export const prerender = false;

export const load: PageLoad = async ({ fetch, depends, url }) => {
	depends(`data:search`);

	const listAnimalsParams: operations['listAnimals']['parameters'] = {
		query: {
			embed: [PathsAnimalsGetParametersQueryEmbed.breeds, PathsAnimalsGetParametersQueryEmbed.tags],
			status: [AnimalStatus.adoptable]
		}
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listAnimalsParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listAnimalsParams.query!.page_size = pageSize;
		}
	}
	const animalTypeIdParam = url.searchParams.get('animal_type_id');
	if (animalTypeIdParam) {
		const animalTypeId = Number.parseInt(animalTypeIdParam);
		if (!Number.isNaN(animalTypeId)) {
			listAnimalsParams.query!.animal_type_id = [animalTypeId];
		}
	}
	const animalSpecieIdParam = url.searchParams.get('animal_specie_id');
	if (animalSpecieIdParam) {
		const animalSpecieId = Number.parseInt(animalSpecieIdParam);
		if (!Number.isNaN(animalSpecieId)) {
			listAnimalsParams.query!.animal_specie_id = [animalSpecieId];
		}
	}
	const animalBreedsIdParams = url.searchParams.getAll('animal_breed_id');
	if (animalBreedsIdParams.length > 0) {
		const animalBreedIds = animalBreedsIdParams.map(Number).filter(Number.isFinite);
		if (animalBreedIds.length > 0) {
			listAnimalsParams.query!.animal_breed_id = animalBreedIds;
		}
	}
	const ageParams = url.searchParams.getAll('age');
	if (ageParams?.length > 0) {
		const ages = ageParams.filter(x => animalAgeValues.includes(x as AnimalAge)) as AnimalAge[];
		if (ages.length > 0) {
			listAnimalsParams.query!.age = ages;
		}
	}
	const sizeParams = url.searchParams.getAll('size');
	if (sizeParams?.length > 0) {
		const sizes = sizeParams.filter(x => animalSizeValues.includes(x as AnimalSize)) as AnimalSize[];
		if (sizes.length > 0) {
			listAnimalsParams.query!.size = sizes;
		}
	}
	const genderParams = url.searchParams.getAll('gender');
	if (genderParams?.length > 0) {
		const genders = genderParams.filter(x => animalGenderValues.includes(x as AnimalGender)) as AnimalGender[];
		if (genders.length > 0) {
			listAnimalsParams.query!.gender = genders;
		}
	}
	const hermaphroditeParam = url.searchParams.get('hermaphrodite');
	if (hermaphroditeParam === 'true' || hermaphroditeParam === 'false') {
		listAnimalsParams.query!.hermaphrodite = hermaphroditeParam === 'true';
	}
	const microchipParam = url.searchParams.get('microchip');
	if (microchipParam === 'true' || microchipParam === 'false') {
		listAnimalsParams.query!.microchip = microchipParam === 'true';
	}
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listAnimalsParams.query!.name = nameParam;
	}
	const tagParams = url.searchParams.getAll('tag');
	if (tagParams.length > 0) {
		listAnimalsParams.query!.tag = tagParams;
	}
	const daysLtParam = url.searchParams.get('days_lt');
	if (daysLtParam) {
		const daysLt = Number.parseInt(daysLtParam);
		if (!Number.isNaN(daysLt)) {
			listAnimalsParams.query!.days_lt = daysLt;
		}
	}
	const daysGtParam = url.searchParams.get('days_gt');
	if (daysGtParam) {
		const daysGt = Number.parseInt(daysGtParam);
		if (!Number.isNaN(daysGt)) {
			listAnimalsParams.query!.days_gt = daysGt;
		}
	}
	const locationParam = url.searchParams.get('location');
	if (locationParam) {
		listAnimalsParams.query!.location = locationParam;
	}
	const latParam = url.searchParams.get('lat');
	if (latParam) {
		const lat = Number.parseFloat(latParam);
		if (!Number.isNaN(lat)) {
			listAnimalsParams.query!.lat = lat;
		}
	}
	const lonParam = url.searchParams.get('lon');
	if (lonParam) {
		const lon = Number.parseFloat(lonParam);
		if (!Number.isNaN(lon)) {
			listAnimalsParams.query!.lon = lon;
		}
	}
	const radiusMparam = url.searchParams.get('radius_m');
	if (radiusMparam) {
		const radiusM = Number.parseFloat(radiusMparam);
		if (!Number.isNaN(radiusM)) {
			listAnimalsParams.query!.radius_m = radiusM;
		}
	}
	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listAnimalsParams.query!.sort = sortParams;
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
		listAnimalsParams.query!.properties = properties;
	}

	try {
		const animalsResult = await fluffly.GET('/animals', {
			fetch,
			params: listAnimalsParams,
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
				goto('/');
			}
		}

		const animalTypesResult = await fluffly.GET('/animal_types', {
			fetch
		});
		if (animalTypesResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		const animalSpeciesResult = await fluffly.GET('/animal_species', {
			fetch,
			params: { query: { sort: ['name'] } }
		});
		if (animalSpeciesResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
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
				goto('/');
			}
		}

		return {
			animalsResult,
			animalTypesResult,
			animalSpeciesResult,
			animalBreedsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
