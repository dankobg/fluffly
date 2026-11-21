<script lang="ts">
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import EyeOffIcon from '@lucide/svelte/icons/eye-off';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import { cn } from '#lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import { useHeaderContext } from './app-table';

	let {
		title,
		class: className,
		...restProps
	}: {
		title: string;
	} & HTMLAttributes<HTMLDivElement> = $props();

	const header = useHeaderContext();
	const column = $derived(header.column);

	const displayTitle = $derived(column.columnDef.meta?.label ?? title ?? column.id);
</script>

{#if !header.column?.getCanSort()}
	<div class={className} {...restProps}>
		{title}
	</div>
{:else}
	<div class={cn('flex items-center', className)} {...restProps}>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button {...props} variant="ghost" size="sm" class="-ms-3 h-8 data-[state=open]:bg-accent">
						<span>
							{displayTitle}
						</span>
						{#if header.column.getIsSorted() === 'desc'}
							<ArrowDownIcon />
						{:else if header.column.getIsSorted() === 'asc'}
							<ArrowUpIcon />
						{:else}
							<ChevronsUpDownIcon />
						{/if}
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="start">
				<DropdownMenu.Item onclick={() => header.column.toggleSorting(false)}>
					<ArrowUpIcon class="me-2 size-3.5 text-muted-foreground/70" />
					Asc
				</DropdownMenu.Item>
				<DropdownMenu.Item onclick={() => header.column.toggleSorting(true)}>
					<ArrowDownIcon class="me-2 size-3.5 text-muted-foreground/70" />
					Desc
				</DropdownMenu.Item>
				<DropdownMenu.Separator />
				<DropdownMenu.Item onclick={() => header.column.toggleVisibility(false)}>
					<EyeOffIcon class="me-2 size-3.5 text-muted-foreground/70" />
					Hide
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>
{/if}
