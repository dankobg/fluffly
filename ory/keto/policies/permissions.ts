//@ts-nocheck

import { Namespace, SubjectSet, Context } from "@ory/keto-namespace-types";

class Group implements Namespace {
  related: {
    members: Identity[];
  };
}

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
