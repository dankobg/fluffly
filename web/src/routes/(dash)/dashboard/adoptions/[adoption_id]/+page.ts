import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { PathsAdoptionsIdGetParametersQueryEmbed } from '$lib/gen/fluffly_openapi';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends, fetch, params }) => {
	depends(`data:dashboard-adoptions-${params.adoption_id}`);

	try {
		const adoptionResult = await fluffly.GET('/adoptions/{id}', {
			fetch,
			params: {
				path: { id: Number(params.adoption_id) },
				query: { embed: [PathsAdoptionsIdGetParametersQueryEmbed.organization] }
			}
		});

		if (adoptionResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			adoptionResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
