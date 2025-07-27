local claims = {
  email_verified: true
} + std.extVar('claims');

{
  identity: {
    traits: {
      [if "email" in claims && claims.email_verified then "email" else null]: claims.email,
      [if "given_name" in claims && claims.email_verified then "first_name" else null]: claims.given_name,
      [if "family_name" in claims && claims.email_verified then "last_name" else null]: claims.family_name,
      [if "picture" in claims && claims.email_verified then "avatar_url" else null]: claims.picture,
    },
  },
}
