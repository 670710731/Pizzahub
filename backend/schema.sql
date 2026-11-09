CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    display_name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'customer'
        CHECK (role IN ('customer', 'admin')),
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'suspended')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX users_email_lower_unique ON users (LOWER(email));

CREATE TABLE categories (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pizzas (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'archived')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX pizzas_published_created_idx
    ON pizzas (created_at DESC, id DESC)
    WHERE status = 'published';

CREATE TABLE pizza_categories (
    pizza_id BIGINT NOT NULL REFERENCES pizzas (id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (pizza_id, category_id)
);

CREATE INDEX pizza_categories_category_idx ON pizza_categories (category_id);

CREATE TABLE pizza_images (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pizza_id BIGINT NOT NULL REFERENCES pizzas (id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    alt_text TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (pizza_id, sort_order)
);

CREATE TABLE pizza_variants (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pizza_id BIGINT NOT NULL REFERENCES pizzas (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    size_label TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'THB'
        CHECK (currency ~ '^[A-Z]{3}$'),
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (pizza_id, size_label)
);

CREATE INDEX pizza_variants_available_idx
    ON pizza_variants (pizza_id, sort_order)
    WHERE is_available = TRUE;

CREATE TABLE reviews (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pizza_id BIGINT NOT NULL REFERENCES pizzas (id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title TEXT NOT NULL DEFAULT '',
    body TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'published'
        CHECK (status IN ('pending', 'published', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (pizza_id, user_id)
);

CREATE INDEX reviews_published_pizza_created_idx
    ON reviews (pizza_id, created_at DESC, id DESC)
    WHERE status = 'published';

CREATE INDEX reviews_user_created_idx ON reviews (user_id, created_at DESC);

CREATE TABLE review_images (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    review_id BIGINT NOT NULL REFERENCES reviews (id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    alt_text TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (review_id, sort_order)
);
