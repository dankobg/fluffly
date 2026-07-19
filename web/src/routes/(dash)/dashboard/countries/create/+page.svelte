<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import CreateCountry from './create-country.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Create country</title>
</svelte:head>

<CreateCountry {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
