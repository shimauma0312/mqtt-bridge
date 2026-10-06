-- サンプル

SET NAMES utf8mb4;
SET time_zone = '+00:00';
SET @base = FROM_UNIXTIME(FLOOR(UNIX_TIMESTAMP() / 120) * 120);

-- gps online
INSERT IGNORE INTO records (device_id, label, recorded_at, payload) VALUES
    ('dev-001', 'gps', @base - INTERVAL 300 SECOND, JSON_OBJECT('lat', 35.6812, 'lon', 139.7671)),
    ('dev-001', 'gps', @base - INTERVAL 270 SECOND, JSON_OBJECT('lat', 35.6832, 'lon', 139.7683)),
    ('dev-001', 'gps', @base - INTERVAL 240 SECOND, JSON_OBJECT('lat', 35.6852, 'lon', 139.7695)),
    ('dev-001', 'gps', @base - INTERVAL 210 SECOND, JSON_OBJECT('lat', 35.6872, 'lon', 139.7707)),
    ('dev-001', 'gps', @base - INTERVAL 180 SECOND, JSON_OBJECT('lat', 35.6892, 'lon', 139.7719)),
    ('dev-001', 'gps', @base - INTERVAL 150 SECOND, JSON_OBJECT('lat', 35.6912, 'lon', 139.7731)),
    ('dev-001', 'gps', @base - INTERVAL 120 SECOND, JSON_OBJECT('lat', 35.6932, 'lon', 139.7743)),
    ('dev-001', 'gps', @base - INTERVAL 90 SECOND, JSON_OBJECT('lat', 35.6952, 'lon', 139.7755)),
    ('dev-001', 'gps', @base - INTERVAL 60 SECOND, JSON_OBJECT('lat', 35.6972, 'lon', 139.7767)),
    ('dev-001', 'gps', @base - INTERVAL 30 SECOND, JSON_OBJECT('lat', 35.6992, 'lon', 139.7779)),
    ('dev-001', 'gps', @base - INTERVAL 0 SECOND, JSON_OBJECT('lat', 35.7012, 'lon', 139.7791)),
    ('dev-002', 'gps', @base - INTERVAL 300 SECOND, JSON_OBJECT('lat', 35.658, 'lon', 139.7016)),
    ('dev-002', 'gps', @base - INTERVAL 270 SECOND, JSON_OBJECT('lat', 35.66, 'lon', 139.7028)),
    ('dev-002', 'gps', @base - INTERVAL 240 SECOND, JSON_OBJECT('lat', 35.662, 'lon', 139.704)),
    ('dev-002', 'gps', @base - INTERVAL 210 SECOND, JSON_OBJECT('lat', 35.664, 'lon', 139.7052)),
    ('dev-002', 'gps', @base - INTERVAL 180 SECOND, JSON_OBJECT('lat', 35.666, 'lon', 139.7064)),
    ('dev-002', 'gps', @base - INTERVAL 150 SECOND, JSON_OBJECT('lat', 35.668, 'lon', 139.7076)),
    ('dev-002', 'gps', @base - INTERVAL 120 SECOND, JSON_OBJECT('lat', 35.67, 'lon', 139.7088)),
    ('dev-002', 'gps', @base - INTERVAL 90 SECOND, JSON_OBJECT('lat', 35.672, 'lon', 139.71)),
    ('dev-002', 'gps', @base - INTERVAL 60 SECOND, JSON_OBJECT('lat', 35.674, 'lon', 139.7112)),
    ('dev-002', 'gps', @base - INTERVAL 30 SECOND, JSON_OBJECT('lat', 35.676, 'lon', 139.7124)),
    ('dev-002', 'gps', @base - INTERVAL 0 SECOND, JSON_OBJECT('lat', 35.678, 'lon', 139.7136)),
    ('dev-003', 'gps', @base - INTERVAL 300 SECOND, JSON_OBJECT('lat', 35.7295, 'lon', 139.7109)),
    ('dev-003', 'gps', @base - INTERVAL 270 SECOND, JSON_OBJECT('lat', 35.7315, 'lon', 139.7121)),
    ('dev-003', 'gps', @base - INTERVAL 240 SECOND, JSON_OBJECT('lat', 35.7335, 'lon', 139.7133)),
    ('dev-003', 'gps', @base - INTERVAL 210 SECOND, JSON_OBJECT('lat', 35.7355, 'lon', 139.7145)),
    ('dev-003', 'gps', @base - INTERVAL 180 SECOND, JSON_OBJECT('lat', 35.7375, 'lon', 139.7157)),
    ('dev-003', 'gps', @base - INTERVAL 150 SECOND, JSON_OBJECT('lat', 35.7395, 'lon', 139.7169)),
    ('dev-003', 'gps', @base - INTERVAL 120 SECOND, JSON_OBJECT('lat', 35.7415, 'lon', 139.7181)),
    ('dev-003', 'gps', @base - INTERVAL 90 SECOND, JSON_OBJECT('lat', 35.7435, 'lon', 139.7193)),
    ('dev-003', 'gps', @base - INTERVAL 60 SECOND, JSON_OBJECT('lat', 35.7455, 'lon', 139.7205)),
    ('dev-003', 'gps', @base - INTERVAL 30 SECOND, JSON_OBJECT('lat', 35.7475, 'lon', 139.7217)),
    ('dev-003', 'gps', @base - INTERVAL 0 SECOND, JSON_OBJECT('lat', 35.7495, 'lon', 139.7229));

