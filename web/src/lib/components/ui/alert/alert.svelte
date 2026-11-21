<script lang="ts" module>
	import { type VariantProps, tv } from 'tailwind-variants';
	import IconTriangleAlert from '@lucide/svelte/icons/triangle-alert';
	import IconCircleAlert from '@lucide/svelte/icons/circle-alert';
	import IconInfo from '@lucide/svelte/icons/info';

	export const alertVariants = tv({
		base: "grid gap-0.5 rounded-lg border px-2.5 py-2 text-left text-sm has-data-[slot=alert-action]:relative has-data-[slot=alert-action]:pr-18 has-[>svg]:grid-cols-[auto_1fr] has-[>svg]:gap-x-2 *:[svg]:row-span-2 *:[svg]:translate-y-0.5 *:[svg]:text-current *:[svg:not([class*='size-'])]:size-4 group/alert relative w-full",
		variants: {
			variant: {
				default: 'bg-card text-card-foreground',
				destructive:
					'text-destructive bg-card *:data-[slot=alert-description]:text-destructive/90 *:[svg]:text-current',
				error: 'border-rose-600/50 text-rose dark:border-rose [&>svg]:text-gree bg-rose-500/30',
				warn: 'border-orange-600/50 text-orange dark:border-orange [&>svg]:text-orange bg-orange-500/30',
				info: 'border-sky-600/50 text-sky dark:border-sky [&>svg]:text-sky bg-sky-500/30',
				success: 'border-green-600/50 text-green dark:border-green [&>svg]:text-gree bg-green-500/30'
			}
		},
		defaultVariants: {
			variant: 'default'
		}
	});

	export type AlertVariant = VariantProps<typeof alertVariants>['variant'];
</script>

<script lang="ts">
	import { cn, type WithElementRef } from '#lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		variant = 'default',
		icon = false,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		variant?: AlertVariant;
		icon: boolean;
	} = $props();
</script>

<div bind:this={ref} data-slot="alert" class={cn(alertVariants({ variant }), className)} {...restProps} role="alert">
	{#if variant !== 'default' && icon}
		{#if variant === 'destructive' || variant === 'error'}
			<IconTriangleAlert class="h-4 w-4" />
		{:else if variant === 'warn'}
			<IconCircleAlert class="h-4 w-4" />
		{:else if variant === 'info'}
			<IconInfo class="h-4 w-4" />
		{/if}
	{/if}
	{@render children?.()}
</div>
