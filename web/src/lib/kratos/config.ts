const baseUrl = import.meta.env['VITE_PUBLIC_SITE_URL'] as string;
const publicUrl = import.meta.env['VITE_PUBLIC_KRATOS_PUBLIC_URL'] as string;
const adminUrl = import.meta.env['VITE_PUBLIC_KRATOS_ADMIN_URL'] as string;

export const config = {
	kratos: {
		baseUrl,
		publicUrl,
		adminUrl
	},

	routes: {
		registration: {
			path: '/auth/registration',
			selfServiceUrl: `${publicUrl}/self-service/registration/browser`
		},
		login: {
			path: '/auth/login',
			selfServiceUrl: `${publicUrl}/self-service/login/browser`
		},
		logout: {
			path: '/auth/logout',
			selfServiceUrl: `${publicUrl}/self-service/logout/browser`
		},
		verification: {
			path: '/auth/verification',
			selfServiceUrl: `${publicUrl}/self-service/verification/browser`
		},
		recovery: {
			path: '/auth/recovery',
			selfServiceUrl: `${publicUrl}/self-service/recovery/browser`
		},
		settings: {
			path: '/dashboard/account',
			selfServiceUrl: `${publicUrl}/self-service/settings/browser`
		},
		dashboard: {
			path: '/dashboard'
		},
		refresh: {
			selfServiceUrl: (returnTo?: string): string =>
				`${publicUrl}/self-service/browser/flows/login?refresh=true&return_to=${returnTo || baseUrl}`
		}
	},

	labels: {
		to_verify: {
			label: 'Email',
			priority: 100
		},
		csrf_token: {
			label: '',
			priority: 90
		},
		'traits.email': {
			label: 'Email',
			priority: 80
		},
		email: {
			label: 'Email',
			priority: 70
		},
		identifier: {
			label: 'Email',
			priority: 60
		},
		password: {
			label: 'Password',
			priority: 50
		},
		'traits.username': {
			label: 'Username',
			priority: 40
		},
		'traits.first_name': {
			label: 'First name',
			priority: 30
		},
		'traits.last_name': {
			label: 'Last name',
			priority: 20
		}
	}
};
