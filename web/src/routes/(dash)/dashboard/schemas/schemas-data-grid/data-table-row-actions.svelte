<script lang="ts">
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconCopy from '@lucide/svelte/icons/copy';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.ts';
	import Button from '#lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { useCellContext } from '#lib/components/data-grid-shared/app-table.ts';

	const cell = useCellContext();
	const row = $derived(cell.row);

	const hasId = $derived(
		typeof row.original === 'object' && !!row.original && 'id' in row.original && typeof row.original.id === 'string'
	);

	function copyIdToClipboard() {
		try {
			navigator.clipboard.writeText(row.original.id).then(() => {
				toast.success('id coppied');
			});
		} catch (error) {
			if (Error.isError(error)) {
				toast.error('failed to copy id: ' + error.message);
			}
		}
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="ghost" class="flex h-8 w-8 p-0 data-[state=open]:bg-muted">
				<IconEllipsis />
				<span class="sr-only">Open Menu</span>
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		{#if hasId}
			<DropdownMenu.Item onclick={copyIdToClipboard}>
				<IconCopy />
				Copy ID to clipboard
			</DropdownMenu.Item>
		{/if}
		<a href="/dashboard/schemas/{row.getValue('id')}">
			<DropdownMenu.Item class="cursor-pointer">
				<IconEye />
				View
			</DropdownMenu.Item>
		</a>
	</DropdownMenu.Content>
</DropdownMenu.Root>
