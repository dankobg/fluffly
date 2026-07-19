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

	let createdCountry: components['schemas']['Country'] | null = $state(null);
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	export const createContryFormSchema = v.object({
		name: v.pipe(v.string(), v.minLength(1, 'name is required')),
		iso_alpha2: v.pipe(v.string(), v.length(2, 'iso_alpha2 must have 2 chars')),
		iso_alpha3: v.pipe(v.string(), v.length(3, 'iso_alpha3 must have 3 chars')),
		iso_numeric: v.pipe(v.string(), v.length(3, 'iso_numeric must have 3 chars'))
	});

	type CreateCountryFormInput = v.InferInput<typeof createContryFormSchema>;
	type CreateCountryFormOutput = v.InferOutput<typeof createContryFormSchema>;

	const initialCreateCountryForm: CreateCountryFormInput = {
		name: '',
		iso_alpha2: '',
		iso_alpha3: '',
		iso_numeric: ''
	};

	const supForm = superForm(initialCreateCountryForm, {
		id: 'create_country',
		validators: valibot(createContryFormSchema),
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
				const countryResult = await fluffly.POST('/countries', {
					body: form.data as CreateCountryFormOutput
				});
				if (countryResult.data) {
					createdCountry = countryResult.data;
					toast.success('Country created');
					reset();
				}
				if (countryResult.error) {
					if (countryResult.error.code === 'validation') {
						const details = countryResult.error.details as Array<components['schemas']['ValidationDetail']>;
						const fieldErrors: ValidationErrors<CreateCountryFormInput> = {};
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
			<h2 class="text-2xl font-bold tracking-tight">Create a country</h2>
			<!-- <p class="text-muted-foreground">Create new country</p> -->
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
				<Form.Button class="font-bold">Create a country</Form.Button>
			</div>
		</form>
	</div>

	{#if createdCountry}
		<Card.Root class="max-w-6xl">
			<Card.Header>
				<Card.Title>View created country</Card.Title>
			</Card.Header>

			<Card.Content>
				<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">ID</span>
						<a class="font-medium text-sky-500 underline" href="/dashboard/countries/{createdCountry.id}">
							{createdCountry.id}
						</a>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Name</span>
						<span class="font-medium">{createdCountry.name}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Iso alpha 2</span>
						<span class="font-medium">{createdCountry.iso_alpha2}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Iso alpha 3</span>
						<span class="font-medium">{createdCountry.iso_alpha3}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Iso numeric</span>
						<span class="font-medium">{createdCountry.iso_numeric}</span>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Created time</span>
						<time class="font-medium">{fmt.format(new Date(createdCountry.created_at))}</time>
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Updated time</span>
						<time class="font-medium">{fmt.format(new Date(createdCountry.updated_at))}</time>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
</div>
