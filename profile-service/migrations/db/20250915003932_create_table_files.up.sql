DROP TABLE IF EXISTS public.files;

CREATE TABLE public.files (
    id UUID NOT NULL,
    fileUri VARCHAR(255),
    fileThumbnailUri VARCHAR(255),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT files_pkey PRIMARY KEY (id)
);
