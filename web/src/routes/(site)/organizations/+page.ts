import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { OrganizationStatus, type operations } from '$lib/gen/fluffly_openapi';
import type { PageLoad } from './$types';

export const prerender = false;

export const load: PageLoad = async ({ fetch, depends, url }) => {
	depends(`data:organizations`);

	const listOrganizationsParams: operations['listOrganizations']['parameters'] = {
		query: {
			status: [OrganizationStatus.approved]
		}
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
	const idParams = url.searchParams.getAll('id');
	if (idParams.length > 0) {
		const organizationIds = idParams.map(Number).filter(Number.isFinite);
		if (organizationIds.length > 0) {
			listOrganizationsParams.query!.id = organizationIds;
		}
	}

	try {
		const organizationsResult = await fluffly.GET('/organizations', {
			fetch,
			params: listOrganizationsParams
		});
		if (organizationsResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		return {
			organizationsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
