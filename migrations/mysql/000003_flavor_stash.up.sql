CREATE TABLE IF NOT EXISTS flavor_stashes
(
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NULL,
    updated_at datetime(3) DEFAULT NULL,
    deleted_at datetime(3) DEFAULT NULL,
    flavor_id bigint(20) unsigned NOT NULL,
    owner_id bigint(20) unsigned NOT NULL,
    on_hand_m bigint(20) unsigned,
    density_m mediumint(6) unsigned,
    vg tinyint(1) default 0,
    rating tinyint(1) default null,
    PRIMARY KEY (id),
    INDEX idx_flavor_stashes_deleted_at (deleted_at),
    INDEX idx_flavor_stashes_created_at (created_at),
    INDEX idx_flavor_stashes_owner_id (owner_id),
    UNIQUE INDEX flavor_owner (flavor_id, owner_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE flavor_stashes ADD CONSTRAINT fk_flavor_stashes_flavor_id FOREIGN KEY (flavor_id) REFERENCES flavors (id);
ALTER TABLE flavor_stashes ADD CONSTRAINT fk_flavor_stashes_owner_id FOREIGN KEY (owner_id) REFERENCES users (id);
