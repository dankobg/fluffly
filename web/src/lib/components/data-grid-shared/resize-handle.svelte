<script lang="ts">
	import { cn } from '#lib/utils.ts';
	import { useHeaderContext } from './app-table';

	const header = useHeaderContext();
</script>

{#if header.column.getCanResize()}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		ondblclick={() => header.column.resetSize()}
		onmousedown={header.getResizeHandler()}
		ontouchstart={header.getResizeHandler()}
		class={cn(
			'absolute top-1/2 right-[-2px] z-10 h-6 w-[3px] -translate-y-1/2 cursor-e-resize touch-none rounded-md transition-colors select-none before:absolute before:top-0 before:right-[-4px] before:left-[-4px] before:h-full before:content-[""] hover:bg-blue-600',
			header.column.getIsResizing() && 'bg-blue-600'
		)}
	></div>
{/if}
