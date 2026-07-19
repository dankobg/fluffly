<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import UpdateBreed from './update-breed.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Update breed</title>
</svelte:head>

<UpdateBreed {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
