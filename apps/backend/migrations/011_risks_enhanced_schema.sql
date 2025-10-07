-- Migration 011: Enhanced risks schema according to ASSETS_RISKS.md Sprint 1
-- Remove old statuses and add new fields

-- Update risks table structure
ALTER TABLE risks 
    ADD COLUMN IF NOT EXISTS owner_user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS methodology VARCHAR(50),
    ADD COLUMN IF NOT EXISTS strategy VARCHAR(50),
    ADD COLUMN IF NOT EXISTS due_date DATE;

-- Update existing owner_id to owner_user_id if owner_id exists
UPDATE risks 
SET owner_user_id = owner_id 
WHERE owner_id IS NOT NULL AND owner_user_id IS NULL;

-- Drop old owner_id column
ALTER TABLE risks DROP COLUMN IF EXISTS owner_id;

-- Update likelihood and impact to use 1-4 scale instead of 1-5
-- First, update the constraints
ALTER TABLE risks DROP CONSTRAINT IF EXISTS risks_likelihood_check;
ALTER TABLE risks DROP CONSTRAINT IF EXISTS risks_impact_check;
ALTER TABLE risks DROP CONSTRAINT IF EXISTS risks_level_check;

-- Add new constraints for 1-4 scale
ALTER TABLE risks ADD CONSTRAINT risks_likelihood_check CHECK (likelihood >= 1 AND likelihood <= 4);
ALTER TABLE risks ADD CONSTRAINT risks_impact_check CHECK (impact >= 1 AND impact <= 4);

-- Update the computed level column to use new formula
ALTER TABLE risks DROP COLUMN IF EXISTS level;
ALTER TABLE risks ADD COLUMN level INTEGER GENERATED ALWAYS AS (likelihood * impact) STORED;

-- Update status constraints - remove old statuses, add new ones
ALTER TABLE risks DROP CONSTRAINT IF EXISTS risks_status_check;
ALTER TABLE risks ADD CONSTRAINT risks_status_check CHECK (status IN ('new', 'in_analysis', 'in_treatment', 'accepted', 'transferred', 'mitigated', 'closed'));

-- Update existing risks to use new status 'new' instead of 'draft'
UPDATE risks SET status = 'new' WHERE status = 'draft';
UPDATE risks SET status = 'in_analysis' WHERE status = 'registered';
UPDATE risks SET status = 'in_treatment' WHERE status = 'analysis';

-- Set default status to 'new'
ALTER TABLE risks ALTER COLUMN status SET DEFAULT 'new';

-- Add methodology and strategy constraints
ALTER TABLE risks ADD CONSTRAINT risks_methodology_check CHECK (methodology IN ('ISO27005', 'NIST', 'COSO', 'Custom'));
ALTER TABLE risks ADD CONSTRAINT risks_strategy_check CHECK (strategy IN ('accept', 'mitigate', 'transfer', 'avoid'));

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_risks_owner_user_id ON risks(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_risks_methodology ON risks(methodology);
CREATE INDEX IF NOT EXISTS idx_risks_strategy ON risks(strategy);
CREATE INDEX IF NOT EXISTS idx_risks_due_date ON risks(due_date);
CREATE INDEX IF NOT EXISTS idx_risks_level ON risks(level);
CREATE INDEX IF NOT EXISTS idx_risks_status ON risks(status);




