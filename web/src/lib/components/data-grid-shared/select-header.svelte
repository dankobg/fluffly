<script lang="ts">
	import { Checkbox } from '#lib/components/ui/checkbox/index.ts';
	import type { Checkbox as CheckboxPrimitive, WithoutChildrenOrChild } from 'bits-ui';
	import { useHeaderContext } from './app-table';

	let props: WithoutChildrenOrChild<CheckboxPrimitive.RootProps> = $props();

	const header = useHeaderContext();
	const table = $derived(header.table);
</script>

<Checkbox
	aria-label="Select all"
	bind:checked={() => table.getIsAllPageRowsSelected(), value => table.toggleAllPageRowsSelected(value)}
	indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
	{...props}
/>
