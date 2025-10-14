import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, depends }) => {
	depends(`data:identity-${params.identity_id}`);
	try {
		const schemasResult = await fluffly.GET('/schemas', {
			params: {
				query: { page_size: 1_000 }
			}
		});
		const identityResult = await fluffly.GET('/identities/{id}', {
			params: {
				path: { id: params.identity_id }
			}
		});
		return {
			schemasResult,
			identityResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
