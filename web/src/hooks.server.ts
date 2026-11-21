import type { Handle } from '@sveltejs/kit/hooks';

export const handle: Handle = async ({ event, resolve }) => {
	const response = await resolve(event, {
		filterSerializedResponseHeaders() {
			return true;
		}
	});
	return response;
};
