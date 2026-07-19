const positionErrorCodes = new Map([
	[1, 'PERMISSION_DENIED'],
	[2, 'POSITION_UNAVAILABLE'],
	[3, 'TIMEOUT']
]);

async function checkGeolocationPermission() {
	const result = await navigator.permissions.query({ name: 'geolocation' });
	return result.state;
}

export async function getCurrentPosition(): Promise<[number, number]> {
	if (import.meta.env.DEV) {
		return [44.7866, 20.4489];
	}

	if (!('geolocation' in navigator)) {
		throw new Error('geolocation not supported');
	}

	if ('permissions' in navigator) {
		const state = await checkGeolocationPermission();
		if (state === 'denied') {
			throw new Error('user denied geolocation');
		}
	}

	return new Promise((resolve, reject) => {
		navigator.geolocation.getCurrentPosition(
			pos => {
				resolve([pos.coords.latitude, pos.coords.longitude]);
			},
			err => {
				reject(new Error('dsada', { cause: positionErrorCodes.get(err.code) }));
			},
			{ timeout: 5_000, maximumAge: 30_000 }
		);
	});
}
