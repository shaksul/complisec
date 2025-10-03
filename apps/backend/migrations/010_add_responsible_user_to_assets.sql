-- Add responsible_user_id field to assets table
-- This field represents the user who is responsible for the asset
-- while owner_id represents the organization that owns the asset

ALTER TABLE assets 
ADD COLUMN responsible_user_id UUID REFERENCES users(id);

-- Add index for better performance
CREATE INDEX idx_assets_responsible_user_id ON assets(responsible_user_id);

-- Add comment to clarify the difference between owner and responsible user
COMMENT ON COLUMN assets.owner_id IS 'Organization that owns the asset (always organization)';
COMMENT ON COLUMN assets.responsible_user_id IS 'User who is responsible for the asset (can be individual user)';
