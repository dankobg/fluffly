import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-identities-${params.identity_id}`);

	try {
		const identityResult = await fluffly.GET('/identities/{id}', {
			fetch,
			params: {
				path: { id: params.identity_id }
			}
		});

		if (identityResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			identityResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
