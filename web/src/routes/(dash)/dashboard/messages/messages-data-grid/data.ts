import IconMailX from '@lucide/svelte/icons/mail-x';
import IconRefreshCw from '@lucide/svelte/icons/refresh-cw';
import IconLayers from '@lucide/svelte/icons/layers';
import IconMailCheck from '@lucide/svelte/icons/mail-check';
import IconMail from '@lucide/svelte/icons/mail';
import IconSmartphone from '@lucide/svelte/icons/smartphone';
import IconShieldCheck from '@lucide/svelte/icons/shield-check';
import IconShieldMinus from '@lucide/svelte/icons/shield-minus';
import IconShieldPlus from '@lucide/svelte/icons/shield-plus';
import IconShieldQuestion from '@lucide/svelte/icons/shield-question';
import IconShieldUser from '@lucide/svelte/icons/shield-user';
import IconShieldX from '@lucide/svelte/icons/shield-x';
import IconAsterisk from '@lucide/svelte/icons/asterisk';
import { CourierMessageStatus, CourierMessageTemplateType, CourierMessageType } from '$lib/gen/fluffly_openapi';

export const statusIcons = new Map([
	[CourierMessageStatus.abandoned, IconMailX],
	[CourierMessageStatus.processing, IconRefreshCw],
	[CourierMessageStatus.queued, IconLayers],
	[CourierMessageStatus.sent, IconMailCheck]
]);
export const statuses = Object.values(CourierMessageStatus).map(value => ({
	label: value,
	value,
	icon: statusIcons.get(value)
}));

export const typeIcons = new Map([
	[CourierMessageType.email, IconMail],
	[CourierMessageType.phone, IconSmartphone]
]);
export const types = Object.values(CourierMessageType).map(value => ({
	label: value,
	value,
	icon: typeIcons.get(value)
}));

export const templateTypeIcons = new Map([
	[CourierMessageTemplateType.recovery_invalid, IconShieldX],
	[CourierMessageTemplateType.recovery_valid, IconShieldCheck],
	[CourierMessageTemplateType.recovery_code_invalid, IconShieldMinus],
	[CourierMessageTemplateType.recovery_code_valid, IconShieldPlus],
	[CourierMessageTemplateType.verification_invalid, IconShieldX],
	[CourierMessageTemplateType.verification_valid, IconShieldCheck],
	[CourierMessageTemplateType.verification_code_invalid, IconShieldMinus],
	[CourierMessageTemplateType.verification_code_valid, IconShieldPlus],
	[CourierMessageTemplateType.stub, IconAsterisk],
	[CourierMessageTemplateType.login_code_valid, IconShieldQuestion],
	[CourierMessageTemplateType.registration_code_valid, IconShieldUser]
]);
export const templateTypes = Object.values(CourierMessageTemplateType).map(value => ({
	label: value,
	value,
	icon: templateTypeIcons.get(value)
}));
