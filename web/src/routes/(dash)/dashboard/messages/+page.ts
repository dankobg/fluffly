import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const messagesRes = await fluffly.GET('/courier/messages', {
			params: {
				query: { page_size: 1_000 }
			}
		});
		return {
			messages: messagesRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
