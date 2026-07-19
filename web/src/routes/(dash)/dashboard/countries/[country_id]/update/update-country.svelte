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

	let updateCountryError = $state<components['schemas']['APIError']>();

	export const updateContryFormSchema = v.object({
		name: v.pipe(v.string(), v.minLength(1, 'name is required')),
		iso_alpha2: v.pipe(v.string(), v.length(2, 'iso_alpha2 must have 2 chars')),
		iso_alpha3: v.pipe(v.string(), v.length(3, 'iso_alpha3 must have 3 chars')),
		iso_numeric: v.pipe(v.string(), v.length(3, 'iso_numeric must have 3 chars'))
	});

	type UpdateCountryFormInput = v.InferInput<typeof updateContryFormSchema>;
	type UpdateCountryFormOutput = v.InferOutput<typeof updateContryFormSchema>;

	const initialUpdateCountryForm: UpdateCountryFormInput = {
		name: data.countryResult?.data?.name ?? '',
		iso_alpha2: data.countryResult?.data?.iso_alpha2 ?? '',
		iso_alpha3: data.countryResult?.data?.iso_alpha3 ?? '',
		iso_numeric: data.countryResult?.data?.iso_numeric ?? ''
	};

	const supForm = superForm(initialUpdateCountryForm, {
		id: 'update_country',
		validators: valibot(updateContryFormSchema),
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
			if (!data.countryResult?.data) {
				return;
			}

			const changed = getChangedFormFields(form.data, isTainted) as Partial<UpdateCountryFormOutput>;

			try {
				const countryResult = await fluffly.PATCH('/countries/{id}', {
					params: {
						path: { id: data.countryResult.data.id }
					},
					body: changed
				});
				updateCountryError = countryResult.error;

				if (countryResult.data) {
					toast.success('Country updated');
				}

				if (countryResult.error) {
					if (countryResult.error.code === 'validation') {
						const details = countryResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<UpdateCountryFormInput> = {};
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
			<h2 class="text-2xl font-bold tracking-tight">Update country</h2>
			<p class="text-muted-foreground">Update country information</p>
		</div>
	</div>
	<div class="max-w-lg">
		<form method="POST" use:enhance class="grid gap-4 md:grid-cols-2">
			{#if updateCountryError}
				<Alert.Root variant="error" icon class="col-span-full">
					<Alert.Title>{updateCountryError.message}</Alert.Title>
					{#if updateCountryError.reason}
						<Alert.Description>{updateCountryError.reason}</Alert.Description>
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

			<div class="grid gap-2">
				<Form.Field form={supForm} name="iso_alpha2">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>ISO-2</Form.Label>
							<Input {...props} bind:value={$form.iso_alpha2} />
						{/snippet}
					</Form.Control>
					<Form.Description />
					<Form.FieldErrors />
				</Form.Field>
			</div>

			<div class="grid gap-2">
				<Form.Field form={supForm} name="iso_alpha3">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>ISO-3</Form.Label>
							<Input {...props} bind:value={$form.iso_alpha3} />
						{/snippet}
					</Form.Control>
					<Form.Description />
					<Form.FieldErrors />
				</Form.Field>
			</div>

			<div class="grid gap-2">
				<Form.Field form={supForm} name="iso_numeric">
					<Form.Control>
						{#snippet children({ props })}
							<Form.Label>ISO numeric</Form.Label>
							<Input {...props} bind:value={$form.iso_numeric} />
						{/snippet}
					</Form.Control>
					<Form.Description />
					<Form.FieldErrors />
				</Form.Field>
			</div>

			<div>
				<Form.Button class="font-bold">Update country</Form.Button>
			</div>
		</form>
	</div>
</div>
