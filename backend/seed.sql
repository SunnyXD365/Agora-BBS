-- 仅用于本地开发/课堂演示。执行前 API 应已完成全部迁移。
-- PowerShell: Get-Content backend/seed.sql -Raw | docker exec -i agora-postgres psql -U agora_user -d agora_db

INSERT INTO categories(name,slug,description,sort_order,is_active,requires_review) VALUES
  ('综合讨论','general','默认交流与综合话题讨论区',0,TRUE,FALSE),
  ('技术实践','technology','软件工程与技术实践',10,TRUE,FALSE),
  ('公共议题','public-issues','需要更审慎表达的高争议讨论区',20,TRUE,TRUE)
ON CONFLICT(slug) DO UPDATE SET name=EXCLUDED.name,description=EXCLUDED.description,sort_order=EXCLUDED.sort_order,is_active=TRUE,requires_review=EXCLUDED.requires_review;

-- 以下账号的本地演示密码均为 password；生产环境禁止执行本文件。
INSERT INTO users(username,email,password_hash,role,status) VALUES
  ('demo_admin','demo-admin@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','admin','active'),
  ('demo_l0','demo-l0@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l1','demo-l1@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l2','demo-l2@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l3_a','demo-l3-a@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l3_b','demo-l3-b@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l3_c','demo-l3-c@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active')
ON CONFLICT(username) DO UPDATE SET email=EXCLUDED.email,password_hash=EXCLUDED.password_hash,role=EXCLUDED.role,status='active',updated_at=CURRENT_TIMESTAMP;

UPDATE user_trust_profiles tp SET
  unlock_level=CASE u.username WHEN 'demo_l0' THEN 0 WHEN 'demo_l1' THEN 1 WHEN 'demo_l2' THEN 2 ELSE 3 END,
  trust_score=CASE WHEN u.username='demo_l0' THEN 0 ELSE 20 END,
  verified_read_seconds=CASE u.username WHEN 'demo_l0' THEN 0 WHEN 'demo_l1' THEN 60 WHEN 'demo_l2' THEN 180 ELSE 300 END,
  compliant_interactions=CASE WHEN u.username IN ('demo_admin','demo_l3_a','demo_l3_b','demo_l3_c') THEN 5 WHEN u.username='demo_l2' THEN 2 ELSE 0 END,
  audit_probability=0.05,
  updated_at=CURRENT_TIMESTAMP
FROM users u WHERE tp.user_id=u.id AND u.username IN ('demo_admin','demo_l0','demo_l1','demo_l2','demo_l3_a','demo_l3_b','demo_l3_c');

UPDATE user_profiles p SET
  onboarding_statement='本账号用于本地课堂演示，承诺基于事实、尊重他人并说明判断依据。',
  background_tag=CASE u.username WHEN 'demo_l3_a' THEN 'engineering' WHEN 'demo_l3_b' THEN 'humanities' WHEN 'demo_l3_c' THEN 'design' ELSE 'education' END,
  onboarding_status=CASE WHEN u.username='demo_l0' THEN 'not_submitted' ELSE 'approved' END,
  updated_at=CURRENT_TIMESTAMP
FROM users u WHERE p.user_id=u.id AND u.username IN ('demo_admin','demo_l0','demo_l1','demo_l2','demo_l3_a','demo_l3_b','demo_l3_c');

INSERT INTO topics(category_id,user_id,title,content,structured_content,status)
SELECT c.id,u.id,'欢迎来到 Agora-BBS','这是用于本地演示的公开主题。',jsonb_build_object('claim','这是用于本地演示的公开主题。','evidence','项目通过阅读感知、冷静期和匿名盲审改善讨论过程。','uncertainty','欢迎通过结构化回复补充不同视角。'),'published'
FROM categories c CROSS JOIN users u
WHERE c.slug='general' AND u.username='demo_admin'
  AND NOT EXISTS(SELECT 1 FROM topics WHERE title='欢迎来到 Agora-BBS');
