import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends }) => {
	depends('data:identities');
	try {
		const identitiesResult = await fluffly.GET('/identities', {
			fetch,
			params: {
				query: { page_size: 500 }
			}
		});
		return {
			identitiesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
