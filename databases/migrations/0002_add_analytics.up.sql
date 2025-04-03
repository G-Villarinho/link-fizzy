-- migrations/0002_add_analytics.up.sql
CREATE TABLE IF NOT EXISTS link_clicks (
    id TEXT PRIMARY KEY,
    link_id TEXT NOT NULL,
    clicked_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (link_id) REFERENCES links(id)
);

CREATE INDEX IF NOT EXISTS idx_link_clicks_link_id ON link_clicks (link_id);
CREATE INDEX IF NOT EXISTS idx_link_clicks_clicked_at ON link_clicks (clicked_at);