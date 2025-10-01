CREATE TYPE attendance_status AS ENUM ('present', 'absent', 'canceled');

CREATE TABLE IF NOT EXISTS attendance (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  attended attendance_status NOT NULL,
  date    DATE NOT NULL DEFAULT CURRENT_DATE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (date, user_id, event_id)
);
