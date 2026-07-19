<script lang="ts">
	import IconEye from '@lucide/svelte/icons/eye';
	import Button from '../ui/button/button.svelte';
	let { value, href }: { value?: string; href?: string } = $props();
</script>

<div class="flex w-[320px] items-center gap-2">
	{#if href}
		<Button variant="ghost" size="icon" {href}>
			<IconEye class="text-sky-400" />
		</Button>
	{/if}
	<span>{value ?? ''}</span>
</div>
