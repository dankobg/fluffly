<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import UpdateCountry from './update-country.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Update country</title>
</svelte:head>

<UpdateCountry {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
