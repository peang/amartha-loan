CREATE TABLE investments (
  id BIGSERIAL PRIMARY KEY,
  loan_id BIGINT NOT NULL,
  investor_id BIGINT NOT NULL,
  amount NUMERIC(10,2) NOT NULL,
  roi NUMERIC(10,2) NOT NULL,
  send_aggreement_email BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_investments_loan_id ON investments (loan_id);
CREATE INDEX idx_investments_investor_id ON investments (investor_id);
CREATE INDEX idx_investments_send_agreement_email ON investments (send_aggreement_email);