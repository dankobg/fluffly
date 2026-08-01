<script lang="ts">
	import { fluffly } from '$lib/fluffly/client';
	import { AnimalStatus, OrganizationStatus, type components } from '$lib/gen/fluffly_openapi';
	import type { PageProps } from './$types';
	import { toast } from 'svelte-sonner';
	import useEmblaCarousel from 'embla-carousel-svelte';
	import { type EmblaCarouselType } from 'embla-carousel';
	import { capitalize, getErrorFallbackAnimalImage } from '$lib/utils';
	import { onDestroy } from 'svelte';
	import IconMapPin from '@lucide/svelte/icons/map-pin';
	import IconHeart from '@lucide/svelte/icons/heart';
	import IconHeartOff from '@lucide/svelte/icons/heart-off';
	import IconHandHeart from '@lucide/svelte/icons/hand-heart';
	import IconMail from '@lucide/svelte/icons/mail';
	import IconPhone from '@lucide/svelte/icons/phone';
	import IconClock from '@lucide/svelte/icons/clock';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { invalidate } from '$app/navigation';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import IconCircleQuestionMark from '@lucide/svelte/icons/circle-question-mark';
	import { MapLibre, DefaultMarker, Popup } from 'svelte-maplibre';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import Label from '$lib/components/ui/label/label.svelte';

	let { data }: PageProps = $props();

	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	let myOrganizationsResult = $state<Awaited<ReturnType<typeof listMyOrganizations>>>();
	let adoptingAs = $state<'individual' | 'organization'>('individual');
	let selectedOrganizationId = $state<number | undefined>();

	let selectedOrganizationName = $derived(
		myOrganizationsResult?.data?.data.find(x => x.id === selectedOrganizationId)?.name
	);

	$effect(() => {
		if (adoptingAs === 'organization') {
			(async () => {
				myOrganizationsResult = await listMyOrganizations();
			})();
		}
	});

	let slideEmblaApi: EmblaCarouselType | undefined;
	let thumbEmblaApi: EmblaCarouselType | undefined;

	let fallbackImageUrl = $derived(getErrorFallbackAnimalImage(data?.animalResult?.data?.type.name ?? ''));

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

	let carouselPhotos: components['schemas']['AnimalPhoto'][] = $derived.by(() => {
		if (!data?.animalResult?.data?.photos) {
			return [
				{
					id: Math.random(),
					animal_id: data?.animalResult?.data?.id as number,
					small_url: data?.animalResult?.data?.image_full_url as string,
					medium_url: data?.animalResult?.data?.image_full_url as string,
					large_url: data?.animalResult?.data?.image_full_url as string,
					full_url: data?.animalResult?.data?.image_full_url as string,
					created_at: '',
					updated_at: ''
				}
			];
		}
		return data?.animalResult?.data?.photos;
	});

	type GroupedFields = Record<string, FieldInfo[]>;
	type FieldInfo = {
		key: string;
		title: string;
		value: unknown;
		type: string;
		enum?: string[];
	};

	let groupedFields: GroupedFields = {};

	for (const [key, val] of Object.entries(data?.animalSpecieResult?.data?.properties_schema?.['properties'] ?? {})) {
		const group = val['x-group'] || 'Other';
		if (!groupedFields[group]) {
			groupedFields[group] = [];
		}

		groupedFields[group].push({
			key,
			title: val.title || key,
			value: data?.animalResult?.data?.properties?.[key],
			type: val.type,
			enum: val.enum
		});
	}

	const groupOrder = ['Traits', 'Behaviour', 'Health'];

	async function listMyOrganizations() {
		try {
			const myOrganizationsResult = await fluffly.GET('/me/organizations', {
				params: {
					query: { status: [OrganizationStatus.approved] }
				}
			});
			if (myOrganizationsResult.error) {
				toast.error(
					[myOrganizationsResult.error.message, myOrganizationsResult.error.reason].filter(Boolean).join(', ')
				);
			}
			return myOrganizationsResult;
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		}
	}

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
			invalidate(`data:animal:${id}`);
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
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
			invalidate(`data:animal:${id}`);
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		}
	}

	async function applyForAdoption(id: number, organizationId?: number) {
		try {
			const likeResult = await fluffly.POST('/animals/{id}/adoptions', {
				params: {
					path: { id }
				},
				body: { organization_id: organizationId }
			});
			if (likeResult.error) {
				toast.error([likeResult.error.message, likeResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Adoption request submitted');
			invalidate(`data:animal:${id}`);
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		}
	}

	onDestroy(() => {
		slideEmblaApi?.destroy();
		thumbEmblaApi?.destroy();
	});
</script>

{#if data?.animalResult?.data}
	{#if carouselPhotos}
		<section class="mt-8 px-4 md:px-8 lg:px-12">
			<div class="embla">
				<h1 class="mb-4 text-center text-4xl">{capitalize(data.animalResult.data.name)}</h1>
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
										alt="animal gallery {index + 1}"
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
								{#each data?.animalResult?.data?.photos ?? [] as photo, index (photo.id)}
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
												alt="animal thumbnail {index + 1}"
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
			<section class="rounded-xl bg-card p-4">
				<div class="flex items-start justify-between gap-4">
					<div class="grid gap-2">
						<h2 class="text-3xl">
							About {capitalize(data.animalResult.data.name)}
							<span class="text-xl">— the {data.animalResult.data.specie.name}</span>
						</h2>
						<Badge class="block">
							{data.animalResult.data.status}
						</Badge>
						{#if data?.animalResult?.data?.organization_id}
							<div class="inline-flex items-center gap-2">
								<IconMapPin />
								<span>
									{data?.animalResult?.data?.organization?.contact?.address.country.name} - {data?.animalResult?.data
										?.organization?.contact?.address.city}
								</span>
							</div>
						{/if}
					</div>
					<div>
						{#if data.animalResult.data?.liked}
							<Button
								size="icon"
								aria-label="Like"
								class="rounded-full"
								onclick={() => unlikeAnimal(data.animalResult.data!.id)}
							>
								<IconHeartOff class="fill-red-500" />
							</Button>
						{:else}
							<Button
								size="icon"
								aria-label="Like"
								class="rounded-full"
								onclick={() => likeAnimal(data.animalResult.data!.id)}
							>
								<IconHeart />
							</Button>
						{/if}
						{#if data.animalResult.data.likes > 0}
							<Badge variant="outline" class="h-5 min-w-5 rounded-full px-1 font-mono tabular-nums">
								{data.animalResult.data.likes}
							</Badge>
						{/if}
					</div>
				</div>

				<div class="mt-4 grid gap-4">
					<span>
						{#if data.animalResult.data.breeds}
							<strong>Breed: </strong> {data.animalResult.data.breeds.map(x => x.name).join(', ')}
						{/if}
					</span>
					<span><strong>Age: </strong> {data.animalResult.data.age}</span>
					<span><strong>Size: </strong> {data.animalResult.data.size}</span>
					<span><strong>Gender: </strong> {data.animalResult.data.gender === 'm' ? 'Male' : 'Female'}</span>
					<span><strong>Hermaphrodite: </strong> {data.animalResult.data.hermaphrodite ? 'Yes' : 'No'}</span>
					<span><strong>Description: </strong> {data.animalResult.data.description}</span>
				</div>

				{#if data.animalResult.data.tags}
					<div class="mt-4 flex flex-wrap items-center gap-2">
						{#each data.animalResult.data.tags as tag (tag.id)}
							<Badge class="bg-pink-500/70 text-fuchsia-50">{tag.name}</Badge>
						{/each}
					</div>
				{/if}

				<h3 class="mt-4 text-xl">Characteristics</h3>
				<div class="mt-4 grid gap-x-8 gap-y-8 text-sm">
					{#each groupOrder as group (group)}
						{#if groupedFields[group]}
							<div>
								<p class="text-lg">{group}</p>
								<ul>
									{#each groupedFields[group] as field (field.title)}
										<li class="mt-3 flex items-center gap-2">
											{#if field.type === 'boolean'}
												{#if field.value === true}
													<IconCheck class="text-green-600" />
												{:else if field.value === false}
													<IconX class="text-red-600" />
												{:else}
													<IconCircleQuestionMark class="text-gray-500" />
												{/if}
											{/if}
											<strong>{field.title}:</strong>
											{#if field.type === 'boolean'}
												{#if field.value === true}
													Yes
												{:else if field.value === false}
													No
												{:else}
													Unknown
												{/if}
											{:else if field.enum}
												{field.value}
											{:else}
												{field.value}
											{/if}
										</li>
									{/each}
								</ul>
							</div>
						{/if}
					{/each}

					{#each Object.keys(groupedFields).filter(g => !groupOrder.includes(g)) as group (group)}
						<div>
							<p class="text-lg">{group}</p>
							<ul>
								{#each groupedFields[group] as field (field.title)}
									{#if field.value}
										<li class="mt-3 flex items-center gap-2">
											{#if field.type === 'boolean'}
												{#if field.value === true}
													<IconCheck class="text-green-600" />
												{:else if field.value === false}
													<IconX class="text-red-600" />
												{:else}
													<IconCircleQuestionMark class="text-gray-500" />
												{/if}
											{/if}
											<strong>{field.title}:</strong>
											{#if field.type === 'boolean'}
												{#if field.value === true}
													Yes
												{:else if field.value === false}
													No
												{:else}
													Unknown
												{/if}
											{:else if field.enum}
												{field.value}
											{:else}
												{field.value}
											{/if}
										</li>
									{/if}
								{/each}
							</ul>
						</div>
					{/each}
				</div>

				{#if data?.animalResult?.data?.microchip}
					<div class="mt-8 grid">
						<h4 class="text-lg">Microchip</h4>
						<div class="mt-4 grid gap-1">
							<span><strong>Number: </strong> {data.animalResult.data.microchip.number}</span>
							{#if data.animalResult.data.microchip.brand}
								<span><strong>Brand: </strong> {data.animalResult.data.microchip.brand}</span>
							{/if}
							{#if data.animalResult.data.microchip.location}
								<span><strong>Location: </strong> {data.animalResult.data.microchip.location}</span>
							{/if}
							{#if data.animalResult.data.microchip.description}
								<span><strong>Description: </strong> {data.animalResult.data.microchip.description}</span>
							{/if}
						</div>
					</div>
				{/if}
			</section>
		</div>

		<div>
			{#if data?.animalResult?.data?.status !== AnimalStatus.pending}
				<div class="flex h-fit flex-col items-center justify-center gap-3 rounded-xl bg-green-600/30 p-4">
					{#if data?.animalResult?.data?.status === AnimalStatus.adoptable}
						<div class="text-xl">Consider adopting {capitalize(data.animalResult.data.name)}?</div>
					{:else if data?.animalResult?.data?.status === AnimalStatus.reserved}
						<div class="text-xl">{capitalize(data.animalResult.data.name)} is currently in process for adoption</div>
					{:else if data?.animalResult?.data?.status === AnimalStatus.adopted}
						<div class="text-xl">{capitalize(data.animalResult.data.name)} is adopted!</div>
						{#if data?.animalResult?.data?.adopted_at}
							<div class="text-lg text-muted">Adopted at {fmt.format(new Date(data.animalResult.data.adopted_at))}</div>
						{/if}
					{/if}
					{#if data?.animalResult?.data?.status === AnimalStatus.adoptable}
						<Dialog.Root>
							<form>
								<Dialog.Trigger type="button" class={buttonVariants({ variant: 'default' })}>
									<IconHandHeart /> Start adoption
								</Dialog.Trigger>
								<Dialog.Content class="sm:max-w-[425px]">
									<Dialog.Header>
										<Dialog.Title>How do you want to adopt {capitalize(data?.animalResult?.data?.name)}?</Dialog.Title>
										<Dialog.Description>
											You can adopt as an individual or as part of an organization.
										</Dialog.Description>
									</Dialog.Header>
									<div class="grid gap-4">
										<Field.Group>
											<Field.Set>
												<Field.Label for="compute-environment-p8w">
													I am adopting {capitalize(data?.animalResult?.data?.name)} as
												</Field.Label>
												<RadioGroup.Root bind:value={adoptingAs}>
													<Field.Label for="adopt-as-individual">
														<Field.Field orientation="horizontal">
															<Field.Content>
																<Field.Title>Individual</Field.Title>
																<Field.Description>Adopting as a regular user</Field.Description>
															</Field.Content>
															<RadioGroup.Item value="individual" id="adopt-as-individual" />
														</Field.Field>
													</Field.Label>
													<Field.Label for="adopt-as-organization">
														<Field.Field orientation="horizontal">
															<Field.Content>
																<Field.Title>Organization</Field.Title>
																<Field.Description>Adopting as part of an organization</Field.Description>
															</Field.Content>
															<RadioGroup.Item value="organization" id="adopt-as-organization" />
														</Field.Field>
													</Field.Label>
												</RadioGroup.Root>
											</Field.Set>
										</Field.Group>
									</div>

									{#if adoptingAs === 'organization'}
										{#if myOrganizationsResult?.data?.data && myOrganizationsResult?.data?.data?.length > 0}
											<div class="grid gap-2">
												<Label>Organization</Label>
												<Select.Root
													type="single"
													onValueChange={v => {
														selectedOrganizationId = Number(v);
													}}
													value={String(selectedOrganizationId)}
												>
													<Select.Trigger>
														{selectedOrganizationName || 'Select organization'}
													</Select.Trigger>
													<Select.Content>
														{#each myOrganizationsResult?.data?.data ?? [] as org (org.id)}
															<Select.Item value={`${org.id}`} label={org.name} />
														{/each}
													</Select.Content>
												</Select.Root>
											</div>
										{:else}
											<div>
												<p>You are not part of any approved organizations</p>
												<a href="/dashboard/organizations/apply" class="underline">Apply for an organization here</a>
											</div>
										{/if}
									{/if}

									<Dialog.Footer>
										<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}>Cancel</Dialog.Close>
										<Button
											type="button"
											onclick={() => {
												if (adoptingAs === 'organization' && !selectedOrganizationId) {
													toast.error('No organization selected, select one or try again as individual');
													return;
												}
												applyForAdoption(data?.animalResult?.data!.id, selectedOrganizationId);
											}}
										>
											Continue with adoption
										</Button>
									</Dialog.Footer>
								</Dialog.Content>
							</form>
						</Dialog.Root>
					{/if}
				</div>
			{/if}

			{#if data?.animalResult?.data?.organization_id && data?.animalResult?.data?.organization}
				<div class="mt-4 grid gap-2 rounded-xl bg-card p-4">
					<div class="text-xl">
						<strong>{capitalize(data?.animalResult?.data?.name)}</strong> is from
						<strong>{capitalize(data?.animalResult?.data?.organization?.name)}</strong>
					</div>
					{#if data?.animalResult?.data?.organization?.contact?.email}
						<div class="flex items-center gap-2">
							<IconMail />{data?.animalResult.data.organization.contact.email}
						</div>
					{/if}
					{#if data?.animalResult?.data?.organization?.contact?.phone}
						<div class="flex items-center gap-2">
							<IconPhone />{data?.animalResult.data.organization.contact.phone}
						</div>
					{/if}

					<div class="flex items-center gap-2">
						<IconMapPin />
						<div>
							<address>
								{[
									data?.animalResult?.data?.organization?.contact?.address?.street_address,
									data?.animalResult?.data?.organization?.contact?.address?.city,
									data?.animalResult?.data?.organization?.contact?.address?.postal_code,
									data?.animalResult?.data?.organization?.contact?.address?.country?.name
								]
									.filter(Boolean)
									.join(', ')}
							</address>
						</div>
					</div>

					{#if data?.animalResult?.data?.organization.work_hour}
						<div class="mt-4">
							<div class="mb-2 flex items-center gap-2 text-lg">
								<IconClock /> Work hours
							</div>
							<div>
								Monday: {data?.animalResult?.data?.organization.work_hour?.monday || 'unknown'}
							</div>
							<div>
								Tuesday: {data?.animalResult?.data?.organization.work_hour?.tuesday || 'unknown'}
							</div>
							<div>
								Wednesday: {data?.animalResult?.data?.organization.work_hour?.wednesday || 'unknown'}
							</div>
							<div>
								Thursday: {data?.animalResult?.data?.organization.work_hour?.thursday || 'unknown'}
							</div>
							<div>
								Friday: {data?.animalResult?.data?.organization.work_hour?.friday || 'unknown'}
							</div>
							<div>
								Saturday: {data?.animalResult?.data?.organization.work_hour?.saturday || 'unknown'}
							</div>
							<div>
								Sunday: {data?.animalResult?.data?.organization.work_hour?.sunday || 'unknown'}
							</div>
						</div>
					{/if}

					{#if data?.animalResult?.data?.organization?.contact?.address.lat && data?.animalResult?.data?.organization?.contact?.address.lon}
						<MapLibre
							class="relative mt-4 aspect-9/16 max-h-[70vh] w-full sm:aspect-video sm:max-h-full"
							style="https://tiles.openfreemap.org/styles/liberty"
							standardControls
							zoom={7}
							center={[
								data?.animalResult.data.organization.contact.address.lon,
								data?.animalResult.data.organization.contact.address.lat
							]}
						>
							{#if data?.animalResult?.data?.organization?.contact?.address?.lon && data?.animalResult?.data?.organization?.contact?.address?.lat}
								<DefaultMarker
									lngLat={[
										data?.animalResult.data.organization.contact.address.lon,
										data?.animalResult.data.organization.contact.address.lat
									]}
								>
									<Popup offset={[0, -10]}>
										<div class="text-lg font-bold text-black">
											{data?.animalResult?.data?.organization?.name}
										</div>
									</Popup>
								</DefaultMarker>
							{/if}
						</MapLibre>
					{/if}
				</div>
			{/if}
		</div>
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
