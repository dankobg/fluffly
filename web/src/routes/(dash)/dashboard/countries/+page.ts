import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import type { operations } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-countries');

	const listCountriesParams: operations['listCountries']['parameters'] = {
		query: {}
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listCountriesParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listCountriesParams.query!.page_size = pageSize;
		}
	}
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listCountriesParams.query!.name = nameParam;
	}
	const isoAlpha2Param = url.searchParams.get('iso_alpha2');
	if (isoAlpha2Param) {
		listCountriesParams.query!.iso_alpha2 = isoAlpha2Param;
	}
	const isoAlpha3Param = url.searchParams.get('iso_alpha3');
	if (isoAlpha3Param) {
		listCountriesParams.query!.iso_alpha3 = isoAlpha3Param;
	}

	try {
		const countriesResult = await fluffly.GET('/countries', {
			fetch,
			params: listCountriesParams
		});

		if (countriesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			countriesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
