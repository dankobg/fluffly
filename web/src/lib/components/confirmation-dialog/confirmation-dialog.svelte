<script lang="ts" module>
	type T = unknown;
</script>

<script lang="ts" generics="T">
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index';
	import { buttonVariants } from '../ui/button';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
</script>

<AlertDialog.Root bind:open={confirmation.open}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>{confirmation.title}</AlertDialog.Title>
			<AlertDialog.Description>
				{@render confirmation?.description?.()}
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={confirmation.onCancel}>{confirmation.cancelText}</AlertDialog.Cancel>
			{#if confirmation.destructive}
				<AlertDialog.Action class={buttonVariants({ variant: 'destructive' })} onclick={confirmation.onConfirm}>
					{confirmation.confirmText}
				</AlertDialog.Action>
			{:else}
				<AlertDialog.Action onclick={confirmation.onConfirm}>
					{confirmation.confirmText}
				</AlertDialog.Action>
			{/if}
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
