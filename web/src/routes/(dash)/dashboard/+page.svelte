<script lang="ts">
	import type { PageProps } from './$types';
	import Dashboard from './dashboard.svelte';

	let { data, params }: PageProps = $props();
</script>

<Dashboard {data} {params} />
