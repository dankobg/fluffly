import IconCheck from '@lucide/svelte/icons/check';
import IconX from '@lucide/svelte/icons/x';
import IconShieldUser from '@lucide/svelte/icons/shield-user';
import IconUser from '@lucide/svelte/icons/user';
import { IdentityState } from '$lib/gen/fluffly_openapi';

export const stateIcons = new Map([
	[IdentityState.active, IconCheck],
	[IdentityState.inactive, IconX]
]);
export const states = Object.values(IdentityState).map(value => ({
	label: value,
	value,
	icon: stateIcons.get(value)
}));

export const schemaIdIcons = new Map([
	['customer', IconUser],
	['developer', IconShieldUser]
]);
export const schemaIds = ['customer', 'developer'].map(value => ({
	label: value,
	value,
	icon: schemaIdIcons.get(value)
}));
