CREATE TABLE IF NOT EXISTS translator.words_ru (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `word__idx` (`word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS translator.words_en (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `word__idx` (`word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS translator.ru_en (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ru_id` bigint unsigned NOT NULL,
  `en_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `ru_id_en_id__idx` (`ru_id`, `en_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS translator.tags (
   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
   `name` varchar(50) NOT NULL,
   PRIMARY KEY (`id`),
   UNIQUE INDEX `name__idx` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;