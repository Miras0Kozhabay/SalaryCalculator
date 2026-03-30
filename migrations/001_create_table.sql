CREATE TABLE IF NOT EXISTS calculations (
    id          BIGSERIAL PRIMARY KEY,
    mode        VARCHAR(5)     NOT NULL CHECK (mode IN ('gross', 'net')),
    gross_salary NUMERIC(15,2) NOT NULL,
    net_salary   NUMERIC(15,2) NOT NULL,
    opv          NUMERIC(15,2) NOT NULL,
    vosms        NUMERIC(15,2) NOT NULL,
    ipn          NUMERIC(15,2) NOT NULL,
    so           NUMERIC(15,2) NOT NULL,
    oosms        NUMERIC(15,2) NOT NULL,
    sn           NUMERIC(15,2) NOT NULL,
    employer_total NUMERIC(15,2) NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Composite index for the most common query pattern (ordering by created_at, filtering)
CREATE INDEX idx_created_at_desc ON calculations(created_at DESC);

-- Index for filtering by mode
CREATE INDEX idx_mode ON calculations(mode);

-- Composite index for history queries with both mode and created_at
CREATE INDEX idx_mode_created_at ON calculations(mode, created_at DESC);
 