import {
	AdoptionStatus,
	AnimalAge,
	AnimalGender,
	AnimalSize,
	AnimalStatus,
	OrganizationStatus
} from './gen/fluffly_openapi';

export const organizationStatusValues = Object.values(OrganizationStatus);

export const animalSizeValues = Object.values(AnimalSize);
export const animalAgeValues = Object.values(AnimalAge);
export const animalGenderValues = Object.values(AnimalGender);
export const animalStatusValues = Object.values(AnimalStatus);
export const adoptionStatusValues = Object.values(AdoptionStatus);

export const daysLessThanOptions = [
	{ value: '1', label: '1 Days' },
	{ value: '3', label: '3 Days' },
	{ value: '5', label: '5 Days' },
	{ value: '7', label: '7 Days' },
	{ value: '14', label: '14 Days' },
	{ value: '30', label: '30 Days' },
	{ value: '60', label: '60 Days' }
];

export const daysGreaterThanOptions = [
	{ value: '60', label: '1 Month' },
	{ value: '90', label: '3 Months' },
	{ value: '180', label: '6 Months' },
	{ value: '270', label: '9 Months' },
	{ value: '365', label: '1 Year' }
];

export const radiusMetersOptions = [
	{ value: '10000', label: '10 km' },
	{ value: '20000', label: '20 km' },
	{ value: '30000', label: '30 km' },
	{ value: '50000', label: '50 km' },
	{ value: '75000', label: '75 km' },
	{ value: '100000', label: '100 km' },
	{ value: '200000', label: '200 km' },
	{ value: '0', label: 'Anywhere' }
];
