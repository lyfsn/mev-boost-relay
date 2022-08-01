package database

import (
	"github.com/flashbots/boost-relay/common"
)

var (
	tableBase = common.GetEnv("DB_TABLE_PREFIX", "dev")

	// TableEvent                 = tableBase + "_event"
	TableValidatorRegistration = tableBase + "_validator_registration"
	// TableEpochSummary           = tableBase + "_epoch_summary"
	// TableSlotSummary            = tableBase + "_slot_summary"
	TableBuilderBlockSubmission = tableBase + "_builder_block_submission"
	TableBuilderBlockSimResult  = tableBase + "_builder_block_sim_result"
	TableDeliveredPayload       = tableBase + "_payload_delivered"
)

var schema = `
CREATE TABLE IF NOT EXISTS ` + TableValidatorRegistration + ` (
	id          bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	pubkey        varchar(98) NOT NULL UNIQUE,
	fee_recipient varchar(42) NOT NULL,
	timestamp     bigint NOT NULL,
	gas_limit     bigint NOT NULL,
	signature     text NOT NULL
);

CREATE TABLE IF NOT EXISTS ` + TableBuilderBlockSubmission + ` (
	id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	epoch bigint NOT NULL,
	slot  bigint NOT NULL,

	builder_pubkey  text NOT NULL,
	proposer_pubkey text NOT NULL,
	proposer_fee_recipient text NOT NULL,

	parent_hash  text NOT NULL,
	block_hash   text NOT NULL,
	block_number bigint NOT NULL,
	num_tx       int NOT NULL,
	value        NUMERIC(48, 0),

	gas_used  bigint NOT NULL,
	gas_limit bigint NOT NULL,

	payload json NOT NULL
);

CREATE TABLE IF NOT EXISTS ` + TableBuilderBlockSimResult + ` (
	id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	block_submission_id bigint references ` + TableBuilderBlockSubmission + `(id) ON DELETE CASCADE,
	success boolean NOT NULL,
	error   text NOT NULL
);

CREATE TABLE IF NOT EXISTS ` + TableDeliveredPayload + ` (
	id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	epoch bigint NOT NULL,
	slot  bigint NOT NULL,

	builder_pubkey  text NOT NULL,
	proposer_pubkey text NOT NULL,
	proposer_fee_recipient text NOT NULL,

	parent_hash  text NOT NULL,
	block_hash   text NOT NULL,
	block_number bigint NOT NULL,
	num_tx       int NOT NULL,
	value        NUMERIC(48, 0),

	gas_used  bigint NOT NULL,
	gas_limit bigint NOT NULL,

	execution_payload     json NOT NULL,
	bid_trace             json NOT NULL,
	bid_trace_builder_sig text NOT NULL,
	signed_builder_bid    json NOT NULL,
	signed_blinded_beacon_block json NOT NULL
);
`

/*
CREATE TABLE IF NOT EXISTS ` + TableEpochSummary + ` (
	id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	epoch      bigint NOT NULL UNIQUE,
	slot_first bigint NOT NULL,
	slot_last  bigint NOT NULL,
	slot_first_processed bigint NOT NULL,
	slot_last_processed  bigint NOT NULL,

	validators_known_total          bigint NOT NULL,
	validator_registrations_total   bigint NOT NULL,
	validator_registrations_saved   bigint NOT NULL,
	validator_registrations_received_unverified  bigint NOT NULL,

	num_register_validator_requests bigint NOT NULL,
	num_get_header_requests         bigint NOT NULL,
	num_get_payload_requests        bigint NOT NULL,

	num_header_sent_ok       bigint NOT NULL,
	num_header_sent_204      bigint NOT NULL,
	num_payload_sent         bigint NOT NULL,
	num_builder_bid_received bigint NOT NULL,

	is_complete boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS ` + TableEvent + ` (
	id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	slot       bigint NOT NULL,
	epoch      bigint NOT NULL,
	event_type varchar(255) NOT NULL,
	event_data json NOT NULL
);

CREATE TABLE IF NOT EXISTS ` + TableSlotSummary + ` (
	id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	epoch  bigint NOT NULL UNIQUE,
	slot   bigint NOT NULL,

	validators_known_total          bigint NOT NULL,
	validator_registrations_total   bigint NOT NULL,

	num_get_header_requests         bigint NOT NULL,
	num_get_payload_requests        bigint NOT NULL,

	num_header_sent          bigint NOT NULL,
	num_header_no_content    bigint NOT NULL,
	num_payload_sent         bigint NOT NULL,
	num_builder_bid_received bigint NOT NULL,
	highest_bid_value NUMERIC(48, 0),

	error text NOT NULL
);
*/
