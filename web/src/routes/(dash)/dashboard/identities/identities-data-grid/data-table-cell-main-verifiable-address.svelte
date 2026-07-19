<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';

	let { address, verified }: { address?: string; verified: boolean } = $props();

	let color = $derived(verified ? 'text-green-400' : 'text-yellow-600');
</script>

<div class="flex items-center gap-2">
	<span>{address ?? ''}</span> -
	<Badge variant="outline" class="border {color} w-fit">
		{verified ? 'verified' : 'unverified'}
	</Badge>
</div>
