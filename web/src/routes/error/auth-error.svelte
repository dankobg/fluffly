<script lang="ts">
	import Header from '$lib/components/shell/header.svelte';
	import SiteSidebar from '$lib/components/shell/site-sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Button } from '$lib/components/ui/button/index';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider class="flex flex-col" open={false}>
		<Header user={data.auth.user} logoutUrl={data.logoutUrl} />
		<div class="flex flex-1">
			<SiteSidebar user={data?.auth?.user} logoutUrl={data?.logoutUrl} />
			<Sidebar.Inset class="mx-auto max-w-[120rem]">
				<div class="mt-16 text-center">
					<h1 class="text-primary-200 text-4xl font-black">Error</h1>
					<p class="text-primary-900 text-2xl font-bold tracking-tight sm:text-4xl">Authentication error occurred.</p>
					{#if data?.flow?.id}
						<p class="text-primary-500 mt-4">Error id: {data.flow.id}</p>
					{/if}
					<Button href="/" class="mt-4 font-bold">Go back home</Button>
				</div>
			</Sidebar.Inset>
		</div>
	</Sidebar.Provider>
</div>
