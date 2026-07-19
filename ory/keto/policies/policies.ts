//@ts-nocheck

import { Namespace, SubjectSet, Context } from "@ory/keto-namespace-types";

// Group membership
class Group implements Namespace {
  related: {
    members: Identity[];
  };
}

// Health represents health check resources
class Health implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean => this.related.viewers.includes(ctx.subject),
  };
}

// Schemas represents all identity schemas resource
class Schemas implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// Schema represents individual identity schema resource
class Schema implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Schemas[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// Sessions represents all sessions resource
class Sessions implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// Session represents individual session resource
class Session implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Sessions[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// Identities represents all identities resource
class Identities implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// Identity represents individual identity resource
class Identity implements Namespace {
  related: {
    owners: Identity[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Identities[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// CourierMessages represents all courier messages resource
class CourierMessages implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// CourierMessage represents individual courier message resource
class CourierMessage implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: CourierMessages[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// Countries represents all countries resource
class Countries implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// Country represents individual country resource
class Country implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Countries[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// Animals represents all animals resource
class Animals implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
    submit: (ctx: Context): boolean => this.permits.view(ctx),
  };
}

// Animal represents individual animal resource
class Animal implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Animals[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
    like: (ctx: Context): boolean => this.permits.view(ctx),
    unlike: (ctx: Context): boolean => this.permits.view(ctx),
    adopt: (ctx: Context): boolean =>
      !this.related.owners.includes(ctx.subject) && this.permits.view(ctx),
  };
}

// AnimalTypes represents all animal types resource
class AnimalTypes implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// AnimalType represents individual animal type resource
class AnimalType implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Countries[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// AnimalSpecies represents all animal species resource
class AnimalSpecies implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// AnimalSpecie represents individual animal specie resource
class AnimalSpecie implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: AnimalSpecies[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// Organizations represents all organizations resource
class Organizations implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
    apply: (ctx: Context): boolean => this.permits.view(ctx),
  };
}

// Organization represents individual organization resource
class Organization implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Organizations[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// Breeds represents all breeds resource
class Breeds implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) || this.permits.manage(ctx),
    manage: (ctx: Context): boolean =>
      this.related.managers.includes(ctx.subject),
  };
}

// Breed represents individual breed resource
class Breed implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    managers: (Identity | SubjectSet<Group, "members">)[];
    parents: Breeds[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.permits.manage(ctx) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
    manage: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.managers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.manage(ctx)),
  };
}

// Analytics represents all analytics resource
class Analytics implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean => this.related.viewers.includes(ctx.subject),
  };
}

// Analytic represents individual analytic resource
class Analytic implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    parents: Analytics[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.related.owners.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
  };
}

// Adoptions represents all adoptions resource
class Adoptions implements Namespace {
  related: {
    viewers: (Identity | SubjectSet<Group, "members">)[];
  };

  permits = {
    view: (ctx: Context): boolean => this.related.viewers.includes(ctx.subject),
  };
}

// Adoption represents individual adoption resource
class Adoption implements Namespace {
  related: {
    owners: (Identity | SubjectSet<Group, "members">)[];
    viewers: (Identity | SubjectSet<Group, "members">)[];
    parents: Adoptions[];
  };

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.related.owners.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
  };
}
