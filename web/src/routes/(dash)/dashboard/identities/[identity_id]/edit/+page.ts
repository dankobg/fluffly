import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:identity-${params.identity_id}`);
	try {
		const schemasResult = await fluffly.GET('/schemas', {
			fetch,
			params: {
				query: { page_size: 500 }
			}
		});
		const identityResult = await fluffly.GET('/identities/{id}', {
			fetch,
			params: {
				path: { id: params.identity_id }
			}
		});

		if (schemasResult.error?.status_code === 403 || identityResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			schemasResult,
			identityResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
