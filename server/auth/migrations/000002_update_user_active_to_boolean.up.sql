-- Step 1: Drop any default
ALTER TABLE users
ALTER COLUMN user_active DROP DEFAULT;

-- Step 2: Convert 't' or 1 → true, else false
ALTER TABLE users
ALTER COLUMN user_active TYPE BOOLEAN
USING CASE
    WHEN user_active = '1' OR user_active = 't' THEN true
    ELSE false
END;

-- Step 3: Set a new default
ALTER TABLE users
ALTER COLUMN user_active SET DEFAULT true;
