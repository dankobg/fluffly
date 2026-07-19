<script lang="ts">
	import type { PageProps } from './$types';
	import UpdateIdentity from './update-identity.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Update identity</title>
</svelte:head>

<UpdateIdentity {data} {params} />
