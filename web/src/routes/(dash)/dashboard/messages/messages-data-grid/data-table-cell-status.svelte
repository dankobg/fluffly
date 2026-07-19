<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { CourierMessageStatus } from '$lib/gen/fluffly_openapi';
	import { statusIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(statusIcons.get(value as CourierMessageStatus));
	let color = $derived.by(() => {
		switch (value as CourierMessageStatus) {
			case CourierMessageStatus.abandoned:
				return 'text-red-400';
			case CourierMessageStatus.processing:
				return 'text-purple-400';
			case CourierMessageStatus.queued:
				return 'text-yellow-400';
			case CourierMessageStatus.sent:
				return 'text-green-400';
			default:
				return '';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	{#if Icon}
		<Icon />
	{/if}
	<span>{value ?? ''}</span>
</Badge>
