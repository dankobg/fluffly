--ListUsers
select * from "user";

--GetUserById
select * from "user" where "id" = $1;