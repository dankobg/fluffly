import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const identityRes = await fluffly.GET('/identities/{id}', {
			params: {
				path: { id: params.identity_id }
			}
		});
		return {
			identity: identityRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
