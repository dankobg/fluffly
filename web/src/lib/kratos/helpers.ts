import {
	type ErrorBrowserLocationChangeRequired,
	type LoginFlow,
	type RecoveryFlow,
	type RegistrationFlow,
	type SettingsFlow,
	type VerificationFlow
} from '@ory/client-fetch';

type Provider = {
	name: string;
	label: string;
};
export const providers: Provider[] = [
	{ name: 'google', label: 'Google' },
	{ name: 'github', label: 'GitHub' },
	{ name: 'facebook', label: 'Facebook' },
	{ name: 'discord', label: 'Discord' },
	{ name: 'twitch', label: 'Twitch' },
	{ name: 'slack', label: 'Slack' },
	{ name: 'spotify', label: 'Spotify' }
];

export type KratosFlow = LoginFlow | RegistrationFlow | RecoveryFlow | VerificationFlow | SettingsFlow;

export function extractCSRFToken(flow: KratosFlow | null): string {
	if (!flow) {
		return '';
	}
	const csrfAttributes = flow?.ui?.nodes?.find(node => {
		return node.attributes.node_type === 'input' && node.group === 'default' && node.attributes.name === 'csrf_token';
	})?.attributes;
	return csrfAttributes?.node_type === 'input' ? csrfAttributes.value : '';
}

export function isBrowserLocationChangeRequiredResponse(
	response: unknown
): response is { error: ErrorBrowserLocationChangeRequired } {
	return (
		typeof response === 'object' &&
		!!response &&
		'error' in response &&
		typeof response.error === 'object' &&
		!!response.error &&
		'id' in response.error
	);
}

export function isErrorIdSessionAlreadyAvailable(id: string | undefined): boolean {
	return id === 'session_already_available';
}
export function isErrorIdSecurityCsrfViolation(id: string | undefined): boolean {
	return id === 'security_csrf_violation';
}
export function isErrorIdSessionAal1Required(id: string | undefined): boolean {
	return id === 'session_aal1_required';
}
export function isErrorIdSessionAal2Required(id: string | undefined): boolean {
	return id === 'session_aal2_required';
}
export function isErrorIdSecurityIdentityMismatch(id: string | undefined): boolean {
	return id === 'security_identity_mismatch';
}
export function isErrorIdSelfServiceFlowExpired(id: string | undefined): boolean {
	return id === 'self_service_flow_expired';
}
export function isErrorIdSessionRefreshRequired(id: string | undefined): boolean {
	return id === 'session_refresh_required';
}
export function isErrorIdSelfServiceFlowDisabled(id: string | undefined): boolean {
	return id === 'self_service_flow_disabled';
}
export function isErrorIdBrowserLocationChangeRequired(id: string | undefined): boolean {
	return id === 'browser_location_change_required';
}
export function isErrorIdSelfServiceFlowReplaced(id: string | undefined): boolean {
	return id === 'self_service_flow_replaced';
}
export function isErrorIdSessionVerifiedAddressRequired(id: string | undefined): boolean {
	return id === 'session_verified_address_required';
}
export function isErrorIdSessionAalAlreadyFulfilled(id: string | undefined): boolean {
	return id === 'session_aal_already_fulfilled';
}
export function isErrorIdSessionInactive(id: string | undefined): boolean {
	return id === 'session_inactive';
}
export function isErrorIdSelfServiceFlowReturnToForbidden(id: string | undefined): boolean {
	return id === 'self_service_flow_return_to_forbidden';
}
