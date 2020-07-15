# goApi

## Go Api with Postgres backend storing JSON

## Database Schema:
```sql
CREATE TABLE posts (
  ID serial NOT NULL PRIMARY KEY, 
  post_info jsonb NOT NULL
);
```

test data
```sql
INSERT INTO new_posts (post_info) VALUES 
  ('{"name": "Go Interfaces", "desc": "This is another test post from the database...", "tags": {"posted": "07.07.2020", "type": "coding"}}'),
  ('{"name": "Go Pointers", "desc": "This is a test post from the database...", "tags": {"posted": "07.06.2020", "type": "coding"}}');
```
