<script lang="ts">
	import * as FormPrimitive from 'formsnap';
	import { Label } from '$lib/components/ui/label/index.js';
	import { cn, type WithoutChild } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		children,
		class: className,
		mandatory = false,
		...restProps
	}: WithoutChild<FormPrimitive.LabelProps> & { mandatory?: boolean } = $props();

	const mandatoryText = '*';
</script>

<FormPrimitive.Label {...restProps} bind:ref>
	{#snippet child({ props })}
		<Label {...props} data-slot="form-label" class={cn('data-[fs-error]:text-destructive', className)}>
			{@render children?.()}
			{#if mandatory}<span class="text-destructive">{mandatoryText}</span>{/if}
		</Label>
	{/snippet}
</FormPrimitive.Label>
