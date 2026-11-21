const apiEndpoint = import.meta.env['VITE_PUBLIC_API_ENDPOINT'] as string;
import createClient from 'openapi-fetch';
import type { paths } from '#lib/gen/fluffly_openapi.ts';

export const fluffly = createClient<paths>({
	baseUrl: apiEndpoint,
	credentials: 'include',
	headers: {
		Accept: 'application/json'
	}
});
