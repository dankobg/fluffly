import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	try {
		const schemaResult = await fluffly.GET('/schemas/{id}', {
			fetch,
			params: {
				path: { id: params.schema_id }
			}
		});
		return {
			id: params.schema_id,
			schemaResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
