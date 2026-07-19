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

	let updateAnimalTypeError = $state<components['schemas']['APIError']>();

	export const updateAnimalTypeFormSchema = v.object({
		name: v.pipe(v.string(), v.minLength(1, 'name is required'))
	});

	type UpdateAnimalTypeFormInput = v.InferInput<typeof updateAnimalTypeFormSchema>;
	type UpdateAnimalTypeFormOutput = v.InferOutput<typeof updateAnimalTypeFormSchema>;

	const initialUpdateAnimalTypeForm: UpdateAnimalTypeFormInput = {
		name: data.animalTypeResult?.data?.name ?? ''
	};

	const supForm = superForm(initialUpdateAnimalTypeForm, {
		id: 'update_animal_type',
		validators: valibot(updateAnimalTypeFormSchema),
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
			if (!data.animalTypeResult?.data) {
				return;
			}

			const changed = getChangedFormFields(form.data, isTainted) as Partial<UpdateAnimalTypeFormOutput>;

			try {
				const animalTypeResult = await fluffly.PATCH('/animal_types/{id}', {
					params: {
						path: { id: data.animalTypeResult.data.id }
					},
					body: changed
				});
				updateAnimalTypeError = animalTypeResult.error;

				if (animalTypeResult.data) {
					toast.success('Animal type updated');
				}

				if (animalTypeResult.error) {
					if (animalTypeResult.error.code === 'validation') {
						const details = animalTypeResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<UpdateAnimalTypeFormInput> = {};
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
			<h2 class="text-2xl font-bold tracking-tight">Update animal type</h2>
			<p class="text-muted-foreground">Update animal type information</p>
		</div>
	</div>
	<div class="max-w-lg">
		<form method="POST" use:enhance class="grid gap-4 md:grid-cols-2">
			{#if updateAnimalTypeError}
				<Alert.Root variant="error" icon class="col-span-full">
					<Alert.Title>{updateAnimalTypeError.message}</Alert.Title>
					{#if updateAnimalTypeError.reason}
						<Alert.Description>{updateAnimalTypeError.reason}</Alert.Description>
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
				<Form.Button class="font-bold">Update animal type</Form.Button>
			</div>
		</form>
	</div>
</div>
