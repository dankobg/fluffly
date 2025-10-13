import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const messageRes = await fluffly.GET('/courier/messages/{id}', {
			params: {
				path: { id: params.message_id }
			}
		});
		return {
			message: messageRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