-- gps offline
INSERT IGNORE INTO records (device_id, label, recorded_at, payload) VALUES
    ('dev-004', 'gps', @base - INTERVAL 14700 SECOND, JSON_OBJECT('lat', 35.6965, 'lon', 139.8147)),
    ('dev-004', 'gps', @base - INTERVAL 14670 SECOND, JSON_OBJECT('lat', 35.6985, 'lon', 139.8135)),
    ('dev-004', 'gps', @base - INTERVAL 14640 SECOND, JSON_OBJECT('lat', 35.7005, 'lon', 139.8123)),
    ('dev-004', 'gps', @base - INTERVAL 14610 SECOND, JSON_OBJECT('lat', 35.7025, 'lon', 139.8111)),
    ('dev-004', 'gps', @base - INTERVAL 14580 SECOND, JSON_OBJECT('lat', 35.7045, 'lon', 139.8099)),
    ('dev-004', 'gps', @base - INTERVAL 14550 SECOND, JSON_OBJECT('lat', 35.7065, 'lon', 139.8087)),
    ('dev-004', 'gps', @base - INTERVAL 14520 SECOND, JSON_OBJECT('lat', 35.7085, 'lon', 139.8075)),
    ('dev-004', 'gps', @base - INTERVAL 14490 SECOND, JSON_OBJECT('lat', 35.7105, 'lon', 139.8063)),
    ('dev-004', 'gps', @base - INTERVAL 14460 SECOND, JSON_OBJECT('lat', 35.7125, 'lon', 139.8051)),
    ('dev-004', 'gps', @base - INTERVAL 14430 SECOND, JSON_OBJECT('lat', 35.7145, 'lon', 139.8039)),
    ('dev-004', 'gps', @base - INTERVAL 14400 SECOND, JSON_OBJECT('lat', 35.7165, 'lon', 139.8027));

-- battery
INSERT IGNORE INTO records (device_id, label, recorded_at, payload) VALUES
    ('dev-001', 'battery', @base - INTERVAL 300 SECOND, JSON_OBJECT('pct', 87)),
    ('dev-001', 'battery', @base - INTERVAL 0 SECOND, JSON_OBJECT('pct', 86)),
    ('dev-002', 'battery', @base - INTERVAL 300 SECOND, JSON_OBJECT('pct', 64)),
    ('dev-002', 'battery', @base - INTERVAL 0 SECOND, JSON_OBJECT('pct', 63)),
    ('dev-003', 'battery', @base - INTERVAL 300 SECOND, JSON_OBJECT('pct', 92)),
    ('dev-003', 'battery', @base - INTERVAL 0 SECOND, JSON_OBJECT('pct', 91)),
    ('dev-004', 'battery', @base - INTERVAL 4 HOUR, JSON_OBJECT('pct', 15));

-- gps walker-01
INSERT IGNORE INTO records (device_id, label, recorded_at, payload) VALUES
    ('walker-01', 'gps', @base - INTERVAL 300 SECOND, JSON_OBJECT('lat', 35.6896, 'lon', 139.7006)),
    ('walker-01', 'gps', @base - INTERVAL 240 SECOND, JSON_OBJECT('lat', 35.6898, 'lon', 139.7007)),
    ('walker-01', 'gps', @base - INTERVAL 180 SECOND, JSON_OBJECT('lat', 35.69, 'lon', 139.7008)),
    ('walker-01', 'gps', @base - INTERVAL 120 SECOND, JSON_OBJECT('lat', 35.6902, 'lon', 139.7009)),
    ('walker-01', 'gps', @base - INTERVAL 60 SECOND, JSON_OBJECT('lat', 35.6904, 'lon', 139.701)),
    ('walker-01', 'gps', @base - INTERVAL 0 SECOND, JSON_OBJECT('lat', 35.6906, 'lon', 139.7011));

-- gps lat/lon 欠落・非数値
INSERT IGNORE INTO records (device_id, label, recorded_at, payload) VALUES
    ('dev-001', 'gps', @base - INTERVAL 10 SECOND, JSON_OBJECT('lon', 139.7791)),
    ('dev-002', 'gps', @base - INTERVAL 40 SECOND, JSON_OBJECT('lat', 'abc', 'lon', 139.7136)),
    ('dev-003', 'gps', @base - INTERVAL 70 SECOND, JSON_OBJECT('lat', NULL, 'lon', 139.7229));
