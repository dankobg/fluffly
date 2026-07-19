import * as v from 'valibot';

export const emptyStringToUndefined = v.pipe(
	v.literal(''),
	v.transform(() => undefined)
);

export const emptyStringToNull = v.pipe(
	v.literal(''),
	v.transform(() => null)
);

export const zeroToUndefined = v.pipe(
	v.literal(0),
	v.transform(() => undefined)
);

export const zeroToNull = v.pipe(
	v.literal(0),
	v.transform(() => null)
);
