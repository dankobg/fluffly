<script lang="ts">
	import Header from '$lib/components/header/header.svelte';

	let { children } = $props();
</script>

<Header />
{@render children()}
