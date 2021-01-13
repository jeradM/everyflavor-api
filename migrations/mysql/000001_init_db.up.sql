CREATE TABLE IF NOT EXISTS batches
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    batch_size_m bigint(20) unsigned DEFAULT NULL,
    batch_strength smallint(3) unsigned DEFAULT NULL,
    batch_vg_m mediumint(6) unsigned DEFAULT NULL,
    max_vg tinyint(1) DEFAULT 0,
    nic_strength smallint(4) unsigned DEFAULT NULL,
    nic_vg_m mediumint(6) unsigned DEFAULT NULL,
    recipe_id bigint(20) unsigned DEFAULT NULL,
    owner_id bigint(20) unsigned DEFAULT NULL,
    use_nic tinyint(1) DEFAULT 1,
    PRIMARY KEY (id),
    INDEX idx_batches_deleted_at (deleted_at),
    INDEX idx_batches_recipe_id (recipe_id),
    INDEX idx_batches_owner_id (owner_id)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS flavors
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    vendor_id bigint(20) unsigned DEFAULT NULL,
    name varchar(191) DEFAULT NULL,
    aliases longtext DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_flavors_deleted_at (deleted_at),
    INDEX idx_flavors_vendor_id (vendor_id),
    INDEX idx_flavors_name (name)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS flavor_ratings
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    flavor_id bigint(20) unsigned DEFAULT NULL,
    rating bigint(20) unsigned DEFAULT NULL,
    owner_id bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_flavor_ratings_flavor_id (flavor_id),
    INDEX idx_flavor_ratings_owner_id (owner_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS flavor_reviews
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    rating_id bigint(20) unsigned DEFAULT NULL,
    content longtext DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_flavor_reviews_deleted_at (deleted_at),
    INDEX idx_flavor_reviews_rating_id (rating_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS users
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    username varchar(191) NOT NULL,
    email varchar(191) DEFAULT NULL,
    password longtext DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX username (username),
    UNIQUE INDEX email (email),
    INDEX idx_users_deleted_at (deleted_at)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recipes
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    owner_id bigint(20) unsigned DEFAULT NULL,
    current tinyint(1) DEFAULT NULL,
    description longtext DEFAULT NULL,
    public tinyint(1) DEFAULT 0,
    snv tinyint(1) DEFAULT 0,
    steep_days bigint(20) unsigned DEFAULT NULL,
    temp_f mediumint(3) unsigned DEFAULT NULL,
    title varchar(191) DEFAULT NULL,
    uuid varchar(191) DEFAULT NULL,
    version bigint(20) unsigned DEFAULT NULL,
    vg_percent_m mediumint(6) unsigned DEFAULT NULL,
    wip tinyint(1) unsigned default 0,
    PRIMARY KEY (id),
    INDEX idx_recipes_public (public),
    INDEX idx_recipes_title (title),
    INDEX idx_recipes_uuid (uuid),
    INDEX idx_recipes_deleted_at (deleted_at),
    INDEX idx_recipes_owner_id (owner_id)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recipe_collaborators
(
    recipe_id bigint(20) unsigned NOT NULL,
    user_id bigint(20) unsigned NOT NULL,
    PRIMARY KEY (recipe_id,user_id),
    INDEX fk_recipe_collaborators_user (user_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recipe_comments
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    content text DEFAULT NULL,
    recipe_id bigint(20) unsigned DEFAULT NULL,
    owner_id bigint(20) unsigned DEFAULT NULL,
    reply_to_id bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_recipe_comments_deleted_at (deleted_at),
    INDEX idx_recipe_comments_recipe_id (recipe_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recipe_flavors
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    percent_m mediumint(6) unsigned DEFAULT NULL,
    flavor_id bigint(20) unsigned DEFAULT NULL,
    recipe_id bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY (flavor_id, recipe_id),
    INDEX idx_recipe_flavors_flavor_id (flavor_id),
    INDEX idx_recipe_flavors_recipe_id (recipe_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recipe_flavor_substitutions
(
    flavor_id bigint(20) unsigned NOT NULL,
    recipe_flavor_id bigint(20) unsigned NOT NULL,
    PRIMARY KEY (flavor_id, recipe_flavor_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recipe_ratings
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    rating tinyint(1) unsigned DEFAULT NULL,
    recipe_id bigint(20) unsigned DEFAULT NULL,
    owner_id bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_recipe_ratings_recipe_id (recipe_id),
    INDEX idx_recipe_ratings_owner_id (owner_id)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recipe_tags
(
    recipe_id bigint(20) unsigned NOT NULL,
    tag_id bigint(20) unsigned NOT NULL,
    PRIMARY KEY (recipe_id, tag_id),
    INDEX fk_recipe_tag_tag (tag_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS roles
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    authority longtext DEFAULT NULL,
    PRIMARY KEY (id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_roles
(
    user_id bigint(20) unsigned NOT NULL,
    role_id bigint(20) unsigned NOT NULL,
    PRIMARY KEY (user_id,role_id),
    INDEX fk_user_roles_role (role_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS vendors
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    abbreviation varchar(191) DEFAULT NULL,
    name longtext DEFAULT NULL,
    aliases longtext DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_vendors_deleted_at (deleted_at),
    INDEX idx_vendors_abbreviation (abbreviation)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS tags
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NOW(),
    updated_at datetime(3) DEFAULT NOW(),
    deleted_at datetime(3) DEFAULT NULL,
    tag varchar(40) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_tags_tag (tag)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Foreign Keys
ALTER TABLE batches ADD CONSTRAINT fk_batches_owner_id FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE batches ADD CONSTRAINT fk_batches_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE SET NULL;
ALTER TABLE flavors ADD CONSTRAINT fk_flavors_vendor_id FOREIGN KEY (vendor_id) REFERENCES vendors (id) ON DELETE SET NULL;
ALTER TABLE flavor_ratings ADD CONSTRAINT fk_flavor_ratings_flavor_id FOREIGN KEY (flavor_id) REFERENCES flavors (id) ON DELETE CASCADE;
ALTER TABLE flavor_ratings ADD CONSTRAINT fk_flavor_ratings_owner_id FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE flavor_reviews ADD CONSTRAINT fk_flavors_reviews_rating_id FOREIGN KEY (rating_id) REFERENCES flavor_ratings (id) ON DELETE CASCADE;
ALTER TABLE recipes ADD CONSTRAINT fk_recipes_owner_id FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE recipe_collaborators ADD CONSTRAINT fk_recipe_collaborators_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE;
ALTER TABLE recipe_collaborators ADD CONSTRAINT fk_recipes_collaborators_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE recipe_comments ADD CONSTRAINT fk_recipes_comments_owner_id FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE recipe_comments ADD CONSTRAINT fk_recipes_comments_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE;
ALTER TABLE recipe_comments ADD CONSTRAINT fk_recipes_comments_reply_to_id FOREIGN KEY (reply_to_id) REFERENCES recipe_comments (id) ON DELETE CASCADE;
ALTER TABLE recipe_flavors ADD CONSTRAINT fk_recipes_flavors_flavor_id FOREIGN KEY (flavor_id) REFERENCES flavors (id) ON DELETE CASCADE;
ALTER TABLE recipe_flavors ADD CONSTRAINT fk_recipe_flavors_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE;
ALTER TABLE recipe_ratings ADD CONSTRAINT fk_recipe_ratings_owner_id FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE recipe_ratings ADD CONSTRAINT fk_recipe_ratings_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE recipe_tags ADD CONSTRAINT fk_recipe_tags_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE;
ALTER TABLE recipe_tags ADD CONSTRAINT fk_recipe_tags_tag_id FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE;
ALTER TABLE recipe_flavor_substitutions ADD CONSTRAINT fk_recipe_flavor_sub_flavor_id FOREIGN KEY (flavor_id) REFERENCES flavors (id) ON DELETE CASCADE;
ALTER TABLE recipe_flavor_substitutions ADD CONSTRAINT fk_recipe_flavor_sub_rf_id FOREIGN KEY (recipe_flavor_id) REFERENCES recipe_flavors (id) ON DELETE CASCADE;
