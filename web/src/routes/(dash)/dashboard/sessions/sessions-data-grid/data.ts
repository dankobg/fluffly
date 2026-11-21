import { AuthenticatorAssuranceLevel } from '#lib/gen/fluffly_openapi.ts';

export const aals = Object.values(AuthenticatorAssuranceLevel).map(value => ({
	label: value,
	value
}));
