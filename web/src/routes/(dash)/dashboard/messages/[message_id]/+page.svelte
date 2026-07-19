<script lang="ts">
	import type { PageProps } from './$types';
	import ViewMessage from './view-message.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View courier message</title>
</svelte:head>

<ViewMessage {data} {params} />
