import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { fluffly } from '$lib/fluffly/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends, parent }) => {
	depends(`data:dashboard`);

	const data = await parent();

	if (data.auth.user?.isDeveloper) {
		try {
			const analyticsStatsResult = await fluffly.GET('/analytics/stats', {
				fetch
			});

			if (analyticsStatsResult.error?.status_code === 403) {
				if (browser) {
					goto('/');
				}
			}

			return {
				analyticsStatsResult
			};
		} catch (error) {
			console.log('err', error);
		}
	} else {
		try {
			const myAnalyticsStatsResult = await fluffly.GET('/me/analytics/stats', {
				fetch
			});

			if (myAnalyticsStatsResult.error?.status_code === 403) {
				if (browser) {
					goto('/');
				}
			}

			return {
				myAnalyticsStatsResult
			};
		} catch (error) {
			console.log('err', error);
		}
	}
};
