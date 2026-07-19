import type { PageLoad } from './$types';
import { type FlowError } from '@ory/client-fetch';
import { kratos } from '$lib/kratos/client';
import { browser } from '$app/environment';

export const load: PageLoad = (async ({ url }) => {
	const errorIdParam = browser && url.searchParams.get('id');

	let flow: FlowError | null = null;

	if (errorIdParam) {
		try {
			const flowError = await kratos.getFlowError({
				id: errorIdParam
			});
			flow = flowError;
		} catch (error: unknown) {
			console.error('error flow', error);
		}
	}

	return {
		flow
	};
}) satisfies PageLoad;
