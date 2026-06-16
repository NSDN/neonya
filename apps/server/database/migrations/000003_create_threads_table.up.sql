CREATE TABLE IF NOT EXISTS threads(
  id BIGSERIAL PRIMARY KEY,
  author_id UUID NOT NULL REFERENCES users(uid),
  board_id VARCHAR(100) NOT NULL REFERENCES boards(id),
  title TEXT NOT NULL,
  thread_type VARCHAR(20) NOT NULL CHECK (thread_type IN ('comic', 'article')),
  thumbnail_link TEXT NOT NULL,
  tag VARCHAR(20)[],
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

INSERT INTO threads (
  author_id,
  board_id,
  title,
  thread_type,
  thumbnail_link,
  tag,
  created_at,
  updated_at
) VALUES
  (
    (SELECT uid FROM users WHERE username = 'reimu'),
    'chat',
    '青年在选择职业时的考虑',
    'article',
    '',
    ARRAY['马克思'],
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    (SELECT uid FROM users WHERE username = 'reimu'),
    'chat', 
    '德意志意识形态', 
    'article',
    '',
    ARRAY['马克思', '恩格斯'], 
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP
  ),
  (
    (SELECT uid FROM users WHERE username = 'reimu'),
    'chat', 
    '共产党宣言', 
    'article',
    '',
    ARRAY['马克思', '陈望道'], 
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP
  );
