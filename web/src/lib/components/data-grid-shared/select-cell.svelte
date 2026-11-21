<script lang="ts">
	import { Checkbox } from '#lib/components/ui/checkbox/index.ts';
	import type { Checkbox as CheckboxPrimitive, WithoutChildrenOrChild } from 'bits-ui';
	import { useCellContext } from './app-table';

	let props: WithoutChildrenOrChild<CheckboxPrimitive.RootProps> = $props();

	const cell = useCellContext();
	const row = $derived(cell.row);
</script>

<Checkbox
	aria-label="Select row"
	bind:checked={() => row.getIsSelected(), value => row.toggleSelected(value)}
	{...props}
/>
