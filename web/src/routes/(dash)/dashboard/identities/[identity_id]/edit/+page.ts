import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, depends }) => {
	depends(`data:identity-${params.identity_id}`);
	try {
		const schemasRes = await fluffly.GET('/schemas', {
			params: {
				query: { page_size: 1_000 }
			}
		});
		const identityRes = await fluffly.GET('/identities/{id}', {
			params: {
				path: { id: params.identity_id }
			}
		});
		return {
			schemas: schemasRes.data,
			identity: identityRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
