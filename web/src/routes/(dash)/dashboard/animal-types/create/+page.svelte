<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import CreateAnimalType from './create-animal-type.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Create animal type</title>
</svelte:head>

<CreateAnimalType {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
