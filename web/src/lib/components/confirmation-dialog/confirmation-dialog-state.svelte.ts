import type { Snippet } from 'svelte';

type ConfirmationDialogOptions = {
	title?: string;
	cancelText?: string;
	confirmText?: string;
	description?: Snippet;
	destructive?: boolean;
	onCancel?: () => void;
	onConfirm?: () => void;
};

export class Confirmation {
	open: boolean = $state(false);
	title: string = $state('Are you sure?');
	cancelText: string = $state('Cancel');
	confirmText: string = $state('Continue');
	description?: Snippet;
	destructive?: boolean = $state(false);
	onCancel?: () => void = () => {
		this.open = false;
	};
	onConfirm?: () => void;

	openDialog(obj: ConfirmationDialogOptions): void {
		this.open = true;
		this.title = obj.title ?? 'Are you sure?';
		this.cancelText = obj.cancelText ?? 'Cancel';
		this.confirmText = obj.confirmText ?? 'Continue';
		if (obj.description) {
			this.description = obj.description;
		}
		if (obj.destructive) {
			this.destructive = obj.destructive;
		}
		if (obj.onCancel) {
			this.onCancel = obj.onCancel;
		}
		if (obj.onConfirm) {
			this.onConfirm = obj.onConfirm;
		}
	}

	closeDialog(): void {
		this.open = false;
	}
}

export const confirmation = new Confirmation();
