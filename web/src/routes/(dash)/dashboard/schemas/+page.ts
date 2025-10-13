import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const schemasRes = await fluffly.GET('/schemas', {
			params: {
				query: { page_size: 1_000 }
			}
		});
		return {
			schemas: schemasRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
