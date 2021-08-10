// Code generated from source setup_postgres.sql via go generate. DO NOT EDIT.

package postgres

const setup_postgres_sql = `-- This file is setup_postgres.sql which gets compiled into go source using a go:generate statement in postgres.go
--
-- TODO? replace all 'addr bytea' with 'addr_id bigint' and a mapping table? makes addrs an 8 byte int that fits in a register instead of a 32 byte string
CREATE SCHEMA public;

CREATE TABLE IF NOT EXISTS block_header (
round bigint PRIMARY KEY,
realtime timestamp without time zone NOT NULL,
rewardslevel bigint NOT NULL,
header jsonb NOT NULL
);

-- For looking round by timestamp. We could replace this with a round-to-timestamp algorithm, it should be extremely
-- efficient since there is such a high correlation between round and time.
CREATE INDEX IF NOT EXISTS block_header_time ON block_header (realtime);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS txn (
round bigint NOT NULL,
intra smallint NOT NULL,
typeenum smallint NOT NULL,
asset bigint NOT NULL, -- 0=Algos, otherwise AssetIndex
txid bytea NOT NULL, -- base32 of [32]byte hash
txnbytes bytea NOT NULL, -- msgpack encoding of signed txn with apply data
txn jsonb NOT NULL, -- json encoding of signed txn with apply data
note_type VARCHAR NOT NULL,
note_txid UUID NOT NULL,
note jsonb NOT NULL,
extra jsonb,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
PRIMARY KEY ( round, intra )
);

CREATE TRIGGER set_timestamp
BEFORE
UPDATE ON txn
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- For transaction lookup
CREATE INDEX IF NOT EXISTS txn_by_tixid ON txn ( txid );

--For note field lookup
CREATE INDEX IF NOT EXISTS txn_by_note ON txn (note_type, note_txid);

-- Optional, to make txn queries by asset fast:
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS txn_asset ON txn (asset, round, intra);

CREATE TABLE IF NOT EXISTS txn_participation (
addr bytea NOT NULL,
round bigint NOT NULL,
intra smallint NOT NULL
);

-- For query account transactions
CREATE UNIQUE INDEX IF NOT EXISTS txn_participation_i ON txn_participation ( addr, round DESC, intra DESC );

-- expand data.basics.AccountData
CREATE TABLE IF NOT EXISTS account (
  addr bytea primary key,
  microalgos bigint NOT NULL, -- okay because less than 2^54 Algos
  rewardsbase bigint NOT NULL,
  rewards_total bigint NOT NULL,
  deleted bool NOT NULL, -- whether or not it is currently deleted
  created_at bigint NOT NULL DEFAULT 0, -- round that the account is first used
  closed_at bigint, -- round that the account was last closed
  keytype varchar(8), -- sig,msig,lsig
  account_data jsonb -- trimmed AccountData that only contains auth addr and keyreg info
);

-- data.basics.AccountData Assets[asset id] AssetHolding{}
CREATE TABLE IF NOT EXISTS account_asset (
  addr bytea NOT NULL, -- [32]byte
  assetid bigint NOT NULL,
  amount numeric(20) NOT NULL, -- need the full 18446744073709551615
  frozen boolean NOT NULL,
  deleted bool NOT NULL, -- whether or not it is currently deleted
  created_at bigint NOT NULL DEFAULT 0, -- round that the asset was added to an account
  closed_at bigint, -- round that the asset was last removed from the account
  PRIMARY KEY (addr, assetid)
);

-- For account lookup
CREATE INDEX IF NOT EXISTS account_asset_by_addr ON account_asset ( addr );

-- Optional, to make queries of all asset balances fast /v2/assets/<assetid>/balances
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS account_asset_asset ON account_asset (assetid, addr ASC);

-- data.basics.AccountData AssetParams[index] AssetParams{}
CREATE TABLE IF NOT EXISTS asset (
  index bigint PRIMARY KEY,
  creator_addr bytea NOT NULL,
  params jsonb NOT NULL, -- data.basics.AssetParams -- TODO index some fields?
  deleted bool NOT NULL, -- whether or not it is currently deleted
  created_at bigint NOT NULL DEFAULT 0, -- round that the asset was created
  closed_at bigint -- round that the asset was closed; cannot be recreated because the index is unique
);

-- For account lookup
CREATE INDEX IF NOT EXISTS asset_by_creator_addr ON asset ( creator_addr );

-- subsumes ledger/accountdb.go accounttotals and acctrounds
-- "state":{online, onlinerewardunits, offline, offlinerewardunits, notparticipating, notparticipatingrewardunits, rewardslevel, round bigint}
CREATE TABLE IF NOT EXISTS metastate (
  k text primary key,
  v jsonb
);

-- per app global state
-- roughly go-algorand/data/basics/userBalance.go AppParams
CREATE TABLE IF NOT EXISTS app (
  index bigint PRIMARY KEY,
  creator bytea, -- account address
  params jsonb,
  deleted bool NOT NULL, -- whether or not it is currently deleted
  created_at bigint NOT NULL DEFAULT 0, -- round that the asset was created
  closed_at bigint -- round that the app was deleted; cannot be recreated because the index is unique
);

-- For account lookup
CREATE INDEX IF NOT EXISTS app_by_creator ON app ( creator );

-- per-account app local state
CREATE TABLE IF NOT EXISTS account_app (
  addr bytea,
  app bigint,
  localstate jsonb,
  deleted bool NOT NULL, -- whether or not it is currently deleted
  created_at bigint NOT NULL DEFAULT 0, -- round that the app was added to an account
  closed_at bigint, -- round that the account_app was last removed from the account
  PRIMARY KEY (addr, app)
);

-- For storing per-transaction closing balance
CREATE TABLE IF NOT EXISTS txn_closingbalance (
  note_txid UUID NOT NULL primary key, 
  receiver_closingbalance numeric(20), 
  sender_closingbalance numeric(20), 
  assetid bigint NOT NULL, 
  receiver_addr bytea NOT NULL, 
  sender_addr bytea NOT NULL
);

-- For account lookup
CREATE INDEX IF NOT EXISTS account_app_by_addr ON account_app ( addr );

-- Adding Redemptions table via foreign data wrapper
CREATE EXTENSION postgres_fdw;

CREATE SERVER redemptions FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'redemptions-staging.postgres.database.azure.com', dbname 'postgres', port '5432');

CREATE USER MAPPING FOR CURRENT_USER
SERVER redemptions
OPTIONS (user 'xfinite@redemptions-staging', password 'Xinaam@123');

CREATE SCHEMA redemption;

IMPORT FOREIGN SCHEMA public FROM SERVER redemptions INTO redemption;

-- Adding Rewards table via foreign data wrapper
CREATE SERVER balances FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'reward-engine-stagin.postgres.database.azure.com', dbname 'postgres', port '5432' );

CREATE USER MAPPING FOR CURRENT_USER
SERVER balances
OPTIONS (user 'xfinite@reward-engine-stagin', password 'Xinaam@123');

CREATE SCHEMA balances;

IMPORT FOREIGN SCHEMA public FROM SERVER balances INTO balances;

create type result_type as (
  amount varchar,
  closingbalance varchar
);
  
create or replace function get_transaction_closing_balance(nid uuid, user_addr varchar, uid varchar) returns result_type as 
$$                                                                                                                    
DECLARE                                                                                                                 
transaction_row record;
closingbalance decimal:=0 ;
amount decimal:=0;
coin_id varchar:='362b2e89-de10-4974-99aa-ea6a55bf30d3';
result result_type;
begin
select cast("Balances".amount as decimal) into closingbalance from balances."Balances" where user_id=uid  order by created_at desc limit 1;
select txn.note->'meta'->'coin_id' into coin_id from public.txn where txn.note->'meta'?'coin_id' and txn.txn->'txn'->>'type' = 'axfer' and note_txid = nid;
select cast(txn.txn->'txn'->'aamt' as decimal) into amount from public.txn where txn.txn->'txn' ? 'aamt' and txn.txn->'txn'->>'type' = 'axfer' and cast(txn.txn->'txn'->>'xaid' as integer) = (case when coin_id = '362b2e89-de10-4974-99aa-ea6a55bf30d3' then 1 else 1 end) and note_txid = nid;
if closingbalance<0 then closingbalance = 0;
elsif user_addr = txn.txn->'txn'->>'snd' from txn where note_txid = nid then
closingbalance = closingbalance + (select sender_closingbalance from txn_closingbalance where note_txid = nid);
else closingbalance = closingbalance + (select receiver_closingbalance from txn_closingbalance where note_txid = nid);
end if;
result.amount := amount::varchar;
result.closingbalance := closingbalance::varchar;
return result;
end;
$$ LANGUAGE plpgsql;`
