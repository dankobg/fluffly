import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { fluffly } from '#lib/fluffly/client.ts';
import { PathsOrganizationsIdGetParametersQueryEmbed } from '#lib/gen/fluffly_openapi.ts';
import { config } from '#lib/kratos/config.ts';
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
