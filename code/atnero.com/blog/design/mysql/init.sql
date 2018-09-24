drop database blog;

BEGIN;

-- 创建blog数据库
CREATE DATABASE blog;

-- 跳到blog数据库
USE blog;

-- 表 dbversion
CREATE TABLE version(
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,version CHAR(16)
        ,CONSTRAINT pk_version PRIMARY KEY (id)
);

INSERT INTO version(version) VALUES ('0.0.1');

-- 后台管理开关
/*
CREATE TABLE bg_manager_enable (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,enable boolean
        ,CONSTRAINT pk_bg_manager_enable PRIMARY KEY (id)
);

INSERT INTO bg_manager_enable(enable) VALUES (true);
*/

-- 权限选项
CREATE TABLE user_right_item (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) UNIQUE
        ,dsc VARCHAR(500)
        ,enable boolean NOT NULL
        ,CONSTRAINT pk_user_right_item PRIMARY KEY (
                id
        )
);

-- INSERT INTO user_right_item(name, dsc, enable) VALUES ('modify_db', 'right to modify database, limited to superuser', true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'view_others_published_article',
        '查看别人的文章',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'view_others_hidden_article',
        '查看别人隐藏的文章;只有管理员有此权限',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'view_my_published_article',
        '查看自己的文章',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'view_my_hidden_article',
        '查看自己隐藏的文章',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'edit_my_article',
        '修改自己的文章',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'edit_others_article',
        '修改别人的文章;只有管理员有此权限',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'delete_my_article',
        '删除自己的文章',
        true);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'delete_others_article',
        '删除别人的文章;只有管理员有此权限',
        false);
INSERT INTO user_right_item(name, dsc, enable) VALUES (
        'create_article',
        '创建文章',
        true);



-- 权限集
-- 如administrator
-- 如普通用户
CREATE TABLE user_right_set (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) UNIQUE
        ,dsc VARCHAR(500)
        ,CONSTRAINT pk_user_right_set PRIMARY KEY (
                id
        )
);

-- INSERT INTO user_right_set(name, dsc) VALUES ('superuser', 'right to modify database, limited to superuser');
INSERT INTO user_right_set(name, dsc) VALUES (
        'tourist_rightset',
        '游客权限组;只能查看别人发布的文章');
INSERT INTO user_right_set(name, dsc) VALUES (
        'base_user_rightset',
        '普通用户,查看发布的文章，查看自己隐藏的文章，编辑自己的文章，删除自己的文章');

-- 权限集到权限选项一对多映射
CREATE TABLE user_right_set2item_map (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,set_id bigint NOT NULL
        ,item_id bigint NOT NULL
        ,CONSTRAINT pk_user_right_set2itemmap PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_user_right_set2itemmap_setid FOREIGN KEY (set_id) references user_right_set(id) on delete cascade
        ,CONSTRAINT fk_user_right_set2itemmap_itemid FOREIGN KEY (item_id) references user_right_item(id) on delete cascade
);

/* INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'superuser' AND user_right_item.name = 'modify_db' ;
*/
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'tourist_rightset' AND user_right_item.name = 'view_others_published_article';
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'view_others_published_article';
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'view_my_published_article';
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'view_my_hidden_article';
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'edit_my_article';
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'delete_my_article';
INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'base_user_rightset' AND user_right_item.name = 'create_article';


CREATE TABLE default_right_sets(
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) UNIQUE
        ,dsc VARCHAR(500)
        ,right_set_id bigint NOT NULL
        ,CONSTRAINT pk_default_rightsets PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_default_rightsets_rightsetid FOREIGN KEY (right_set_id) references user_right_set(id)
);


-- 用户
CREATE TABLE users(
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) UNIQUE
        ,pwd VARCHAR(40) NOT NULL
        ,email VARCHAR(40) NOT NULL
        ,rightset bigint NOT NULL
        ,reg_time DATETIME NOT null DEFAULT now()
        ,last_time DATETIME NOT null DEFAULT now()
        ,fail_time DATETIME NOT null DEFAULT now()
        ,CONSTRAINT pk_users PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_users_rightset FOREIGN KEY (rightset) references user_right_set(id)
);

/*
INSERT INTO users(name, pwd, email, rightset)
        SELECT 'superuser', md5('superuser'), 'superuser@atnero.com', id FROM user_right_set
                WHERE user_right_set.name = 'superuser';
*/

-- 文章分类
-- 如 原创，转发等
CREATE TABLE article_sorts (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) UNIQUE
        ,CONSTRAINT pk_article_sorts PRIMARY KEY (
                id
        )
);

-- 用户文章类型
-- 如编程开发，数据库，linux之类的
CREATE TABLE article_classes (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) NOT NULL UNIQUE
        ,CONSTRAINT pk_article_classes PRIMARY KEY (
                id
        )
);

-- 用户文章标签
CREATE TABLE article_labels (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,user_id bigint NOT NULL
        ,name VARCHAR(40) NOT NULL
        ,CONSTRAINT pk_article_labels PRIMARY KEY (
                id
        )
        ,CONSTRAINT article_labels_uniquename UNIQUE(user_id, name)
        ,CONSTRAINT fk_article_labels FOREIGN KEY (user_id) references users(id)
);

-- 文章
CREATE TABLE articles (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,user_id bigint NOT NULL
        ,title VARCHAR(100) NOT NULL
        ,content LONGTEXT NOT NULL
        ,sort_id bigint NOT NULL
        ,class_id bigint NOT NULL
        ,published boolean NOT NULL DEFAULT false
        ,create_time DATETIME NOT NULL DEFAULT now()
        ,lastupdate_time DATETIME NOT NULL DEFAULT now()
        ,view_count BIGINT NOT NULL DEFAULT 0
        ,CONSTRAINT pk_articles PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_articles_userid FOREIGN KEY (user_id) references users(id)
        ,CONSTRAINT fk_articles_sortid FOREIGN KEY (sort_id) references article_sorts(id)
        ,CONSTRAINT fk_articles_classid FOREIGN KEY (class_id) references article_classes(id)
);


-- 文章标签
CREATE TABLE article_attached_labels (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        , article_id bigint NOT NULL
        , label_id bigint NOT NULL
        ,CONSTRAINT pk_article_labelattach PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_article_attached_labels_articleid FOREIGN KEY (article_id) references articles(id)
        ,CONSTRAINT fk_article_attached_labels_labelid FOREIGN KEY (label_id) references article_labels(id)
        ,CONSTRAINT article_attached_labels_unique1 UNIQUE (article_id, label_id)
);

-- 回复
CREATE TABLE article_comments (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        , article_id bigint NOT NULL
        , user_name VARCHAR(40) NOT NULL
        , email VARCHAR(40) NOT NULL
        , content VARCHAR(2000) NOT NULL
        , create_time DATETIME NOT null DEFAULT now()
        ,CONSTRAINT pk_article_comments PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_article_comments_articleid FOREIGN KEY (article_id) references articles(id)
);


CREATE TABLE invitation_code(
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,code VARCHAR(40) NOT NULL UNIQUE
        ,used boolean NOT NULL DEFAULT false
        ,create_time DATETIME NOT null DEFAULT now()
        ,expire_time  DATETIME NOT null
        ,CONSTRAINT pk_invitation_code PRIMARY KEY (
                id
        )
);


COMMIT;