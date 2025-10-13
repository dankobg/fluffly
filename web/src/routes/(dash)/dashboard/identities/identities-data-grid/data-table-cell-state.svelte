<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { IdentityState } from '$lib/gen/fluffly_openapi';
	import { stateIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(stateIcons.get(value as IdentityState));
	let color = $derived.by(() => {
		switch (value as IdentityState) {
			case IdentityState.active:
				return 'text-green-400';
			case IdentityState.inactive:
				return 'text-red-400';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	{#if Icon}
		<Icon />
	{/if}
	<span>{value ?? ''}</span>
</Badge>
