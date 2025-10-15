import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';
import { PathsSessionsIdGetParametersQueryExpand, type operations } from '$lib/gen/fluffly_openapi';

export const load: PageLoad = async ({ fetch, params, url }) => {
	try {
		const getSessionParams: operations['getSession']['parameters'] = {
			path: { id: params.session_id },
			query: {
				expand: [PathsSessionsIdGetParametersQueryExpand.identity, PathsSessionsIdGetParametersQueryExpand.devices]
			}
		};
		const expand = url.searchParams.getAll('expand') as PathsSessionsIdGetParametersQueryExpand[];
		if (expand.length > 0) {
			getSessionParams.query!.expand = expand;
		}
		const sessionResult = await fluffly.GET('/sessions/{id}', {
			fetch,
			params: {
				path: { id: params.session_id }
			}
		});
		return {
			sessionResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
