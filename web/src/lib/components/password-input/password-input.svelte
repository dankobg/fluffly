<script lang="ts">
	import Input from '$lib/components/ui/input/input.svelte';
	import Toggle from '$lib/components/ui/toggle/toggle.svelte';
	import { cn } from '$lib/utils';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconEyeOff from '@lucide/svelte/icons/eye-off';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'> & { files?: undefined }>;

	let { ref = $bindable(null), value = $bindable(), class: className, ...restProps }: Props = $props();

	let visible: boolean = $state(false);

	function toggle() {
		visible = !visible;
	}
</script>

<div class="relative">
	<Input type={visible ? 'text' : 'password'} class={cn('pr-10', className)} bind:value bind:ref {...restProps} />
	<Toggle aria-label={visible ? 'hide password' : 'show password'} class="absolute top-0 right-0" onclick={toggle}>
		{#if visible}
			<IconEye />
		{:else}
			<IconEyeOff />
		{/if}
	</Toggle>
</div>
