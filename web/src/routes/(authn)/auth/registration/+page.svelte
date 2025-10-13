<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import Registration from './registration.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Registration</title>
</svelte:head>

<Registration {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
