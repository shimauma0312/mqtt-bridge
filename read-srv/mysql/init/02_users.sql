-- Bridge
CREATE USER 'bridge_writer'@'%' IDENTIFIED BY 'bridge_writer';
GRANT INSERT ON bridge.records TO 'bridge_writer'@'%';

-- Read 側
CREATE USER 'reader'@'%' IDENTIFIED BY 'reader';
GRANT SELECT ON bridge.* TO 'reader'@'%';
