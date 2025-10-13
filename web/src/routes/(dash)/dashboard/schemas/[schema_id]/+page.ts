import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const schemaRes = await fluffly.GET('/schemas/{id}', {
			params: {
				path: { id: params.schema_id }
			}
		});
		return {
			id: params.schema_id,
			schema: schemaRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
