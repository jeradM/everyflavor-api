ALTER TABLE batch_flavors DROP FOREIGN KEY fk_batch_flavors_batch_id;
ALTER TABLE batch_flavors DROP FOREIGN KEY fk_batch_flavors_flavor_id;

DROP TABLE IF EXISTS batch_flavors;
