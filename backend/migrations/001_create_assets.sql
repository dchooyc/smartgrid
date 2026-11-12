-- Assets committed from validated CSV uploads. The constraints mirror the
-- backend's validation rules so bad rows cannot be stored by any route.
CREATE TABLE assets (
    asset_id           TEXT PRIMARY KEY CHECK (asset_id <> ''),
    parent_asset_id    TEXT REFERENCES assets (asset_id),
    asset_type         TEXT NOT NULL CHECK (asset_type IN (
                           'SUBSTATION', 'TRANSFORMER', 'LV_BOARD',
                           'SWITCHBOARD', 'SWITCHBOARD_PANEL')),
    asset_name         TEXT NOT NULL CHECK (asset_name <> ''),
    operational_status TEXT NOT NULL CHECK (operational_status IN (
                           'IN_SERVICE', 'MAINTENANCE', 'OUT_OF_SERVICE')),
    voltage_kv         NUMERIC,
    rating_kva         NUMERIC CHECK (rating_kva >= 0),
    manufacturer       TEXT,
    model              TEXT,
    serial_number      TEXT,
    commissioned_date  DATE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Substations are the roots; every other asset has a parent.
    CONSTRAINT assets_only_substations_are_roots
        CHECK ((asset_type = 'SUBSTATION') = (parent_asset_id IS NULL)),
    CONSTRAINT assets_not_own_parent
        CHECK (parent_asset_id <> asset_id)
);

-- For finding the children of an asset, and for foreign-key checks.
CREATE INDEX assets_parent_asset_id_idx ON assets (parent_asset_id);

-- For looking assets up by serial number. Not unique: duplicates are only
-- a warning during validation.
CREATE INDEX assets_serial_number_idx ON assets (serial_number);
