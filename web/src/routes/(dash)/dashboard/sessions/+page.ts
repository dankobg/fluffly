import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';
import type { operations, PathsSessionsGetParametersQueryExpand } from '$lib/gen/fluffly_openapi';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:sessions');
	try {
		const listSessionsParams: operations['listSessions']['parameters'] = {
			query: { page_size: 500 }
		};
		const active = url.searchParams.get('active');
		if (active) {
			listSessionsParams.query!.active = active === 'true' ? true : active === 'false' ? false : false;
		}
		const expand = url.searchParams.getAll('expand') as PathsSessionsGetParametersQueryExpand[];
		if (expand.length > 0) {
			listSessionsParams.query!.expand = expand;
		}
		const pageToken = url.searchParams.get('page_token');
		if (pageToken) {
			listSessionsParams.query!.page_token = pageToken;
		}
		const sessionsResult = await fluffly.GET('/sessions', {
			fetch,
			params: listSessionsParams
		});
		return {
			sessionsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
