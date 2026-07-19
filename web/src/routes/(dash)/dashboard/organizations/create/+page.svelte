<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import CreateOrganization from './create-organization.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Create organization</title>
</svelte:head>

<CreateOrganization {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
