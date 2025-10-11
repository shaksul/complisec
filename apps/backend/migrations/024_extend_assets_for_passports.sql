-- Extend assets table with passport fields
-- These fields store technical specifications and passport information for assets

-- Add passport/technical specification fields
ALTER TABLE assets ADD COLUMN IF NOT EXISTS serial_number VARCHAR(255);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS pc_number VARCHAR(100);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS model VARCHAR(255);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS cpu VARCHAR(255);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS ram VARCHAR(100);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS hdd_info TEXT;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS network_card VARCHAR(255);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS optical_drive VARCHAR(255);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS ip_address INET;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS mac_address MACADDR;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS manufacturer VARCHAR(255);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS purchase_year INTEGER;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS warranty_until DATE;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Add indexes for commonly searched fields
CREATE INDEX IF NOT EXISTS idx_assets_serial_number ON assets(serial_number);
CREATE INDEX IF NOT EXISTS idx_assets_model ON assets(model);
CREATE INDEX IF NOT EXISTS idx_assets_ip_address ON assets(ip_address);

-- Add comments
COMMENT ON COLUMN assets.serial_number IS 'Serial number (S/N) of the asset';
COMMENT ON COLUMN assets.pc_number IS 'PC/Computer number assigned by organization';
COMMENT ON COLUMN assets.model IS 'Model name/number of the asset';
COMMENT ON COLUMN assets.cpu IS 'CPU/Processor information';
COMMENT ON COLUMN assets.ram IS 'RAM/Memory information';
COMMENT ON COLUMN assets.hdd_info IS 'HDD/Storage information';
COMMENT ON COLUMN assets.network_card IS 'Network card information';
COMMENT ON COLUMN assets.optical_drive IS 'Optical drive information';
COMMENT ON COLUMN assets.ip_address IS 'IP address of the asset';
COMMENT ON COLUMN assets.mac_address IS 'MAC address of the asset';
COMMENT ON COLUMN assets.manufacturer IS 'Manufacturer of the asset';
COMMENT ON COLUMN assets.purchase_year IS 'Year when asset was purchased';
COMMENT ON COLUMN assets.warranty_until IS 'Date until warranty is valid';
COMMENT ON COLUMN assets.metadata IS 'Additional metadata in JSON format';






