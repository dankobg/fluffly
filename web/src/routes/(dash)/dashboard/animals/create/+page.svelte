<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import CreateAnimal from './create-animal.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Create animal</title>
</svelte:head>

<CreateAnimal {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
