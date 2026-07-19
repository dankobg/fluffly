<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import UpdateAnimalType from './update-animal-type.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Update animal type</title>
</svelte:head>

<UpdateAnimalType {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
