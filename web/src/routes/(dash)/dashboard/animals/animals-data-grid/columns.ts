import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import { renderComponent, renderSnippet } from '@tanstack/svelte-table';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';
import { createRawSnippet } from 'svelte';

const columnHelper = createAppColumnHelper<components['schemas']['Animal']>();

export const columns = columnHelper.columns([
	columnHelper.display({
		id: 'select',
		header: ({ header }) => renderComponent(header.SelectHeader),
		cell: ({ cell }) => renderComponent(cell.SelectCell),
		enableSorting: false,
		enableHiding: false,
		enableResizing: false,
		maxSize: 40
	}),
	columnHelper.accessor('id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'ID'
			});
		},
		cell: ({ cell, row }) => {
			return renderComponent(cell.DataTableCellId, {
				value: row.original.id,
				href: `/dashboard/animals/${row.original.id}`
			});
		}
	}),
	columnHelper.accessor('name', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Name'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('type', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Animal type'
			});
		},
		cell: ({ row }) => {
			const typeSnippet = createRawSnippet<[{ type: string | undefined }]>(getType => {
				const { type } = getType();
				return {
					render: () => `<div>${type ?? ''}</div>`
				};
			});
			return renderSnippet(typeSnippet, {
				type: row.original.type.name
			});
		},
		enableSorting: false
	}),
	columnHelper.accessor('specie', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Animal specie'
			});
		},
		cell: ({ row }) => {
			const speciesSnippet = createRawSnippet<[{ species: string | undefined }]>(getSpecies => {
				const { species } = getSpecies();
				return {
					render: () => `<div>${species}</div>`
				};
			});
			return renderSnippet(speciesSnippet, {
				species: row.original.specie.name
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	}),
	columnHelper.accessor('gender', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Gender'
			});
		},
		cell: ({ row }) => {
			const genderSnippet = createRawSnippet<[{ gender: string | undefined }]>(getGender => {
				const { gender } = getGender();
				return {
					render: () => `<div>${gender ? (gender === 'm' ? 'male' : 'female') : ''}</div>`
				};
			});
			return renderSnippet(genderSnippet, {
				gender: row.original.gender
			});
		}
	}),
	columnHelper.accessor('hermaphrodite', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Hermaphrodite'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('age', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Age'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('size', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Size'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('status', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Status'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('description', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Description'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('adoption_id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Adoption ID'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('created_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Create time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		}
	}),
	columnHelper.accessor('updated_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Update time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		}
	}),
	columnHelper.display({
		id: 'actions',
		cell: () => renderComponent(DataTableRowActions),
		enableSorting: false,
		enableHiding: false,
		enableResizing: false,
		maxSize: 40
	})
]);
