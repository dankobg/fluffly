<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { useCellContext } from './app-table';

	let props: HTMLAttributes<HTMLSpanElement> = $props();

	const cell = useCellContext<string>();
</script>

<span {...props}>{cell.getValue()}</span>
