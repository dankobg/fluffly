import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	try {
		const messagesResult = await fluffly.GET('/courier/messages', {
			fetch,
			params: {
				query: { page_size: 500 }
			}
		});
		return {
			messagesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
