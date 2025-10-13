<script lang="ts">
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from '$lib/components/ui/sonner/index';
	import '../app.css';
	// import ConfirmationDialog from '$lib/components/confirmation-dialog/confirmation-dialog.svelte'; // @TODO: SOME BUG SSR SHIT HAPPENED

	let { children } = $props();
</script>

<!-- <ConfirmationDialog /> -->
<ModeWatcher />
<Toaster position="top-right" richColors />
{@render children()}
