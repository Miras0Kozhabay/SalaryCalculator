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
 
CREATE INDEX idx_created_at ON calculations(created_at);
CREATE INDEX idx_mode       ON calculations(mode);
 