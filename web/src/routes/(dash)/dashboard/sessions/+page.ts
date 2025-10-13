import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';
import type { operations, PathsSessionsGetParametersQueryExpand } from '$lib/gen/fluffly_openapi';

export const load: PageLoad = async ({ url, depends }) => {
	depends('data:sessions');
	try {
		const listSessionsParams: operations['listSessions']['parameters'] = {
			query: { page_size: 1_000 }
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
		const sessionsRes = await fluffly.GET('/sessions', {
			params: listSessionsParams
		});
		return {
			sessions: sessionsRes.data
		};
	} catch (error) {
		console.log('err', error);
	}
};
