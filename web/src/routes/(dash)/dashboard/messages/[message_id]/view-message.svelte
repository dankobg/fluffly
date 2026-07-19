<script lang="ts">
	import type { PageProps } from './$types';
	import * as Table from '$lib/components/ui/table/index';
	import { statusIcons, templateTypeIcons, typeIcons } from '../messages-data-grid/data';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import {
		CourierMessageStatus,
		CourierMessageTemplateType,
		CourierMessageType,
		MessageDispatchStatus
	} from '$lib/gen/fluffly_openapi';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	let StatusIcon = $derived(data.messageResult?.data?.status && statusIcons.get(data.messageResult?.data.status));
	let statusIconClasses = $derived.by(() => {
		switch (data.messageResult?.data?.status as CourierMessageStatus) {
			case CourierMessageStatus.abandoned:
				return 'text-red-400';
			case CourierMessageStatus.processing:
				return 'text-purple-400';
			case CourierMessageStatus.queued:
				return 'text-yellow-400';
			case CourierMessageStatus.sent:
				return 'text-green-400';
			default:
				return '';
		}
	});

	let TypeIcon = $derived(data.messageResult?.data?.type && typeIcons.get(data.messageResult?.data.type));
	let typeIconClasses = $derived.by(() => {
		switch (data.messageResult?.data?.type as CourierMessageType) {
			case CourierMessageType.email:
				return 'text-blue-400';
			case CourierMessageType.phone:
				return 'text-purple-400';
			default:
				return '';
		}
	});

	let TemplateTypeIcon = $derived(
		data.messageResult?.data?.template_type && templateTypeIcons.get(data.messageResult?.data.template_type)
	);
	let templateTypeIconClasses = $derived.by(() => {
		switch (data.messageResult?.data?.template_type as CourierMessageTemplateType) {
			case CourierMessageTemplateType.recovery_valid:
			case CourierMessageTemplateType.recovery_code_valid:
			case CourierMessageTemplateType.verification_valid:
			case CourierMessageTemplateType.verification_code_valid:
			case CourierMessageTemplateType.login_code_valid:
			case CourierMessageTemplateType.registration_code_valid:
				return 'text-green-400';
			case CourierMessageTemplateType.verification_code_invalid:
			case CourierMessageTemplateType.verification_invalid:
			case CourierMessageTemplateType.recovery_code_invalid:
			case CourierMessageTemplateType.recovery_invalid:
				return 'text-red-400';
			case CourierMessageTemplateType.stub:
				return 'text-purple-400';
			default:
				return '';
		}
	});
</script>

{#if data.messageResult?.data}
	<h1 class="mb-6 text-2xl font-bold">Courier Message</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.messageResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Recipient</span>
			<span class="font-medium">{data.messageResult?.data.recipient}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Send count</span>
			<span class="font-medium">{data.messageResult?.data.send_count}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Status</span>
			<span class="flex gap-2 font-medium"
				>{data.messageResult?.data.status} <StatusIcon class={statusIconClasses} /></span
			>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Channel</span>
			<span class="font-medium">{data.messageResult?.data.channel}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Type</span>
			<span class="flex gap-2 font-medium">{data.messageResult?.data.type} <TypeIcon class={typeIconClasses} /></span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Template type</span>
			<span class="flex gap-2 font-medium">
				{data.messageResult?.data.template_type}
				<TemplateTypeIcon class={templateTypeIconClasses} />
			</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Subject</span>
			<span class="font-medium">{data.messageResult?.data.subject}</span>
		</div>
		<div class="col-span-1 flex flex-col sm:col-span-2">
			<span class="text-muted-foreground">Body</span>
			<span class="font-medium">{data.messageResult?.data.body}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.messageResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.messageResult?.data.updated_at))}</time>
		</div>
	</div>

	{#if data.messageResult?.data.dispatches && data.messageResult?.data.dispatches.length > 0}
		<p class="mt-8 text-lg">Message dispatches</p>
		<Table.Root>
			<Table.Caption>A list of message dispatches</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>ID</Table.Head>
					<Table.Head>Message ID</Table.Head>
					<Table.Head>Status</Table.Head>
					<Table.Head>Created time</Table.Head>
					<Table.Head>Update time</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.messageResult?.data.dispatches as dispatch (dispatch)}
					<Table.Row>
						<Table.Cell class="font-medium">{dispatch.id}</Table.Cell>
						<Table.Cell>{dispatch.message_id}</Table.Cell>
						<Table.Cell>
							<div class="flex gap-2">
								{dispatch.status}
								{#if dispatch.status === MessageDispatchStatus.success}
									<IconCheck class="text-green-400" />
								{/if}
								{#if dispatch.status === MessageDispatchStatus.failed}
									<IconX class="text-red-400" />
								{/if}
							</div>
						</Table.Cell>
						<Table.Cell>{fmt.format(new Date(dispatch.created_at))}</Table.Cell>
						<Table.Cell>{fmt.format(new Date(dispatch.updated_at))}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
{/if}
