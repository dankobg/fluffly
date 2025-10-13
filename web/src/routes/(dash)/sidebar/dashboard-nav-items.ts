import IconLayoutDashboard from '@lucide/svelte/icons/layout-dashboard';
import IconUser from '@lucide/svelte/icons/user';
import IconUsers from '@lucide/svelte/icons/users';
import IconMail from '@lucide/svelte/icons/mail';
import IconFingerprint from '@lucide/svelte/icons/fingerprint';
import IconNewspaper from '@lucide/svelte/icons/newspaper';

export const dashboardNavItems = [
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
		title: 'Identity',
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
			}
		]
	},
	{
		title: 'Courier',
		url: '#',
		items: [
			{
				title: 'Messages',
				url: '/dashboard/messages',
				icon: IconMail
			}
		]
	}
];
