<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import UpdateAnimalSpecies from './update-animal-species.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Update animal species</title>
</svelte:head>

<UpdateAnimalSpecies {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
