<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';

	let { value }: { value: boolean } = $props();

	let Icon = $derived(value ? IconCheck : IconX);
	let color = $derived(value ? 'text-green-400' : 'text-yellow-400');
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	{#if Icon}
		<Icon />
	{/if}
	<span>{value ?? ''}</span>
</Badge>
