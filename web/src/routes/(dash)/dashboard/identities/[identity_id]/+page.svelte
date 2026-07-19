<script lang="ts">
	import type { PageProps } from './$types';
	import ViewIdentity from './view-identity.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View identity</title>
</svelte:head>

<ViewIdentity {data} {params} />
