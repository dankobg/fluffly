--ListCountries
select * from "country";

--ListUsers
select * from "user";

--GetUserById
select * from "user" where "id" = $1;


--ListOrganizations
select
  o.*,
  coalesce(
    json_agg(distinct jsonb_build_object(
      'id', oph.id,
      'small', oph.small,
      'medium', oph.medium,
      'large', oph.large,
      'full', oph.full
    )) filter (where oph.organization_id is not null),
    '[]'
  ) as photos,
  coalesce(
    json_agg(distinct jsonb_build_object(
      'platform', osc.platform,
      'url', osc.url
    )) filter (where osc.organization_id is not null),
    '[]'
  ) as socials,
  --prefix:organization_work_hour__
  owh.*,
  --prefix:organization_contact__
  oc.*,
  --prefix:organization_contact.address__
  a.*,
  --prefix:organization_contact.address.country__
  c.*
from organization o
left join organization_contact oc on o.id = oc.organization_id
left join address a on a.id = oc.address_id
left join country c on c.id = a.country_id
left join organization_work_hour owh on o.id = owh.organization_id
left join organization_photo oph on o.id = oph.organization_id
left join organization_social osc on o.id = osc.organization_id
group by o.id, oc.id, oc.email, owh.id, a.id, c.id;

--GetOrganizationById
select
  o.*,
  coalesce(
    json_agg(distinct jsonb_build_object(
      'id', oph.id,
      'small', oph.small,
      'medium', oph.medium,
      'large', oph.large,
      'full', oph.full
    )) filter (where oph.organization_id is not null),
    '[]'
  ) as photos,
  coalesce(
    json_agg(distinct jsonb_build_object(
      'platform', osc.platform,
      'url', osc.url
    )) filter (where osc.organization_id is not null),
    '[]'
  ) as socials,
  --prefix:organization_work_hour__
  owh.*,
  --prefix:organization_contact__
  oc.*,
  --prefix:organization_contact.address__
  a.*,
  --prefix:organization_contact.address.country__
  c.*
from organization o
left join organization_contact oc on o.id = oc.organization_id
left join address a on a.id = oc.address_id
left join country c on c.id = a.country_id
left join organization_work_hour owh on o.id = owh.organization_id
left join organization_photo oph on o.id = oph.organization_id
left join organization_social osc on o.id = osc.organization_id
where o."id" = $1
group by o.id, oc.id, oc.email, owh.id, a.id, c.id;