<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import ApplyForOrganization from './apply-for-organization.svelte';

	let { data, params }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Apply for organization</title>
</svelte:head>

<ApplyForOrganization {data} {params} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
