ALTER TABLE public.users
    ADD COLUMN "fileId" UUID,
    ADD COLUMN "bankAccountName" VARCHAR(32) DEFAULT NULL,
    ADD COLUMN "bankAccountHolder" VARCHAR(32) DEFAULT NULL,
    ADD COLUMN "bankAccountNumber" VARCHAR(32) DEFAULT NULL,


ALTER TABLE public.users
    ADD CONSTRAINT fk_users_file
    FOREIGN KEY ("fileId") REFERENCES public.files(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;
