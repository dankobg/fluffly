<script lang="ts" module>
	export const dateFmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});
</script>

<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { useCellContext } from './app-table';

	let props: HTMLAttributes<HTMLSpanElement> = $props();

	const cell = useCellContext<string | undefined>();
	const val = cell.getValue();
</script>

<span {...props}>
	{#if val}
		{dateFmt.format(new Date(val))}
	{/if}
</span>
