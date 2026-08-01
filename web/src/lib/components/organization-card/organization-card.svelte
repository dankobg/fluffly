<script lang="ts">
	import type { components } from '$lib/gen/fluffly_openapi';
	import { cn, getErrorFallbackOrganizationImage } from '$lib/utils';

	let {
		className,
		organization
	}: {
		className?: string;
		organization: components['schemas']['Organization'];
	} = $props();
</script>

<div
	class={cn(
		'__dark:bg-card relative row-span-4 grid grid-rows-subgrid gap-2 rounded-xl border bg-fuchsia-100 dark:border-0 dark:bg-emerald-950',
		className
	)}
>
	<img
		class="h-48 w-full rounded-t-xl object-cover"
		src={getErrorFallbackOrganizationImage()}
		alt={`${organization.name} ${organization.name}`}
	/>
	<a class="mt-2 px-4 text-2xl hover:underline" href={`/organizations/${organization.id}`}>
		{organization?.name}
	</a>
	<span class="px-4">{organization?.mission_statement ?? ''}</span>

	<a href={`/search?organization_id=${organization.id}`} class="inline-block p-4">
		<button
			class="transition-bg w-full cursor-pointer rounded-xl bg-green-500 p-2 font-semibold text-black duration-200 hover:bg-green-500/80"
		>
			View animals
		</button>
	</a>
</div>
