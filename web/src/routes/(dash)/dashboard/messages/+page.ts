import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const messagesResult = await fluffly.GET('/courier/messages', {
			params: {
				query: { page_size: 1_000 }
			}
		});
		return {
			messagesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
