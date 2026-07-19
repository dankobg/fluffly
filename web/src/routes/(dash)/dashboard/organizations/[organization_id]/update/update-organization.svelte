<script lang="ts">
	import type { PageProps } from './$types';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Form from '$lib/components/ui/form';
	import { fluffly } from '$lib/fluffly/client';
	import type { components, OrganizationStatus } from '$lib/gen/fluffly_openapi';
	import * as Alert from '$lib/components/ui/alert';
	import IconPlus from '@lucide/svelte/icons/plus';
	import IconX from '@lucide/svelte/icons/x';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import IconPen from '@lucide/svelte/icons/pen';
	import IconGlobe from '@lucide/svelte/icons/globe';
	import IconImage from '@lucide/svelte/icons/image';
	import IconVideo from '@lucide/svelte/icons/video';
	import { onDestroy, onMount } from 'svelte';
	import { Uppy } from '@uppy/core';
	import UppyDashboard from '@uppy/dashboard';
	import UppyImageEditor from '@uppy/image-editor';
	import Button from '$lib/components/ui/button/button.svelte';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { invalidate } from '$app/navigation';
	import * as Item from '$lib/components/ui/item/index.js';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import { capitalize, getChangedFormFields } from '$lib/utils';
	import { emptyStringToNull } from '$lib/validation/common';
	import { organizationStatusValues } from '$lib/enum-values';

	let { data }: PageProps = $props();

	const statusOptions = organizationStatusValues.map(x => ({
		value: x,
		label: capitalize(x)
	}));

	let uppyPhotos: Uppy | null = null;
	let uppyVideos: Uppy | null = null;

	let selectedPhotos: { id: number; key: string }[] = $state([]);
	let selectedVideos: { id: number; key: string }[] = $state([]);

	onMount(() => {
		uppyPhotos = new Uppy({
			id: 'update-organization-photos',
			restrictions: {
				maxNumberOfFiles: 20,
				allowedFileTypes: ['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/avif']
			},
			autoProceed: false
		});
		uppyPhotos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#update-organization-photos',
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
			id: 'update-organization-videos',
			restrictions: {
				maxNumberOfFiles: 5,
				allowedFileTypes: ['video/mp4', 'video/ogg', 'video/mpeg', 'video/webm']
			},
			autoProceed: false
		});
		uppyVideos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#update-organization-videos',
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
		uppyPhotos?.destroy?.();
		uppyVideos?.destroy?.();
	});

	let updateOrganizationError = $state<components['schemas']['APIError']>();

	let selectedStatusName = $derived(statusOptions.find(x => x.value === $form.status)?.label ?? 'Select status');

	let selectedCountryName = $derived(
		data.countriesResult?.data?.data.find(c => c.id === $form.contact.address.country_id)?.name
	);

	export const updateOrganizationFormSchema = v.object({
		name: v.pipe(v.string(), v.minLength(1, 'name is required')),
		website: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'website is required'))])),
		mission_statement: v.nullish(
			v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'mission_statement is required'))])
		),
		adoption_policy: v.nullish(
			v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'adoption_policy is required'))])
		),
		adoption_url: v.nullish(
			v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'adoption_url is required'))])
		),
		status: v.pipe(
			v.nullable(v.picklist(organizationStatusValues)),
			v.picklist(organizationStatusValues, 'status is required')
		),
		work_hour: v.pipe(
			v.nullish(
				v.object({
					monday: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'monday is required'))])),
					tuesday: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'tuesday is required'))])),
					wednesday: v.nullish(
						v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'wednesday is required'))])
					),
					thursday: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'thursday is required'))])),
					friday: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'friday is required'))])),
					saturday: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'saturday is required'))])),
					sunday: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'sunday is required'))]))
				})
			),
			v.transform(x => (Object.values(x ?? {}).some(x => x) ? x : null))
		),
		contact: v.object({
			phone: v.pipe(v.string(), v.minLength(1, 'phone is required')),
			email: v.pipe(v.string(), v.minLength(1, 'email is required'), v.email('email must be a valid email')),
			address: v.object({
				unit_number: v.nullish(
					v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'unit_number is required'))])
				),
				street_number: v.nullish(
					v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'street_number is required'))])
				),
				street_address: v.pipe(v.string(), v.minLength(1, 'street_address is required')),
				city: v.pipe(v.string(), v.minLength(1, 'city is required')),
				region: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'region is required'))])),
				postal_code: v.nullish(
					v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'postal_code is required'))])
				),
				note: v.nullish(v.union([emptyStringToNull, v.pipe(v.string(), v.minLength(1, 'note is required'))])),
				country_id: v.pipe(v.number(), v.minValue(1, 'country_id is required'))
			})
		}),
		socials: v.pipe(
			v.nullish(
				v.array(
					v.object({
						platform: v.pipe(v.string(), v.minLength(1, 'platform is required')),
						url: v.pipe(v.string(), v.url('url must be a valid url'))
					})
				)
			),
			v.transform(x => (x?.length === 0 ? null : x))
		),
		photos: v.pipe(
			v.nullish(
				v.array(
					v.object({
						url: v.pipe(v.string(), v.url('url must be a valid url'))
					})
				)
			),
			v.transform(x => (x?.length === 0 ? null : x))
		),
		videos: v.pipe(
			v.nullish(
				v.array(
					v.object({
						url: v.pipe(v.string(), v.url('url must be a valid url'))
					})
				)
			),
			v.transform(x => (x?.length === 0 ? null : x))
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
	});
	type UpdateOrganizationFormInput = v.InferInput<typeof updateOrganizationFormSchema>;
	type UpdateOrganizationFormOutput = v.InferOutput<typeof updateOrganizationFormSchema>;

	const initialUpdateOrganizationForm: UpdateOrganizationFormInput = {
		name: data.organizationResult?.data?.name ?? '',
		website: data.organizationResult?.data?.website,
		mission_statement: data.organizationResult?.data?.mission_statement,
		adoption_policy: data.organizationResult?.data?.adoption_policy,
		adoption_url: data.organizationResult?.data?.adoption_url,
		status: data.organizationResult?.data?.status ?? null,
		work_hour: data.organizationResult?.data?.work_hour
			? {
					monday: data.organizationResult?.data?.work_hour?.monday,
					tuesday: data.organizationResult?.data?.work_hour?.tuesday,
					wednesday: data.organizationResult?.data?.work_hour?.wednesday,
					thursday: data.organizationResult?.data?.work_hour?.thursday,
					friday: data.organizationResult?.data?.work_hour?.friday,
					saturday: data.organizationResult?.data?.work_hour?.saturday,
					sunday: data.organizationResult?.data?.work_hour?.sunday
				}
			: undefined,
		contact: {
			phone: data.organizationResult?.data?.contact?.phone ?? '',
			email: data.organizationResult?.data?.contact?.email ?? '',
			address: {
				note: data.organizationResult?.data?.contact?.address?.note,
				postal_code: data.organizationResult?.data?.contact?.address?.postal_code,
				region: data.organizationResult?.data?.contact?.address?.region,
				street_number: data.organizationResult?.data?.contact?.address?.street_number,
				unit_number: data.organizationResult?.data?.contact?.address?.unit_number,
				street_address: data.organizationResult?.data?.contact?.address?.street_address ?? '',
				city: data.organizationResult?.data?.contact?.address?.city ?? '',
				country_id: data.organizationResult?.data?.contact?.address?.country.id ?? 198
			}
		},
		socials:
			data.organizationResult?.data?.socials && data.organizationResult?.data?.socials?.length > 0
				? data.organizationResult.data.socials.map(x => ({ id: x.id, platform: x.platform, url: x.url }))
				: undefined
	};

	const supForm = superForm(initialUpdateOrganizationForm, {
		id: 'update_organization',
		validators: valibot(updateOrganizationFormSchema),
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
			if (!data.organizationResult?.data) {
				return;
			}

			const changed = getChangedFormFields(form.data, isTainted) as Partial<UpdateOrganizationFormOutput>;

			try {
				const organizationResult = await fluffly.PATCH('/organizations/{id}', {
					params: {
						path: { id: data.organizationResult.data.id }
					},
					body: {
						data: changed
					},
					bodySerializer(body) {
						const fd = new FormData();
						const data = body?.data as Partial<UpdateOrganizationFormOutput> | undefined;
						if (!data) {
							throw new Error('no data');
						}
						const { photos_files, videos_files, ...dataWithoutFiles } = data;
						const orgDoc = JSON.stringify(dataWithoutFiles);
						fd.append('data', new Blob([orgDoc], { type: 'application/json' }));
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
				updateOrganizationError = organizationResult.error;

				if (organizationResult.data) {
					toast.success('Organization updated');
					uppyPhotos?.clear();
					uppyVideos?.clear();
				}

				if (organizationResult.error) {
					toast.error([organizationResult.error.message, organizationResult.error.reason].filter(Boolean).join(', '));

					if (organizationResult.error.code === 'validation') {
						const details = organizationResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<UpdateOrganizationFormInput> = {};
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

	const { form, enhance, errors, isTainted } = supForm;

	async function deleteOrganizationPhoto(organizationId: number, photoId: number) {
		try {
			const deleteOrganizationPhotoResult = await fluffly.DELETE('/organizations/{id}/photos/{photo_id}', {
				params: {
					path: { id: organizationId, photo_id: photoId }
				}
			});
			if (deleteOrganizationPhotoResult.error) {
				toast.error(
					[deleteOrganizationPhotoResult.error.message, deleteOrganizationPhotoResult.error.reason]
						.filter(x => x)
						.join(',')
				);
				return;
			}
			toast.success('Organization photo deleted');
			invalidate(`data:dashboard-organizations-${organizationId}-update`);
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

	async function deleteOrganizationPhotos(organizationId: number, photoIds: number[]) {
		try {
			const deleteOrganizationPhotosResult = await fluffly.DELETE('/organizations/{id}/photos', {
				params: {
					path: {
						id: organizationId
					}
				},
				body: { ids: photoIds }
			});
			if (deleteOrganizationPhotosResult.error) {
				toast.error(
					[deleteOrganizationPhotosResult.error.message, deleteOrganizationPhotosResult.error.reason]
						.filter(x => x)
						.join(',')
				);
				return;
			}
			toast.success('Organization photos deleted');
			invalidate(`data:dashboard-organizations-${organizationId}-update`);
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

	async function deleteOrganizationVideo(organizationId: number, videoId: number) {
		try {
			const deleteOrganizationVideoResult = await fluffly.DELETE('/organizations/{id}/videos/{video_id}', {
				params: {
					path: { id: organizationId, video_id: videoId }
				}
			});
			if (deleteOrganizationVideoResult.error) {
				toast.error(
					[deleteOrganizationVideoResult.error.message, deleteOrganizationVideoResult.error.reason]
						.filter(x => x)
						.join(',')
				);
				return;
			}
			toast.success('Organization video deleted');
			invalidate(`data:dashboard-organizations-${organizationId}-update`);
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

	async function deleteOrganizationVideos(organizationId: number, videoIds: number[]) {
		try {
			const deleteOrganizationVideosResult = await fluffly.DELETE('/organizations/{id}/videos', {
				params: {
					path: { id: organizationId }
				},
				body: {
					ids: videoIds
				}
			});
			if (deleteOrganizationVideosResult.error) {
				toast.error(
					[deleteOrganizationVideosResult.error.message, deleteOrganizationVideosResult.error.reason]
						.filter(x => x)
						.join(',')
				);
				return;
			}
			toast.success('Organization videos deleted');
			invalidate(`data:dashboard-organizations-${organizationId}-update`);
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

{#if data.organizationResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Update organization</h2>
				<p class="text-muted-foreground">Update organization information</p>
			</div>
		</div>
		<div>
			<form method="POST" use:enhance enctype="multipart/form-data" class="grid gap-4">
				{#if updateOrganizationError}
					<Alert.Root variant="error" icon class="col-span-full">
						<Alert.Title>{updateOrganizationError.message}</Alert.Title>
						{#if updateOrganizationError.reason}
							<Alert.Description>{updateOrganizationError.reason}</Alert.Description>
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
									<Input {...props} bind:value={$form.name} minlength={1} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="website">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Website</Form.Label>
									<Input {...props} bind:value={$form.website} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="mission_statement">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Mission statement</Form.Label>
									<Input {...props} bind:value={$form.mission_statement} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="adoption_policy">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Adoption policy</Form.Label>
									<Input {...props} bind:value={$form.adoption_policy} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="adoption_url">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Adoption URL</Form.Label>
									<Input {...props} bind:value={$form.adoption_url} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
				</section>

				<div class="grid gap-2">
					<Form.Field form={supForm} name="status">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label mandatory>Status</Form.Label>
								<Select.Root
									type="single"
									onValueChange={v => {
										$form.status = (v as OrganizationStatus) || null;
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

				<section class="mt-8 grid gap-4 md:grid-cols-2">
					<div class="col-span-full flex items-center gap-4">
						<h3 class="text-muted-foreground">Work hours information</h3>
						{#if Object.values($form.work_hour ?? {}).some(x => x)}
							<Button variant="outline" onclick={() => ($form.work_hour = null)}>Unset work hours</Button>
						{/if}
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="work_hour.monday">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Monday</Form.Label>
									<Input
										{...props}
										bind:value={
											() => $form.work_hour?.monday ?? '',
											v => {
												$form.work_hour ||= {};
												$form.work_hour.monday = v;
											}
										}
									/>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="work_hour.tuesday">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Tuesday</Form.Label>
									<Input
										{...props}
										bind:value={
											() => $form.work_hour?.tuesday ?? '',
											v => {
												$form.work_hour ||= {};
												$form.work_hour.tuesday = v;
											}
										}
									/>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="work_hour.wednesday">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Wednesday</Form.Label>
									<Input
										{...props}
										bind:value={
											() => $form.work_hour?.wednesday ?? '',
											v => {
												$form.work_hour ||= {};
												$form.work_hour.wednesday = v;
											}
										}
									/>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="work_hour.thursday">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Thursday</Form.Label>
									<Input
										{...props}
										bind:value={
											() => $form.work_hour?.thursday ?? '',
											v => {
												$form.work_hour ||= {};
												$form.work_hour.thursday = v;
											}
										}
									/>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="work_hour.friday">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Friday</Form.Label>
									<Input
										{...props}
										bind:value={
											() => $form.work_hour?.friday ?? '',
											v => {
												$form.work_hour ||= {};
												$form.work_hour.friday = v;
											}
										}
									/>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="work_hour.saturday">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Saturday</Form.Label>
									<Input
										{...props}
										bind:value={
											() => $form.work_hour?.saturday ?? '',
											v => {
												$form.work_hour ||= {};
												$form.work_hour.saturday = v;
											}
										}
									/>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="work_hour.sunday">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Sunday</Form.Label>
									<Input
										{...props}
										bind:value={
											() => $form.work_hour?.sunday ?? '',
											v => {
												$form.work_hour ||= {};
												$form.work_hour.sunday = v;
											}
										}
									/>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
				</section>

				<section class="mt-8 grid gap-4 md:grid-cols-2">
					<h4 class="col-span-full text-muted-foreground">Contact information</h4>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.phone">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label mandatory>Phone</Form.Label>
									<Input {...props} bind:value={$form.contact.phone} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.email">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label mandatory>E-Mail</Form.Label>
									<Input {...props} type="email" bind:value={$form.contact.email} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.address.country_id">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label mandatory>Address country</Form.Label>
									<Select.Root
										type="single"
										onValueChange={v => {
											$form.contact.address.country_id = Number(v);
										}}
										value={String($form.contact.address.country_id)}
									>
										<Select.Trigger {...props}>
											{selectedCountryName || 'Select a country'}
										</Select.Trigger>
										<Select.Content>
											{#each data.countriesResult?.data?.data ?? [] as country (country.id)}
												<Select.Item value={`${country.id}`} label={country.name} />
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
						<Form.Field form={supForm} name="contact.address.unit_number">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Address unit number</Form.Label>
									<Input {...props} bind:value={$form.contact.address.unit_number} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.address.street_number">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Address street number</Form.Label>
									<Input {...props} bind:value={$form.contact.address.street_number} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.address.street_address">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label mandatory>Address street address</Form.Label>
									<Input {...props} bind:value={$form.contact.address.street_address} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.address.city">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label mandatory>Address city</Form.Label>
									<Input {...props} bind:value={$form.contact.address.city} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.address.region">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Address region</Form.Label>
									<Input {...props} bind:value={$form.contact.address.region} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.address.postal_code">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Address postal code</Form.Label>
									<Input {...props} bind:value={$form.contact.address.postal_code} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="contact.address.note">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Address note</Form.Label>
									<Input {...props} bind:value={$form.contact.address.note} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
				</section>

				<section class="grid gap-4 md:grid-cols-2">
					<Item.Root variant="muted" class="col-span-full">
						<Item.Media variant="icon">
							<IconGlobe />
						</Item.Media>
						<Item.Content>
							<Item.Title>Social platforms</Item.Title>
							<Item.Description>Add social platforms: fb, instagram etc</Item.Description>
						</Item.Content>
						<Item.Actions>
							<Button
								size="sm"
								variant="outline"
								onclick={() => {
									$form.socials ||= [];
									$form.socials = [...$form.socials, { platform: '', url: '' }];
								}}
							>
								Add
								<IconPlus />
							</Button>
						</Item.Actions>
					</Item.Root>
					<Item.Group class="col-span-full grid gap-2">
						{#if $form.socials}
							{#each $form.socials as _, i (i)}
								<Item.Root variant="outline">
									<Item.Content>
										<Item.Title class="mb-2">
											Social platform {i + 1}
										</Item.Title>
										<Form.Field form={supForm} name="socials[{i}].platform">
											<Form.Control>
												{#snippet children({ props })}
													<Form.Label mandatory>Name</Form.Label>
													<Input {...props} bind:value={$form.socials![i]!.platform} />
												{/snippet}
											</Form.Control>
											<Form.Description />
											<Form.FieldErrors />
										</Form.Field>
										<Form.Field form={supForm} name="socials[{i}].url">
											<Form.Control>
												{#snippet children({ props })}
													<Form.Label mandatory>URL</Form.Label>
													<Input {...props} bind:value={$form.socials![i]!.url} />
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
												if ($form.socials) {
													$form.socials = $form.socials.filter((_, idx) => idx !== i);
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
					<h5 class="mb-2 font-medium">Organization photos</h5>
					<div id="update-organization-photos"></div>
					{#each $errors.photos_files?._errors ?? [] as e (e)}
						<div class="text-sm font-medium text-destructive">{e}</div>
					{/each}
				</div>

				<div class="col-span-full mt-8">
					<h5 class="mb-2 font-medium">Organization videos</h5>
					<div id="update-organization-videos"></div>
					{#each $errors.videos_files?._errors ?? [] as e (e)}
						<div class="text-sm font-medium text-destructive">{e}</div>
					{/each}
				</div>

				<div class="mt-8">
					<Form.Button class="font-bold">Update organization</Form.Button>
				</div>

				{#if selectedPhotos.length > 0}
					<div class="mt-8 flex gap-4 self-start">
						<Button
							variant="destructive"
							onclick={() => {
								confirmation.openDialog({
									title: `Delete organization photos (${selectedPhotos.length})?`,
									destructive: true,
									async onConfirm() {
										const orgId = data.organizationResult.data?.id;
										if (orgId) {
											await deleteOrganizationPhotos(
												orgId,
												selectedPhotos.map(x => x.id)
											);
											invalidate(`data:dashboard-organizations-${orgId}-update`);
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

				{#if data.organizationResult.data?.photos && data.organizationResult.data.photos.length > 0}
					<div class="mt-8 grid grid-cols-1 text-sm">
						<span class="mb-2 text-muted-foreground">Photos</span>
						<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
							{#each data.organizationResult.data.photos as photo, i (photo.full_url)}
								<div class="group grid h-64 overflow-hidden [grid-template-areas:'stack']">
									<img
										src={photo.full_url}
										alt="organization photo {i + 1}"
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
															title: `Delete organization photo: ${photo.id}?`,
															destructive: true,
															async onConfirm() {
																const orgId = data?.organizationResult?.data?.id;
																if (orgId) {
																	await deleteOrganizationPhoto(orgId, photo.id);
																	invalidate(`data:dashboard-organizations-${orgId}-update`);
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
									title: `Delete organization videos (${selectedVideos.length})?`,
									destructive: true,
									async onConfirm() {
										const orgId = data.organizationResult.data?.id;
										if (orgId) {
											await deleteOrganizationVideos(
												orgId,
												selectedVideos.map(x => x.id)
											);
											invalidate(`data:dashboard-organizations-${orgId}-update`);
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

				{#if data.organizationResult.data?.videos && data.organizationResult.data.videos.length > 0}
					<div class="mt-8 grid grid-cols-1 text-sm">
						<span class="mb-2 text-muted-foreground">Videos</span>
						<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
							{#each data.organizationResult.data.videos as video (video.url)}
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
															title: `Delete organization video: ${video.id}?`,
															destructive: true,
															async onConfirm() {
																const orgId = data?.organizationResult?.data?.id;
																if (orgId) {
																	await deleteOrganizationVideo(orgId, video.id);
																	invalidate(`data:dashboard-organizations-${orgId}-update`);
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
			</form>
		</div>
	</div>
{/if}
