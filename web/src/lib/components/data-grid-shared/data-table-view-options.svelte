<script lang="ts">
	import Settings2Icon from '@lucide/svelte/icons/settings-2';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.ts';
	import { buttonVariants } from '#lib/components/ui/button/index.ts';
	import { useTableContext } from './app-table';

	const table = useTableContext();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={buttonVariants({
			variant: 'outline',
			size: 'sm',
			class: 'ms-auto hidden h-8 lg:flex'
		})}
	>
		<Settings2Icon />
		View
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Group>
			<DropdownMenu.Label>Toggle columns</DropdownMenu.Label>
			<DropdownMenu.Separator />
			{#each table
				.getAllColumns()
				.filter(col => typeof col.accessorFn !== 'undefined' && col.getCanHide()) as column (column.id)}
				<DropdownMenu.CheckboxItem
					bind:checked={() => column.getIsVisible(), v => column.toggleVisibility(!!v)}
					class="capitalize"
				>
					{column.id}
				</DropdownMenu.CheckboxItem>
			{/each}
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
