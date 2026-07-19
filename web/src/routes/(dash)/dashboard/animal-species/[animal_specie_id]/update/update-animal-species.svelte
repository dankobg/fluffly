<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Form from '$lib/components/ui/form';
	import { fluffly } from '$lib/fluffly/client';
	import type { components } from '$lib/gen/fluffly_openapi';
	import * as Alert from '$lib/components/ui/alert';
	import { getChangedFormFields } from '$lib/utils';

	let {
		data,
		// eslint-disable-next-line no-useless-assignment
		capture = $bindable(),
		// eslint-disable-next-line no-useless-assignment
		restore = $bindable()
	}: PageProps & { capture: Snapshot['capture']; restore: Snapshot['restore'] } = $props();

	let updateAnimalSpeciesError = $state<components['schemas']['APIError']>();

	export const updateAnimalSpeciesFormSchema = v.object({
		animal_type_id: v.pipe(
			v.nullable(v.number()),
			v.number('animal_type_id is required'),
			v.minValue(1, 'animal_type_id is required')
		),
		name: v.pipe(v.string(), v.minLength(1, 'name is required'))
	});

	type UpdateAnimalSpeciesFormInput = v.InferInput<typeof updateAnimalSpeciesFormSchema>;
	type UpdateAnimalSpeciesFormOutput = v.InferOutput<typeof updateAnimalSpeciesFormSchema>;

	const initialUpdateAnimalSpeciesForm: UpdateAnimalSpeciesFormInput = {
		animal_type_id: data.animalSpeciesResult?.data?.animal_type_id ?? null,
		name: data.animalSpeciesResult?.data?.name ?? ''
	};

	const supForm = superForm(initialUpdateAnimalSpeciesForm, {
		id: 'update_animal_species',
		validators: valibot(updateAnimalSpeciesFormSchema),
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
			if (!data.animalSpeciesResult?.data) {
				return;
			}

			const changed = getChangedFormFields(form.data, isTainted) as Partial<UpdateAnimalSpeciesFormOutput>;

			try {
				const animalSpeciesResult = await fluffly.PATCH('/animal_species/{id}', {
					params: {
						path: { id: data.animalSpeciesResult.data.id }
					},
					body: changed
				});
				updateAnimalSpeciesError = animalSpeciesResult.error;

				if (animalSpeciesResult.data) {
					toast.success('Animal species updated');
				}

				if (animalSpeciesResult.error) {
					if (animalSpeciesResult.error.code === 'validation') {
						const details = animalSpeciesResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<UpdateAnimalSpeciesFormInput> = {};
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

	const { form, enhance, errors, isTainted } = supForm;
	// eslint-disable-next-line no-useless-assignment
	capture = supForm.capture;
	// eslint-disable-next-line no-useless-assignment
	restore = supForm.restore;
</script>

<div class="flex h-full flex-1 flex-col space-y-8 p-8">
	<div class="flex items-center justify-between space-y-2">
		<div>
			<h2 class="text-2xl font-bold tracking-tight">Update animal species</h2>
			<p class="text-muted-foreground">Update animal species information</p>
		</div>
	</div>
	<div class="max-w-lg">
		<form method="POST" use:enhance class="grid gap-4 md:grid-cols-2">
			{#if updateAnimalSpeciesError}
				<Alert.Root variant="error" icon class="col-span-full">
					<Alert.Title>{updateAnimalSpeciesError.message}</Alert.Title>
					{#if updateAnimalSpeciesError.reason}
						<Alert.Description>{updateAnimalSpeciesError.reason}</Alert.Description>
					{/if}
				</Alert.Root>
			{/if}

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
				<Form.Button class="font-bold">Update animal species</Form.Button>
			</div>
		</form>
	</div>
</div>
