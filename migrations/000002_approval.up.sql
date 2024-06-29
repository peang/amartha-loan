CREATE TABLE approvals (
  id BIGSERIAL PRIMARY KEY,
  field_validator_id BIGINT NOT NULL,
  approval_file_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE INDEX idx_approval_field_validator_id ON approvals (field_validator_id);
