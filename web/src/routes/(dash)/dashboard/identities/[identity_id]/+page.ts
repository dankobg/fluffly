import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	try {
		const identityResult = await fluffly.GET('/identities/{id}', {
			fetch,
			params: {
				path: { id: params.identity_id }
			}
		});
		return {
			identityResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
