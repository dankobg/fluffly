import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends }) => {
	depends('data:identities');
	try {
		const identitiesRes = await fluffly.GET('/identities', {
			params: {
				query: { page_size: 1_000 }
			}
		});
		return {
			identities: identitiesRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
