import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const identityResult = await fluffly.GET('/identities/{id}', {
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
