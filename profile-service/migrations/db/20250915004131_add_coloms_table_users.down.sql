ALTER TABLE public.users DROP CONSTRAINT IF EXISTS fk_users_file;

ALTER TABLE public.users
    DROP COLUMN IF EXISTS "fileId",
    DROP COLUMN IF EXISTS "bankAccountName",
    DROP COLUMN IF EXISTS "bankAccountHolder",
    DROP COLUMN IF EXISTS "bankAccountNumber",
