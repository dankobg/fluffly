import type { Column, RowData } from '@tanstack/svelte-table';
import type { AppFeatures } from './app-table';

export function getCommonPinningStyles<T extends RowData>(column: Column<AppFeatures, T>, isSelected = false): string {
	const isPinned = column.getIsPinned();
	const isLastLeftPinnedColumn = isPinned === 'start' && column.getIsLastColumn('start');
	const isFirstRightPinnedColumn = isPinned === 'end' && column.getIsFirstColumn('end');

	const styles: Array<string> = [];
	if (isLastLeftPinnedColumn) {
		styles.push('box-shadow: -4px 0 4px -4px var(--border) inset');
	} else if (isFirstRightPinnedColumn) {
		styles.push('box-shadow: 4px 0 4px -4px var(--border) inset');
	}
	if (isPinned === 'start') {
		styles.push(`inset-inline-start: ${column.getStart('start')}px`);
	}
	if (isPinned === 'end') {
		styles.push(`inset-inline-end: ${column.getAfter('end')}px`);
	}
	styles.push(`position: ${isPinned ? 'sticky' : 'relative'}`);
	if (isSelected) {
		styles.push('background: var(--muted)');
	} else if (isPinned) {
		styles.push('background: var(--background)');
	}
	styles.push(`z-index: ${isPinned ? 1 : 0}`);
	return styles.join('; ');
}
