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
CREATE TABLE bg_manager_enable (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,enable boolean
        ,CONSTRAINT pk_bg_manager_enable PRIMARY KEY (id)
);

INSERT INTO bg_manager_enable(enable) VALUES (true);

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

INSERT INTO user_right_item(name, dsc, enable) VALUES ('modify_db', 'right to modify database, limited to superuser', true);

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

INSERT INTO user_right_set(name, dsc) VALUES ('superuser', 'right to modify database, limited to superuser');

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

INSERT INTO user_right_set2item_map(set_id, item_id)
        SELECT user_right_set.id, user_right_item.id FROM user_right_set, user_right_item
                WHERE user_right_set.name = 'superuser' AND user_right_item.name = 'modify_db' ;

-- 用户
CREATE TABLE users(
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) UNIQUE
        ,pwd VARCHAR(40) NOT NULL
        ,email VARCHAR(40) NOT NULL
        ,rightset bigint NOT NULL
        ,reg_time DATETIME NOT null DEFAULT now()
        ,CONSTRAINT pk_users PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_users_rightset FOREIGN KEY (rightset) references user_right_set(id)
);

INSERT INTO users(name, pwd, email, rightset)
        SELECT 'superuser', md5('superuser'), 'superuser@atnero.com', id FROM user_right_set
                WHERE user_right_set.name = 'superuser';

-- 文章分类
-- 如 原创，转发等
CREATE TABLE article_sorts (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,name VARCHAR(40) UNIQUE
        ,dsc VARCHAR(500)
        ,CONSTRAINT pk_article_sorts PRIMARY KEY (
                id
        )
);

-- 用户文章类型
-- 如编程开发，数据库，linux之类的
CREATE TABLE article_classes (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,user_id bigint NOT NULL
        ,name VARCHAR(40) NOT NULL
        ,dsc VARCHAR(500)
        ,CONSTRAINT pk_article_classes PRIMARY KEY (
                id
        )
        ,CONSTRAINT article_classes_uniquename UNIQUE(user_id, name)
        ,CONSTRAINT fk_article_classes FOREIGN KEY (user_id) references users(id)
);

-- 用户文章标签
CREATE TABLE article_labels (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        ,user_id bigint NOT NULL
        ,name VARCHAR(40) NOT NULL
        ,dsc VARCHAR(500)
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
        ,short_dsc VARCHAR(500)
        ,content LONGTEXT NOT NULL
        ,sort_id bigint NOT NULL
        ,class_id bigint NOT NULL
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
);

-- 回复
CREATE TABLE article_comments (
        id bigint NOT NULL AUTO_INCREMENT UNIQUE
        , article_id bigint NOT NULL
        , content VARCHAR(2000) NOT NULL
        ,CONSTRAINT pk_article_comments PRIMARY KEY (
                id
        )
        ,CONSTRAINT fk_article_comments_articleid FOREIGN KEY (article_id) references articles(id)
);

COMMIT;