import IconLayoutDashboard from '@lucide/svelte/icons/layout-dashboard';
import IconUser from '@lucide/svelte/icons/user';
import IconUsers from '@lucide/svelte/icons/users';
import IconMail from '@lucide/svelte/icons/mail';
import IconFingerprint from '@lucide/svelte/icons/fingerprint';
import IconNewspaper from '@lucide/svelte/icons/newspaper';
import IconEarth from '@lucide/svelte/icons/earth';
import IconBuilding1 from '@lucide/svelte/icons/building';
import IconBuilding2 from '@lucide/svelte/icons/building-2';
import IconPawPrint from '@lucide/svelte/icons/paw-print';
import IconHeartHandshake from '@lucide/svelte/icons/heart-handshake';

import type { Component } from 'svelte';

export type NavItem = {
	title: string;
	url?: string;
	isActive?: boolean;
	items?: NavItem[];
	icon?: Component;
};

export const mainNavItems: NavItem[] = [
	{ url: '#', title: 'About' },
	{ url: '/search', title: 'Search' },
	{ url: '#', title: 'Blog' },
	{ url: '#', title: 'Contact' }
];

export const customerDashboardNavItems: NavItem[] = [
	{
		title: 'App',
		url: '#',
		items: [
			{
				title: 'Dashboard',
				url: '/dashboard',
				icon: IconLayoutDashboard
			}
		]
	},
	{
		title: 'User',
		url: '#',
		items: [
			{
				title: 'Account',
				url: '/dashboard/account',
				icon: IconUser
			}
		]
	},
	{
		title: 'Actions',
		url: '#',
		items: [
			{
				title: 'Apply for organization',
				url: '/dashboard/organizations/apply',
				icon: IconBuilding1
			},
			{
				title: 'Submit animal for adoption',
				url: '/dashboard/animals/submit',
				icon: IconBuilding1
			}
		]
	},
	{
		title: 'Pet center',
		url: '#',
		items: [
			{
				title: 'My organizations',
				url: '/dashboard/organizations/mine',
				icon: IconPawPrint
			},
			{
				title: 'My posted',
				url: '/dashboard/animals/posted',
				icon: IconPawPrint
			},
			{
				title: 'My favorites',
				url: '/dashboard/animals/favorite',
				icon: IconPawPrint
			},
			{
				title: 'My adopted',
				url: '/dashboard/animals/adopted',
				icon: IconPawPrint
			}
		]
	}
];

export const developerDashboardNavItems: NavItem[] = [
	{
		title: 'App',
		url: '#',
		items: [
			{
				title: 'Dashboard',
				url: '/dashboard',
				icon: IconLayoutDashboard
			}
		]
	},
	{
		title: 'User',
		url: '#',
		items: [
			{
				title: 'Account',
				url: '/dashboard/account',
				icon: IconUser
			}
		]
	},
	{
		title: 'Actions',
		url: '#',
		items: [
			{
				title: 'Apply for organization',
				url: '/dashboard/organizations/apply',
				icon: IconBuilding1
			},
			{
				title: 'Submit animal for adoption',
				url: '/dashboard/animals/submit',
				icon: IconBuilding1
			}
		]
	},
	{
		title: 'Pet center',
		url: '#',
		items: [
			{
				title: 'My organizations',
				url: '/dashboard/organizations/mine',
				icon: IconPawPrint
			},
			{
				title: 'My posted',
				url: '/dashboard/animals/posted',
				icon: IconPawPrint
			},
			{
				title: 'My favorites',
				url: '/dashboard/animals/favorite',
				icon: IconPawPrint
			},
			{
				title: 'My adopted',
				url: '/dashboard/animals/adopted',
				icon: IconPawPrint
			}
		]
	},
	{
		title: 'Auth',
		url: '#',
		items: [
			{
				title: 'Schemas',
				url: '/dashboard/schemas',
				icon: IconNewspaper
			},
			{
				title: 'Identities',
				url: '/dashboard/identities',
				icon: IconUsers
			},
			{
				title: 'Sessions',
				url: '/dashboard/sessions',
				icon: IconFingerprint
			},
			{
				title: 'Courier messages',
				url: '/dashboard/messages',
				icon: IconMail
			}
		]
	},
	{
		title: 'Managment',
		url: '#',
		items: [
			{
				title: 'Adoptions',
				url: '/dashboard/adoptions',
				icon: IconHeartHandshake
			},
			{
				title: 'Organizations',
				url: '/dashboard/organizations',
				icon: IconEarth
			},
			{
				title: 'Animals',
				url: '/dashboard/animals',
				icon: IconPawPrint
			},
			{
				title: 'Animal types',
				url: '/dashboard/animal-types',
				icon: IconPawPrint
			},
			{
				title: 'Animal species',
				url: '/dashboard/animal-species',
				icon: IconPawPrint
			},
			{
				title: 'Animal breeds',
				url: '/dashboard/breeds',
				icon: IconPawPrint
			},
			{
				title: 'Countries',
				url: '/dashboard/countries',
				icon: IconBuilding2
			}
		]
	}
];
