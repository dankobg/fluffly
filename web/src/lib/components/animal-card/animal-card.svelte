<script lang="ts">
	import type { components } from '$lib/gen/fluffly_openapi';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import IconHeart from '@lucide/svelte/icons/heart';
	import IconHeartOff from '@lucide/svelte/icons/heart-off';
	import { capitalize, cn, getErrorFallbackAnimalImage } from '$lib/utils';
	import Button from '$lib/components/ui/button/button.svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';

	let {
		className,
		animal,
		showLikeUnlikeButton = false,
		onLiked,
		onUnliked,
		showPostedBy = false
	}: {
		className?: string;
		animal: components['schemas']['Animal'];
		showLikeUnlikeButton?: boolean;
		onLiked?: () => void;
		onUnliked?: () => void;
		showPostedBy?: boolean;
	} = $props();

	async function likeAnimal(id: number) {
		try {
			const likeResult = await fluffly.POST('/animals/{id}/like', {
				params: {
					path: { id }
				}
			});
			if (likeResult.error) {
				toast.error([likeResult.error.message, likeResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Animal liked');
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			onLiked?.();
		}
	}

	async function unlikeAnimal(id: number) {
		try {
			const likeResult = await fluffly.POST('/animals/{id}/unlike', {
				params: {
					path: { id }
				}
			});
			if (likeResult.error) {
				toast.error([likeResult.error.message, likeResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Animal unliked');
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			onUnliked?.();
		}
	}
</script>

<div
	class={cn(
		'relative row-span-7 grid grid-rows-subgrid gap-2 rounded-xl border bg-fuchsia-100 dark:border-0 dark:bg-card',
		className
	)}
>
	{#if showLikeUnlikeButton}
		<div class="absolute top-1 right-1 flex flex-col items-center gap-1">
			{#if animal?.liked}
				<Button size="icon" aria-label="Like" class="rounded-full" onclick={() => unlikeAnimal(animal.id)}>
					<IconHeartOff class="fill-red-500" />
				</Button>
			{:else}
				<Button size="icon" aria-label="Like" class="rounded-full" onclick={() => likeAnimal(animal.id)}>
					<IconHeart />
				</Button>
			{/if}
			{#if animal.likes > 0}
				<Badge
					variant="outline"
					class="h-5 min-w-5 rounded-full bg-background px-1 font-mono tabular-nums dark:bg-background/40"
				>
					{animal.likes}
				</Badge>
			{/if}
		</div>
	{/if}

	{#if showPostedBy}
		<Badge variant="default" class="absolute top-45 left-4 bg-sky-900 text-white">
			{#if animal.organization_id}
				From Org
			{:else}
				From User
			{/if}
		</Badge>
	{/if}

	<img
		class="h-48 w-full rounded-t-xl object-cover"
		src={animal.image_full_url}
		alt={`${animal.type.name} ${animal.name}`}
		onerror={e => {
			const url = getErrorFallbackAnimalImage(animal.type.name);
			(e.currentTarget as HTMLImageElement).src = url;
		}}
	/>
	<span class="mt-2 px-4 text-2xl">{animal.name}</span>
	<span class="px-4">{animal.specie.name}</span>
	<span class="px-4">
		{#if animal.breeds}
			breed: {animal.breeds.map(x => x.name).join(', ')}
		{/if}
	</span>

	<ul class="flex flex-wrap items-center gap-2 px-4">
		<li>{capitalize(animal.age)}</li>
		<li class="before:mx-2 before:text-current before:content-['•']">
			{capitalize(animal.size)}
		</li>
		{#if animal.gender}
			<li class="before:mx-2 before:text-current before:content-['•']">
				{animal.gender === 'm' ? 'Male' : 'Female'}
			</li>
		{/if}
	</ul>

	<div class="flex flex-wrap items-center gap-2 px-4">
		{#if animal.tags}
			{#each animal.tags as tag (tag.id)}
				<Badge class="bg-pink-600/70 text-fuchsia-50">{tag.name}</Badge>
			{/each}
		{/if}
	</div>

	<a href={`/animal/${animal.id}`} class="inline-block p-4">
		<button
			class="transition-bg w-full cursor-pointer rounded-xl bg-purple-400 p-2 font-semibold text-black duration-200 hover:bg-purple-400/80"
		>
			View profile
		</button>
	</a>
</div>
