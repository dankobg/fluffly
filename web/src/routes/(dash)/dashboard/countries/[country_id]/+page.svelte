<script lang="ts">
	import type { PageProps } from './$types';
	

	let { data, params }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});


</script>

{#if data.countryResult?.data}
	<h1 class="mb-6 text-2xl font-bold">Country</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.countryResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Name</span>
			<span class="font-medium">{data.countryResult?.data.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Iso alpha 2</span>
			<span class="font-medium">{data.countryResult?.data.iso_alpha2}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Iso alpha 3</span>
			<span class="font-medium">{data.countryResult?.data.iso_alpha3}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Iso numeric</span>
			<span class="font-medium">{data.countryResult?.data.iso_numeric}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.countryResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.countryResult?.data.updated_at))}</time>
		</div>
	</div>
{/if}
