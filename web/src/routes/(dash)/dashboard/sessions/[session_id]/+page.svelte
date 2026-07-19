<script lang="ts">
	import type { PageProps } from './$types';
	import ViewSession from './view-session.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View identity session</title>
</svelte:head>

<ViewSession {data} {params} />
