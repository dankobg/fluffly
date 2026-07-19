<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { AdoptionStatus } from '$lib/gen/fluffly_openapi';
	import { statusIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(statusIcons.get(value as AdoptionStatus));
	let color = $derived.by(() => {
		switch (value as AdoptionStatus) {
			case AdoptionStatus.pending:
				return 'text-yellow-400';
			case AdoptionStatus.approved:
				return 'text-green-400';
			case AdoptionStatus.rejected:
				return 'text-red-400';
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
