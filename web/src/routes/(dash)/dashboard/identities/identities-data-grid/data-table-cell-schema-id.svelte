<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { schemaIdIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(value && schemaIdIcons.get(value));
	let color = $derived.by(() => {
		switch (value) {
			case 'customer':
				return 'text-blue-400';
			case 'employee':
				return 'text-purple-400';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	{#if Icon}
		<Icon />
	{/if}
	<span>{value ?? ''}</span>
</Badge>
