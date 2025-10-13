import { clsx, type ClassValue } from 'clsx';
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
