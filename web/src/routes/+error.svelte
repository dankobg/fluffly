<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index';
</script>

<main>
	<div class="text-center">
		<h1 class="text-primary-200 text-4xl font-black">Error {page.status}</h1>
		<p class="text-primary-900 text-2xl font-bold tracking-tight sm:text-4xl">Uh-oh!</p>
		<p class="text-primary-500 mt-4">{page.error?.message}</p>
		<Button href="/" class="mt-4 font-bold">Go back home</Button>
	</div>
</main>
