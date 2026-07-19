<script lang="ts">
	import type { PageProps } from './$types';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Form from '$lib/components/ui/form';
	import { fluffly } from '$lib/fluffly/client';
	import { AnimalAge, AnimalGender, AnimalSize, type components } from '$lib/gen/fluffly_openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Item from '$lib/components/ui/item/index.js';
	import IconPlus from '@lucide/svelte/icons/plus';
	import IconX from '@lucide/svelte/icons/x';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import IconPen from '@lucide/svelte/icons/pen';
	import IconImage from '@lucide/svelte/icons/image';
	import IconVideo from '@lucide/svelte/icons/video';
	import { onDestroy, onMount } from 'svelte';
	import { Uppy } from '@uppy/core';
	import * as Alert from '$lib/components/ui/alert';
	import UppyDashboard from '@uppy/dashboard';
	import UppyImageEditor from '@uppy/image-editor';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import { animalAgeValues, animalGenderValues, animalSizeValues } from '$lib/enum-values';
	import { capitalize, getChangedFormFields } from '$lib/utils';
	import TagsInput from '$lib/components/tags-input/tags-input.svelte';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { invalidate } from '$app/navigation';
	import { emptyStringToNull } from '$lib/validation/common';

	let { data }: PageProps = $props();

	let uppyImage: Uppy | null = null;
	let uppyPhotos: Uppy | null = null;
	let uppyVideos: Uppy | null = null;

	let selectedPhotos: { id: number; key: string }[] = $state([]);
	let selectedVideos: { id: number; key: string }[] = $state([]);

	onMount(() => {
		uppyImage = new Uppy({
			id: 'update-animal-photo',
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
			target: '#update-animal-image',
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
			id: 'update-animal-photos',
			restrictions: {
				maxNumberOfFiles: 20,
				allowedFileTypes: ['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/avif']
			},
			autoProceed: false
		});
		uppyPhotos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#update-animal-photos',
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
			id: 'update-animal-videos',
			restrictions: {
				maxNumberOfFiles: 5,
				allowedFileTypes: ['video/mp4', 'video/ogg', 'video/mpeg', 'video/webm']
			},
			autoProceed: false
		});
		uppyVideos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#update-animal-videos',
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

	let updateAnimalError = $state<components['schemas']['APIError']>();

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

	let selectedOrganizationName = $derived(
		data.organizationsResult?.data?.data.find(x => x.id === $form.organization_id)?.name
	);

	let selectedAnimalTypeName = $derived(
		data.animalTypesResult?.data?.data.find(x => x.id === $form.animal_type_id)?.name
	);

	let selectedAnimalSpecie = $derived(data.animalSpeciesResult?.data?.data.find(x => x.id === $form.animal_specie_id));

	let selectedAnimalAge = $derived(ageOptions.find(x => x.value === $form.age)?.label);

	let selectedAnimalSize = $derived(sizeOptions.find(x => x.value === $form.size)?.label);

	let selectedAnimalGender = $derived(genderOptions.find(x => x.value === $form.gender)?.label);

	let selectedBreedIds = $derived.by(() => $form.animal_breeds?.map(x => String(x.breed_id)));

	let selectedBreedNames = $derived.by(() => {
		if (selectedBreedIds?.length === 0) {
			return;
		}
		return validBreeds
			?.filter(x => selectedBreedIds?.includes(String(x.id)))
			?.map(x => x.name)
			?.join(', ');
	});

	let validSpecies = $derived.by(() => {
		if ($form.animal_type_id) {
			return speciesByTypeId.get($form.animal_type_id);
		}
		if ($form.animal_breeds && $form.animal_breeds?.length > 0) {
			const out: components['schemas']['AnimalSpecie'][] = [];
			const breeds = $form.animal_breeds?.map(b => breedById.get(b.breed_id)).filter(Boolean);
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

	let hasMicrochip = $derived(Boolean(data?.animalResult?.data?.microchip));
	let ownerIsOrganization = $state<boolean>(false);

	export const updateAnimalFormSchema = v.pipe(
		v.object({
			user_id: v.optional(v.pipe(v.string(), v.uuid())),
			organization_id: v.optional(v.number()),
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
			animal_breeds: v.nullish(
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
			hermaphrodite: v.optional(v.boolean()),
			image_url: v.pipe(
				v.optional(v.pipe(v.string(), v.url('image_url must be a valid url'))),
				v.transform(x => x || undefined)
			),
			description: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.nonEmpty('description is required'))])),
			properties: v.optional(v.record(v.string(), v.any())),
			tags: v.nullish(
				v.pipe(
					v.array(
						v.object({
							name: v.pipe(v.string(), v.minLength(1, 'tag is required'))
						})
					),
					v.transform(x => (x.length === 0 ? null : x))
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
				v.nullish(
					v.object({
						number: v.pipe(v.string(), v.minLength(1, 'number is required')),
						brand: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.nonEmpty('brand is required'))])),
						description: v.nullish(
							v.union([emptyStringToNull, v.pipe(v.string(), v.nonEmpty('description is required'))])
						),
						location: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.nonEmpty('location is required'))]))
					})
				),
				v.transform(x => (Object.values(x ?? {}).some(x => x) ? x : null))
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
				[['user_id'], ['organization_id']],
				input => {
					const hasUser = input.user_id !== undefined;
					const hasOrg = input.organization_id !== undefined;
					return (hasUser && !hasOrg) || (!hasUser && hasOrg);
				},
				'provide either user_id or organization_id, but not both'
			),
			['user_id']
		),
		v.forward(
			v.partialCheck(
				[['user_id'], ['organization_id']],
				input => {
					const hasUser = input.user_id !== undefined;
					const hasOrg = input.organization_id !== undefined;
					return (hasUser && !hasOrg) || (!hasUser && hasOrg);
				},
				'provide either user_id or organization_id, but not both'
			),
			['organization_id']
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

	type UpdateAnimalFormInput = v.InferInput<typeof updateAnimalFormSchema>;
	type UpdateAnimalFormOutput = v.InferOutput<typeof updateAnimalFormSchema>;

	const initialUpdateAnimalForm: UpdateAnimalFormInput = {
		user_id: data?.auth?.user?.id, // @TODO: make nicer later
		// organization_id: null,
		animal_type_id: data.animalResult?.data?.type.id ?? null,
		animal_specie_id: data.animalResult?.data?.specie.id ?? null,
		animal_breeds: (data.animalResult?.data?.breeds ?? []).map(x => ({ breed_id: x.id, primary: x.primary })),
		name: data.animalResult?.data?.name ?? '',
		age: data.animalResult?.data?.age ?? null,
		size: data.animalResult?.data?.size ?? null,
		hermaphrodite: Boolean(data.animalResult?.data?.hermaphrodite),
		tags: (data.animalResult?.data?.tags ?? []).map(x => ({ name: x.name })),
		description: data?.animalResult?.data?.description,
		gender: data?.animalResult?.data?.gender,
		microchip: data?.animalResult?.data?.microchip
			? {
					number: data?.animalResult?.data?.microchip?.number,
					brand: data?.animalResult?.data?.microchip?.brand,
					description: data?.animalResult?.data?.microchip?.description,
					location: data?.animalResult?.data?.microchip?.location
				}
			: undefined,
		properties: data?.animalResult?.data?.properties
	};

	const supForm = superForm(initialUpdateAnimalForm, {
		id: 'update_animal',
		validators: valibot(updateAnimalFormSchema),
		SPA: true,
		dataType: 'json',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		resetForm: false,
		async onUpdate({ form }) {
			if (!form.valid) {
				toast.error('Invalid form, please fix errors and try again');
				return;
			}
			if (!data.animalResult?.data) {
				return;
			}

			const changed = getChangedFormFields(form.data, isTainted) as Partial<UpdateAnimalFormOutput>;

			try {
				const animalResult = await fluffly.PATCH('/animals/{id}', {
					params: {
						path: { id: data.animalResult.data.id }
					},
					body: {
						data: changed
					},
					bodySerializer(body) {
						const fd = new FormData();
						const data = body?.data as Partial<UpdateAnimalFormOutput> | undefined;
						if (!data) {
							throw new Error('no data');
						}
						const { image_file, photos_files, videos_files, ...dataWithoutFiles } = data;
						const animalDoc = JSON.stringify(dataWithoutFiles);
						fd.append('data', new Blob([animalDoc], { type: 'application/json' }));
						if (image_file && image_file instanceof File) {
							fd.append('image', image_file);
						}
						for (const photo of photos_files ?? []) {
							if (photo instanceof File) {
								fd.append('photos', photo);
							}
						}
						for (const video of videos_files ?? []) {
							if (video instanceof File) {
								fd.append('videos', video);
							}
						}
						return fd;
					}
				});
				updateAnimalError = animalResult.error;

				if (animalResult.data) {
					toast.success('Animal updated');
					uppyImage?.clear();
					uppyPhotos?.clear();
					uppyVideos?.clear();
				}

				if (animalResult.error) {
					toast.error([animalResult.error.message, animalResult.error.reason].filter(Boolean).join(', '));
					if (animalResult.error.code === 'validation') {
						const details = animalResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<UpdateAnimalFormInput> = {};
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
				number: data?.animalResult?.data?.microchip?.number ?? '',
				brand: data?.animalResult?.data?.microchip?.brand ?? '',
				description: data?.animalResult?.data?.microchip?.description ?? '',
				location: data?.animalResult?.data?.microchip?.location ?? ''
			};
		} else {
			$form.microchip = null;
			$errors.microchip = undefined;
		}
	});

	const { form, enhance, errors, isTainted } = supForm;

	async function deleteAnimalPhoto(animalId: number, photoId: number) {
		try {
			const deleteAnimalPhotoResult = await fluffly.DELETE('/animals/{id}/photos/{photo_id}', {
				params: {
					path: { id: animalId, photo_id: photoId }
				}
			});
			if (deleteAnimalPhotoResult.error) {
				toast.error(
					[deleteAnimalPhotoResult.error.message, deleteAnimalPhotoResult.error.reason].filter(x => x).join(', ')
				);
				return;
			}
			toast.success('Animal photo deleted');
			invalidate(`data:dashboard-animals-${animalId}-update`);
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			confirmation.closeDialog();
		}
	}

	async function deleteAnimalPhotos(animalId: number, photoIds: number[]) {
		try {
			const deleteAnimalPhotosResult = await fluffly.DELETE('/animals/{id}/photos', {
				params: {
					path: {
						id: animalId
					}
				},
				body: { ids: photoIds }
			});
			if (deleteAnimalPhotosResult.error) {
				toast.error(
					[deleteAnimalPhotosResult.error.message, deleteAnimalPhotosResult.error.reason].filter(x => x).join(', ')
				);
				return;
			}
			toast.success('Animal photos deleted');
			invalidate(`data:dashboard-animals-${animalId}-update`);
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			confirmation.closeDialog();
		}
	}

	async function deleteAnimalVideo(animalId: number, videoId: number) {
		try {
			const deleteAnimalVideoResult = await fluffly.DELETE('/animals/{id}/videos/{video_id}', {
				params: {
					path: { id: animalId, video_id: videoId }
				}
			});
			if (deleteAnimalVideoResult.error) {
				toast.error(
					[deleteAnimalVideoResult.error.message, deleteAnimalVideoResult.error.reason].filter(x => x).join(', ')
				);
				return;
			}
			toast.success('Animal video deleted');
			invalidate(`data:dashboard-animals-${animalId}-update`);
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			confirmation.closeDialog();
		}
	}

	async function deleteAnimalVideos(animalId: number, videoIds: number[]) {
		try {
			const deleteAnimalVideosResult = await fluffly.DELETE('/animals/{id}/videos', {
				params: {
					path: { id: animalId }
				},
				body: {
					ids: videoIds
				}
			});
			if (deleteAnimalVideosResult.error) {
				toast.error(
					[deleteAnimalVideosResult.error.message, deleteAnimalVideosResult.error.reason].filter(x => x).join(', ')
				);
				return;
			}
			toast.success('Animal videos deleted');
			invalidate(`data:dashboard-animals-${animalId}-update`);
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			confirmation.closeDialog();
		}
	}
</script>

{#if data.animalResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h1 class="text-2xl font-bold tracking-tight">Update animal</h1>
				<!-- <p class="text-muted-foreground">Update animal</p> -->
			</div>
		</div>
		<div>
			<form method="POST" use:enhance enctype="multipart/form-data" class="grid gap-4">
				{#if updateAnimalError}
					<Alert.Root variant="error" icon class="col-span-full">
						<Alert.Title>{updateAnimalError.message}</Alert.Title>
						{#if updateAnimalError.reason}
							<Alert.Description>{updateAnimalError.reason}</Alert.Description>
						{/if}
					</Alert.Root>
				{/if}

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
							<Checkbox bind:checked={ownerIsOrganization} />
							<div class="space-y-1 leading-none">
								<Label>Posted by organization</Label>
							</div>
						</div>
					</div>
					{#if ownerIsOrganization}
						<div class="grid gap-2">
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
						</div>
					{/if}

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
											$form.animal_breeds = $form.animal_breeds?.filter(b => {
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
											$form.animal_breeds = $form.animal_breeds?.filter(b => allowed.has(b.breed_id));
										}}
										value={String($form.animal_specie_id)}
									>
										<Select.Trigger {...props}>
											{selectedAnimalSpecie?.name || 'Select animal species'}
										</Select.Trigger>
										<Select.Content>
											{#each validSpecies ?? [] as animalSpecies (animalSpecies.id)}
												<Select.Item
													data-animal-type-id={animalSpecies.animal_type_id}
													value={String(animalSpecies.id)}
													label={animalSpecies.name}
												/>
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
											if (breedIds?.length > 0) {
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
										value={$form.animal_breeds?.map(x => String(x.breed_id)) ?? []}
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
							<Checkbox bind:checked={hasMicrochip} id="update-has-microchip" />
							<div class="space-y-1 leading-none">
								<Label for="update-has-microchip">Has microchip</Label>
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
									<TagsInput
										{...props}
										bind:value={
											() => $form?.tags?.map(x => x.name),
											v => {
												$form.tags = v?.map(x => ({ name: x }));
											}
										}
									/>
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
						<div id="update-animal-image"></div>
						{#each $errors.image_file ?? [] as e (e)}
							<div class="text-sm font-medium text-destructive">{e}</div>
						{/each}
					</div>

					<div class="col-span-full mt-8">
						<h5 class="mb-2 font-medium">Animal photos</h5>
						<div id="update-animal-photos"></div>
						{#each $errors.photos_files?._errors ?? [] as e (e)}
							<div class="text-sm font-medium text-destructive">{e}</div>
						{/each}
					</div>

					<div class="col-span-full mt-8">
						<h5 class="mb-2 font-medium">Animal videos</h5>
						<div id="update-animal-videos"></div>
						{#each $errors.videos_files?._errors ?? [] as e (e)}
							<div class="text-sm font-medium text-destructive">{e}</div>
						{/each}
					</div>

					<div class="mt-4 ml-auto">
						<Form.Button class="font-bold">Update animal</Form.Button>
					</div>

					{#if selectedPhotos.length > 0}
						<div class="mt-8 flex gap-4 self-start">
							<Button
								variant="destructive"
								onclick={() => {
									confirmation.openDialog({
										title: `Delete animal photos (${selectedPhotos.length})?`,
										destructive: true,
										async onConfirm() {
											const orgId = data.animalResult.data?.id;
											if (orgId) {
												await deleteAnimalPhotos(
													orgId,
													selectedPhotos.map(x => x.id)
												);
												invalidate(`data:dashboard-animals-${orgId}-update`);
												const keys = selectedPhotos.map(x => x.key);
												selectedPhotos = selectedPhotos.filter(x => !keys.includes(x.key));
											}
										}
									});
								}}
							>
								Delete selected photos ({selectedPhotos.length})
								<IconTrash />
							</Button>
							<Button variant="ghost" onclick={() => (selectedPhotos = [])}>
								Clear selected photos
								<IconX />
							</Button>
						</div>
					{/if}

					{#if data.animalResult.data?.photos && data.animalResult.data?.photos?.length > 0}
						<div class="mt-8 grid grid-cols-1 text-sm">
							<span class="mb-2 text-muted-foreground">Photos</span>
							<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
								{#each data.animalResult.data.photos as photo, i (photo.full_url)}
									<div class="group grid h-64 overflow-hidden [grid-template-areas:'stack']">
										<img
											src={photo.full_url}
											alt="animal photo {i + 1}"
											class="h-full w-full object-cover object-center [grid-area:stack]"
											loading="lazy"
										/>
										<div
											class={[
												'transition-opacity [grid-area:stack]',
												selectedPhotos.some(x => x.key === photo.full_url)
													? 'opacity-100'
													: 'opacity-0 group-hover:opacity-100'
											]}
										>
											<div class="flex items-center gap-2 bg-black/80 p-4">
												<Checkbox
													checked={selectedPhotos.some(x => x.key === photo.full_url)}
													onCheckedChange={v => {
														if (v) {
															selectedPhotos.push({ id: photo.id, key: photo?.full_url ?? '' });
														} else {
															selectedPhotos = selectedPhotos.filter(x => x.key !== photo.full_url);
														}
													}}
												/>
												<div class="ml-auto flex gap-2">
													<Button
														variant="destructive"
														onclick={() => {
															confirmation.openDialog({
																title: `Delete animal photo: ${photo.id}?`,
																destructive: true,
																async onConfirm() {
																	const orgId = data?.animalResult?.data?.id;
																	if (orgId) {
																		await deleteAnimalPhoto(orgId, photo.id);
																		invalidate(`data:dashboard-animals-${orgId}-update`);
																		selectedPhotos = selectedPhotos.filter(x => x.key !== photo.full_url);
																	}
																}
															});
														}}
													>
														Delete
														<IconTrash />
													</Button>
												</div>
											</div>
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}

					{#if selectedVideos.length > 0}
						<div class="mt-8 flex gap-4 self-start">
							<Button
								variant="destructive"
								onclick={() => {
									confirmation.openDialog({
										title: `Delete animal videos (${selectedVideos.length})?`,
										destructive: true,
										async onConfirm() {
											const orgId = data.animalResult.data?.id;
											if (orgId) {
												await deleteAnimalVideos(
													orgId,
													selectedVideos.map(x => x.id)
												);
												invalidate(`data:dashboard-animals-${orgId}-update`);
												const keys = selectedVideos.map(x => x.key);
												selectedVideos = selectedVideos.filter(x => !keys.includes(x.key));
											}
										}
									});
								}}
							>
								Delete selected videos ({selectedVideos.length})
								<IconTrash />
							</Button>
							<Button variant="ghost" onclick={() => (selectedVideos = [])}>
								Clear selected videos
								<IconX />
							</Button>
						</div>
					{/if}

					{#if data.animalResult.data?.videos && data.animalResult.data?.videos?.length > 0}
						<div class="mt-8 grid grid-cols-1 text-sm">
							<span class="mb-2 text-muted-foreground">Videos</span>
							<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
								{#each data.animalResult.data.videos as video (video.url)}
									<div class="group grid h-64 overflow-hidden [grid-template-areas:'stack']">
										<video controls class="h-full w-full object-cover [grid-area:stack]" muted>
											<track kind="captions" />
											<source src={video.url} />
										</video>
										<div
											class={[
												'pointer-events-none transition-opacity [grid-area:stack]',
												selectedVideos.some(x => x.key === video.url)
													? 'opacity-100'
													: 'opacity-0 group-hover:opacity-100'
											]}
										>
											<div class="pointer-events-auto flex items-center gap-2 bg-black/80 p-4">
												<Checkbox
													checked={selectedVideos.some(x => x.key === video.url)}
													onCheckedChange={v => {
														if (v) {
															selectedVideos.push({ id: video.id, key: video?.url ?? '' });
														} else {
															selectedVideos = selectedVideos.filter(x => x.key !== video.url);
														}
													}}
												/>
												<div class="ml-auto flex gap-2">
													<Button>
														Update
														<IconPen />
													</Button>
													<Button
														variant="destructive"
														onclick={() => {
															confirmation.openDialog({
																title: `Delete animal video: ${video.id}?`,
																destructive: true,
																async onConfirm() {
																	const orgId = data?.animalResult?.data?.id;
																	if (orgId) {
																		await deleteAnimalVideo(orgId, video.id);
																		invalidate(`data:dashboard-animals-${orgId}-update`);
																		selectedVideos = selectedVideos.filter(x => x.key !== video.url);
																	}
																}
															});
														}}
													>
														Delete
														<IconTrash />
													</Button>
												</div>
											</div>
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</section>
			</form>
		</div>
	</div>
{/if}
