import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';
import type { operations } from '$lib/gen/fluffly_openapi';

export const load: PageLoad = async ({ fetch, url, params, depends }) => {
	depends(`data:identity-sessions-${params.identity_id}`);
	try {
		const listIdentitySessionsParams: operations['listIdentitySessions']['parameters'] = {
			path: { id: params.identity_id },
			query: { page_size: 500 }
		};
		const active = url.searchParams.get('active');
		if (active) {
			listIdentitySessionsParams.query!.active = active === 'true' ? true : active === 'false' ? false : false;
		}
		const pageToken = url.searchParams.get('page_token');
		if (pageToken) {
			listIdentitySessionsParams.query!.page_token = pageToken;
		}
		const sessionsResult = await fluffly.GET('/identities/{id}/sessions', {
			fetch,
			params: listIdentitySessionsParams
		});
		const identityResult = await fluffly.GET('/identities/{id}', {
			fetch,
			params: {
				path: { id: params.identity_id }
			}
		});
		return {
			identityResult,
			sessionsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
