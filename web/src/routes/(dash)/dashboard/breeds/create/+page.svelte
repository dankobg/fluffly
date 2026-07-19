<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import CreateBreed from './create-breed.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Create breed</title>
</svelte:head>

<CreateBreed {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
