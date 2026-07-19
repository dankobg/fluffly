<script lang="ts">
	import type { PageProps } from './$types';
	import {
		type UpdateVerificationFlowWithCodeMethod,
		VerificationFlowState,
		instanceOfVerificationFlow,
		ResponseError,
		isGenericErrorResponse,
		FetchError,
		RequiredError
	} from '@ory/client-fetch';
	import { kratos } from '$lib/kratos/client';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Form from '$lib/components/ui/form';
	import * as InputOTP from '$lib/components/ui/input-otp/index';
	import { REGEXP_ONLY_DIGITS } from 'bits-ui';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { config } from '$lib/kratos/config';
	import {
		isErrorIdSecurityCsrfViolation,
		isErrorIdSecurityIdentityMismatch,
		isErrorIdSelfServiceFlowExpired,
		isErrorIdSessionAlreadyAvailable
	} from '$lib/kratos/helpers';

	let { data }: PageProps = $props();

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		data = { ...data, flow: null, csrf: data?.csrf ?? '' };
		if (browser) {
			goto(redirectUrl);
		}
		return;
	}

	const verificationFormSchema = v.object({
		csrf_token: v.pipe(v.string(), v.minLength(1, 'csrf_token is required')),
		method: v.literal('code'),
		code: v.pipe(v.string(), v.length(6, 'Code must have exactly 6 characters'))
	});

	type VerificationFormInput = v.InferInput<typeof verificationFormSchema>;

	const initialVerificationForm: VerificationFormInput = {
		code: '',
		method: 'code',
		csrf_token: data.csrf ?? ''
	};

	const supForm = superForm(initialVerificationForm, {
		id: 'auth_verification',
		validators: valibot(verificationFormSchema),
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

			const url = data.flow?.ui.action;
			const body = form.data as UpdateVerificationFlowWithCodeMethod & { method: 'code' };

			if (url) {
				try {
					const verificationFlow = await kratos.updateVerificationFlow({
						flow: data.flow?.id ?? '',
						updateVerificationFlowBody: body
					});
					if (verificationFlow.state === VerificationFlowState.PassedChallenge) {
						for (const node of verificationFlow.ui.nodes) {
							if (node.group === 'code') {
								if (node.attributes.node_type === 'a') {
									window.sessionStorage.setItem(
										'fluffly_email_verified',
										'Your E-Mail has been verified! You can now log in'
									);
									window.location.href = node.attributes.href;
									return;
								}
							}
						}
					}
				} catch (error: unknown) {
					if (!error || typeof error !== 'object') {
						return;
					}
					if (error instanceof ResponseError) {
						const err = await error.response.json();
						switch (error.response.status) {
							case 400: {
								if (instanceOfVerificationFlow(err)) {
									data = { ...data, flow: err, csrf: data.csrf ?? '' };
									const nodes = err.ui.nodes ?? [];
									const fieldErrors: ValidationErrors<VerificationFormInput> = {};
									for (const node of nodes) {
										const errMsgs: string[] = [];
										if (node.attributes.node_type === 'input') {
											for (const msg of node?.messages ?? []) {
												errMsgs.push(msg.text);
												const fieldName = node.attributes.name;
												set(fieldErrors, fieldName, errMsgs);
											}
										}
									}
									errors.set(fieldErrors);
								}
								break;
							}
							case 410: {
								if (isGenericErrorResponse(err)) {
									if (isErrorIdSessionAlreadyAvailable(err.error?.id)) {
										goto('/');
									} else if (isErrorIdSelfServiceFlowExpired(err.error?.id)) {
										if (browser) {
											goto(`${config.routes.verification.path}?return_to=${window.location.href}`);
										}
									} else if (isErrorIdSecurityCsrfViolation(err.error?.id)) {
										handleFlowErrAction(config.routes.verification.path, err.error.message);
									} else if (isErrorIdSecurityIdentityMismatch(err.error?.id)) {
										goto('/');
									}
								}
								break;
							}
							case 500:
								console.error('unexpected server error');
								break;
							default:
								break;
						}
						return;
					}
					if (error instanceof FetchError) {
						console.error('fetch error: ', error.cause);
						return;
					}
					if (error instanceof RequiredError) {
						console.error('required error: ', error.field);
						return;
					}
					console.error('unexpected error');
				}
			}
		}
	});

	const { form, enhance, errors } = supForm;
</script>

<section class="grid h-[calc(100vh-4rem)] place-content-center gap-4">
	<a href="/" class="justify-self-center">
		<img class="mr-2 h-8 w-8" src="/images/logo.svg" alt="logo" />
	</a>

	<Card.Root class="mx-auto max-w-md">
		<Card.Header>
			<Card.Title class="text-center text-2xl">Verify account</Card.Title>
			<Card.Description>Verify your account</Card.Description>
		</Card.Header>

		<Card.Content>
			<div class="grid gap-4">
				<form method="POST" use:enhance class="grid gap-4">
					{#each data?.flow?.ui?.messages ?? [] as msg (msg.id)}
						<Alert.Root variant={msg.type === '11184809' ? 'info' : msg.type} icon>
							<Alert.Title>{msg.type === 'error' ? 'Unable to verify account' : ''}</Alert.Title>
							<Alert.Description>{msg.text}</Alert.Description>
						</Alert.Root>
					{/each}

					<div class="grid gap-2">
						<Form.Field form={supForm} name="code">
							<Form.Control>
								{#snippet children({ props })}
									<InputOTP.Root
										bind:value={$form.code}
										maxlength={6}
										inputmode="numeric"
										pattern={REGEXP_ONLY_DIGITS}
										class="flex w-full"
										pasteTransformer={text => text.trim()}
										{...props}
									>
										{#snippet children({ cells })}
											<InputOTP.Group class="flex-[5] justify-end">
												{#each cells.slice(0, 3) as cell}
													<InputOTP.Slot {cell} class="h-14 w-full" />
												{/each}
											</InputOTP.Group>
											<InputOTP.Separator class="flex flex-[1] justify-center" />
											<InputOTP.Group class="flex-[5] justify-start">
												{#each cells.slice(3, 6) as cell}
													<InputOTP.Slot {cell} class="h-14 w-full" />
												{/each}
											</InputOTP.Group>
										{/snippet}
									</InputOTP.Root>
								{/snippet}
							</Form.Control>
							<Form.Description>Please enter the one-time password sent to your E-Mail</Form.Description>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<Form.Button class="w-full font-bold">Verify account</Form.Button>
				</form>
			</div>
		</Card.Content>
	</Card.Root>
</section>
