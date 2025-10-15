import type { Session } from '@ory/client-fetch';

export type CustomTraits = {
	email: string;
	first_name?: string;
	last_name?: string;
	avatar_url?: string;
};

export interface User {
	id: string;
	email?: string;
	firstName?: string;
	lastName?: string;
	avatarUrl?: string;
	role?: string;
	isEmployee?: boolean;
	fullName?: string;
}

export type SessionService =
	| {
			status: 'active';
			session: Session;
			user: User;
	  }
	| {
			status: 'inactive';
			session: Session | null;
			user: null;
	  };

export function createSessionService(session: Session | null): SessionService {
	if (!session?.active || !session.identity) {
		return { status: 'inactive', session, user: null };
	}

	return {
		status: 'active',
		session,
		user: {
			id: session.identity.id,
			email: session.identity.traits?.email,
			firstName: session.identity.traits?.first_name,
			lastName: session.identity.traits?.last_name,
			avatarUrl: session.identity.traits?.avatar_url,
			get role(): string {
				switch (session.identity?.schema_id) {
					case 'default':
						return 'default';
					case 'employee':
						return 'employee';
					default:
						return 'default';
				}
			},
			get isEmployee(): boolean {
				return this.role === 'employee';
			},
			get fullName(): string {
				return `${this.firstName || ''}${this.lastName ? ` ${this.lastName}` : ''}`.trim();
			}
		}
	};
}
