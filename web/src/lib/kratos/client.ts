import { Configuration, FrontendApi } from '@ory/client-fetch';
import { config } from './config';

const kratosConfiguration = new Configuration({
	basePath: config.kratos.publicUrl,
	credentials: 'include',
	fetchApi: fetch,
	headers: {
		'Content-Type': 'application/json',
		Accept: 'application/json'
	}
});

export const kratos = new FrontendApi(kratosConfiguration);
