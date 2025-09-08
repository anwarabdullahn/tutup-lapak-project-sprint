-- CreateTable
CREATE TABLE "profiles" (
    "id" TEXT NOT NULL,
    "user_id" TEXT NOT NULL,
    "file_id" TEXT NULL DEFAULT '',
    "file_uri" TEXT NULL DEFAULT '',
    "file_thumbnail_uri" TEXT NOT NULL DEFAULT '',
    "bank_account_name" TEXT NOT NULL DEFAULT '',
    "bank_account_holder" TEXT NOT NULL DEFAULT '',
    "bank_account_number" TEXT NOT NULL DEFAULT '',

    CONSTRAINT "profiles_pkey" PRIMARY KEY ("id")
);