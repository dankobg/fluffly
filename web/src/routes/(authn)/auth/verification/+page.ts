import type { PageLoad } from './$types';
import { type VerificationFlow } from '@ory/client-fetch';
import { kratos } from '#lib/kratos/client.ts';
import { extractCSRFToken } from '#lib/kratos/helpers.ts';
import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import { config } from '#lib/kratos/config.ts';
import { ResponseError } from '@ory/client-fetch';
import { isGenericErrorResponse } from '@ory/client-fetch';
import { isErrorIdSessionAlreadyAvailable } from '#lib/kratos/helpers.ts';
import { isErrorIdSelfServiceFlowExpired } from '#lib/kratos/helpers.ts';
import { FetchError } from '@ory/client-fetch';
import { RequiredError } from '@ory/client-fetch';
import { isErrorIdSecurityCsrfViolation } from '#lib/kratos/helpers.ts';
import { isErrorIdSecurityIdentityMismatch } from '#lib/kratos/helpers.ts';

export const load: PageLoad = (async ({ url, depends }) => {
	depends('data:verification');

	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: VerificationFlow | null = null;

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		flow = null;
		if (browser) {
			goto(redirectUrl);
		}
		return;
	}

	if (flowIdParam) {
		try {
			const verificationFlow = await kratos.getVerificationFlow({
				id: flowIdParam
			});
			flow = verificationFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (error instanceof ResponseError) {
				const err = await error.response.json();
				switch (error.response.status) {
					case 403:
					case 404:
					case 410: {
						if (isGenericErrorResponse(err)) {
							if (isErrorIdSessionAlreadyAvailable(err.error?.id)) {
								goto('/');
							} else if (isErrorIdSelfServiceFlowExpired(err.error?.id)) {
								if (browser) {
									goto(`${config.routes.verification.path}?return_to=${window.location.href}`);
								}
							}
						}
						break;
					}
					case 500:
						console.error('unexpected server error');
						break;
					default:
						break;
				}
				return;
			}
			if (error instanceof FetchError) {
				console.error('fetch error: ', error.cause);
				return;
			}
			if (error instanceof RequiredError) {
				console.error('required error: ', error.field);
				return;
			}
			console.error('unexpected error');
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const verificationFlow = await kratos.createBrowserVerificationFlow({
				returnTo
			});
			flow = verificationFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (error instanceof ResponseError) {
				const err = await error.response.json();
				switch (error.response.status) {
					case 400: {
						if (isGenericErrorResponse(err)) {
							if (isErrorIdSessionAlreadyAvailable(err.error?.id)) {
								goto('/');
							} else if (isErrorIdSecurityCsrfViolation(err.error?.id)) {
								handleFlowErrAction(config.routes.login.path, err.error.message);
							} else if (isErrorIdSecurityIdentityMismatch(err.error?.id)) {
								goto('/');
							}
						}
						break;
					}
					case 500:
						console.error('unexpected server error');
						break;
					default:
						break;
				}
				return;
			}
			if (error instanceof FetchError) {
				console.error('fetch error: ', error.cause);
				return;
			}
			if (error instanceof RequiredError) {
				console.error('required error: ', error.field);
				return;
			}
			console.error('unexpected error');
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf
	};
}) satisfies PageLoad;
