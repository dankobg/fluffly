<script lang="ts">
	import { type components } from '$lib/gen/fluffly_openapi';
	import type { PageProps } from './$types';
	import useEmblaCarousel from 'embla-carousel-svelte';
	import { type EmblaCarouselType } from 'embla-carousel';
	import { capitalize, getErrorFallbackOrganizationImage } from '$lib/utils';
	import { onDestroy } from 'svelte';
	import IconMapPin from '@lucide/svelte/icons/map-pin';
	import IconMail from '@lucide/svelte/icons/mail';
	import IconPhone from '@lucide/svelte/icons/phone';
	import IconClock from '@lucide/svelte/icons/clock';
	import { MapLibre, DefaultMarker, Popup } from 'svelte-maplibre';
	import * as Item from '$lib/components/ui/item/index.js';

	let { data }: PageProps = $props();

	let slideEmblaApi: EmblaCarouselType | undefined;
	let thumbEmblaApi: EmblaCarouselType | undefined;

	let fallbackImageUrl = $derived(getErrorFallbackOrganizationImage());

	let currentIndex = $state<number>(0);

	function onInitSlideEmbla(e: CustomEvent) {
		slideEmblaApi = e.detail;

		if (slideEmblaApi) {
			slideEmblaApi.on('select', emblaApi => {
				const index = emblaApi.selectedScrollSnap();
				currentIndex = index;
				thumbEmblaApi?.scrollTo(index);
			});
		}
	}

	function onInitThumbEmbla(e: CustomEvent) {
		thumbEmblaApi = e.detail;
	}

	let carouselPhotos: components['schemas']['OrganizationPhoto'][] = $derived(
		data?.organizationResult?.data?.photos ?? []
	);

	onDestroy(() => {
		slideEmblaApi?.destroy();
		thumbEmblaApi?.destroy();
	});
</script>

