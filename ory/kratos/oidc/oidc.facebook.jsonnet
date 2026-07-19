local claims = std.extVar('claims');

{
  identity: {
    traits: {
      [if "email" in claims then "email" else null]: claims.email,
      [if "given_name" in claims then "first_name" else null]: claims.given_name,
      [if "family_name" in claims then "last_name" else null]: claims.family_name,
      [if "picture" in claims then "avatar_url" else null]: claims.picture,
    },
  },
}