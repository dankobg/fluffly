import IconBadgeCheck from '@lucide/svelte/icons/badge-check';
import IconClock from '@lucide/svelte/icons/clock';
import IconCircleX from '@lucide/svelte/icons/circle-x';
import { OrganizationStatus } from '$lib/gen/fluffly_openapi';

export const statusIcons = new Map([
	[OrganizationStatus.pending, IconClock],
	[OrganizationStatus.approved, IconBadgeCheck],
	[OrganizationStatus.rejected, IconCircleX]
]);
export const statuses = Object.values(OrganizationStatus).map(value => ({
	label: value,
	value,
	icon: statusIcons.get(value)
}));
