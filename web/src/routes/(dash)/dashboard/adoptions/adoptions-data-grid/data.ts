import IconBadgeCheck from '@lucide/svelte/icons/badge-check';
import IconClock from '@lucide/svelte/icons/clock';
import IconCircleX from '@lucide/svelte/icons/circle-x';
import { AdoptionStatus } from '$lib/gen/fluffly_openapi';

export const statusIcons = new Map([
	[AdoptionStatus.pending, IconClock],
	[AdoptionStatus.approved, IconBadgeCheck],
	[AdoptionStatus.rejected, IconCircleX]
]);
export const statuses = Object.values(AdoptionStatus).map(value => ({
	label: value,
	value,
	icon: statusIcons.get(value)
}));
