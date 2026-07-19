import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { organizationStatusValues } from '$lib/enum-values';
import { fluffly } from '$lib/fluffly/client';
import type { operations, OrganizationStatus } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-organizations');

	const listOrganizationsParams: operations['listOrganizations']['parameters'] = {
		query: { sort: ['-created_at'] }
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listOrganizationsParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listOrganizationsParams.query!.page_size = pageSize;
		}
	}
	const nameParam = url.searchParams.get('name');
	if (nameParam) {
		listOrganizationsParams.query!.name = nameParam;
	}
	const statusParams = url.searchParams.getAll('status');
	if (statusParams?.length > 0) {
		const statuses = statusParams.filter(x =>
			organizationStatusValues.includes(x as OrganizationStatus)
		) as OrganizationStatus[];
		if (statuses.length > 0) {
			listOrganizationsParams.query!.status = statuses;
		}
	}
	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listOrganizationsParams.query!.sort = sortParams;
	}

	try {
		const organizationsResult = await fluffly.GET('/organizations', {
			fetch,
			params: listOrganizationsParams
		});

		if (organizationsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			organizationsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
