<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { fluffly } from '$lib/fluffly/client';
	import type { components } from '$lib/gen/fluffly_openapi';

	let {
		data,
		// eslint-disable-next-line no-useless-assignment
		capture = $bindable(),
		// eslint-disable-next-line no-useless-assignment
		restore = $bindable()
	}: PageProps & { capture: Snapshot['capture']; restore: Snapshot['restore'] } = $props();

	let createdAnimalSpecie: components['schemas']['AnimalSpecie'] | null = $state(null);
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	let selectedAnimalTypeName = $derived(
		data.animalTypesResult?.data?.data.find(x => x.id === $form.animal_type_id)?.name
	);

	export const createAnimalSpeciesFormSchema = v.object({
		animal_type_id: v.pipe(
			v.nullable(v.number()),
			v.number('animal_type_id is required'),
			v.minValue(1, 'animal_type_id is required')
		),
		name: v.pipe(v.string(), v.minLength(1, 'name is required'))
	});

	type CreateAnimalSpeciesFormInput = v.InferInput<typeof createAnimalSpeciesFormSchema>;
	type CreateAnimalSpeciesFormOutput = v.InferOutput<typeof createAnimalSpeciesFormSchema>;

	const initialCreateAnimalSpeciesForm: CreateAnimalSpeciesFormInput = {
		animal_type_id: null,
		name: ''
	};

	const supForm = superForm(initialCreateAnimalSpeciesForm, {
		id: 'create_animal_species',
		validators: valibot(createAnimalSpeciesFormSchema),
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
				const animalSpeciesResult = await fluffly.POST('/animal_species', {
					body: form.data as CreateAnimalSpeciesFormOutput
				});
				if (animalSpeciesResult.data) {
					createdAnimalSpecie = animalSpeciesResult.data;
					toast.success('Animal species created');
					reset();
				}
				if (animalSpeciesResult.error) {
					if (animalSpeciesResult.error.code === 'validation') {
						const details = animalSpeciesResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<CreateAnimalSpeciesFormInput> = {};
						for (const detail of details) {
							if (detail.in === 'body') {
								const path = detail.pointer.replaceAll('/', '.').slice(1);
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
			<h2 class="text-2xl font-bold tracking-tight">Create animal species</h2>
			<!-- <p class="text-muted-foreground">Create new animal species</p> -->
		</div>
	</div>
	<div class="max-w-lg">
		<form method="POST" use:enhance class="grid gap-4 md:grid-cols-2">
			<div class="grid gap-2">
				<Form.Field form={supForm} name="animal_type_id">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label mandatory>Animal type</Form.Label>
							<Select.Root
								type="single"
								onValueChange={v => {
									$form.animal_type_id = Number(v);
								}}
								value={String($form.animal_type_id)}
							>
								<Select.Trigger {...props}>
									{selectedAnimalTypeName || 'Select animal type'}
								</Select.Trigger>
								<Select.Content>
									{#each data.animalTypesResult?.data?.data ?? [] as animalType (animalType.name)}
										<Select.Item value={`${animalType.id}`} label={animalType.name} />
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
				<Form.Field form={supForm} name="name">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>Name</Form.Label>
							<Input {...props} bind:value={$form.name} />
						{/snippet}
					</Form.Control>
					<Form.Description />
					<Form.FieldErrors />
				</Form.Field>
			</div>

			<div>
				<Form.Button class="font-bold">Create animal species</Form.Button>
			</div>
		</form>
	</div>

	{#if createdAnimalSpecie}
		<Card.Root class="max-w-6xl">
			<Card.Header>
				<Card.Title>View created animal species</Card.Title>
			</Card.Header>

			<Card.Content>
				<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">ID</span>
						<a class="font-medium text-sky-500 underline" href="/dashboard/animal-species/{createdAnimalSpecie.id}">
							{createdAnimalSpecie.id}
						</a>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Animal type ID</span>
						<span class="font-medium">{createdAnimalSpecie.animal_type_id}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Name</span>
						<span class="font-medium">{createdAnimalSpecie.name}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Created time</span>
						<time class="font-medium">{fmt.format(new Date(createdAnimalSpecie.created_at))}</time>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Updated time</span>
						<time class="font-medium">{fmt.format(new Date(createdAnimalSpecie.updated_at))}</time>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
</div>
