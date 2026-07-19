<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
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

	let createdAnimalType: components['schemas']['AnimalType'] | null = $state(null);
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	export const createAnimalTypeFormSchema = v.object({
		name: v.pipe(v.string(), v.minLength(1, 'name is required'))
	});

	type CreateAnimalTypeFormInput = v.InferInput<typeof createAnimalTypeFormSchema>;
	type CreateAnimalTypeFormOutput = v.InferOutput<typeof createAnimalTypeFormSchema>;

	const initialCreateAnimalTypeForm: CreateAnimalTypeFormInput = {
		name: ''
	};

	const supForm = superForm(initialCreateAnimalTypeForm, {
		id: 'create_animal_type',
		validators: valibot(createAnimalTypeFormSchema),
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
				const animalTypeResult = await fluffly.POST('/animal_types', {
					body: form.data as CreateAnimalTypeFormOutput
				});
				if (animalTypeResult.data) {
					createdAnimalType = animalTypeResult.data;
					toast.success('Animal type created');
					reset();
				}
				if (animalTypeResult.error) {
					if (animalTypeResult.error.code === 'validation') {
						const details = animalTypeResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<CreateAnimalTypeFormInput> = {};
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
				<Form.Button class="font-bold">Create animal type</Form.Button>
			</div>
		</form>
	</div>

	{#if createdAnimalType}
		<Card.Root class="max-w-6xl">
			<Card.Header>
				<Card.Title>View created animal type</Card.Title>
			</Card.Header>

			<Card.Content>
				<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">ID</span>
						<a class="font-medium text-sky-500 underline" href="/dashboard/animal-types/{createdAnimalType.id}">
							{createdAnimalType.id}
						</a>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Name</span>
						<span class="font-medium">{createdAnimalType.name}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Created time</span>
						<time class="font-medium">{fmt.format(new Date(createdAnimalType.created_at))}</time>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Updated time</span>
						<time class="font-medium">{fmt.format(new Date(createdAnimalType.updated_at))}</time>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
</div>