{#if data?.organizationResult?.data}
	{#if carouselPhotos}
		<section class="mt-8 px-4 md:px-8 lg:px-12">
			<div class="embla">
				<h1 class="mb-4 text-center text-4xl">{capitalize(data.organizationResult.data.name)}</h1>
				<div
					class="embla-viewport"
					onemblaInit={onInitSlideEmbla}
					use:useEmblaCarousel={{
						options: {
							loop: true,
							watchDrag: carouselPhotos?.length > 1
						},
						plugins: []
					}}
				>
					<div class="embla-container">
						{#each carouselPhotos ?? [] as photo, index (photo.id)}
							<div class="embla-slide" style="--embla-slide-blur-img: url('{photo.full_url}')">
								<div>
									<img
										src={photo.full_url}
										alt="organization gallery {index + 1}"
										class="embla-slide-image"
										onerror={e => {
											(e.currentTarget as HTMLImageElement).src = fallbackImageUrl;
										}}
									/>
								</div>
							</div>
						{/each}
					</div>
				</div>

				{#if carouselPhotos.length > 1}
					<div class="embla-thumbnails">
						<div
							class="embla-thumbnails-viewport"
							onemblaInit={onInitThumbEmbla}
							use:useEmblaCarousel={{
								options: {
									containScroll: 'keepSnaps',
									dragFree: true
								},
								plugins: []
							}}
						>
							<div class="embla-thumbnails-container">
								{#each data?.organizationResult?.data?.photos ?? [] as photo, index (photo.id)}
									<div class="embla-thumbnails-slide">
										<button
											class="embla-thumbnails-btn"
											onclick={() => {
												currentIndex = index;
												slideEmblaApi?.scrollTo(index);
												thumbEmblaApi?.scrollTo(index);
											}}
										>
											<img
												src={photo.full_url}
												alt="organization thumbnail {index + 1}"
												class="embla-thumbnails-image {currentIndex === index && 'current'}"
												onerror={e => {
													(e.currentTarget as HTMLImageElement).src = (e.currentTarget as HTMLImageElement).src =
														fallbackImageUrl;
												}}
											/>
										</button>
									</div>
								{/each}
							</div>
						</div>
					</div>
				{/if}
			</div>
		</section>
	{/if}

	<div class="mx-auto my-10 grid gap-4 px-4 md:px-8 lg:grid-cols-[1fr_22rem] lg:px-12 xl:grid-cols-[1fr_32rem]">
		<div class="max-w-[80ch]">
			<section class="rounded-xl bg-emerald-950 p-4">
				<div class="flex items-start justify-between gap-4">
					<div class="grid gap-2">
						<h2 class="text-3xl">
							About {capitalize(data.organizationResult.data.name)}
						</h2>
						{#if data?.organizationResult?.data?.id}
							<div class="inline-flex items-center gap-2">
								<IconMapPin />
								<span>
									{data?.organizationResult?.data?.contact?.address.country.name} - {data?.organizationResult?.data
										?.contact?.address.city}
								</span>
							</div>
						{/if}
						<p>
							{data.organizationResult.data.mission_statement}
						</p>
						{#if data?.organizationResult?.data?.website}
							<p>
								Website: <a
									class="text-blue-400 visited:text-purple-700 hover:text-blue-500 hover:underline"
									href={data.organizationResult.data.website}
									target="_blank"
								>
									{data.organizationResult.data.website}
								</a>
							</p>
						{/if}
						{#if data?.organizationResult?.data?.adoption_policy}
							<p>
								Adoption policy: <a
									class="text-blue-400 visited:text-purple-700 hover:text-blue-500 hover:underline"
									href={data.organizationResult.data.adoption_policy}
									target="_blank"
								>
									{data.organizationResult.data.adoption_policy}
								</a>
							</p>
						{/if}
						{#if data?.organizationResult?.data?.adoption_policy}
							<p>
								Adoption URL: <a
									class="text-blue-400 visited:text-purple-700 hover:text-blue-500 hover:underline"
									href={data.organizationResult.data.adoption_url}
									target="_blank"
								>
									{data.organizationResult.data.adoption_url}
								</a>
							</p>
						{/if}
					</div>
				</div>

				{#if data.organizationResult.data?.socials?.length && data.organizationResult.data?.socials?.length > 0}
					<div class="mt-8 grid grid-cols-1 text-sm">
						<h5 class="text-lg">Social platforms</h5>
						<Item.Group>
							{#each data.organizationResult.data?.socials as social (social.id + social.platform)}
								<Item.Root class="p-2">
									<Item.Content class="gap-0">
										<Item.Title>{social.platform}</Item.Title>
										<Item.Description>
											<a href={social.url}>{social.url}</a>
										</Item.Description>
									</Item.Content>
								</Item.Root>
							{/each}
						</Item.Group>
					</div>
				{/if}

				<a href={`/search?organization_id=${data?.organizationResult?.data?.id}`} class="mt-8 inline-block">
					<button
						class="transition-bg w-full cursor-pointer rounded-xl bg-green-500 p-2 font-semibold text-black duration-200 hover:bg-green-500/80"
					>
						View our animals
					</button>
				</a>
			</section>
		</div>

		{#if data?.organizationResult?.data && data?.organizationResult?.data?.id}
			<div class="grid gap-2 rounded-xl bg-emerald-950 p-4">
				<div class="text-xl">
					<strong>{capitalize(data?.organizationResult?.data?.name)}</strong> contact info
				</div>
				{#if data?.organizationResult?.data?.contact?.email}
					<div class="flex items-center gap-2">
						<IconMail />{data?.organizationResult.data.contact.email}
					</div>
				{/if}
				{#if data?.organizationResult?.data?.contact?.phone}
					<div class="flex items-center gap-2">
						<IconPhone />{data?.organizationResult.data.contact.phone}
					</div>
				{/if}

				<div class="flex items-center gap-2">
					<IconMapPin />
					<div>
						<address>
							{[
								data?.organizationResult?.data?.contact?.address?.street_address,
								data?.organizationResult?.data?.contact?.address?.city,
								data?.organizationResult?.data?.contact?.address?.postal_code,
								data?.organizationResult?.data?.contact?.address?.country?.name
							]
								.filter(Boolean)
								.join(', ')}
						</address>
					</div>
				</div>

				{#if data?.organizationResult?.data?.work_hour}
					<div class="mt-4">
						<div class="mb-2 flex items-center gap-2 text-lg">
							<IconClock /> Work hours
						</div>
						<div>
							Monday: {data?.organizationResult?.data?.work_hour?.monday || 'unknown'}
						</div>
						<div>
							Tuesday: {data?.organizationResult?.data?.work_hour?.tuesday || 'unknown'}
						</div>
						<div>
							Wednesday: {data?.organizationResult?.data?.work_hour?.wednesday || 'unknown'}
						</div>
						<div>
							Thursday: {data?.organizationResult?.data?.work_hour?.thursday || 'unknown'}
						</div>
						<div>
							Friday: {data?.organizationResult?.data?.work_hour?.friday || 'unknown'}
						</div>
						<div>
							Saturday: {data?.organizationResult?.data?.work_hour?.saturday || 'unknown'}
						</div>
						<div>
							Sunday: {data?.organizationResult?.data?.work_hour?.sunday || 'unknown'}
						</div>
					</div>
				{/if}

				{#if data?.organizationResult?.data?.contact?.address.lat && data?.organizationResult?.data?.contact?.address.lon}
					<MapLibre
						class="relative mt-4 aspect-9/16 max-h-[70vh] w-full sm:aspect-video sm:max-h-full"
						style="https://tiles.openfreemap.org/styles/liberty"
						standardControls
						zoom={7}
						center={[
							data?.organizationResult.data.contact.address.lon,
							data?.organizationResult.data.contact.address.lat
						]}
					>
						{#if data?.organizationResult?.data?.contact?.address?.lon && data?.organizationResult?.data?.contact?.address?.lat}
							<DefaultMarker
								lngLat={[
									data?.organizationResult.data.contact.address.lon,
									data?.organizationResult.data.contact.address.lat
								]}
							>
								<Popup offset={[0, -10]}>
									<div class="text-lg font-bold text-black">
										{data?.organizationResult?.data?.name}
									</div>
								</Popup>
							</DefaultMarker>
						{/if}
					</MapLibre>
				{/if}
			</div>
		{/if}
	</div>
{/if}

<style>
	.embla {
		max-width: 60rem;
		margin-inline: auto;
	}

	.embla-viewport {
		overflow: hidden;
	}

	.embla-container {
		display: flex;
		touch-action: pan-y pinch-zoom;
	}

	.embla-slide {
		flex: 0 0 100%;
		min-width: 0;
		height: 31rem;
	}

	.embla-slide {
		position: relative;
		height: 31rem;
		overflow: hidden;
	}
	.embla-slide::before {
		content: '';
		position: absolute;
		inset: 0;
		background-image: var(--embla-slide-blur-img);
		background-size: cover;
		background-position: center;
		filter: blur(10px);
	}

	.embla-slide-image {
		height: 31rem;
		width: 100%;
		object-fit: contain;
		object-position: center;
		isolation: isolate;
	}

	.embla-thumbnails {
		height: 6rem;
		margin-top: 1rem;
	}

	.embla-thumbnails-viewport {
		overflow: hidden;
		height: 100%;
	}

	.embla-thumbnails-container {
		height: 100%;
		display: flex;
		gap: 1rem;
	}

	.embla-thumbnails-slide {
		height: 100%;
		min-width: 0;
		flex: 0 0 auto;
	}

	.embla-thumbnails-image {
		width: 100%;
		height: 6rem;
		object-fit: contain;
		object-position: center;
		display: block;
	}

	.embla-thumbnails-image.current {
		border: 3px solid gold;
	}
</style>
