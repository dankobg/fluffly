const backendBaseUrl = import.meta.env['VITE_PUBLIC_BACKEND_BASE'] as string;
import { clsx, type ClassValue } from 'clsx';
import type { FormPath } from 'sveltekit-superforms';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, 'child'> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, 'children'> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };

export function getInitials(name?: string): string {
	if (!name) {
		return '';
	}
	const parts = name.trim().toUpperCase().split(/\s+/);
	const first = parts[0]?.charAt(0) || '';
	if (parts.length === 1) {
		return first;
	}
	const last = parts.at(-1)?.charAt(0);
	return first + last;
}

export function getChangedFormFields<T extends Record<string, unknown>>(
	form: T,
	isTainted: (path?: FormPath<T> | Record<string, unknown> | boolean | undefined) => boolean,
	pathPrefix = ''
): Partial<T> {
	const result: Record<string, unknown> = {};

	for (const key in form) {
		const value = form[key];
		const path = [pathPrefix, key].filter(Boolean).join('.');

		if (value instanceof File) {
			if (isTainted(path as FormPath<T>)) {
				result[key] = value;
			}
			continue;
		}

		if (value !== null && typeof value === 'object' && !Array.isArray(value)) {
			const nested = getChangedFormFields(value as T, isTainted, path);
			if (Object.keys(nested).length > 0) {
				result[key] = nested;
			}
		} else {
			if (isTainted(path as FormPath<T>)) {
				result[key] = value;
			}
		}
	}

	return result as Partial<T>;
}

export function capitalize(s: string) {
	if (!s) return '';
	return s[0]?.toUpperCase() + s.slice(1);
}

export function getErrorFallbackAnimalImage(animalTypeName: string): string {
	let max = 5;
	let name = animalTypeName.toLowerCase();
	switch (animalTypeName) {
		case 'Small & Furry':
			name = 'rabbit';
			break;
		case 'Scales, Fins & Other':
			name = 'aquatic';
			break;
		case 'Barnyard':
		case 'Horse':
			max = 3;
			name = 'farm';
			break;
	}
	const num = Math.floor(Math.random() * max) + 1;
	return `${backendBaseUrl}/public/images/placeholder/${name}-${num}.svg`;
}
