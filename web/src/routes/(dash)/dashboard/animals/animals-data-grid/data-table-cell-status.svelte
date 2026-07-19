<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { AnimalStatus } from '$lib/gen/fluffly_openapi';
	import { statusIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(statusIcons.get(value as AnimalStatus));
	let color = $derived.by(() => {
		switch (value as AnimalStatus) {
			case AnimalStatus.pending:
				return 'text-yellow-400';
			case AnimalStatus.adopted:
				return 'text-green-400';
			case AnimalStatus.rejected:
				return 'text-red-400';
			case AnimalStatus.adoptable:
				return 'text-purple-400';
			case AnimalStatus.reserved:
				return 'text-sky-400';
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
