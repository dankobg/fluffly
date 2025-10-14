import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const messageResult = await fluffly.GET('/courier/messages/{id}', {
			params: {
				path: { id: params.message_id }
			}
		});
		return {
			messageResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
