<script lang="ts" module>
	type TData = unknown;
	type TValue = unknown;
</script>

<script lang="ts" generics="TData, TValue">
	import IconEyeOff from '@lucide/svelte/icons/eye-off';
	import IconArrowDown from '@lucide/svelte/icons/arrow-down';
	import IconArrowUp from '@lucide/svelte/icons/arrow-up';
	import IconChevronsUpDown from '@lucide/svelte/icons/chevrons-up-down';
	import type { HTMLAttributes } from 'svelte/elements';
	import type { Column } from '@tanstack/table-core';
	import type { WithoutChildren } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import Button from '$lib/components/ui/button/button.svelte';

	type Props = HTMLAttributes<HTMLDivElement> & {
		column: Column<TData, TValue>;
		title: string;
	};

	let { column, class: className, title, ...restProps }: WithoutChildren<Props> = $props();
</script>

{#if !column?.getCanSort()}
	<div class={className} {...restProps}>
		{title}
	</div>
{:else}
	<div class={cn('flex items-center', className)} {...restProps}>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button {...props} variant="ghost" size="sm" class="-ml-3 h-8 data-[state=open]:bg-accent">
						<span>
							{title}
						</span>
						{#if column.getIsSorted() === 'desc'}
							<IconArrowDown />
						{:else if column.getIsSorted() === 'asc'}
							<IconArrowUp />
						{:else}
							<IconChevronsUpDown />
						{/if}
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="start">
				<DropdownMenu.Item onclick={() => column.toggleSorting(false)}>
					<IconArrowUp class="mr-2 size-3.5 text-muted-foreground/70" />
					Asc
				</DropdownMenu.Item>
				<DropdownMenu.Item onclick={() => column.toggleSorting(true)}>
					<IconArrowDown class="mr-2 size-3.5 text-muted-foreground/70" />
					Desc
				</DropdownMenu.Item>
				<DropdownMenu.Separator />
				<DropdownMenu.Item onclick={() => column.toggleVisibility(false)}>
					<IconEyeOff class="mr-2 size-3.5 text-muted-foreground/70" />
					Hide
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>
{/if}
