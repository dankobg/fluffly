<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { OrganizationStatus } from '$lib/gen/fluffly_openapi';
	import { statusIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(statusIcons.get(value as OrganizationStatus));
	let color = $derived.by(() => {
		switch (value as OrganizationStatus) {
			case OrganizationStatus.pending:
				return 'text-yellow-400';
			case OrganizationStatus.approved:
				return 'text-green-400';
			case OrganizationStatus.rejected:
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
