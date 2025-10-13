import type { PageLoad } from './$types';
import { FetchError, isGenericErrorResponse, RequiredError, ResponseError, type LoginFlow } from '@ory/client-fetch';
import { kratos } from '$lib/kratos/client';
import {
	extractCSRFToken,
	isErrorIdSecurityCsrfViolation,
	isErrorIdSecurityIdentityMismatch,
	isErrorIdSelfServiceFlowExpired,
	isErrorIdSessionAal1Required,
	isErrorIdSessionAal2Required,
	isErrorIdSessionAlreadyAvailable
} from '$lib/kratos/helpers';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import { config } from '$lib/kratos/config';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');
	const refreshParam = browser && url.searchParams.get('refresh');
	const aalParam = browser && url.searchParams.get('aal');
	const loginChallengeParam = browser && url.searchParams.get('login_challenge');

	let flow: LoginFlow | null = null;

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
			const loginFlow = await kratos.getLoginFlow({
				id: flowIdParam
			});
			flow = loginFlow;
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
									goto(`${config.routes.login.path}?return_to=${window.location.href}`);
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
		const aal: string | undefined = aalParam ? aalParam.toString() : undefined;
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;
		const refresh: boolean | undefined = refreshParam?.toString().toLowerCase() === 'true' ? true : undefined;
		const loginChallenge: string | undefined = loginChallengeParam ? loginChallengeParam.toString() : undefined;
		try {
			const loginFlow = await kratos.createBrowserLoginFlow({
				aal,
				returnTo,
				refresh,
				loginChallenge
			});
			flow = loginFlow;
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
							} else if (isErrorIdSessionAal1Required(err.error?.id)) {
								goto(`${config.routes.login.path}?aal=aal1&return_to=${window.location.href}`);
							} else if (isErrorIdSecurityIdentityMismatch(err.error?.id)) {
								goto('/');
							} else if (isErrorIdSessionAal2Required(err.error?.id)) {
								if (browser) {
									goto(`${config.routes.login.path}?aal=aal2&return_to=${window.location.href}`);
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
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf
	};
}) satisfies PageLoad;
