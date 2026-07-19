import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { adoptionStatusValues } from '$lib/enum-values';
import { fluffly } from '$lib/fluffly/client';
import type { AdoptionStatus, operations } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-adoptions');

	const listAdoptionsParams: operations['listAdoptions']['parameters'] = {
		query: {}
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listAdoptionsParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listAdoptionsParams.query!.page_size = pageSize;
		}
	}
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listAdoptionsParams.query!.name = nameParam;
	}
	const statusParams = url.searchParams.getAll('status');
	if (statusParams?.length > 0) {
		const statuses = statusParams.filter(x => adoptionStatusValues.includes(x as AdoptionStatus)) as AdoptionStatus[];
		if (statuses.length > 0) {
			listAdoptionsParams.query!.status = statuses;
		}
	}
	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listAdoptionsParams.query!.sort = sortParams;
	}

	try {
		const adoptionsResult = await fluffly.GET('/adoptions', {
			fetch,
			params: listAdoptionsParams
		});

		if (adoptionsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			adoptionsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
