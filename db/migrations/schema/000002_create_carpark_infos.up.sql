CREATE TABLE `carpark_infos` (
    `id` int(11) AUTO_INCREMENT NOT NULL,
    `carpark_id` int(11) NOT NULL,
    `lot_type` varchar(5) NOT NULL,
    `lots_available` int(11) NOT NULL,
    `total_lots` int(11) NOT NULL,
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY `uq_idx_carpark_id_and_lot_type` (`carpark_id`, `lot_type`)
);