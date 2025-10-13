import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent, depends }) => {
	const data = await parent();
	depends('data:dashboard');
	try {
		const mystuff = { rofl: 'hello bozo' };
		return {
			...data,
			mystuff
		};
	} catch (error) {
		console.log('err', error);
		return {
			...data,
			mystuff: null
		};
	}
};
