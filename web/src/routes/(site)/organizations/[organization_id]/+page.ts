import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { PathsOrganizationsIdGetParametersQueryEmbed } from '$lib/gen/fluffly_openapi';
import type { PageLoad } from './$types';

export const prerender = false;

export const load: PageLoad = async ({ fetch, depends, params }) => {
  depends(`data:organization:${params.organization_id}`);

	try {
		const organizationResult = await fluffly.GET('/organizations/{id}', {
			fetch,
			params: {
				path: {
					id: Number(params.organization_id)
				},
				query: {
					embed: [
						PathsOrganizationsIdGetParametersQueryEmbed.work_hour,
						PathsOrganizationsIdGetParametersQueryEmbed.contact,
						PathsOrganizationsIdGetParametersQueryEmbed.socials,
						PathsOrganizationsIdGetParametersQueryEmbed.photos,
						PathsOrganizationsIdGetParametersQueryEmbed.videos
					]
				}
			}
		});

		if (organizationResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		return {
			organizationResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
