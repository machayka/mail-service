CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    -- email właściciela formularza, który opłaca usługę obsługi formularz
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
