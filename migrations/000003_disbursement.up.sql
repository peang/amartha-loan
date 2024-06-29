CREATE TABLE disbursements (
  id BIGSERIAL PRIMARY KEY,
  field_officer_id BIGINT NOT NULL,
  aggreement_file_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE INDEX idx_disbursment_field_officer_id ON disbursements (field_officer_id);
