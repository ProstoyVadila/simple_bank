alter table if exists accounts drop constraint if exists accounts_owner_name_fkey;
alter table if exists accounts drop constraint if exists owner_currency_key;
drop table if exists users;
