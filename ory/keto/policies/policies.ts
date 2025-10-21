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
    viewers: (Animal | SubjectSet<Group, "members">)[];
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
    owners: Schema[];
    viewers: (Schema | SubjectSet<Group, "members">)[];
    managers: (Schema | SubjectSet<Group, "members">)[];
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
    owners: Session[];
    viewers: (Session | SubjectSet<Group, "members">)[];
    managers: (Session | SubjectSet<Group, "members">)[];
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
    owners: CourierMessage[];
    viewers: (CourierMessage | SubjectSet<Group, "members">)[];
    managers: (CourierMessage | SubjectSet<Group, "members">)[];
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

// Country represents individual identity resource
class Country implements Namespace {
  related: {
    owners: Country[];
    viewers: (Country | SubjectSet<Group, "members">)[];
    managers: (Country | SubjectSet<Group, "members">)[];
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

// Animals represents all identities resource
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
  };
}

// Animal represents individual identity resource
class Animal implements Namespace {
  related: {
    owners: Animal[];
    viewers: (Animal | SubjectSet<Group, "members">)[];
    managers: (Animal | SubjectSet<Group, "members">)[];
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
  };
}

// Organization represents individual organization resource
class Organization implements Namespace {
  related: {
    owners: Organization[];
    viewers: (Organization | SubjectSet<Group, "members">)[];
    managers: (Organization | SubjectSet<Group, "members">)[];
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
