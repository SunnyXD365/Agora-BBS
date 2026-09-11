#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "=== 1. 分类与公开数据测试 ==="
echo "[GET] 获取板块列表:"
curl -s -X GET "${BASE_URL}/categories" | jq .

echo -e "\n=== 2. 账号注册与登录测试 ==="
echo "[POST] 注册新用户:"
curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gopher_test",
    "password": "password123",
    "email": "gopher@example.com"
  }' | jq .

echo -e "\n[POST] 用户登录 (提取 Token):"
LOGIN_RESP=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gopher_test",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESP | jq -r '.data.token')
echo "获取到的 Token: ${TOKEN}"

echo -e "\n=== 3. 受保护个人信息测试 ==="
echo "[GET] 获取当前登录用户信息:"
curl -s -X GET "${BASE_URL}/users/me" \
  -H "Authorization: Bearer ${TOKEN}" | jq .

echo -e "\n=== 4. 主题帖 CRUD 测试 ==="
echo "[POST] 发布新帖:"
CREATE_TOPIC_RESP=$(curl -s -X POST "${BASE_URL}/topics" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "title": "Go 高并发论坛架构讨论",
    "content": "Slice 1 已经全链路贯通！"
  }')
echo $CREATE_TOPIC_RESP | jq .

TOPIC_ID=$(echo $CREATE_TOPIC_RESP | jq -r '.data.topic_id')

echo -e "\n[GET] 分页拉取主题帖列表:"
curl -s -X GET "${BASE_URL}/topics?category_id=1&page=1&page_size=10" | jq .

echo -e "\n[GET] 查看单个帖子详情 (ID: ${TOPIC_ID}):"
curl -s -X GET "${BASE_URL}/topics/${TOPIC_ID}" | jq .

echo -e "\n=== 5. 帖子回复测试 ==="
echo "[POST] 在帖子下发表回复:"
CREATE_POST_RESP=$(curl -s -X POST "${BASE_URL}/topics/${TOPIC_ID}/posts" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "写得真好，前排挤挤！"
  }')
echo $CREATE_POST_RESP | jq .

echo -e "\n[GET] 查看帖子下的楼层回复列表:"
curl -s -X GET "${BASE_URL}/topics/${TOPIC_ID}/posts?page=1&page_size=10" | jq .

echo -e "\n=== 6. 点赞互动测试 ==="
echo "[POST] 对主题帖点赞:"
curl -s -X POST "${BASE_URL}/likes" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"target_type\": \"topic\",
    \"target_id\": ${TOPIC_ID}
  }" | jq .

echo -e "\n[DELETE] 取消点赞:"
curl -s -X DELETE "${BASE_URL}/likes" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"target_type\": \"topic\",
    \"target_id\": ${TOPIC_ID}
  }" | jq .