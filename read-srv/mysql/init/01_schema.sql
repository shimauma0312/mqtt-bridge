-- 汎用レコード

CREATE TABLE records (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    device_id   VARCHAR(128)    NOT NULL,
    label       VARCHAR(64)     NOT NULL,
    recorded_at DATETIME(3)     NOT NULL,
    payload     JSON            NOT NULL,
    PRIMARY KEY (id),
    -- 重複排除
    UNIQUE KEY uq_record (device_id, label, recorded_at),
    KEY idx_label_time (label, recorded_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
