import type { PageLoad } from './$types';
import {
	FetchError,
	instanceOfGenericError,
	isGenericErrorResponse,
	RequiredError,
	ResponseError,
	type SettingsFlow
} from '@ory/client-fetch';
import { kratos } from '$lib/kratos/client';
import {
	extractCSRFToken,
	isErrorIdSecurityCsrfViolation,
	isErrorIdSecurityIdentityMismatch,
	isErrorIdSessionInactive
} from '$lib/kratos/helpers';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import { config } from '$lib/kratos/config';

export const load: PageLoad = (async ({ url, depends }) => {
	depends('data:account');
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: SettingsFlow | null = null;

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
			const settingsFlow = await kratos.getSettingsFlow({
				id: flowIdParam
			});
			flow = settingsFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (instanceOfGenericError(error)) {
				if (error.id === 'session_inactive' || error.id === 'session_refresh_required') {
					handleFlowErrAction(
						config.routes.login.path + `?return_to=${encodeURIComponent(window.location.href)}`,
						error.message
					);
				} else if (error.id === 'security_csrf_violation' || error.id === 'security_identity_mismatch') {
					handleFlowErrAction(config.routes.settings.path, error.message);
				}
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const settingsFlow = await kratos.createBrowserSettingsFlow({
				returnTo
			});
			flow = settingsFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (error instanceof ResponseError) {
				const err = await error.response.json();
				switch (error.response.status) {
					case 401:
					case 403:
					case 404:
					case 410: {
						if (isGenericErrorResponse(err)) {
							if (isErrorIdSecurityCsrfViolation(err.error?.id)) {
								handleFlowErrAction(config.routes.settings.path, err.error.message);
							} else if (isErrorIdSessionInactive(err.error?.id)) {
								handleFlowErrAction(
									config.routes.login.path + `?return_to=${encodeURIComponent(window.location.href)}`,
									err.error.message
								);
							} else if (isErrorIdSecurityIdentityMismatch(err?.error.id)) {
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
