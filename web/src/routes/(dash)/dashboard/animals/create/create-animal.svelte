<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { fluffly } from '$lib/fluffly/client';
	import { AnimalAge, AnimalGender, AnimalSize, AnimalStatus, type components } from '$lib/gen/fluffly_openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Item from '$lib/components/ui/item/index.js';
	import IconPlus from '@lucide/svelte/icons/plus';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import IconImage from '@lucide/svelte/icons/image';
	import IconVideo from '@lucide/svelte/icons/video';
	import IconChevronsUpDown from '@lucide/svelte/icons/chevrons-up-down';
	import IconCheck from '@lucide/svelte/icons/check';
	import { onDestroy, onMount, tick } from 'svelte';
	import { Uppy } from '@uppy/core';
	import UppyDashboard from '@uppy/dashboard';
	import UppyImageEditor from '@uppy/image-editor';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import { animalAgeValues, animalGenderValues, animalSizeValues, animalStatusValues } from '$lib/enum-values';
	import { capitalize, cn } from '$lib/utils';
	import TagsInput from '$lib/components/tags-input/tags-input.svelte';
	import { emptyStringToUndefined, zeroToUndefined } from '$lib/validation/common';
	import { useDebounce } from 'runed';

	let {
		data,
		// eslint-disable-next-line no-useless-assignment
		capture = $bindable(),
		// eslint-disable-next-line no-useless-assignment
		restore = $bindable()
	}: PageProps & { capture: Snapshot['capture']; restore: Snapshot['restore'] } = $props();

	let uppyImage: Uppy | null = null;
	let uppyPhotos: Uppy | null = null;
	let uppyVideos: Uppy | null = null;

	onMount(() => {
		uppyImage = new Uppy({
			id: 'create-animal-photo',
			restrictions: {
				minNumberOfFiles: 1,
				maxNumberOfFiles: 1,
				allowedFileTypes: ['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/avif']
			},
			autoProceed: false
		});
		uppyImage.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#create-animal-image',
			note: 'Image only, 1 file (jpg, png, gif)',
			width: '100%',
			height: '300px',
			hideUploadButton: true,
			locale: {
				strings: {
					dropPasteBoth: 'Drop main photo file here, %{browseFiles} or %{browseFolders}',
					dropPasteFiles: 'Drop main photo file here or %{browseFiles}',
					dropPasteFolders: 'Drop main photo file here or %{browseFolders}',
					dropPasteImportBoth: 'Drop main photo file here, %{browseFiles}, %{browseFolders} or import from:',
					dropPasteImportFiles: 'Drop main photo file here, %{browseFiles} or import from:',
					dropPasteImportFolders: 'Drop main photo file here, %{browseFolders} or import from:'
				}
			}
		});
		uppyImage.use(UppyImageEditor);
		uppyImage.on('file-added', file => {
			if (file.data instanceof File) {
				$form.image_file = file.data;
			}
		});
		uppyImage.on('file-removed', file => {
			if (file.data instanceof File) {
				$form.image_file = undefined;
			}
		});

		uppyPhotos = new Uppy({
			id: 'create-animal-photos',
			restrictions: {
				maxNumberOfFiles: 20,
				allowedFileTypes: ['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/avif']
			},
			autoProceed: false
		});
		uppyPhotos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#create-animal-photos',
			note: 'Images only, up to 20 files (jpg, png, gif)',
			width: '100%',
			height: '300px',
			hideUploadButton: true,
			locale: {
				strings: {
					dropPasteBoth: 'Drop photo files here, %{browseFiles} or %{browseFolders}',
					dropPasteFiles: 'Drop photo files here or %{browseFiles}',
					dropPasteFolders: 'Drop photo files here or %{browseFolders}',
					dropPasteImportBoth: 'Drop photo files here, %{browseFiles}, %{browseFolders} or import from:',
					dropPasteImportFiles: 'Drop photo files here, %{browseFiles} or import from:',
					dropPasteImportFolders: 'Drop photo files here, %{browseFolders} or import from:'
				}
			}
		});
		uppyPhotos.use(UppyImageEditor);
		uppyPhotos.on('files-added', files => {
			for (const file of files) {
				if (file.data instanceof File) {
					$form.photos_files ||= [];
					$form.photos_files = [...$form.photos_files, file.data];
				}
			}
		});
		uppyPhotos.on('file-removed', file => {
			if (file.data instanceof File) {
				if ($form?.photos_files) {
					$form.photos_files = $form.photos_files.filter(x => x.name !== file.name);
				}
			}
		});

		uppyVideos = new Uppy({
			id: 'create-animal-videos',
			restrictions: {
				maxNumberOfFiles: 5,
				allowedFileTypes: ['video/mp4', 'video/ogg', 'video/mpeg', 'video/webm']
			},
			autoProceed: false
		});
		uppyVideos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#create-animal-videos',
			note: 'Videos only, up to 5 files',
			width: '100%',
			height: '300px',
			hideUploadButton: true,
			locale: {
				strings: {
					dropPasteBoth: 'Drop video files here, %{browseFiles} or %{browseFolders}',
					dropPasteFiles: 'Drop video files here or %{browseFiles}',
					dropPasteFolders: 'Drop video files here or %{browseFolders}',
					dropPasteImportBoth: 'Drop video files here, %{browseFiles}, %{browseFolders} or import from:',
					dropPasteImportFiles: 'Drop video files here, %{browseFiles} or import from:',
					dropPasteImportFolders: 'Drop video files here, %{browseFolders} or import from:'
				}
			}
		});
		uppyVideos.on('files-added', files => {
			for (const file of files) {
				if (file.data instanceof File) {
					$form.videos_files ||= [];
					$form.videos_files = [...$form.videos_files, file.data];
				}
			}
		});
		uppyVideos.on('file-removed', file => {
			if (file.data instanceof File) {
				if ($form?.videos_files) {
					$form.videos_files = $form.videos_files.filter(x => x.name !== file.name);
				}
			}
		});
	});

	onDestroy(() => {
		uppyImage?.destroy();
		uppyPhotos?.destroy();
		uppyVideos?.destroy();
	});

	let createdAnimal: components['schemas']['Animal'] | null = $state(null);
	let createdAnimalTypeName = $state<string>('');
	let createdAnimalSpeciesName = $state<string>('');
	let createdAnimalBreedNames = $state<string>('');
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	const typeById = new Map(data?.animalTypesResult?.data?.data.map(t => [t.id, t]));
	let speciesById = new Map<number, components['schemas']['AnimalSpecie']>();
	let breedById = new Map<number, components['schemas']['Breed']>();
	let speciesByTypeId = new Map<number, components['schemas']['AnimalSpecie'][]>();
	let breedsBySpecieId = new Map<number, components['schemas']['Breed'][]>();

	for (const s of data?.animalSpeciesResult?.data?.data ?? []) {
		speciesById.set(s.id, s);

		const list = speciesByTypeId.get(s.animal_type_id);
		if (list) {
			list.push(s);
		} else {
			speciesByTypeId.set(s.animal_type_id, [s]);
		}
	}
	for (const b of data?.animalBreedsResult?.data?.data ?? []) {
		breedById.set(b.id, b);

		const list = breedsBySpecieId.get(b.animal_specie_id);
		if (list) {
			list.push(b);
		} else {
			breedsBySpecieId.set(b.animal_specie_id, [b]);
		}
	}

	const ageOptions = animalAgeValues.map(x => ({
		value: x,
		label: capitalize(x)
	}));

	const sizeOptions = animalSizeValues.map(x => ({
		value: x,
		label: capitalize(x)
	}));

	const genderOptions = animalGenderValues.map(x => ({
		value: x,
		label: x === 'm' ? 'Male' : 'Female'
	}));

	const statusOptions = animalStatusValues.map(x => ({
		value: x,
		label: capitalize(x)
	}));

	let organizations: components['schemas']['Organization'][] = $state([]);

	let selectedOrganizationName = $derived(organizations.find(x => x.id === $form.organization_id)?.name);

	let selectedAnimalTypeName = $derived(
		data.animalTypesResult?.data?.data.find(x => x.id === $form.animal_type_id)?.name
	);

	let selectedAnimalSpecie = $derived(data.animalSpeciesResult?.data?.data.find(x => x.id === $form.animal_specie_id));

	let selectedAnimalAge = $derived(ageOptions.find(x => x.value === $form.age)?.label);

	let selectedAnimalSize = $derived(sizeOptions.find(x => x.value === $form.size)?.label);

	let selectedAnimalGender = $derived(genderOptions.find(x => x.value === $form.gender)?.label);

	let selectedBreedIds = $derived.by(() => ($form.animal_breeds ?? []).map(x => String(x.breed_id)));

	let selectedStatusName = $derived(statusOptions.find(x => x.value === $form.status)?.label ?? 'Select status');

	let selectedBreedNames = $derived.by(() => {
		if (selectedBreedIds?.length === 0) {
			return;
		}
		return validBreeds
			?.filter(x => selectedBreedIds.includes(String(x.id)))
			?.map(x => x.name)
			?.join(', ');
	});

	let validSpecies = $derived.by(() => {
		if ($form.animal_type_id) {
			return speciesByTypeId.get($form.animal_type_id);
		}
		if ($form.animal_breeds && $form.animal_breeds.length > 0) {
			const out: components['schemas']['AnimalSpecie'][] = [];
			const breeds = $form.animal_breeds.map(b => breedById.get(b.breed_id)).filter(Boolean);
			for (const breed of breeds) {
				if (breed?.animal_specie_id) {
					const specie = speciesById.get(breed.animal_specie_id);
					if (specie) {
						out.push(specie);
					}
				}
			}
			return out;
		}
		return data?.animalSpeciesResult?.data?.data;
	});

	let validBreeds = $derived.by(() => {
		if (!$form.animal_specie_id && $form.animal_type_id) {
			const species = speciesByTypeId.get($form.animal_type_id);
			const out: components['schemas']['Breed'][] = [];
			for (const s of species ?? []) {
				const breeds = breedsBySpecieId.get(s.id);
				if (breeds) {
					out.push(...breeds);
				}
			}
			return out;
		}
		if ($form.animal_specie_id) {
			return breedsBySpecieId.get($form.animal_specie_id);
		}
		return data?.animalBreedsResult?.data?.data;
	});

	let hasMicrochip = $state<boolean>(false);
	let ownerIsOrganization = $state<boolean>(false);

	export const createAnimalFormSchema = v.pipe(
		v.object({
			user_id: v.pipe(v.string(), v.minLength(1, 'user_id is required')),
			organization_id: v.optional(
				v.union([zeroToUndefined, v.pipe(v.number(), v.minValue(1, 'organization_id is required'))])
			),
			animal_type_id: v.pipe(
				v.nullable(v.number()),
				v.number('animal_type_id is required'),
				v.minValue(1, 'animal_type_id is required')
			),
			animal_specie_id: v.pipe(
				v.nullable(v.number()),
				v.number('animal_specie_id is required'),
				v.minValue(1, 'animal_specie_id is required')
			),
			animal_breeds: v.optional(
				v.pipe(
					v.array(
						v.object({
							breed_id: v.pipe(v.number(), v.minValue(1, 'breed_id is required')),
							primary: v.boolean()
						})
					),
					v.transform(x => (x.length === 0 ? undefined : x))
				)
			),
			name: v.pipe(v.string(), v.minLength(1, 'name is required')),
			age: v.pipe(v.nullable(v.picklist(animalAgeValues)), v.picklist(animalAgeValues, 'age is required')),
			size: v.pipe(v.nullable(v.picklist(animalSizeValues)), v.picklist(animalSizeValues, 'size is required')),
			gender: v.optional(v.picklist(animalGenderValues)),
			status: v.optional(v.picklist(animalStatusValues)),
			hermaphrodite: v.optional(v.boolean()),
			image_url: v.optional(
				v.union([emptyStringToUndefined, v.pipe(v.string(), v.url('image_url must be a valid url'))])
			),
			description: v.optional(
				v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'description is required'))])
			),
			properties: v.optional(v.record(v.string(), v.any())),
			tags: v.optional(
				v.pipe(
					v.array(v.pipe(v.string(), v.minLength(1, 'tag is required'))),
					v.transform(x => (x.length === 0 ? undefined : x))
				)
			),
			photos: v.optional(
				v.array(
					v.object({
						url: v.pipe(v.string(), v.url('url must be a valid url'))
					})
				)
			),
			videos: v.optional(
				v.array(
					v.object({
						url: v.pipe(v.string(), v.url('url must be a valid url'))
					})
				)
			),
			microchip: v.pipe(
				v.optional(
					v.object({
						number: v.pipe(v.string(), v.minLength(1, 'number is required')),
						brand: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'brand is required'))])
						),
						description: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'description is required'))])
						),
						location: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'location is required'))])
						)
					})
				),
				v.transform(x => (Object.values(x ?? {}).some(x => x) ? x : undefined))
			),
			image_file: v.optional(
				v.pipe(
					v.file('main image is required'),
					v.mimeType(
						['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/avif'],
						'please select a JPEG or PNG file'
					)
				)
			),
			photos_files: v.optional(
				v.array(
					v.pipe(
						v.file('invalid file'),
						v.mimeType(
							['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/avif'],
							'please select JPEG or PNG files'
						)
					)
				)
			),
			videos_files: v.optional(
				v.array(
					v.pipe(
						v.file('invalid file'),
						v.mimeType(['video/mp4', 'video/ogg', 'video/mpeg', 'video/webm'], 'please select MP4 or OGG files')
					)
				)
			)
		}),
		v.forward(
			v.partialCheck(
				[['image_url'], ['image_file']],
				input => {
					const hasImageUrl = input.image_url !== undefined;
					const hasPhotoFile = input.image_file !== undefined;
					return hasImageUrl || hasPhotoFile;
				},
				'provide either image_url or image_file'
			),
			['image_url']
		),
		v.forward(
			v.partialCheck(
				[['image_url'], ['image_file']],
				input => {
					const hasImageUrl = input.image_url !== undefined;
					const hasPhotoFile = input.image_file !== undefined;
					return hasImageUrl || hasPhotoFile;
				},
				'provide either image_url or image_file'
			),
			['image_file']
		),
		v.forward(
			v.partialCheck(
				[['image_file'], ['photos_files'], ['videos_files']],
				input => {
					const files: File[] = [
						...(input.image_file ? [input.image_file] : []),
						...(input.photos_files ?? []),
						...(input.videos_files ?? [])
					];
					let total = 0;
					for (const f of files) {
						total += f.size;
					}
					return total < 100 << 20;
				},
				'total combined files size must be less than 100 MiB'
			),
			['image_file']
		),
		v.forward(
			v.partialCheck(
				[['image_file'], ['photos_files'], ['videos_files']],
				input => {
					const files: File[] = [
						...(input.image_file ? [input.image_file] : []),
						...(input.photos_files ?? []),
						...(input.videos_files ?? [])
					];
					let total = 0;
					for (const f of files) {
						total += f.size;
					}
					return total < 100 << 20;
				},
				'total combined files size must be less than 100 MiB'
			),
			['photos_files']
		),
		v.forward(
			v.partialCheck(
				[['image_file'], ['photos_files'], ['videos_files']],
				input => {
					const files: File[] = [
						...(input.image_file ? [input.image_file] : []),
						...(input.photos_files ?? []),
						...(input.videos_files ?? [])
					];
					let total = 0;
					for (const f of files) {
						total += f.size;
					}
					return total < 100 << 20;
				},
				'total combined files size must be less than 100 MiB'
			),
			['videos_files']
		)
	);

	type CreateAnimalFormInput = v.InferInput<typeof createAnimalFormSchema>;
	type CreateAnimalFormOutput = v.InferOutput<typeof createAnimalFormSchema>;

	const initialCreateAnimalForm: CreateAnimalFormInput = {
		user_id: data?.auth?.user?.id ?? '', // @TODO: make nicer later
		// organization_id: null,
		animal_type_id: null,
		animal_specie_id: null,
		animal_breeds: [],
		name: '',
		age: null,
		size: null,
		status: AnimalStatus.adoptable,
		hermaphrodite: false,
		tags: []
	};

	const supForm = superForm(initialCreateAnimalForm, {
		id: 'create_animal',
		validators: valibot(createAnimalFormSchema),
		SPA: true,
		dataType: 'json',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		resetForm: true,
		async onUpdate({ form }) {
			if (!form.valid) {
				toast.error('Invalid form, please fix errors and try again');
				return;
			}

			try {
				const animalResult = await fluffly.POST('/animals', {
					body: {
						data: form.data as CreateAnimalFormOutput
					},
					bodySerializer(body) {
						const fd = new FormData();
						const data = body?.data as CreateAnimalFormOutput | undefined;
						if (!data) {
							throw new Error('no data');
						}
						const dataWithoutFiles = {
							...data,
							image_file: undefined,
							photos_files: undefined,
							videos_files: undefined
						};
						const animalDoc = JSON.stringify(dataWithoutFiles);
						fd.append('data', new Blob([animalDoc], { type: 'application/json' }));
						if (data.image_file instanceof File) {
							fd.append('image', data.image_file);
						}
						for (const photo of data.photos_files ?? []) {
							if (photo instanceof File) {
								fd.append('photos', photo);
							}
						}
						for (const video of data.videos_files ?? []) {
							if (video instanceof File) {
								fd.append('videos', video);
							}
						}
						return fd;
					}
				});
				if (animalResult.data) {
					createdAnimal = animalResult.data;
					createdAnimalTypeName = selectedAnimalTypeName || '';
					createdAnimalSpeciesName = selectedAnimalSpecie?.name || '';
					createdAnimalBreedNames = selectedBreedNames || '';

					toast.success('Animal created');
					reset();
					uppyImage?.clear();
					uppyPhotos?.clear();
					uppyVideos?.clear();
				}
				if (animalResult.error) {
					toast.error([animalResult.error.message, animalResult.error.reason].filter(Boolean).join(', '));

					if (animalResult.error.code === 'validation') {
						const details = animalResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<CreateAnimalFormInput> = {};
						for (const detail of details) {
							if (detail.in === 'body') {
								const path = detail.pointer.replaceAll('/data', '').replaceAll('/', '.').slice(1);
								set(fieldErrors, path, [detail.reason]);
							}
						}
						errors.set(fieldErrors);
					}
				}
			} catch (error: unknown) {
				if (error instanceof TypeError) {
					toast.error(`${error}`);
					return;
				}
				console.error('unexpected error');
			}
		}
	});

	$effect(() => {
		if (hasMicrochip) {
			$form.microchip = {
				number: '',
				brand: '',
				description: '',
				location: ''
			};
		} else {
			$form.microchip = undefined;
			$errors.microchip = undefined;
		}
	});

	const { form, enhance, errors, reset } = supForm;
	// eslint-disable-next-line no-useless-assignment
	capture = supForm.capture;
	// eslint-disable-next-line no-useless-assignment
	restore = supForm.restore;

	let organizationsNameFilter = $state('');
	let orgComboboxOpen = $state(false);
	let triggerRef = $state<HTMLButtonElement>(null!);

	function closeAndFocusTrigger() {
		orgComboboxOpen = false;
		tick().then(() => {
			triggerRef.focus();
		});
	}

	const debounceOrganizations = useDebounce(
		async () => {
			if (organizationsNameFilter.length >= 3) {
				const organizationsResult = await fluffly.GET('/organizations', {
					params: { query: { page_size: 100, name: organizationsNameFilter || undefined } }
				});
				if (organizationsResult.data) {
					organizations = organizationsResult.data.data;
				}
			}
		},
		() => 700
	);
