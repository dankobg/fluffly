import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboad-schemas-${params.schema_id}`);

	try {
		const schemaResult = await fluffly.GET('/schemas/{id}', {
			fetch,
			params: {
				path: { id: params.schema_id }
			}
		});

		if (schemaResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			id: params.schema_id,
			schemaResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
