CREATE TABLE generated_strings (
    id TEXT PRIMARY KEY,          -- The generated unique string itself (e.g., "sadiq-hikma-masjid-123456789")
    used_at INTEGER,              -- Unix timestamp when the string was first used (NULL if not used)
    used_by TEXT,                 -- Identifier of the user who used the string (NULL if not used)
    created_at INTEGER NOT NULL  -- Unix timestamp when the string was generated and inserted
);

-- Indexes for efficient lookups
CREATE INDEX idx_used_at ON generated_strings (used_at);
CREATE INDEX idx_created_at ON generated_strings (created_at);
-- If we frequently query by user, add an index on used_by
-- CREATE INDEX idx_used_by ON generated_strings (used_by);
