import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const schemasResult = await fluffly.GET('/schemas', {
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
