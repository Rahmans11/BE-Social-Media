SELECT avatar
FROM users 
WHERE account_id = 20

select * FROM posts_reactions

SELECT 
    DISTINCT p.id AS post_id,
    p.caption,
    p.image,
    p.created_at AS post_created,
    u.id AS user_id,
    u.first_name,
    u.last_name,
    u.avatar,
    COALESCE(l.like_count, 0) AS total_likes,
    COALESCE(c.comment_count, 0) AS total_comments
FROM follows f
INNER JOIN posts p ON f.followed_id = p.user_id
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS like_count
    FROM posts_reactions
    GROUP BY post_id
) l ON p.id = l.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS comment_count
    FROM posts_reactions
    WHERE comment IS NOT NULL AND comment != ''
    GROUP BY post_id
) c ON p.id = c.post_id
WHERE f.follower_id = 18  
ORDER BY p.created_at DESC;