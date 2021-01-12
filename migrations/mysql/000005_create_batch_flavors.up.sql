CREATE TABLE IF NOT EXISTS batch_flavors (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    percent_m mediumint(6) unsigned DEFAULT NULL,
    vg tinyint(1) unsigned DEFAULT 0,
    flavor_id bigint(20) unsigned DEFAULT NULL,
    batch_id bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_batch_flavors_flavor_id (flavor_id),
    INDEX idx_batch_flavors_recipe_id (batch_id)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE batch_flavors ADD CONSTRAINT fk_batch_flavors_batch_id FOREIGN KEY (batch_id) REFERENCES batches (id) ON DELETE CASCADE ;
ALTER TABLE batch_flavors ADD CONSTRAINT fk_batch_flavors_flavor_id FOREIGN KEY (flavor_id) REFERENCES flavors (id) ON DELETE CASCADE ;
