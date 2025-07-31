--ListUsers
select * from "user";

--GetUserById
select * from "user" where "id" = $1;


--ListOrganizations
select * from "organization";

--GetOrganizationById
select * from "organization" where "id" = $1;
