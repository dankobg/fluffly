<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Form from '$lib/components/ui/form';
	import { fluffly } from '$lib/fluffly/client';
	import type { components } from '$lib/gen/fluffly_openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Item from '$lib/components/ui/item/index.js';
	import IconPlus from '@lucide/svelte/icons/plus';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import IconGlobe from '@lucide/svelte/icons/globe';
	import IconImage from '@lucide/svelte/icons/image';
	import IconVideo from '@lucide/svelte/icons/video';
	import { onDestroy, onMount } from 'svelte';
	import { Uppy } from '@uppy/core';
	import UppyDashboard from '@uppy/dashboard';
	import UppyImageEditor from '@uppy/image-editor';
	import { emptyStringToUndefined } from '$lib/validation/common';

	let {
		data,
		// eslint-disable-next-line no-useless-assignment
		capture = $bindable(),
		// eslint-disable-next-line no-useless-assignment
		restore = $bindable()
	}: PageProps & { capture: Snapshot['capture']; restore: Snapshot['restore'] } = $props();

	let uppyPhotos: Uppy | null = null;
	let uppyVideos: Uppy | null = null;

	onMount(() => {
		uppyPhotos = new Uppy({
			id: 'apply-for-organization-photos',
			restrictions: {
				maxNumberOfFiles: 20,
				allowedFileTypes: ['image/jpg', 'image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/avif']
			},
			autoProceed: false
		});
		uppyPhotos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#apply-for-organization-photos',
			note: 'Images only, up to 20 files (jpg, png, gif)',
			width: '100%',
			height: '300px',
			hideUploadButton: true,
			locale: {
				strings: {
					dropPasteBoth: 'Drop image files here, %{browseFiles} or %{browseFolders}',
					dropPasteFiles: 'Drop image files here or %{browseFiles}',
					dropPasteFolders: 'Drop image files here or %{browseFolders}',
					dropPasteImportBoth: 'Drop image files here, %{browseFiles}, %{browseFolders} or import from:',
					dropPasteImportFiles: 'Drop image files here, %{browseFiles} or import from:',
					dropPasteImportFolders: 'Drop image files here, %{browseFolders} or import from:'
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
			id: 'apply-for-organization-videos',
			restrictions: {
				maxNumberOfFiles: 5,
				allowedFileTypes: ['video/mp4', 'video/ogg', 'video/mpeg', 'video/webm']
			},
			autoProceed: false
		});
		uppyVideos.use(UppyDashboard, {
			inline: true,
			theme: 'auto',
			target: '#apply-for-organization-videos',
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
		uppyPhotos?.destroy();
		uppyVideos?.destroy();
	});

	let selectedCountryName = $derived(
		data.countriesResult?.data?.data.find(c => c.id === $form.contact.address.country_id)?.name
	);

	export const applyForOrganizationFormSchema = v.pipe(
		v.object({
			name: v.pipe(v.string(), v.minLength(1, 'name is required')),
			website: v.optional(v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'website is required'))])),
			mission_statement: v.optional(
				v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'mission_statement is required'))])
			),
			adoption_policy: v.optional(
				v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'adoption_policy is required'))])
			),
			adoption_url: v.optional(
				v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'adoption_url is required'))])
			),
			work_hour: v.pipe(
				v.optional(
					v.object({
						monday: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'monday is required'))])
						),
						tuesday: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'tuesday is required'))])
						),
						wednesday: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'wednesday is required'))])
						),
						thursday: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'thursday is required'))])
						),
						friday: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'friday is required'))])
						),
						saturday: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'saturday is required'))])
						),
						sunday: v.optional(
							v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'sunday is required'))])
						)
					})
				),
				v.transform(x => (Object.values(x ?? {}).some(x => x) ? x : undefined))
			),
			contact: v.object({
				phone: v.pipe(v.string(), v.minLength(1, 'phone is required')),
				email: v.pipe(v.string(), v.minLength(1, 'email is required'), v.email('email must be a valid email')),
				address: v.object({
					unit_number: v.optional(
						v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'unit_number is required'))])
					),
					street_number: v.optional(
						v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'street_number is required'))])
					),
					street_address: v.pipe(v.string(), v.minLength(1, 'street_address is required')),
					city: v.pipe(v.string(), v.minLength(1, 'city is required')),
					region: v.optional(
						v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'region is required'))])
					),
					postal_code: v.optional(
						v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'postal_code is required'))])
					),
					note: v.optional(v.union([emptyStringToUndefined, v.pipe(v.string(), v.minLength(1, 'note is required'))])),
					country_id: v.pipe(v.number(), v.minValue(1, 'country_id is required'))
				})
			}),
			socials: v.optional(
				v.array(
					v.object({
						platform: v.pipe(v.string(), v.minLength(1, 'platform is required')),
						url: v.pipe(v.string(), v.url('url must be a valid url'))
					})
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
				[['photos_files'], ['videos_files']],
				input => {
					const files: File[] = [...(input.photos_files ?? []), ...(input.videos_files ?? [])];
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
				[['photos_files'], ['videos_files']],
				input => {
					const files: File[] = [...(input.photos_files ?? []), ...(input.videos_files ?? [])];
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

	type ApplyForOrganizationFormInput = v.InferInput<typeof applyForOrganizationFormSchema>;
	type ApplyForOrganizationFormOutput = v.InferOutput<typeof applyForOrganizationFormSchema>;

	const initialApplyForOrganizationForm: ApplyForOrganizationFormInput = {
		name: '',
		contact: {
			phone: '',
			email: '',
			address: {
				street_address: '',
				city: '',
				country_id: 198
			}
		}
	};

	const supForm = superForm(initialApplyForOrganizationForm, {
		id: 'apply_for_organization',
		validators: valibot(applyForOrganizationFormSchema),
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
				const organizationResult = await fluffly.POST('/organizations/application', {
					body: {
						data: form.data
					},
					bodySerializer(body) {
						const fd = new FormData();
						const data = body?.data as ApplyForOrganizationFormOutput | undefined;
						if (!data) {
							throw new Error('no data');
						}
						const dataWithoutFiles = {
							...data,
							photos_files: undefined,
							videos_files: undefined
						};
						const orgDoc = JSON.stringify(dataWithoutFiles);
						fd.append('data', new Blob([orgDoc], { type: 'application/json' }));
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
				if (organizationResult.data) {
					toast.success('Applied successfuly');
					reset();
					sessionStorage.removeItem('sveltekit:snapshot');
					uppyPhotos?.clear();
					uppyVideos?.clear();
				}
				if (organizationResult.error) {
					toast.error([organizationResult.error.message, organizationResult.error.reason].filter(Boolean).join(', '));

					if (organizationResult.error.code === 'validation') {
						const details = organizationResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<ApplyForOrganizationFormInput> = {};
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

	const { form, enhance, errors, reset } = supForm;
	// eslint-disable-next-line no-useless-assignment
	capture = supForm.capture;
	// eslint-disable-next-line no-useless-assignment
	restore = supForm.restore;
</script>

<div class="flex h-full flex-1 flex-col space-y-8 p-8">
	<div class="flex items-center justify-between space-y-2">
		<div>
			<h1 class="text-2xl font-bold tracking-tight">Apply for organization</h1>
			<!-- <p class="text-muted-foreground">Apply for organization</p> -->
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

			<section class="mt-8 grid gap-4 md:grid-cols-2">
				<h3 class="col-span-full text-muted-foreground">Work hours information</h3>
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

			<div class="col-span-full">
				<h5 class="mb-2 font-medium">Organization photos</h5>
				<div id="apply-for-organization-photos"></div>
				{#each $errors.photos_files?._errors ?? [] as e (e)}
					<div class="text-sm font-medium text-destructive">{e}</div>
				{/each}
			</div>

			<div class="col-span-full">
				<h5 class="mb-2 font-medium">Organization videos</h5>
				<div id="apply-for-organization-videos"></div>
				{#each $errors.videos_files?._errors ?? [] as e (e)}
					<div class="text-sm font-medium text-destructive">{e}</div>
				{/each}
			</div>

			<div class="mt-4 ml-auto">
				<Form.Button class="font-bold">apply for organization</Form.Button>
			</div>
		</form>
	</div>
</div>
