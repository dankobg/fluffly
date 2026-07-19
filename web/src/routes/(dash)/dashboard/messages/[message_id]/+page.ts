import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-messages-${params.message_id}`);

	try {
		const messageResult = await fluffly.GET('/courier/messages/{id}', {
			fetch,
			params: {
				path: { id: params.message_id }
			}
		});

		if (messageResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			messageResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
