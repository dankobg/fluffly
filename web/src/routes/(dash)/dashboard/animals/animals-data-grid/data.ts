import IconBadgeCheck from '@lucide/svelte/icons/badge-check';
import IconHourglass from '@lucide/svelte/icons/hourglass';
import IconClock from '@lucide/svelte/icons/clock';
import IconCircleX from '@lucide/svelte/icons/circle-x';
import IconHandHeart from '@lucide/svelte/icons/hand-heart';
import { AnimalStatus } from '$lib/gen/fluffly_openapi';

export const statusIcons = new Map([
	[AnimalStatus.pending, IconClock],
	[AnimalStatus.rejected, IconCircleX],
	[AnimalStatus.adoptable, IconHandHeart],
	[AnimalStatus.adopted, IconBadgeCheck],
	[AnimalStatus.reserved, IconHourglass]
]);
export const statuses = Object.values(AnimalStatus).map(value => ({
	label: value,
	value,
	icon: statusIcons.get(value)
}));
