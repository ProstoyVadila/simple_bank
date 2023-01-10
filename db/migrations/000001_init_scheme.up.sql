create extension if not exists "uuid-ossp";

create table accounts (
    id uuid primary key default uuid_generate_v4(),
    owner_name varchar not null,
    balance bigint not null,
    currency varchar not null,
    created_at timestamptz not null default (now())
);
create table entries (
    id uuid primary key default uuid_generate_v4(),
    account_id uuid not null references accounts(id) on delete cascade,
    amount bigint not null,
    created_at timestamptz not null default (now())
);
create table transfers (
    id uuid primary key default uuid_generate_v4(),
    from_account_id uuid not null references accounts(id) on delete cascade,
    to_account_id uuid not null references accounts(id) on delete cascade,
    amount bigint not null,
    created_at timestamptz not null default (now())
);

create index on accounts(owner_name);
create index on entries(account_id);
create index on transfers(from_account_id);
create index on transfers(to_account_id);
create index on transfers(from_account_id, to_account_id);
comment on column entries.amount is 'can be negative or positive';
comment on column transfers.amount is 'must be positive';
