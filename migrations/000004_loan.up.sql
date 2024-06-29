CREATE TABLE loans (
  id BIGSERIAL PRIMARY KEY,
  uuid UUID NOT NULL,
  borrower_id BIGINT NOT NULL,
  approval_id BIGINT,
  proposed_amount NUMERIC(10,2) NOT NULL,
  principal_amount NUMERIC(10,2),
  rate NUMERIC(10,2),
  roi NUMERIC(10,2),
  status INT NOT NULL,
  aggreement_file_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
  CONSTRAINT fk_approval
    FOREIGN KEY(approval_id) 
    REFERENCES approvals(id)
);

CREATE INDEX idx_loan_status ON loans (status);
CREATE INDEX idx_loan_borrower_id ON loans (borrower_id);
CREATE INDEX idx_loan_approval_id ON loans (approval_id);
CREATE INDEX idx_loan_proposed_amount ON loans (proposed_amount);
CREATE INDEX idx_loan_principal_amount ON loans (principal_amount);