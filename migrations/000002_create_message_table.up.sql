begin;

create table if not exists messages (
	id bigserial primary key,
	created_at timestamp(0) with time zone not null default now(),
	user_id bigint not null references users on delete cascade,
    body text not null
);

commit;
