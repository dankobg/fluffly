import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	try {
		const schemasResult = await fluffly.GET('/schemas', {
			fetch,
			params: {
				query: { page_size: 500 }
			}
		});
		return {
			schemasResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