</script>

<div class="flex h-full flex-1 flex-col space-y-8 p-8">
	<div class="flex items-center justify-between space-y-2">
		<div>
			<h1 class="text-2xl font-bold tracking-tight">Create animal</h1>
			<!-- <p class="text-muted-foreground">Create new animal</p> -->
		</div>
	</div>
	<div>
		<form method="POST" use:enhance enctype="multipart/form-data" class="grid gap-4">
			<section class="grid gap-4 md:grid-cols-2">
				<h2 class="col-span-full text-muted-foreground">Main information</h2>
				<div class="grid gap-2">
					<Form.Field form={supForm} name="name">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label mandatory>Name</Form.Label>
								<Input {...props} bind:value={$form.name} />
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<div class="flex items-center gap-2">
						<Checkbox
							id="create-owner-is-org"
							checked={ownerIsOrganization}
							onCheckedChange={flag => {
								ownerIsOrganization = flag;
							}}
						/>
						<div class="space-y-1 leading-none">
							<Label for="create-owner-is-org">Posted by organization</Label>
						</div>
					</div>
				</div>
				{#if ownerIsOrganization}
					<Popover.Root bind:open={orgComboboxOpen}>
						<Popover.Trigger bind:ref={triggerRef}>
							{#snippet child({ props })}
								<Button
									variant="outline"
									class="justify-between"
									{...props}
									role="combobox"
									aria-expanded={orgComboboxOpen}
								>
									{selectedOrganizationName || 'Select organization'}
									<IconChevronsUpDown class="ms-2 size-4 shrink-0 opacity-50" />
								</Button>
							{/snippet}
						</Popover.Trigger>
						<Popover.Content class="p-0">
							<Command.Root shouldFilter={false}>
								<Command.Input
									placeholder="Type min 3 letters..."
									value={organizationsNameFilter}
									oninput={e => {
										organizationsNameFilter = e.currentTarget.value;
										if (organizationsNameFilter.length >= 3) {
											debounceOrganizations();
										}
									}}
								/>
								<Command.List>
									<Command.Empty>No organizations found</Command.Empty>
									<Command.Group>
										{#each organizations as org (org.id)}
											<Command.Item
												value={String($form.organization_id)}
												onSelect={() => {
													$form.organization_id = org.id;
													closeAndFocusTrigger();
												}}
											>
												<IconCheck class={cn('me-2 size-4', $form.organization_id !== org.id && 'text-transparent')} />
												{org.name}
											</Command.Item>
										{/each}
									</Command.Group>
								</Command.List>
							</Command.Root>
						</Popover.Content>
					</Popover.Root>
				{/if}

				<!-- <div class="grid gap-2">
							<Form.Field form={supForm} name="organization_id">
								<Form.Control>
									{#snippet children({ props })}
										<Form.Label mandatory>Organization</Form.Label>
										<Select.Root
											type="single"
											onValueChange={v => {
												$form.organization_id = Number(v);
											}}
											value={String($form.organization_id)}
										>
											<Select.Trigger {...props}>
												{selectedOrganizationName || 'Select organization'}
											</Select.Trigger>
											<Select.Content>
												{#each data.organizationsResult?.data?.data ?? [] as org (org.id)}
													<Select.Item value={`${org.id}`} label={org.name} />
												{/each}
											</Select.Content>
										</Select.Root>
									{/snippet}
								</Form.Control>
								<Form.Description />
								<Form.FieldErrors />
							</Form.Field>
						</div> -->

				<div class="grid gap-2">
					<Form.Field form={supForm} name="animal_type_id">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label mandatory>Animal type</Form.Label>
								<Select.Root
									type="single"
									onValueChange={v => {
										const animalTypeId = Number(v);
										$form.animal_type_id = animalTypeId;
										const species = speciesByTypeId.get(animalTypeId) ?? [];
										if (species.length === 1) {
											$form.animal_specie_id = species[0]!.id;
										} else {
											$form.animal_specie_id = null;
										}
										const allowed = new Set(species.map(s => s.id));
										$form.animal_breeds = ($form.animal_breeds ?? []).filter(b => {
											const breed = breedById.get(b.breed_id);
											return breed && allowed.has(breed.animal_specie_id);
										});
									}}
									value={String($form.animal_type_id)}
								>
									<Select.Trigger {...props}>
										{selectedAnimalTypeName || 'Select animal type'}
									</Select.Trigger>
									<Select.Content>
										{#each data?.animalTypesResult?.data?.data ?? [] as animalType (animalType.id)}
											<Select.Item value={String(animalType.id)} label={animalType.name} />
										{/each}
									</Select.Content>
								</Select.Root>
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="animal_specie_id">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label mandatory>Animal species</Form.Label>
								<Select.Root
									type="single"
									onValueChange={v => {
										const animalSpecieId = Number(v);
										$form.animal_specie_id = animalSpecieId;
										const species = speciesById.get(animalSpecieId);
										if (species && $form.animal_type_id !== species.animal_type_id) {
											$form.animal_type_id = species.animal_type_id;
										}
										const allowed = new Set(breedsBySpecieId.get(animalSpecieId)?.map(b => b.id) ?? []);
										$form.animal_breeds = ($form.animal_breeds ?? []).filter(b => allowed.has(b.breed_id));
									}}
									value={String($form.animal_specie_id)}
								>
									<Select.Trigger {...props}>
										{selectedAnimalSpecie?.name || 'Select animal species'}
									</Select.Trigger>
									<Select.Content>
										{#each validSpecies ?? [] as animalSpecies (animalSpecies.id)}
											<Select.Item value={String(animalSpecies.id)} label={animalSpecies.name} />
										{/each}
									</Select.Content>
								</Select.Root>
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="animal_breeds">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label>Animal breeds</Form.Label>
								<Select.Root
									type="multiple"
									onValueChange={v => {
										const breedIds = v.map(Number);
										$form.animal_breeds = breedIds.map((x, i) => ({ breed_id: x, primary: i === 0 }));
										if (breedIds.length > 0) {
											const breed = breedById.get(breedIds[0]!);
											if (!breed) {
												return;
											}
											const species = speciesById.get(breed.animal_specie_id);
											if (!species) {
												return;
											}
											const type = typeById.get(species.animal_type_id);
											if (!type) {
												return;
											}
											if ($form.animal_specie_id !== species.id) {
												$form.animal_specie_id = species.id;
											}
											if ($form.animal_type_id !== type.id) {
												$form.animal_type_id = type.id;
											}
										}
									}}
									value={($form.animal_breeds ?? []).map(x => String(x.breed_id))}
								>
									<Select.Trigger {...props}>
										{selectedBreedNames || 'Select animal breeds'}
									</Select.Trigger>
									<Select.Content>
										{#each validBreeds ?? [] as animalBreed (animalBreed.id)}
											<Select.Item value={String(animalBreed.id)} label={animalBreed.name} />
										{/each}
									</Select.Content>
								</Select.Root>
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="age">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label mandatory>Age</Form.Label>
								<Select.Root
									type="single"
									onValueChange={v => {
										$form.age = (v as AnimalAge) ?? null;
									}}
									value={$form.age ?? ''}
								>
									<Select.Trigger {...props}>
										{selectedAnimalAge || 'Select age'}
									</Select.Trigger>
									<Select.Content>
										{#each ageOptions as age (age.value)}
											<Select.Item value={`${age.value}`} label={age.label} />
										{/each}
									</Select.Content>
								</Select.Root>
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="size">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label mandatory>Size</Form.Label>
								<Select.Root
									type="single"
									onValueChange={v => {
										$form.size = (v as AnimalSize) ?? null;
									}}
									value={$form.size ?? ''}
								>
									<Select.Trigger {...props}>
										{selectedAnimalSize || 'Select size'}
									</Select.Trigger>
									<Select.Content>
										{#each sizeOptions as size (size.value)}
											<Select.Item value={`${size.value}`} label={size.label} />
										{/each}
									</Select.Content>
								</Select.Root>
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="gender">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label>Gender</Form.Label>
								<Select.Root
									type="single"
									onValueChange={v => {
										$form.gender = (v as AnimalGender) ?? undefined;
									}}
									value={$form.gender}
								>
									<Select.Trigger {...props}>
										{selectedAnimalGender || 'Select gender'}
									</Select.Trigger>
									<Select.Content>
										{#each genderOptions as gender (gender.value)}
											<Select.Item value={`${gender.value}`} label={gender.label} />
										{/each}
									</Select.Content>
								</Select.Root>
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="status">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label mandatory>Status</Form.Label>
								<Select.Root
									type="single"
									onValueChange={v => {
										$form.status = (v as AnimalStatus) || null;
									}}
									value={$form.status ?? undefined}
								>
									<Select.Trigger {...props}>
										{selectedStatusName}
									</Select.Trigger>
									<Select.Content>
										{#each statusOptions as status (status.value)}
											<Select.Item value={`${status.value}`} label={status.label} />
										{/each}
									</Select.Content>
								</Select.Root>
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="description">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label>Description</Form.Label>
								<Textarea {...props} bind:value={$form.description} />
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="hermaphrodite">
						<Form.Control>
							{#snippet children({ props })}
								<div class="flex items-center gap-2">
									<Checkbox {...props} bind:checked={$form.hermaphrodite} />
									<div class="space-y-1 leading-none">
										<Form.Label>Hermaphrodite</Form.Label>
										<Form.Description />
									</div>
								</div>
							{/snippet}
						</Form.Control>
					</Form.Field>
				</div>

				<div class="grid gap-2">
					<div class="flex items-center gap-2">
						<Checkbox bind:checked={hasMicrochip} id="create-has-microchip" />
						<div class="space-y-1 leading-none">
							<Label for="create-has-microchip">Has microchip</Label>
						</div>
					</div>
				</div>
				{#if hasMicrochip && $form.microchip}
					<div class="grid gap-2">
						<Form.Field form={supForm} name="microchip.number">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label mandatory>Microchip number</Form.Label>
									<Input {...props} bind:value={$form.microchip!.number} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="microchip.brand">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Microchip brand</Form.Label>
									<Input {...props} bind:value={$form.microchip!.brand} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="microchip.description">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Microchip description</Form.Label>
									<Input {...props} bind:value={$form.microchip!.description} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="microchip.location">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Microchip location</Form.Label>
									<Input {...props} bind:value={$form.microchip!.location} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
				{/if}

				<div class="grid gap-2">
					<Form.Field form={supForm} name="tags">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label>Tags</Form.Label>
								<TagsInput {...props} bind:value={$form.tags} />
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				{#if selectedAnimalSpecie?.properties_schema?.['properties']}
					{#each Object.entries(selectedAnimalSpecie.properties_schema?.['properties']) as [key, val] (key)}
						{#if typeof val === 'object'}
							{@const value = val as Record<string, unknown>}

							{#if value?.['type'] === 'boolean'}
								<div class="grid gap-2">
									<Form.Field form={supForm} name={`properties.${key}`}>
										<Form.Control>
											{#snippet children({ props })}
												<Label>{value?.['title']}</Label>
												<RadioGroup.Root
													{...props}
													class="flex items-center gap-4"
													bind:value={
														() => {
															const v = $form?.properties?.[key];
															if (v === true) {
																return 'yes';
															}
															if (v === false) {
																return 'no';
															}
															return 'unknown';
														},
														v => {
															$form.properties ||= {};
															if (v === 'yes') {
																$form.properties[key] = true;
															} else if (v === 'no') {
																$form.properties[key] = false;
															} else {
																$form.properties[key] = undefined;
															}
														}
													}
												>
													<div class="flex items-center space-x-2">
														<RadioGroup.Item value="unknown" id={`${key}-unknown`} />
														<Label for={`${key}-unknown`}>Unknown</Label>
													</div>
													<div class="flex items-center space-x-2">
														<RadioGroup.Item value="yes" id={`${key}-yes`} />
														<Label for={`${key}-yes`}>Yes</Label>
													</div>
													<div class="flex items-center space-x-2">
														<RadioGroup.Item value="no" id={`${key}-no`} />
														<Label for={`${key}-no`}>No</Label>
													</div>
												</RadioGroup.Root>
											{/snippet}
										</Form.Control>
									</Form.Field>
								</div>
							{/if}

							{#if value?.['type'] === 'string'}
								{#if value?.['enum']}
									<div class="grid gap-2">
										<Form.Field form={supForm} name={`properties.${key}`}>
											<Form.Control>
												{#snippet children({ props })}
													<Form.Label>{value?.['title']}</Form.Label>
													<Select.Root
														type="single"
														bind:value={
															() => $form?.properties?.[key] ?? '',
															v => {
																$form.properties ||= {};
																$form.properties[key] = v as string;
															}
														}
													>
														<Select.Trigger {...props}>
															{$form.properties?.[key] ?? `Select ${value?.['title']}`}
														</Select.Trigger>
														<Select.Content>
															{#each value['enum'] as string[] as item (item)}
																<Select.Item value={`${item}`} label={capitalize(item)} />
															{/each}
														</Select.Content>
													</Select.Root>
												{/snippet}
											</Form.Control>
											<Form.Description />
											<Form.FieldErrors />
										</Form.Field>
									</div>
								{:else}
									<div class="grid gap-2">
										<Form.Field form={supForm} name={`properties.${key}`}>
											<Form.Control>
												{#snippet children({ props })}
													<Form.Label>{value?.['title']}</Form.Label>
													<Input
														{...props}
														bind:value={
															() => $form?.properties?.[key] ?? '',
															v => {
																$form.properties ||= {};
																$form.properties[key] = v as string;
															}
														}
													/>
												{/snippet}
											</Form.Control>
											<Form.Description />
											<Form.FieldErrors />
										</Form.Field>
									</div>
								{/if}
							{/if}
						{/if}
					{/each}
				{/if}

				<div class="grid gap-2">
					<Form.Field form={supForm} name="image_url">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label>Image URL <span class="text-orange-400">*</span></Form.Label>
								<Input {...props} bind:value={$form.image_url} />
							{/snippet}
						</Form.Control>
						<Form.Description />
						<Form.FieldErrors />
					</Form.Field>
				</div>

				<section class="grid gap-4 md:grid-cols-2">
					<Item.Root variant="muted" class="col-span-full">
						<Item.Media variant="icon">
							<IconImage />
						</Item.Media>
						<Item.Content>
							<Item.Title>Photo URLs</Item.Title>
							<Item.Description>Add photo URLs</Item.Description>
						</Item.Content>
						<Item.Actions>
							<Button
								size="sm"
								variant="outline"
								onclick={() => {
									$form.photos ||= [];
									$form.photos = [...$form.photos, { url: '' }];
								}}
							>
								Add
								<IconPlus />
							</Button>
						</Item.Actions>
					</Item.Root>
					<Item.Group class="col-span-full grid gap-2">
						{#if $form.photos}
							{#each $form.photos as _, i (i)}
								<Item.Root variant="outline">
									<Item.Content>
										<Item.Title class="mb-2">
											Photo {i + 1}
										</Item.Title>
										<Form.Field form={supForm} name="photos[{i}].url">
											<Form.Control>
												{#snippet children({ props })}
													<Form.Label mandatory>URL</Form.Label>
													<Input {...props} bind:value={$form.photos![i]!.url} />
												{/snippet}
											</Form.Control>
											<Form.Description />
											<Form.FieldErrors />
										</Form.Field>
									</Item.Content>
									<Item.Actions>
										<Button
											size="sm"
											variant="outline"
											onclick={() => {
												if ($form.photos) {
													$form.photos = $form.photos.filter((_, idx) => idx !== i);
												}
											}}
										>
											Remove
											<IconTrash />
										</Button>
									</Item.Actions>
								</Item.Root>
							{/each}
						{/if}
					</Item.Group>
				</section>

				<section class="grid gap-4 md:grid-cols-2">
					<Item.Root variant="muted" class="col-span-full">
						<Item.Media variant="icon">
							<IconVideo />
						</Item.Media>
						<Item.Content>
							<Item.Title>Video URLs</Item.Title>
							<Item.Description>Add video URLs</Item.Description>
						</Item.Content>
						<Item.Actions>
							<Button
								size="sm"
								variant="outline"
								onclick={() => {
									$form.videos ||= [];
									$form.videos = [...$form.videos, { url: '' }];
								}}
							>
								Add
								<IconPlus />
							</Button>
						</Item.Actions>
					</Item.Root>
					<Item.Group class="col-span-full grid gap-2">
						{#if $form.videos}
							{#each $form.videos as _, i (i)}
								<Item.Root variant="outline">
									<Item.Content>
										<Item.Title class="mb-2">
											Video {i + 1}
										</Item.Title>
										<Form.Field form={supForm} name="videos[{i}].url">
											<Form.Control>
												{#snippet children({ props })}
													<Form.Label mandatory>URL</Form.Label>
													<Input {...props} bind:value={$form.videos![i]!.url} />
												{/snippet}
											</Form.Control>
											<Form.Description />
											<Form.FieldErrors />
										</Form.Field>
									</Item.Content>
									<Item.Actions>
										<Button
											size="sm"
											variant="outline"
											onclick={() => {
												if ($form.videos) {
													$form.videos = $form.videos.filter((_, idx) => idx !== i);
												}
											}}
										>
											Remove
											<IconTrash />
										</Button>
									</Item.Actions>
								</Item.Root>
							{/each}
						{/if}
					</Item.Group>
				</section>

				<div class="col-span-full mt-8">
					<h5 class="mb-2 font-medium">Animal main photo <span class="text-orange-400">*</span></h5>
					<div id="create-animal-image"></div>
					{#each $errors.image_file ?? [] as e (e)}
						<div class="text-sm font-medium text-destructive">{e}</div>
					{/each}
				</div>

				<div class="col-span-full mt-8">
					<h5 class="mb-2 font-medium">Animal photos</h5>
					<div id="create-animal-photos"></div>
					{#each $errors.photos_files?._errors ?? [] as e (e)}
						<div class="text-sm font-medium text-destructive">{e}</div>
					{/each}
				</div>

				<div class="col-span-full mt-8">
					<h5 class="mb-2 font-medium">Animal videos</h5>
					<div id="create-animal-videos"></div>
					{#each $errors.videos_files?._errors ?? [] as e (e)}
						<div class="text-sm font-medium text-destructive">{e}</div>
					{/each}
				</div>

				<div class="mt-4 ml-auto">
					<Form.Button class="font-bold">Create animal</Form.Button>
				</div>
			</section>
		</form>
	</div>

	{#if createdAnimal}
		<Card.Root class="max-w-6xl">
			<Card.Header>
				<Card.Title>View created animal</Card.Title>
			</Card.Header>

			<Card.Content>
				<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">ID</span>
						<a class="font-medium text-sky-500 underline" href="/dashboard/animals/{createdAnimal.id}">
							{createdAnimal.id}
						</a>
					</div>

					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Name</span>
						<span class="font-medium">{createdAnimal.name}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Type</span>
						<span class="font-medium">{createdAnimalTypeName}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Species</span>
						<span class="font-medium">{createdAnimalSpeciesName}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Breeds</span>
						<span class="font-medium">{createdAnimalBreedNames}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Gender</span>
						<span class="font-medium">{createdAnimal.gender}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Hermaphrodite</span>
						<span class="font-medium">{createdAnimal.hermaphrodite}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Age</span>
						<span class="font-medium">{createdAnimal.age}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Size</span>
						<span class="font-medium">{createdAnimal.size}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Status</span>
						<span class="font-medium">{createdAnimal.status}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Description</span>
						<span class="font-medium">{createdAnimal.description}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Created time</span>
						<time class="font-medium">{fmt.format(new Date(createdAnimal.created_at))}</time>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Updated time</span>
						<time class="font-medium">{fmt.format(new Date(createdAnimal.updated_at))}</time>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
</div>
