<script lang="ts">
	import Header from '$lib/components/header/header.svelte';
	import { Button } from '$lib/components/ui/button/index';
</script>

<Header />

<main class="grid h-[calc(100vh-4rem)] place-content-center px-4">
	<div class="text-center">
		<h1 class="text-primary-200 text-9xl font-black">404</h1>
		<p class="text-primary-900 text-2xl font-bold tracking-tight sm:text-4xl">Uh-oh!</p>
		<p class="text-primary-500 mt-4">We can't find that page.</p>
		<Button href="/" class="mt-4 font-bold">Go back home</Button>
	</div>
</main>
