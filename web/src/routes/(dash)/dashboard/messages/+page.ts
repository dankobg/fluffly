import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { fluffly } from '#lib/fluffly/client.ts';
import { config } from '#lib/kratos/config.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends }) => {
	depends('data:dashboard-messages');

	try {
		const messagesResult = await fluffly.GET('/courier/messages', {
			fetch,
			params: {
				query: { page_size: 500 }
			}
		});

		if (messagesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			messagesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
