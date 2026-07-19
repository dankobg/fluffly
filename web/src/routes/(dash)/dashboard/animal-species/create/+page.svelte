<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import CreateAnimalSpecies from './create-animal-species.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Create animal species</title>
</svelte:head>

<CreateAnimalSpecies {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
