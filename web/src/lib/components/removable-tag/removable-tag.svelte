<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import IconX from '@lucide/svelte/icons/x';
	import type { ComponentProps } from 'svelte';

	let {
		class: className,
		onClose,
		children,
		...restProps
	}: ComponentProps<typeof Badge> & {
		onClose?: () => void;
	} = $props();
</script>

<Badge variant="secondary" class={cn('flex items-center gap-2 px-3 py-1 text-sm', className)} {...restProps}>
	{@render children?.()}

	<Button
		type="button"
		size="icon"
		variant="ghost"
		class="h-5 w-5 rounded-[inherit] p-0 hover:bg-black/20"
		onclick={onClose}
	>
		<IconX class="h-5 w-5" />
	</Button>
</Badge>
