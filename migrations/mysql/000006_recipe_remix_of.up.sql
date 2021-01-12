ALTER TABLE recipes
ADD remix_of_id BIGINT(20) UNSIGNED DEFAULT NULL;

ALTER TABLE recipes ADD CONSTRAINT fk_recipes_remix_of_id FOREIGN KEY (remix_of_id) references recipes (id) ON DELETE SET NULL;