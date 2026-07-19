import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { animalAgeValues, animalGenderValues, animalSizeValues, animalStatusValues } from '$lib/enum-values';
import { fluffly } from '$lib/fluffly/client';
import type { AnimalAge, AnimalGender, AnimalSize, AnimalStatus, operations } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-animals');

	const listAnimalsParams: operations['listAnimals']['parameters'] = {
		query: { sort: ['-created_at'] }
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
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listAnimalsParams.query!.name = nameParam;
	}
	const animalTypeIdParams = url.searchParams.getAll('animal_type_id');
	if (animalTypeIdParams.length > 0) {
		const animalTypeIds = animalTypeIdParams.map(Number).filter(Number.isFinite);
		if (animalTypeIds.length > 0) {
			listAnimalsParams.query!.animal_type_id = animalTypeIds;
		}
	}
	const animalSpecieIdParams = url.searchParams.getAll('animal_specie_id');
	if (animalSpecieIdParams.length > 0) {
		const animalSpecieIds = animalSpecieIdParams.map(Number).filter(Number.isFinite);
		if (animalSpecieIds.length > 0) {
			listAnimalsParams.query!.animal_specie_id = animalSpecieIds;
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
	const statusParams = url.searchParams.getAll('status');
	if (statusParams?.length > 0) {
		const statuses = statusParams.filter(x => animalStatusValues.includes(x as AnimalStatus)) as AnimalStatus[];
		if (statuses.length > 0) {
			listAnimalsParams.query!.status = statuses;
		}
	}
	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listAnimalsParams.query!.sort = sortParams;
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
			fetch
		});

		if (animalSpeciesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const animalsResult = await fluffly.GET('/animals', {
			fetch,
			params: listAnimalsParams
		});

		if (animalsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			animalTypesResult,
			animalSpeciesResult,
			animalsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
