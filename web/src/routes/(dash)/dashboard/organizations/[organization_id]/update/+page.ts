import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { PathsOrganizationsIdGetParametersQueryEmbed } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-organizations-${params.organization_id}-update`);

	try {
		const organizationResult = await fluffly.GET('/organizations/{id}', {
			fetch,
			params: {
				path: { id: Number(params.organization_id) },
				query: {
					embed: [
						PathsOrganizationsIdGetParametersQueryEmbed.contact,
						PathsOrganizationsIdGetParametersQueryEmbed.work_hour,
						PathsOrganizationsIdGetParametersQueryEmbed.photos,
						PathsOrganizationsIdGetParametersQueryEmbed.videos,
						PathsOrganizationsIdGetParametersQueryEmbed.socials
					]
				}
			}
		});

		const countriesResult = await fluffly.GET('/countries', {
			fetch,
			params: {
				query: { page_size: 1_000 }
			}
		});

		if (organizationResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		if (countriesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			organizationResult,
			countriesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
