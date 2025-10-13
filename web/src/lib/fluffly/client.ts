const apiEndpoint = import.meta.env['VITE_PUBLIC_API_ENDPOINT'] as string;
import createClient from 'openapi-fetch';
import type { paths } from '$lib/gen/fluffly_openapi';

export const fluffly = createClient<paths>({
	baseUrl: apiEndpoint,
	credentials: 'include',
	headers: {
		'Content-Type': 'application/json',
		Accept: 'application/json'
	}
});
